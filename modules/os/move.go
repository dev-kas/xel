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

func _moveFile(src, dst string) error {
	if err := _copyFile(src, dst); err != nil {
		return err
	}
	return os.Remove(src)
}

func _moveDir(src, dst string) error {
	if err := _copyDir(src, dst); err != nil {
		return err
	}
	return os.RemoveAll(src)
}

var move = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{Message: "move() takes exactly 2 arguments"}
	}
	if args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "move() expects string as first argument"}
	}

	fromPath, err := resolvePath(filepath.Dir(xShared.XelRootDebugger.CurrentFile), args[0].Value.(string))
	if err != nil {
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	toPath, err := resolvePath(filepath.Dir(xShared.XelRootDebugger.CurrentFile), args[1].Value.(string))
	if err != nil {
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	fromStats, err := os.Stat(fromPath)
	if err != nil {
		return nil, &errors.RuntimeError{Message: "Source does not exist"}
	}

	if _, err := os.Stat(toPath); err == nil {
		return nil, &errors.RuntimeError{Message: "Destination already exists"}
	}

	if fromStats.IsDir() {
		err = _moveDir(fromPath, toPath)
	} else {
		err = _moveFile(fromPath, toPath)
	}

	if err != nil {
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	nilVal := values.MK_NIL()
	return &nilVal, nil
})
