package strings

import (
	"strings"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var toArray = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 || args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "toArray() function requires exactly 1 argument: a string to convert to an array of characters."}
	}

	str := args[0].Value.(string)
	var result []shared.RuntimeValue
	for _, s := range strings.Split(str, "") {
		result = append(result, values.MK_STRING(s))
	}

	retVal := values.MK_ARRAY(result)
	return &retVal, nil
})
