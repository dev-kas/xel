package object

import (
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var omit = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{Message: "omit() takes exactly 2 arguments"}
	}

	if args[0].Type != shared.Object {
		return nil, &errors.RuntimeError{Message: "omit() expects an object as first argument"}
	}

	obj := args[0]

	if args[1].Type != shared.Array {
		return nil, &errors.RuntimeError{Message: "omit() expects an array as second argument"}
	}

	keys := args[1].Value.([]shared.RuntimeValue)

	excludedKeys := make(map[string]bool)
	for _, item := range keys {
		if item.Type != shared.String {
			return nil, &errors.RuntimeError{Message: "omit() expects a string as key"}
		}
		excludedKeys[item.Value.(string)] = true
	}

	newMap := make(map[string]*shared.RuntimeValue)

	for k, v := range obj.Value.(map[string]*shared.RuntimeValue) {
		if _, found := excludedKeys[k]; !found {
			newMap[k] = v
		}
	}

	objVal := values.MK_OBJECT(newMap)
	return &objVal, nil
})
