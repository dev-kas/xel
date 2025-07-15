package os

import (
	"os"
	"path/filepath"

	xShared "github.com/dev-kas/xel/shared"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var mkdir = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{Message: "mkdir() takes exactly 1 argument"}
	}
	if args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "mkdir() expects string as first argument"}
	}

	path, err := resolvePath(filepath.Dir(xShared.XelRootDebugger.CurrentFile), args[0].Value.(string))
	if err != nil {
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	nilVal := values.MK_NIL()

	err = os.Mkdir(path, 0755)
	if err != nil {
		if os.IsExist(err) {
			return &nilVal, nil
		}
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	return &nilVal, nil
})
