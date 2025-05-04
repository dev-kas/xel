package globals

import (
	"fmt"

	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/shared"
	"github.com/dev-kas/VirtLang-Go/values"
)

var Typeof = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{
			Message: "typeof() takes exactly one argument",
		}
	}
	result := values.MK_STRING(fmt.Sprintf("\"%s\"", shared.Stringify(args[0].Type)))
	return &result, nil
})
