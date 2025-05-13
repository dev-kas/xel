package globals

import (
	"fmt"
	"xel/helpers"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
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
	result := values.MK_NIL()
	return &result, nil
})
