package array

import (
	"fmt"

	"xel/helpers"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var includes = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{Message: "includes() takes exactly 2 arguments"}
	}

	if args[0].Type != shared.Array {
		return nil, &errors.RuntimeError{Message: fmt.Sprintf("includes() expects array as first argument, got %s", shared.Stringify(args[0].Type))}
	}

	array := args[0].Value.([]shared.RuntimeValue)
	search := args[1]

	for _, el := range array {
		if helpers.EqualRuntimeValues(&el, &search) {
			trueVal := values.MK_BOOL(true)
			return &trueVal, nil
		}
	}

	falseVal := values.MK_BOOL(false)
	return &falseVal, nil
})
