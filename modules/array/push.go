package array

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

var push = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 2 {
		return nil, &errors.RuntimeError{
			Message: "push() takes at least two arguments",
		}
	}

	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("push() takes an array as an argument, but got %s", shared.Stringify(args[0].Type)),
		}
	}

	arr := args[0]

	arr.Value = append(arr.Value.([]shared.RuntimeValue), args[1:]...)
	return &arr, nil
})
