package os

import (
	"os"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var cwd = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, &errors.RuntimeError{Message: err.Error()}
	}
	retVal := values.MK_STRING(cwd)
	return &retVal, nil
})
