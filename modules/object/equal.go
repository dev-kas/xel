package object

import (
	"github.com/dev-kas/xel/helpers"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var equal = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 2 {
		return nil, &errors.RuntimeError{Message: "equal() takes at least 2 arguments"}
	}

	if args[0].Type != shared.Object {
		return nil, &errors.RuntimeError{Message: "equal() expects an object as first argument"}
	}

	obj := args[0]

	if args[1].Type != shared.Object {
		return nil, &errors.RuntimeError{Message: "equal() expects an object as second argument"}
	}

	obj2 := args[1]

	if len(obj.Value.(map[string]*shared.RuntimeValue)) != len(obj2.Value.(map[string]*shared.RuntimeValue)) {
		boolVal := values.MK_BOOL(false)
		return &boolVal, nil
	}

	isEqual := helpers.EqualRuntimeValues(&obj, &obj2)
	boolVal := values.MK_BOOL(isEqual)
	return &boolVal, nil
})
