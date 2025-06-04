package globals

import (
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var Len = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{
			Message: "len() takes exactly one argument",
		}
	}

	if args[0].Type != shared.Array && args[0].Type != shared.String {
		return nil, &errors.RuntimeError{
			Message: "len() takes an array or string as an argument",
		}
	}

	result := values.MK_NIL()

	if args[0].Type == shared.String {
		result = values.MK_NUMBER(float64(len(args[0].Value.(string))) - 2)
	} else {
		result = values.MK_NUMBER(float64(len(args[0].Value.([]shared.RuntimeValue))))
	}

	return &result, nil
})
