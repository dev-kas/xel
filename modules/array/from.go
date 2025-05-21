package array

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var from = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{Message: "from() takes exactly 1 argument"}
	}

	iterable := args[0]

	var result []shared.RuntimeValue

	switch iterable.Type {
	case shared.String:
		str := iterable.Value.(string)
		str = str[1 : len(str)-1]
		for i := 0; i < len(str); i++ {
			result = append(result, values.MK_STRING(string(str[i])))
		}
	case shared.Array:
		arr := iterable.Value.([]shared.RuntimeValue)
		result = append(result, arr...)
	default:
		return nil, &errors.RuntimeError{Message: fmt.Sprintf("from() expects a string or array, but got %s", shared.Stringify(iterable.Type))}
	}

	res := values.MK_ARRAY(result)
	return &res, nil
})
