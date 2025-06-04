package strings

import (
	"strings"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var replaceAll = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 3 || args[0].Type != shared.String || args[1].Type != shared.String || args[2].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "replaceAll(string, search, replace)"}
	}

	str := args[0].Value.(string)
	search := args[1].Value.(string)
	replace := args[2].Value.(string)

	result := values.MK_STRING(strings.ReplaceAll(str, search, replace))
	return &result, nil
})
