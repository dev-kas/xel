package math

import (
	"fmt"
	"math"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var asin = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{
			Message: "asin() takes exactly one argument",
		}
	}

	if args[0].Type != shared.Number {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("asin() takes a number as an argument, but got %s", shared.Stringify(args[0].Type)),
		}
	}

	result := math.Asin(args[0].Value.(float64))
	res := values.MK_NUMBER(result)
	return &res, nil
})
