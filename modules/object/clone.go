package object

import (
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var clone = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{Message: "clone() takes exactly 1 argument"}
	}

	if args[0].Type != shared.Object {
		return nil, &errors.RuntimeError{Message: "clone() expects an object as argument"}
	}

	obj := args[0]

	newObj := map[string]*shared.RuntimeValue{}

	for key := range obj.Value.(map[string]*shared.RuntimeValue) {
		newObj[key] = obj.Value.(map[string]*shared.RuntimeValue)[key]
	}

	objVal := values.MK_OBJECT(newObj)
	return &objVal, nil
})
