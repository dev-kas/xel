package strings

import (
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var equal = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 2 {
		return nil, &errors.RuntimeError{Message: "equal(str, str, ...) expects at least two strings"}
	}
	for _, arg := range args {
		if arg.Type != shared.String {
			return nil, &errors.RuntimeError{Message: "equal(str, str, ...) expects only strings"}
		}
	}

	for i := 1; i < len(args); i++ {
		if args[i].Value.(string) != args[i-1].Value.(string) {
			boolVal := values.MK_BOOL(false)
			return &boolVal, nil
		}
	}

	boolVal := values.MK_BOOL(true)
	return &boolVal, nil
})
