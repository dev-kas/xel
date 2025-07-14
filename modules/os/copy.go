package os

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	xShared "xel/shared"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func _copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func _copyDir(srcDir, dstDir string) error {
	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dstDir, relPath)
		if d.IsDir() {
			return os.MkdirAll(targetPath, os.ModePerm)
		}
		return _copyFile(path, targetPath)
	})
}

var copy = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{Message: "copy() takes exactly 2 arguments"}
	}
	if args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "copy() expects string as first argument"}
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
		err = _copyDir(fromPath, toPath)
	} else {
		err = _copyFile(fromPath, toPath)
	}

	if err != nil {
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	nilVal := values.MK_NIL()
	return &nilVal, nil
})
