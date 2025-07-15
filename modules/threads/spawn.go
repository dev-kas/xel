package threads

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dev-kas/xel/helpers"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var spawn = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 1 {
		return nil, &errors.RuntimeError{Message: "spawn() takes at least 1 argument"}
	}
	if args[0].Type != shared.Function {
		return nil, &errors.RuntimeError{Message: fmt.Sprintf("spawn() expects function as argument, got %s", shared.Stringify(args[0].Type))}
	}

	fn := args[0]
	fnArgs := args[1:]

	threadLimiter <- struct{}{}

	ctx, cancel := context.WithCancel(context.Background())
	nilVal := values.MK_NIL()

	thread := &Thread{
		ID:          -1,
		Result:      make(chan *shared.RuntimeValue, 1),
		Error:       make(chan *errors.RuntimeError, 1),
		ReturnValue: &nilVal,
		StartedAt:   time.Now(),
		Status:      RUNNING,
		Cancel:      cancel,
		Context:     ctx,
		mu:          sync.Mutex{},
	}

	threadsMutex.Lock()
	thread.ID = len(threads)
	threads = append(threads, thread)
	threadsMutex.Unlock()

	threadsGroup.Add(1)

	go func(t *Thread, fnVal shared.RuntimeValue) {
		defer func() {
			select {
			case t.Result <- t.ReturnValue:
			default:
			}
			select {
			case t.Error <- nil:
			default:
			}

			<-threadLimiter
			threadsGroup.Done()
		}()

		var result *shared.RuntimeValue
		var _ *errors.RuntimeError

		result, _ = helpers.EvalCancellableFnVal(t.Context, &fnVal, fnArgs, env)

		t.mu.Lock()
		if t.Status == RUNNING {
			t.FinishedAt = time.Now()
			t.Status = FINISHED
		}
		t.ReturnValue = result
		t.mu.Unlock()

	}(thread, fn)

	idVal := values.MK_NUMBER(float64(thread.ID))

	joinVal := values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
		thread.mu.Lock()
		status := thread.Status
		thread.mu.Unlock()

		switch status {
		case KILLED:
			return nil, &errors.RuntimeError{Message: "Thread has been killed"}
		case FINISHED:
			thread.mu.Lock()
			res := thread.ReturnValue
			thread.mu.Unlock()
			return res, nil
		}

		result := <-thread.Result
		err := <-thread.Error

		if err != nil {
			return nil, err
		}
		return result, nil
	})

	timeVal := values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
		thread.mu.Lock()
		finishedAt := thread.FinishedAt
		startedAt := thread.StartedAt
		thread.mu.Unlock()

		if finishedAt.IsZero() {
			return nil, &errors.RuntimeError{Message: "Thread has not finished execution yet"}
		}
		duration := finishedAt.Sub(startedAt).Milliseconds()
		timeVal := values.MK_NUMBER(float64(duration))
		return &timeVal, nil
	})

	killVal := values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
		thread.mu.Lock()
		if thread.Status == RUNNING {
			thread.Status = KILLED
			thread.Cancel()
			thread.FinishedAt = time.Now()
		}
		thread.mu.Unlock()
		return &nilVal, nil
	})

	getStatusVal := values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
		thread.mu.Lock()
		status := thread.Status
		thread.mu.Unlock()
		statusVal := values.MK_STRING(string(status))
		return &statusVal, nil
	})

	getResultVal := values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
		thread.mu.Lock()
		res := thread.ReturnValue
		thread.mu.Unlock()
		return res, nil
	})

	retVal := values.MK_OBJECT(map[string]*shared.RuntimeValue{
		"id":        &idVal,
		"join":      &joinVal,
		"time":      &timeVal,
		"kill":      &killVal,
		"status":    &getStatusVal,
		"getResult": &getResultVal,
	})

	return &retVal, nil
})
