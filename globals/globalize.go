package globals

import (
	"fmt"
	"math"
	"os"

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

	pwd, _ := os.Getwd()
	env.DeclareVar("__dirname__", values.MK_STRING(fmt.Sprintf("\"%s\"", pwd)), true)
	env.DeclareVar("__filename__", values.MK_STRING("\"\""), true)

	env.DeclareVar("proc", Proc(), false)
}
