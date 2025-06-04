// format.go
package strings

import (
	"fmt"
	"regexp"
	"strings"
	"xel/helpers"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var format = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 1 || args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "format() requires at least one argument: the format string."}
	}

	formatStr := args[0].Value.(string)

	vCount := strings.Count(formatStr, "%v")
	re := regexp.MustCompile(`%[^v%]`)
	formatStr = re.ReplaceAllStringFunc(formatStr, func(s string) string {
		if s == "%v" {
			return s
		}
		return strings.ReplaceAll(s, "%", "%%")
	})

	var valuesToInsert []interface{}
	for _, arg := range args[1:] {
		valuesToInsert = append(valuesToInsert, helpers.Stringify(arg, false))
	}
	if len(valuesToInsert) > vCount {
		valuesToInsert = valuesToInsert[:vCount]
	}

	formattedStr := fmt.Sprintf(formatStr, valuesToInsert...)
	retVal := values.MK_STRING(formattedStr)
	return &retVal, nil
})
