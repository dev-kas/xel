package threads

import (
	"time"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var killAll = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 0 {
		return nil, &errors.RuntimeError{
			Message: "killAll() takes no arguments",
		}
	}

	threadsMutex.Lock()
	for _, thread := range threads {
		thread.mu.Lock()
		if thread.Status == RUNNING {
			thread.Status = KILLED
			thread.Cancel()
			thread.FinishedAt = time.Now()
		}
		thread.mu.Unlock()
	}
	threadsMutex.Unlock()

	nilVal := values.MK_NIL()
	return &nilVal, nil
})
