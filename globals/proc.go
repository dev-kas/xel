package globals

import (
	xShared "xel/shared"

	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func Proc() shared.RuntimeValue {
	retVal := values.MK_OBJECT(map[string]*shared.RuntimeValue{})

	runtime_version := values.MK_STRING(xShared.RuntimeVersion)
	retVal.Value.(map[string]*shared.RuntimeValue)["runtime_version"] = &runtime_version

	engine_version := values.MK_STRING(xShared.EngineVersion)
	retVal.Value.(map[string]*shared.RuntimeValue)["engine_version"] = &engine_version

	return retVal
}
