package classes

import (
	"fmt"
	"reflect"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var instanceOf = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{
			Message: "instanceOf() takes exactly two arguments",
		}
	}

	if args[0].Type != shared.Class {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("instanceOf() takes a class as an argument, but got %s", shared.Stringify(args[0].Type)),
		}
	}

	class := args[0].Value.(values.ClassValue)

	if args[1].Type != shared.ClassInstance {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("instanceOf() takes a class-instance as an argument, but got %s", shared.Stringify(args[1].Type)),
		}
	}

	instance := args[1].Value.(values.ClassInstanceValue)

	res := values.MK_BOOL(reflect.DeepEqual(class, instance.Class))
	return &res, nil
})
