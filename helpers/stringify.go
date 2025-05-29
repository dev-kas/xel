package helpers

import (
	"fmt"
	"reflect"

	"github.com/dev-kas/virtlang-go/v3/shared"
)

func Stringify(value shared.RuntimeValue, internal bool) string {
	return stringifyWithVisited(value, internal, make(map[uintptr]bool))
}

func stringifyWithVisited(value shared.RuntimeValue, internal bool, visited map[uintptr]bool) string {
	output := ""

	if value.Type == shared.Object || value.Type == shared.Array || value.Type == shared.ClassInstance {
		if value.Value == nil {
			return "null"
		}
		ptr := reflect.ValueOf(value.Value).Pointer()
		if visited[ptr] {
			return "[Circular]"
		}
		visited[ptr] = true
		defer delete(visited, ptr)
	}

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
		output += "{\r\n"
		first := true
		for key, val := range value.Value.(map[string]*shared.RuntimeValue) {
			if !first {
				output += ",\r\n"
			}
			first = false
			output += fmt.Sprintf("\t%s: %s", key, stringifyWithVisited(*val, internal, visited)) + "\r\n"
		}
		output += "\r\n}"
	case shared.Array:
		arr := value.Value.([]shared.RuntimeValue)
		output += "[\r\n"
		for i, val := range arr {
			if i > 0 {
				output += ",\r\n"
			}
			output += "\t" + stringifyWithVisited(val, internal, visited)
		}
		output += "\r\n]"
	case shared.Function:
		output += "<function>"
	case shared.NativeFN:
		output += "<native function>"
	case shared.Class:
		output += "<class>"
	case shared.ClassInstance:
		output += "[ClassInstance]"

	default:
		output += fmt.Sprintf("Unknown - %+v", value.Value)
	}
	return output
}
