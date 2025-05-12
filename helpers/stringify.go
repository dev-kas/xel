package helpers

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

func Stringify(value shared.RuntimeValue, internal bool) string {
	output := ""
	switch value.Type {
	case shared.String:
		if internal {
			output += value.Value.(string)
		} else {
			output += value.Value.(string)[1 : len(value.Value.(string))-1]
		}
	case shared.Number:
		output += fmt.Sprintf("%g", value.Value.(float64))
	case shared.Boolean:
		if value.Value.(bool) {
			output += "true"
		} else {
			output += "false"
		}
	case shared.Nil:
		output += "nil"
	case shared.Object:
		output += "{\n"
		for key, val := range value.Value.(map[string]*shared.RuntimeValue) {
			output += "  " + key + ": " + Stringify(*val, true) + "\n"
		}
		output += "}"
	case shared.Array:
		output += "["
		for i, val := range value.Value.([]shared.RuntimeValue) {
			output += Stringify(val, true)
			if i != len(value.Value.([]shared.RuntimeValue))-1 {
				output += ", "
			}
		}
		output += "]"
	case shared.Function:
		output += "<function>"
	case shared.NativeFN:
		output += "<native function>"
	case shared.Class:
		output += "<class>"
	case shared.ClassInstance:
		keyValuePair := map[string]*shared.RuntimeValue{}
		ciVal := value.Value.(values.ClassInstanceValue)
		for key, val := range ciVal.Data.Variables {
			if ciVal.Publics[key] {
				keyValuePair[key] = &val
			}
		}
		output += Stringify(values.MK_OBJECT(keyValuePair), true)

	default:
		output += fmt.Sprintf("Unknown - %+v", value.Value)
	}
	return output
}
