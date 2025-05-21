package strings

import (
	"fmt"
	"strings"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var replace = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 3 || args[0].Type != shared.String || args[1].Type != shared.String || args[2].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "replace(string, search, replace)"}
	}

	str := args[0].Value.(string)
	str = str[1 : len(str)-1]
	search := args[1].Value.(string)
	search = search[1 : len(search)-1]
	replace := args[2].Value.(string)
	replace = replace[1 : len(replace)-1]

	result := values.MK_STRING(fmt.Sprintf("\"%s\"", strings.Replace(str, search, replace, 1)))
	return &result, nil
})
