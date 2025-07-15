package strings

import (
	"strings"

	"github.com/dev-kas/xel/helpers"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var concat = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	var sb strings.Builder
	for _, arg := range args {
		strVal := helpers.Stringify(arg, false)
		sb.WriteString(strVal)
	}
	result := values.MK_STRING(sb.String())
	return &result, nil
})
