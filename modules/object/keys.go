package object

import (
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var keys = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{Message: "keys() takes exactly 1 argument"}
	}

	if args[0].Type != shared.Object {
		return nil, &errors.RuntimeError{Message: "keys() expects an object as argument"}
	}

	obj := args[0]

	var keys []shared.RuntimeValue
	for key := range obj.Value.(map[string]*shared.RuntimeValue) {
		keys = append(keys, values.MK_STRING(key))
	}

	retVal := values.MK_ARRAY(keys)
	return &retVal, nil
})
