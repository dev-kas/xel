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

var stat = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{Message: "stat() takes exactly 1 argument"}
	}
	if args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "stat() expects string as first argument"}
	}

	path, err := resolvePath(filepath.Dir(xShared.XelRootDebugger.CurrentFile), args[0].Value.(string))
	if err != nil {
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	stats, err := os.Stat(path)
	if err != nil {
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	nameVal := values.MK_STRING(stats.Name())
	sizeVal := values.MK_NUMBER(float64(stats.Size()))
	permVal := values.MK_STRING(stats.Mode().Perm().String())
	modTimeVal := values.MK_NUMBER(float64(stats.ModTime().UnixNano()) / 1e6)
	isDirVal := values.MK_BOOL(stats.IsDir())

	retVal := values.MK_OBJECT(map[string]*shared.RuntimeValue{
		"name":    &nameVal,
		"size":    &sizeVal,
		"perm":    &permVal,
		"modTime": &modTimeVal,
		"isDir":   &isDirVal,
	})
	return &retVal, nil
})
