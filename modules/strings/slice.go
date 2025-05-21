package strings

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var slice = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 1 || len(args) > 3 {
		return nil, &errors.RuntimeError{Message: "slice() takes 1 to 3 arguments"}
	}
	if args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "slice() expects a string as the first argument"}
	}

	str := args[0].Value.(string)
	str = str[1 : len(str)-1]
	runes := []rune(str)
	length := len(runes)

	start := 0
	end := length

	if len(args) >= 2 {
		if args[1].Type != shared.Number {
			return nil, &errors.RuntimeError{Message: "slice() start must be a number"}
		}
		start = int(args[1].Value.(float64))
		if start < 0 {
			start += length
		}
		if start < 0 {
			start = 0
		}
		if start > length {
			start = length
		}
	}
	if len(args) == 3 {
		if args[2].Type != shared.Number {
			return nil, &errors.RuntimeError{Message: "slice() end must be a number"}
		}
		end = int(args[2].Value.(float64))
		if end < 0 {
			end += length
		}
		if end < 0 {
			end = 0
		}
		if end > length {
			end = length
		}
	}

	if start > end {
		start = end
	}

	result := values.MK_STRING(fmt.Sprintf("\"%s\"", string(runes[start:end])))
	return &result, nil
})
