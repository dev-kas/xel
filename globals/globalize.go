package globals

import (
	"math"

	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/values"
)

func Globalize(env *environment.Environment) {
	env.DeclareVar("true", values.MK_BOOL(true), true)
	env.DeclareVar("false", values.MK_BOOL(false), true)
	env.DeclareVar("nil", values.MK_NIL(), true)
	env.DeclareVar("NaN", values.MK_NUMBER(int(math.NaN())), true)
	env.DeclareVar("inf", values.MK_NUMBER(int(math.Inf(1))), true)

	env.DeclareVar("print", Print, false)
	env.DeclareVar("len", Len, false)
	env.DeclareVar("typeof", Typeof, false)
}
