package strings

import (
	"strings"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var upper = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 || args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "upper(string) takes one string argument"}
	}
	str := args[0].Value.(string)
	result := values.MK_STRING(strings.ToUpper(str))
	return &result, nil
})
