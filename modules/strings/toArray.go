package strings

import (
	"fmt"
	"strings"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

var toArray = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 || args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "toArray() function requires exactly 1 argument: a string to convert to an array of characters."}
	}

	str := args[0].Value.(string)
	str = str[1 : len(str)-1]
	var result []shared.RuntimeValue
	for _, s := range strings.Split(str, "") {
		result = append(result, values.MK_STRING(fmt.Sprintf("\"%s\"", s)))
	}

	retVal := values.MK_ARRAY(result)
	return &retVal, nil
})
