package strings

import (
	"math"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var substr = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 2 || len(args) > 3 || args[0].Type != shared.String || args[1].Type != shared.Number {
		return nil, &errors.RuntimeError{Message: "substr(string, start, [length])"}
	}

	str := args[0].Value.(string)
	start := int(args[1].Value.(float64))
	length := len(str)
	if len(args) == 3 {
		if args[2].Type != shared.Number {
			return nil, &errors.RuntimeError{Message: "substr expects a number for length"}
		}
		length = int(args[2].Value.(float64))
	}

	if start < 0 {
		start = len(str) + start
	}
	if start < 0 {
		start = 0
	}
	end := int(math.Min(float64(start+length), float64(len(str))))

	if start > len(str) {
		start = len(str)
	}
	if end < start {
		end = start
	}

	result := values.MK_STRING(str[start:end])
	return &result, nil
})
