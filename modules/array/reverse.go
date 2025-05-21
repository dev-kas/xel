package array

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var reverse = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{Message: "reverse() expects exactly one argument"}
	}

	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("reverse() expects an array, got %s", shared.Stringify(args[0].Type)),
		}
	}

	original := args[0].Value.([]shared.RuntimeValue)
	n := len(original)
	result := make([]shared.RuntimeValue, n)

	for i := 0; i < n; i++ {
		result[i] = original[n-1-i]
	}

	return &shared.RuntimeValue{
		Type:  shared.Array,
		Value: result,
	}, nil
})
