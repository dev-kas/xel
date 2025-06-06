package strings

import (
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var substring = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 2 || len(args) > 3 || args[0].Type != shared.String || args[1].Type != shared.Number {
		return nil, &errors.RuntimeError{Message: "substring(string, start, [end])"}
	}

	str := args[0].Value.(string)
	start := int(args[1].Value.(float64))
	end := len(str)
	if len(args) == 3 {
		if args[2].Type != shared.Number {
			return nil, &errors.RuntimeError{Message: "substring expects a number for end"}
		}
		end = int(args[2].Value.(float64))
	}

	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start > end {
		start, end = end, start
	}
	if start > len(str) {
		start = len(str)
	}
	if end > len(str) {
		end = len(str)
	}

	result := values.MK_STRING(str[start:end])
	return &result, nil
})
