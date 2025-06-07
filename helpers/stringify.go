package helpers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dev-kas/virtlang-go/v4/shared"
)

const (
	maxLineLength = 80
	itemsPerLine  = 12
)

func Stringify(value shared.RuntimeValue, internal bool) string {
	return stringifyWithVisited(value, internal, make(map[uintptr]bool), 0)
}

func formatObject(obj map[string]*shared.RuntimeValue, internal bool, visited map[uintptr]bool, indentLevel int) string {
	if len(obj) == 0 {
		return "{}"
	}

	indent := strings.Repeat("  ", indentLevel)
	nextIndent := strings.Repeat("  ", indentLevel+1)
	var items []string

	for key, val := range obj {
		item := fmt.Sprintf("%s: %s", key, stringifyWithVisited(*val, internal, visited, indentLevel+1))
		items = append(items, item)
	}

	oneLine := "{ " + strings.Join(items, ", ") + " }"
	if len(oneLine) <= maxLineLength {
		return oneLine
	}

	return "{\n" + nextIndent + strings.Join(items, ",\n"+nextIndent) + "\n" + indent + "}"
}

func isSimpleValue(value shared.RuntimeValue) bool {
	switch value.Type {
	case shared.Number, shared.Boolean, shared.Nil, shared.String:
		return true
	default:
		return false
	}
}

func formatArray(arr []shared.RuntimeValue, internal bool, visited map[uintptr]bool, indentLevel int) string {
	if len(arr) == 0 {
		return "[]"
	}

	indent := strings.Repeat("  ", indentLevel)
	nextIndent := strings.Repeat("  ", indentLevel+1)

	allSimple := true
	for _, item := range arr {
		if !isSimpleValue(item) {
			allSimple = false
			break
		}
	}

	if allSimple {
		singleLine := "[ "
		for i, item := range arr {
			if i > 0 {
				singleLine += ", "
			}
			singleLine += stringifyWithVisited(item, internal, visited, 0)
		}
		singleLine += " ]"

		if len(singleLine) <= maxLineLength {
			return singleLine
		}

		var lines []string
		for i := 0; i < len(arr); i += itemsPerLine {
			end := i + itemsPerLine
			if end > len(arr) {
				end = len(arr)
			}

			var lineItems []string
			for j := i; j < end; j++ {
				lineItems = append(lineItems, stringifyWithVisited(arr[j], internal, visited, 0))
			}

			line := nextIndent + strings.Join(lineItems, ", ")
			if i+itemsPerLine < len(arr) {
				line += ","
			}
			lines = append(lines, line)
		}

		return "[\n" + strings.Join(lines, "\n") + "\n" + indent + "]"
	}

	var items []string
	for _, item := range arr {
		items = append(items, nextIndent+stringifyWithVisited(item, internal, visited, indentLevel+1))
	}
	return "[\n" + strings.Join(items, ",\n") + "\n" + indent + "]"
}

func stringifyWithVisited(value shared.RuntimeValue, internal bool, visited map[uintptr]bool, indentLevel int) string {
	if value.Type == shared.Object || value.Type == shared.Array || value.Type == shared.ClassInstance {
		if value.Value == nil {
			return "null"
		}
		// Only track pointers for non-nil values that can be addressed
		val := reflect.ValueOf(value.Value)
		if val.Kind() == reflect.Ptr && !val.IsNil() {
			ptr := val.Pointer()
			if visited[ptr] {
				return "[Circular]"
			}
			visited[ptr] = true
			defer delete(visited, ptr)
		}
	}

	output := ""
	switch value.Type {
	case shared.String:
		if internal {
			output += fmt.Sprintf("\"%s\"", value.Value.(string))
		} else {
			output += value.Value.(string)
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
		output += formatObject(value.Value.(map[string]*shared.RuntimeValue), true, visited, indentLevel)
	case shared.Array:
		output += formatArray(value.Value.([]shared.RuntimeValue), true, visited, indentLevel)
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
