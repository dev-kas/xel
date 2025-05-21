package array

import (
	"fmt"
	"xel/helpers"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var reduce = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 3 {
		return nil, &errors.RuntimeError{Message: "reduce() takes exactly 3 arguments"}
	}

	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{Message: fmt.Sprintf("reduce() expects array as first argument, got %s", shared.Stringify(args[0].Type))}
	}
	if args[1].Type != shared.Function {
		return nil, &errors.RuntimeError{Message: fmt.Sprintf("reduce() expects function as second argument, got %s", shared.Stringify(args[1].Type))}
	}

	array := args[0].Value.([]shared.RuntimeValue)
	reducer := args[1]
	acc := args[2]

	// Clone arr for passing to each call
	copyArr := make([]shared.RuntimeValue, len(array))
	copy(copyArr, array)
	virtArr := shared.RuntimeValue{Type: shared.Array, Value: copyArr}

	for i, val := range array {
		callArgs := []shared.RuntimeValue{
			acc,                          // acc
			val,                          // current value
			values.MK_NUMBER(float64(i)), // index
			virtArr,                      // array
		}
		res, err := helpers.EvalFnVal(&reducer, callArgs, env)
		if err != nil {
			return nil, err
		}
		acc = *res
	}

	return &acc, nil
})
