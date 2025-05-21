package array

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var create = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{Message: "create() takes exactly 2 arguments"}
	}

	if args[0].Type != shared.Number {
		return nil, &errors.RuntimeError{Message: fmt.Sprintf("create() expects a number as the first argument, got %s", shared.Stringify(args[0].Type))}
	}
	if args[1].Type == shared.Array {
		return nil, &errors.RuntimeError{Message: "create() cannot use an array as the second argument"}
	}

	length := int(args[0].Value.(float64))
	value := args[1]

	var result []shared.RuntimeValue
	for i := 0; i < length; i++ {
		result = append(result, value)
	}

	res := values.MK_ARRAY(result)
	return &res, nil
})
