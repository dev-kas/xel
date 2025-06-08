package object

import (
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var set = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) > 3 || len(args) < 2 {
		return nil, &errors.RuntimeError{Message: "set() takes 2 or 3 arguments"}
	}

	if args[0].Type != shared.Object {
		return nil, &errors.RuntimeError{Message: "set() expects an object as first argument"}
	}

	obj := args[0]

	if args[1].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "set() expects a string as second argument"}
	}

	key := args[1].Value.(string)

	if len(args) == 3 {
		obj.Value.(map[string]*shared.RuntimeValue)[key] = &args[2]
	} else {
		nilVal := values.MK_NIL()
		obj.Value.(map[string]*shared.RuntimeValue)[key] = &nilVal
	}

	return &obj, nil
})
