package math

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var mean = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{
			Message: "mean() takes exactly one argument",
		}
	}

	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("mean() takes an array as an argument, but got %s", shared.Stringify(args[0].Type)),
		}
	}

	result := 0.0
	for _, arg := range args[0].Value.([]shared.RuntimeValue) {
		if arg.Type != shared.Number {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("mean() takes a number as an argument, but got %s", shared.Stringify(arg.Type)),
			}
		}
		result += arg.Value.(float64)
	}
	res := values.MK_NUMBER(result / float64(len(args[0].Value.([]shared.RuntimeValue))))
	return &res, nil
})
