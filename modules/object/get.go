package object

import (
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var get = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{Message: "get() takes exactly 2 arguments"}
	}

	if args[0].Type != shared.Object {
		return nil, &errors.RuntimeError{Message: "get() expects an object as first argument"}
	}

	obj := args[0]

	if args[1].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "get() expects a string as second argument"}
	}

	key := args[1].Value.(string)

	if obj.Value.(map[string]*shared.RuntimeValue)[key] == nil {
		nilVal := values.MK_NIL()
		return &nilVal, nil
	}

	return obj.Value.(map[string]*shared.RuntimeValue)[key], nil
})
