package math

import (
	"fmt"
	"math"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var degToRad = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{
			Message: "degToRad() takes exactly one argument",
		}
	}

	if args[0].Type != shared.Number {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("degToRad() takes a number as an argument, but got %s", shared.Stringify(args[0].Type)),
		}
	}

	result := args[0].Value.(float64) * math.Pi / 180
	res := values.MK_NUMBER(result)
	return &res, nil
})
