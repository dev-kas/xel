package strings

import (
	"fmt"
	"unicode/utf8"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var charAt = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 || args[0].Type != shared.String || args[1].Type != shared.Number {
		return nil, &errors.RuntimeError{Message: "charAt(str, index) expects a string and a number"}
	}
	str := args[0].Value.(string)
	index := int(args[1].Value.(float64))
	if index < 0 || index >= utf8.RuneCountInString(str) {
		nilVal := values.MK_NIL()
		return &nilVal, nil
	}
	for i, r := range str {
		if i == index {
			strVal := values.MK_STRING(fmt.Sprintf("%c", r))
			return &strVal, nil
		}
	}
	nilVal := values.MK_NIL()
	return &nilVal, nil
})
