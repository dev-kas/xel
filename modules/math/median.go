package math

import (
	"fmt"
	"sort"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var median = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{
			Message: "median() takes exactly one argument",
		}
	}

	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("median() takes an array as an argument, but got %s", shared.Stringify(args[0].Type)),
		}
	}

	result := []float64{}
	for _, arg := range args[0].Value.([]shared.RuntimeValue) {
		if arg.Type != shared.Number {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("median() takes a number as an argument, but got %s", shared.Stringify(arg.Type)),
			}
		}
		result = append(result, arg.Value.(float64))
	}

	sort.Float64s(result)

	var res float64
	if len(result)%2 == 0 {
		mid1 := result[len(result)/2-1]
		mid2 := result[len(result)/2]
		res = (mid1 + mid2) / 2
	} else {
		res = result[len(result)/2]
	}
	resultValue := values.MK_NUMBER(res)
	return &resultValue, nil
})
