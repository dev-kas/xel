package os

import (
	"os"
	"path/filepath"

	xShared "xel/shared"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var read = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{Message: "read() takes exactly 1 argument"}
	}
	if args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "read() expects string as first argument"}
	}

	path, err := resolvePath(filepath.Dir(xShared.XelRootDebugger.CurrentFile), args[0].Value.(string))
	if err != nil {
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	retVal := values.MK_STRING(string(data))
	return &retVal, nil
})
