package globals

import (
	"fmt"

	xShared "xel/shared"

	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

func Proc() shared.RuntimeValue {
	retVal := values.MK_OBJECT(map[string]*shared.RuntimeValue{})

	runtime_version := values.MK_STRING(fmt.Sprintf("\"%s\"", xShared.RuntimeVersion))
	retVal.Value.(map[string]*shared.RuntimeValue)["runtime_version"] = &runtime_version

	engine_version := values.MK_STRING(fmt.Sprintf("\"%s\"", xShared.EngineVersion))
	retVal.Value.(map[string]*shared.RuntimeValue)["engine_version"] = &engine_version

	return retVal
}
