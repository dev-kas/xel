package object

import (
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var pick = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{Message: "pick() takes exactly 2 arguments"}
	}

	if args[0].Type != shared.Object {
		return nil, &errors.RuntimeError{Message: "pick() expects an object as first argument"}
	}

	obj := args[0]

	if args[1].Type != shared.Array {
		return nil, &errors.RuntimeError{Message: "pick() expects an array as second argument"}
	}

	keys := args[1]

	newObj := map[string]*shared.RuntimeValue{}

	for _, key := range keys.Value.([]shared.RuntimeValue) {
		if key.Type != shared.String {
			return nil, &errors.RuntimeError{Message: "pick() expects a string as key"}
		}

		newObj[key.Value.(string)] = obj.Value.(map[string]*shared.RuntimeValue)[key.Value.(string)]
	}

	objVal := values.MK_OBJECT(newObj)
	return &objVal, nil
})
