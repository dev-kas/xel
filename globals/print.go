package globals

import (
	"fmt"
	"github.com/dev-kas/xel/helpers"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var Print = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	out := ""
	for i, arg := range args {
		out += helpers.Stringify(arg, false)

		if len(args)-1 != i {
			out += " "
		}
	}
	fmt.Printf("%s\r\n", out)
	result := values.MK_NIL()
	return &result, nil
})
