package globals

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var Throw = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{Message: "throw() takes exactly one argument"}
	}

	if args[0].Type != shared.String {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("throw() takes a string as an argument, but got %s", shared.Stringify(args[0].Type)),
		}
	}

	out := args[0].Value.(string)
	return nil, &errors.RuntimeError{Message: out}
})
