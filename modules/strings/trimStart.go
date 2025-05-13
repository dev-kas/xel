package strings

import (
	"fmt"
	"strings"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

var trimStart = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 || args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "trimStart(string) takes one string argument"}
	}
	str := args[0].Value.(string)
	str = str[1 : len(str)-1]
	result := values.MK_STRING(fmt.Sprintf("\"%s\"", strings.TrimLeft(str, " \t\r\n")))
	return &result, nil
})
