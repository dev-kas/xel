package strings

import (
	"strings"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var lastIndexOf = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 2 || args[0].Type != shared.String || args[1].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "lastIndexOf expects (string, substring, ?end)"}
	}
	str := args[0].Value.(string)
	str = str[1 : len(str)-1]
	search := args[1].Value.(string)
	search = search[1 : len(search)-1]
	end := len(str)
	if len(args) > 2 {
		if args[2].Type != shared.Number {
			return nil, &errors.RuntimeError{Message: "lastIndexOf end must be a number"}
		}
		end = int(args[2].Value.(float64))
		end = max(0, min(end, len(str)))
	}
	i := strings.LastIndex(str[:end], search)
	result := values.MK_NUMBER(float64(i))
	return &result, nil
})
