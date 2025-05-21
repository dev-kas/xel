package strings

import (
	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var startsWith = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 2 || len(args) > 3 || args[0].Type != shared.String || args[1].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "startsWith(str, search, start?) expects 2 or 3 arguments: (string, string, number?)"}
	}
	str := args[0].Value.(string)
	str = str[1 : len(str)-1]
	search := args[1].Value.(string)
	search = search[1 : len(search)-1]
	start := 0
	if len(args) == 3 {
		if args[2].Type != shared.Number {
			return nil, &errors.RuntimeError{Message: "startsWith start index must be a number"}
		}
		start = int(args[2].Value.(float64))
	}
	if start < 0 || start > len(str) || start+len(search) > len(str) {
		boolVal := values.MK_BOOL(false)
		return &boolVal, nil
	}
	boolVal := values.MK_BOOL(str[start:start+len(search)] == search)
	return &boolVal, nil
})
