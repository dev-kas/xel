package array

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

var slice = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 1 || len(args) > 3 {
		return nil, &errors.RuntimeError{Message: "slice() expects 1 to 3 arguments"}
	}

	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("slice() expects first argument to be array, got %s", shared.Stringify(args[0].Type)),
		}
	}

	rawArr := args[0].Value.([]shared.RuntimeValue)
	length := len(rawArr)
	start := 0
	end := length

	if len(args) >= 2 {
		if args[1].Type != shared.Number {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("slice() expects start index to be number, got %s", shared.Stringify(args[1].Type)),
			}
		}
		start = int(args[1].Value.(float64))
		if start < 0 {
			start = length + start
		}
	}

	if len(args) == 3 {
		if args[2].Type != shared.Number {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("slice() expects end index to be number, got %s", shared.Stringify(args[2].Type)),
			}
		}
		end = int(args[2].Value.(float64))
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

	newArr := make([]shared.RuntimeValue, end-start)
	copy(newArr, rawArr[start:end])

	return &shared.RuntimeValue{
		Type:  shared.Array,
		Value: newArr,
	}, nil
})
