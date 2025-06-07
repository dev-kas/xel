package os

import (
	"os"
	user_ "os/user"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var user = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 0 {
		return nil, &errors.RuntimeError{Message: "user() takes exactly 0 arguments"}
	}

	u, err := user_.Current()
	if err != nil {
		return nil, &errors.RuntimeError{Message: err.Error()}
	}

	nameVal := values.MK_STRING(u.Name)
	uidVal := values.MK_STRING(u.Uid)
	gidVal := values.MK_STRING(u.Gid)
	homeVal := values.MK_STRING(u.HomeDir)
	shellVal := values.MK_STRING(os.Getenv("SHELL"))

	retVal := values.MK_OBJECT(map[string]*shared.RuntimeValue{
		"name": &nameVal,
		"uid":  &uidVal,
		"gid":  &gidVal,
		"home": &homeVal,
		"shell": &shellVal,
	})
	return &retVal, nil
})
