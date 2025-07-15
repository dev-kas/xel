package globals

import (
	"math"
	"path/filepath"

	xShared "github.com/dev-kas/xel/shared"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func Globalize(env *environment.Environment) {
	env.DeclareVar("true", values.MK_BOOL(true), true)
	env.DeclareVar("false", values.MK_BOOL(false), true)
	env.DeclareVar("nil", values.MK_NIL(), true)
	env.DeclareVar("NaN", values.MK_NUMBER(math.NaN()), true)
	env.DeclareVar("inf", values.MK_NUMBER(math.Inf(1)), true)

	env.DeclareVar("print", Print, false)
	env.DeclareVar("len", Len, false)
	env.DeclareVar("typeof", Typeof, false)
	env.DeclareVar("import", Import, false)

	env.DeclareVar("__dirname__", values.MK_STRING(filepath.Dir(xShared.XelRootDebugger.CurrentFile)), true)
	env.DeclareVar("__filename__", values.MK_STRING(xShared.XelRootDebugger.CurrentFile), true)

	env.DeclareVar("proc", Proc(), false)
	env.DeclareVar("throw", Throw, false)
}
