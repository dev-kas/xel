package globals

import (
	"github.com/dev-kas/VirtLang-Go/environment"
)

func Globalize(env *environment.Environment) {
	env.DeclareVar("print", Print, false)
	env.DeclareVar("len", Len, false)
	env.DeclareVar("typeof", Typeof, false)
}
