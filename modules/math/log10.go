package math

import (
	"fmt"
	"math"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

var log10 = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{
			Message: "log10() takes exactly one argument",
		}
	}

	if args[0].Type != shared.Number {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("log10() takes a number as an argument, but got %s", shared.Stringify(args[0].Type)),
		}
	}

	result := math.Log10(args[0].Value.(float64))
	res := values.MK_NUMBER(result)
	return &res, nil
})
