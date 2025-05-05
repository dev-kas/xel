package globals

import (
	"fmt"
	"xel/helpers"

	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/shared"
	"github.com/dev-kas/VirtLang-Go/values"
)

var Print = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	out := ""
	for i, arg := range args {
		out += helpers.Stringify(arg, false)

		if len(args)-1 != i {
			out += " "
		}
	}
	fmt.Println(out)
	return nil, nil
})
