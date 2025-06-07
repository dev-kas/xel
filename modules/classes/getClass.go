package classes

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var getClass = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{
			Message: "class() takes exactly one argument",
		}
	}

	if args[0].Type != shared.ClassInstance {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("class() takes a class-instance as an argument, but got %s", shared.Stringify(args[0].Type)),
		}
	}

	instance := args[0].Value.(values.ClassInstanceValue)
	class := instance.Class
	return &shared.RuntimeValue{
		Type:  shared.Class,
		Value: class,
	}, nil
})
