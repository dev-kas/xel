package math

import (
	"fmt"
	"math"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

var atan2 = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{
			Message: "atan2() takes exactly two arguments",
		}
	}

	if args[0].Type != shared.Number {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("atan2() takes a number as an argument, but got %s", shared.Stringify(args[0].Type)),
		}
	}

	if args[1].Type != shared.Number {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("atan2() takes a number as an argument, but got %s", shared.Stringify(args[1].Type)),
		}
	}

	result := math.Atan2(args[0].Value.(float64), args[1].Value.(float64))
	res := values.MK_NUMBER(result)
	return &res, nil
})
