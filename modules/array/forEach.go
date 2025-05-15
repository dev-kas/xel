package array

import (
	"fmt"
	"xel/helpers"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

var forEach = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{Message: "forEach() takes exactly 2 arguments"}
	}

	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("forEach() expects array, got %s", shared.Stringify(args[0].Type)),
		}
	}
	if args[1].Type != shared.Function {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("forEach() expects function, got %s", shared.Stringify(args[1].Type)),
		}
	}

	array := args[0].Value.([]shared.RuntimeValue)
	fn := args[1]

	for i, val := range array {
		callArgs := []shared.RuntimeValue{
			val,
			values.MK_NUMBER(float64(i)),
			args[0],
		}
		_, err := helpers.EvalFnVal(&fn, callArgs, env)
		if err != nil {
			return nil, err
		}
	}

	nilVal := values.MK_NIL()
	return &nilVal, nil
})
