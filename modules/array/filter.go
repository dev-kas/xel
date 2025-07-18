package array

import (
	"fmt"

	xelHelpers "github.com/dev-kas/xel/helpers"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/helpers"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var filter = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{Message: "filter() takes exactly 2 arguments"}
	}

	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("filter() expects array, got %s", shared.Stringify(args[0].Type)),
		}
	}
	if args[1].Type != shared.Function {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("filter() expects function, got %s", shared.Stringify(args[1].Type)),
		}
	}

	array := args[0].Value.([]shared.RuntimeValue)
	fn := args[1]
	result := []shared.RuntimeValue{}

	for i, val := range array {
		callArgs := []shared.RuntimeValue{
			val,
			values.MK_NUMBER(float64(i)),
			args[0],
		}
		out, err := xelHelpers.EvalFnVal(&fn, callArgs, env)
		if err != nil {
			return nil, err
		}
		if helpers.IsTruthy(out) {
			result = append(result, val)
		}
	}

	return &shared.RuntimeValue{
		Type:  shared.Array,
		Value: result,
	}, nil
})
