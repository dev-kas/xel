package threads

import (
	"context"
	"fmt"
	"time"
	"xel/helpers"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
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

	ctx, cancel := context.WithCancel(context.Background())

	thread := &Thread{
		ID:        len(threads),
		Result:    make(chan *shared.RuntimeValue),
		Error:     make(chan *errors.RuntimeError),
		StartedAt: time.Now(),
		Status:    RUNNING,
		Cancel:    cancel,
		Context:   ctx,
	}

	threadsMutex.Lock()
	threads = append(threads, thread)
	threadsMutex.Unlock()

	go func(t *Thread, fn *shared.RuntimeValue) {
		threadsMutex.Lock()
		result, err := helpers.EvalCancellableFnVal(t.Context, fn, fnArgs, env)
		if t.Status == RUNNING {
			t.FinishedAt = time.Now()
			t.Status = FINISHED
		}
		threadsMutex.Unlock()
		t.Result <- result
		t.Error <- err
	}(thread, &fn)

	idVal := values.MK_NUMBER(float64(thread.ID))
	joinVal := values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
		if thread.Status == KILLED {
			return nil, &errors.RuntimeError{Message: "Thread has been killed"}
		} else if thread.Status == FINISHED {
			return nil, &errors.RuntimeError{Message: "Thread has finished execution"}
		}

		result := <-thread.Result
		err := <-thread.Error
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	timeVal := values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
		if thread.FinishedAt.IsZero() {
			return nil, &errors.RuntimeError{Message: "Thread has not finished execution yet"}
		}
		duration := thread.FinishedAt.Sub(thread.StartedAt).Milliseconds()
		timeVal := values.MK_NUMBER(float64(duration))
		return &timeVal, nil
	})
	killVal := values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
		threadsMutex.Lock()
		if thread.Status == RUNNING {
			thread.Status = KILLED
			thread.Cancel()
			thread.FinishedAt = time.Now()
		}
		threadsMutex.Unlock()
		nilVal := values.MK_NIL()
		return &nilVal, nil
	})
	getStatusVal := values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
		statusVal := values.MK_STRING(fmt.Sprintf("\"%s\"", string(thread.Status)))
		return &statusVal, nil
	})
	retVal := values.MK_OBJECT(map[string]*shared.RuntimeValue{
		"id":     &idVal,
		"join":   &joinVal,
		"time":   &timeVal,
		"kill":   &killVal,
		"status": &getStatusVal,
	})

	return &retVal, nil
})
