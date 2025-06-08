package object

import (
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var create = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) > 1 {
		return nil, &errors.RuntimeError{Message: "create() takes 1 or no arguments"}
	}

	if len(args) == 1 && args[0].Type != shared.Object {
		return nil, &errors.RuntimeError{Message: "create() expects an object as argument"}
	}

	var retVal shared.RuntimeValue

	if len(args) == 0 {
		retVal = values.MK_OBJECT(map[string]*shared.RuntimeValue{})
	} else {
		retVal = args[0]
	}

	return &retVal, nil
})
