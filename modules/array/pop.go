package array

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var pop = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{
			Message: "pop() takes exactly one argument",
		}
	}

	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("pop() takes an array as an argument, but got %s", shared.Stringify(args[0].Type)),
		}
	}

	arr := args[0]
	arr.Value = arr.Value.([]shared.RuntimeValue)[:len(arr.Value.([]shared.RuntimeValue))-1]

	return &arr, nil
})
