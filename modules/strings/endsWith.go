package strings

import (
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var endsWith = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 2 || len(args) > 3 || args[0].Type != shared.String || args[1].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "endsWith(str, search, length?) expects 2 or 3 arguments: (string, string, number?)"}
	}
	str := args[0].Value.(string)
	search := args[1].Value.(string)
	length := len(str)
	if len(args) == 3 {
		if args[2].Type != shared.Number {
			return nil, &errors.RuntimeError{Message: "endsWith length must be a number"}
		}
		length = int(args[2].Value.(float64))
		if length > len(str) {
			length = len(str)
		}
	}
	if length < len(search) {
		boolVal := values.MK_BOOL(false)
		return &boolVal, nil
	}
	boolVal := values.MK_BOOL(str[length-len(search):length] == search)
	return &boolVal, nil
})
