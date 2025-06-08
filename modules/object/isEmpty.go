package object

import (
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var isEmpty = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{Message: "isEmpty() takes exactly 1 argument"}
	}

	if args[0].Type != shared.Object {
		return nil, &errors.RuntimeError{Message: "isEmpty() expects an object as argument"}
	}

	obj := args[0]

	isEmpty := true
	if len(obj.Value.(map[string]*shared.RuntimeValue)) > 0 {
		isEmpty = false
	}

	retVal := values.MK_BOOL(isEmpty)
	return &retVal, nil
})
