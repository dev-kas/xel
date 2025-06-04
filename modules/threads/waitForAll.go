package threads

import (
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var waitForAll = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 0 {
		return nil, &errors.RuntimeError{
			Message: "waitForAll() takes no arguments",
		}
	}

	threadsGroup.Wait()

	nilVal := values.MK_NIL()
	return &nilVal, nil
})
