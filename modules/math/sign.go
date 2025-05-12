package math

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

var sign = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{
			Message: "sign() takes exactly one argument",
		}
	}

	if args[0].Type != shared.Number {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("sign() takes a number as an argument, but got %s", shared.Stringify(args[0].Type)),
		}
	}

	result := 0.0
	if args[0].Value.(float64) > 0 {
		result = 1.0
	} else if args[0].Value.(float64) < 0 {
		result = -1.0
	} else {
		result = 0.0
	}

	res := values.MK_NUMBER(result)
	return &res, nil
})
