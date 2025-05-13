package strings

import (
	"fmt"
	"strings"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

var padStart = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 2 || len(args) > 3 || args[0].Type != shared.String || args[1].Type != shared.Number {
		return nil, &errors.RuntimeError{Message: "padStart(string, targetLength, padString?)"}
	}

	str := args[0].Value.(string)
	str = str[1 : len(str)-1]
	targetLen := int(args[1].Value.(float64))
	padStr := " "
	if len(args) == 3 {
		if args[2].Type != shared.String {
			return nil, &errors.RuntimeError{Message: "padString must be a string"}
		}
		padStr = args[2].Value.(string)
		padStr = padStr[1 : len(padStr)-1]
		if padStr == "" {
			return &args[0], nil
		}
	}

	if len(str) >= targetLen {
		return &args[0], nil
	}

	padLen := targetLen - len(str)
	repeats := (padLen + len(padStr) - 1) / len(padStr)
	padding := (strings.Repeat(padStr, repeats))[:padLen]

	result := values.MK_STRING(fmt.Sprintf("\"%s\"", padding+str))
	return &result, nil
})
