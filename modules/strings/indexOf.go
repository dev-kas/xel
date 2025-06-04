package strings

import (
	"strings"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var indexOf = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 2 || args[0].Type != shared.String || args[1].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "indexOf expects (string, substring, ?start)"}
	}
	str := args[0].Value.(string)
	search := args[1].Value.(string)
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
