package strings

import (
	"fmt"
	"strings"

	"xel/helpers"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

var concat = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	var sb strings.Builder
	for _, arg := range args {
		strVal := helpers.Stringify(arg, false)
		sb.WriteString(strVal)
	}
	result := values.MK_STRING(fmt.Sprintf("\"%s\"", sb.String()))
	return &result, nil
})
