package os

import (
	"path/filepath"
	"xel/modules"

	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func resolvePath(baseDir, inputPath string) (string, error) {
	if filepath.IsAbs(inputPath) {
		return filepath.Clean(inputPath), nil
	}
	return filepath.Clean(filepath.Join(baseDir, inputPath)), nil
}

func module() (*shared.RuntimeValue, *errors.RuntimeError) {
	mod := values.MK_OBJECT(map[string]*shared.RuntimeValue{
		"cwd":      &cwd,
		"get":      &get,
		"list":     &list,
		"exists":   &exists,
		"read":     &read,
		"write":    &write,
		"remove":   &remove,
		"mkdir":    &mkdir,
		"join":     &join,
		"sep":      &sep,
		"platform": &platform,
		"arch":     &arch,
		"tempdir":  &tempdir,
		"user":     &user,
		"exec":     &exec,
		"copy":     &copy,
		"move":     &move,
		"stat":     &stat,
	})

	return &mod, nil
}

func init() {
	modules.RegisterNativeModule("xel:os", module)
}
