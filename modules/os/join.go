package os

import (
	"path/filepath"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var join = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) < 2 {
		return nil, &errors.RuntimeError{Message: "join() takes at least 2 arguments"}
	}

	var path string = args[0].Value.(string)

	for i := 1; i < len(args); i++ {
		if args[i].Type != shared.String {
			return nil, &errors.RuntimeError{Message: "join() expects string as argument"}
		}

		path = filepath.Join(path, args[i].Value.(string))
	}

	retVal := values.MK_STRING(path)
	return &retVal, nil
})
