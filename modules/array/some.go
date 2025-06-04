package array

import (
	"fmt"

	"xel/helpers"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var some = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{Message: "some() takes exactly 2 arguments"}
	}

	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{Message: fmt.Sprintf("some() expects array as first argument, got %s", shared.Stringify(args[0].Type))}
	}
	if args[1].Type != shared.Function {
		return nil, &errors.RuntimeError{Message: fmt.Sprintf("some() expects function as second argument, got %s", shared.Stringify(args[1].Type))}
	}

	array := args[0].Value.([]shared.RuntimeValue)
	predicate := args[1]

	for i, el := range array {
		callArgs := []shared.RuntimeValue{
			el,
			values.MK_NUMBER(float64(i)),
			args[0],
		}
		res, err := helpers.EvalFnVal(&predicate, callArgs, env)
		if err != nil {
			return nil, err
		}
		if res.Type == shared.Boolean && res.Value.(bool) {
			trueVal := values.MK_BOOL(true)
			return &trueVal, nil
		}
	}

	falseVal := values.MK_BOOL(false)
	return &falseVal, nil
})
