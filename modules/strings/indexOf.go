package strings

import (
	"strings"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var indexOf = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 2 || args[0].Type != shared.String || args[1].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "indexOf expects (string, substring, ?start)"}
	}
	str := args[0].Value.(string)
	str = str[1 : len(str)-1]
	search := args[1].Value.(string)
	search = search[1 : len(search)-1]
	start := 0
	if len(args) > 2 {
		if args[2].Type != shared.Number {
			return nil, &errors.RuntimeError{Message: "indexOf start must be a number"}
		}
		start = int(args[2].Value.(float64))
		if start < 0 || start > len(str) {
			result := values.MK_NUMBER(-1)
			return &result, nil
		}
	}
	i := strings.Index(str[start:], search)
	if i == -1 {
		result := values.MK_NUMBER(-1)
		return &result, nil
	}
	result := values.MK_NUMBER(float64(start + i))
	return &result, nil
})
