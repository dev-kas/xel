package strings

import (
	"fmt"
	"strings"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var repeat = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 || args[0].Type != shared.String || args[1].Type != shared.Number {
		return nil, &errors.RuntimeError{Message: "repeat(string, count)"}
	}

	str := args[0].Value.(string)
	str = str[1 : len(str)-1]
	count := int(args[1].Value.(float64))
	if count < 0 {
		return nil, &errors.RuntimeError{Message: "count must be non-negative"}
	}

	result := values.MK_STRING(fmt.Sprintf("\"%s\"", strings.Repeat(str, count)))
	return &result, nil
})
