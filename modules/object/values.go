package object

import (
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var values_ = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{Message: "values() takes exactly 1 argument"}
	}

	if args[0].Type != shared.Object {
		return nil, &errors.RuntimeError{Message: "values() expects an object as argument"}
	}

	obj := args[0]

	var values_ []shared.RuntimeValue
	for key := range obj.Value.(map[string]*shared.RuntimeValue) {
		values_ = append(values_, *obj.Value.(map[string]*shared.RuntimeValue)[key])
	}

	retVal := values.MK_ARRAY(values_)
	return &retVal, nil
})
