package array

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var splice = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 2 {
		return nil, &errors.RuntimeError{Message: "splice() expects at least 2 arguments"}
	}

	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("splice() expects first argument to be array, got %s", shared.Stringify(args[0].Type)),
		}
	}

	original := args[0].Value.([]shared.RuntimeValue)
	length := len(original)

	if args[1].Type != shared.Number {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("splice() expects second argument to be number (start), got %s", shared.Stringify(args[1].Type)),
		}
	}
	start := int(args[1].Value.(float64))
	if start < 0 {
		start = length + start
	}
	if start < 0 {
		start = 0
	}
	if start > length {
		start = length
	}

	deleteCount := length - start
	if len(args) >= 3 {
		if args[2].Type != shared.Number {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("splice() expects third argument to be number (deleteCount), got %s", shared.Stringify(args[2].Type)),
			}
		}
		deleteCount = int(args[2].Value.(float64))
		if deleteCount < 0 {
			deleteCount = 0
		}
		if start+deleteCount > length {
			deleteCount = length - start
		}
	}

	insert := []shared.RuntimeValue{}
	if len(args) > 3 {
		insert = append(insert, args[3:]...)
	}

	modified := append([]shared.RuntimeValue{}, original[:start]...)
	modified = append(modified, insert...)
	modified = append(modified, original[start+deleteCount:]...)

	return &shared.RuntimeValue{
		Type:  shared.Array,
		Value: modified,
	}, nil
})
