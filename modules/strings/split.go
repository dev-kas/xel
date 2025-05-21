package strings

import (
	"fmt"
	"strings"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var split = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 1 || args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "split() function requires at least 1 argument: a string to split."}
	}

	str := args[0].Value.(string)
	str = str[1 : len(str)-1]
	delimiter := ""
	maxSplits := -1
	if len(args) > 1 && args[1].Type == shared.String {
		delimiter = args[1].Value.(string)
		delimiter = delimiter[1 : len(delimiter)-1]
	}
	if len(args) > 2 && args[2].Type == shared.Number {
		maxSplits = int(args[2].Value.(float64))
	}

	var splitStr []shared.RuntimeValue
	if maxSplits == -1 {
		for _, s := range strings.Split(str, delimiter) {
			splitStr = append(splitStr, values.MK_STRING(fmt.Sprintf("\"%s\"", s)))
		}
	} else {
		for _, s := range strings.SplitN(str, delimiter, maxSplits) {
			splitStr = append(splitStr, values.MK_STRING(fmt.Sprintf("\"%s\"", s)))
		}
	}

	retVal := values.MK_ARRAY(splitStr)
	return &retVal, nil
})
