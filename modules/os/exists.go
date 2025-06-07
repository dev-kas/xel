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

var exists = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{Message: "exists() takes exactly 1 argument"}
	}
	if args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "exists() expects string as first argument"}
	}

	path, err := resolvePath(filepath.Dir(xShared.XelRootDebugger.CurrentFile), args[0].Value.(string))
	if err != nil {
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			retVal := values.MK_BOOL(false)
			return &retVal, nil
		}
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	retVal := values.MK_BOOL(true)
	return &retVal, nil
})
