package globals

import (
	"fmt"

	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/shared"
	"github.com/dev-kas/VirtLang-Go/values"
)

var Print = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	out := ""
	for i, arg := range args {
		switch arg.Type {
		case shared.String:
			out += arg.Value.(string)[1 : len(arg.Value.(string))-1]
		case shared.Number:
			out += fmt.Sprintf("%d", arg.Value)
		default:
			out += fmt.Sprintf("%v", arg.Value)
		}

		if len(args)-1 != i {
			out += " "
		}
	}
	fmt.Println(out)
	return nil, nil
})
