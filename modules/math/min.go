package math

import (
	"fmt"
	"math"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var min = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 2 {
		return nil, &errors.RuntimeError{
			Message: "min() takes at least two arguments",
		}
	}

	for _, arg := range args {
		if arg.Type != shared.Number {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("min() takes a number as an argument, but got %s", shared.Stringify(arg.Type)),
			}
		}
	}

	result := args[0].Value.(float64)
	for _, arg := range args[1:] {
		result = math.Min(result, arg.Value.(float64))
	}
	res := values.MK_NUMBER(result)
	return &res, nil
})
