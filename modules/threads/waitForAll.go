package threads

import (
	"runtime"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

var waitForAll = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 0 {
		return nil, &errors.RuntimeError{
			Message: "waitForAll() takes no arguments",
		}
	}

	threadsMutex.Lock()
	defer threadsMutex.Unlock()

	for {
		allDone := true
		for _, thread := range threads {
			if thread.Status == RUNNING {
				allDone = false
				break
			}
		}
		if allDone {
			break
		}
		runtime.Gosched()
	}

	nilVal := values.MK_NIL()
	return &nilVal, nil
})
