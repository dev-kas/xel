package array

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var fill = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 2 || len(args) > 4 {
		return nil, &errors.RuntimeError{Message: "fill() expects 2 to 4 arguments"}
	}

	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("fill() expects first argument to be array, got %s", shared.Stringify(args[0].Type)),
		}
	}

	original := args[0].Value.([]shared.RuntimeValue)
	length := len(original)
	value := args[1]

	start := 0
	end := length

	if len(args) >= 3 {
		if args[2].Type != shared.Number {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("fill() expects start to be number, got %s", shared.Stringify(args[2].Type)),
			}
		}
		start = int(args[2].Value.(float64))
		if start < 0 {
			start = length + start
		}
	}

	if len(args) == 4 {
		if args[3].Type != shared.Number {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("fill() expects end to be number, got %s", shared.Stringify(args[3].Type)),
			}
		}
		end = int(args[3].Value.(float64))
		if end < 0 {
			end = length + end
		}
	}

	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	if start > end {
		start = end
	}

	result := make([]shared.RuntimeValue, length)
	copy(result, original)

	for i := start; i < end; i++ {
		result[i] = value
	}

	return &shared.RuntimeValue{
		Type:  shared.Array,
		Value: result,
	}, nil
})
