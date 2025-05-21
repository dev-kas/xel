package array

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var concat = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{Message: "concat() takes exactly 2 arguments"}
	}

	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{Message: fmt.Sprintf("concat() expects array as first argument, got %s", shared.Stringify(args[0].Type))}
	}
	if args[1].Type != shared.Array {
		return nil, &errors.RuntimeError{Message: fmt.Sprintf("concat() expects array as second argument, got %s", shared.Stringify(args[1].Type))}
	}

	arr1 := args[0].Value.([]shared.RuntimeValue)
	arr2 := args[1].Value.([]shared.RuntimeValue)

	concatenated := append(arr1, arr2...)

	res := values.MK_ARRAY(concatenated)
	return &res, nil
})
