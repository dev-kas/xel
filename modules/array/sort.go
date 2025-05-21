package array

import (
	"fmt"
	sort_ "sort"
	"xel/helpers"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var sort = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 1 || len(args) > 2 {
		return nil, &errors.RuntimeError{Message: "sort() takes 1 or 2 arguments"}
	}
	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("sort() expects array, got %s", shared.Stringify(args[0].Type)),
		}
	}

	array := args[0].Value.([]shared.RuntimeValue)
	result := make([]shared.RuntimeValue, len(array))
	copy(result, array)

	if len(args) == 2 {
		comparator := args[1]
		if comparator.Type != shared.Function {
			return nil, &errors.RuntimeError{Message: "sort() second argument must be a function"}
		}

		sort_.Slice(result, func(i, j int) bool {
			out, err := helpers.EvalFnVal(&comparator, []shared.RuntimeValue{result[i], result[j]}, env)
			if err != nil {
				return false
			}
			if out.Type != shared.Number {
				return false
			}
			return out.Value.(float64) < 0
		})
	} else {
		sort_.Slice(result, func(i, j int) bool {
			a := result[i]
			b := result[j]
			if a.Type != shared.Number || b.Type != shared.Number {
				return false
			}
			return a.Value.(float64) < b.Value.(float64)
		})
	}

	return &shared.RuntimeValue{
		Type:  shared.Array,
		Value: result,
	}, nil
})
