package os

import (
	"os"
	"path/filepath"

	"github.com/dev-kas/xel/helpers"

	xShared "github.com/dev-kas/xel/shared"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var write = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{Message: "write() takes exactly 2 arguments"}
	}
	if args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "write() expects string as first argument"}
	}

	path, err := resolvePath(filepath.Dir(xShared.XelRootDebugger.CurrentFile), args[0].Value.(string))
	if err != nil {
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	err = os.WriteFile(path, []byte(helpers.Stringify(args[1], false)), 0644)
	if err != nil {
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	nilVal := values.MK_NIL()
	return &nilVal, nil
})
