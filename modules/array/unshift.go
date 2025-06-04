package array

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var unshift = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 2 {
		return nil, &errors.RuntimeError{
			Message: "unshift() takes at least two arguments",
		}
	}

	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("unshift() takes an array as an argument, but got %s", shared.Stringify(args[0].Type)),
		}
	}

	arr := args[0]

	arr.Value = append(args[1:], arr.Value.([]shared.RuntimeValue)...)
	return &arr, nil
})
