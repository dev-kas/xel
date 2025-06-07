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

var list = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{Message: "list() takes exactly 1 argument"}
	}
	if args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "list() expects string as first argument"}
	}

	path, err := resolvePath(filepath.Dir(xShared.XelRootDebugger.CurrentFile), args[0].Value.(string))
	if err != nil {
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	arr := []shared.RuntimeValue{}
	for _, file := range files {
		arr = append(arr, values.MK_STRING(file.Name()))
	}

	retVal := values.MK_ARRAY(arr)
	return &retVal, nil
})
