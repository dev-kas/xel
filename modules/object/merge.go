package object

import (
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var merge = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{Message: "merge() takes exactly 2 arguments"}
	}

	if args[0].Type != shared.Object {
		return nil, &errors.RuntimeError{Message: "merge() expects an object as first argument"}
	}

	obj := args[0]

	if args[1].Type != shared.Object {
		return nil, &errors.RuntimeError{Message: "merge() expects an object as second argument"}
	}

	obj2 := args[1]

	newObj := map[string]*shared.RuntimeValue{}

	for key := range obj.Value.(map[string]*shared.RuntimeValue) {
		newObj[key] = obj.Value.(map[string]*shared.RuntimeValue)[key]
	}

	for key := range obj2.Value.(map[string]*shared.RuntimeValue) {
		newObj[key] = obj2.Value.(map[string]*shared.RuntimeValue)[key]
	}

	objVal := values.MK_OBJECT(newObj)
	return &objVal, nil
})
