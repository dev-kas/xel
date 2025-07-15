package array

import (
	"fmt"

	"github.com/dev-kas/xel/helpers"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var join = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 1 || len(args) > 2 {
		return nil, &errors.RuntimeError{Message: "join() takes 1 or 2 arguments"}
	}

	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{Message: fmt.Sprintf("join() expects array as first argument, got %s", shared.Stringify(args[0].Type))}
	}

	separator := ", "
	if len(args) == 2 {
		if args[1].Type != shared.String {
			return nil, &errors.RuntimeError{Message: fmt.Sprintf("join() expects string as second argument, got %s", shared.Stringify(args[1].Type))}
		}

		separator = args[1].Value.(string)
	}

	array := args[0].Value.([]shared.RuntimeValue)
	var result string

	for i, el := range array {
		strValue := helpers.Stringify(el, false)
		if i > 0 {
			result += separator
		}
		result += strValue
	}

	res := values.MK_STRING(result)
	return &res, nil
})
