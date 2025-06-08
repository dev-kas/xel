package helpers

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/dev-kas/virtlang-go/v4/shared"
)

type visitedInfo struct {
	ID   int
	Path string
}

const (
	maxLineLength = 80
	itemsPerLine  = 12
	maxDepth     = 25
)

// Stringify converts a RuntimeValue to its string representation.
// The internal parameter controls whether strings should be quoted (true) or not (false).
// When called internally (from within the interpreter), internal should be true to get quoted strings.
// When called for user output (like in print statements), internal should be false to get unquoted strings.
func Stringify(value shared.RuntimeValue, internal bool) string {
	// Always pass internal=true for recursive calls to ensure consistent behavior
	return stringifyWithVisited(value, internal, make(map[uintptr]visitedInfo), 0, "")
}

func formatObject(obj map[string]*shared.RuntimeValue, visited map[uintptr]visitedInfo, indentLevel int, path string) string {
	if len(obj) == 0 {
		return "{}"
	}

	indent := strings.Repeat("  ", indentLevel)
	nextIndent := indent + "  "

	// Check if we can fit everything on one line
	onOneLine := true
	// lineLength tracks the total length of the line
	// Start with 2 for the opening and closing braces
	_ = 2 // lineLength is currently unused but kept for future use
	var items []string

	for key, value := range obj {
		itemStr := key + ": " + stringifyWithVisited(*value, true, visited, indentLevel+1, path+"."+key)
		items = append(items, itemStr)
		if len(itemStr) > maxLineLength {
			onOneLine = false
		}
	}

	if onOneLine {
		return "{" + strings.Join(items, ", ") + "}"
	}

	var lines []string
	for _, item := range items {
		lines = append(lines, nextIndent+item)
	}
	return "{\n" + strings.Join(lines, ",\n") + "\n" + indent + "}"
}

func isSimpleValue(value shared.RuntimeValue) bool {
	switch value.Type {
	case shared.Number, shared.Boolean, shared.Nil, shared.String:
		return true
	default:
		return false
	}
}

func formatArray(arr []shared.RuntimeValue, visited map[uintptr]visitedInfo, indentLevel int, path string) string {
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
			singleLine += stringifyWithVisited(item, true, visited, 0, fmt.Sprintf("%s[%d]", path, i))
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
				lineItems = append(lineItems, stringifyWithVisited(arr[j], true, visited, 0, fmt.Sprintf("%s[%d]", path, j)))
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
	for i, item := range arr {
		items = append(items, nextIndent+stringifyWithVisited(item, true, visited, indentLevel+1, fmt.Sprintf("%s[%d]", path, i)))
	}
	return "[\n" + strings.Join(items, ",\n") + "\n" + indent + "]"
}

// stringifyWithVisited converts a RuntimeValue to a string, tracking visited references to detect cycles.
// The internal parameter controls string quoting behavior - when true, strings are quoted.
func stringifyWithVisited(value shared.RuntimeValue, internal bool, visited map[uintptr]visitedInfo, indentLevel int, path string) string {
	if indentLevel > maxDepth {
		return "[Maximum depth exceeded]"
	}

	if value.Type == shared.Object || value.Type == shared.Array || value.Type == shared.ClassInstance {
		if value.Value == nil {
			return "nil"
		}

		val := reflect.ValueOf(value.Value)
		var ptr uintptr

		switch val.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.UnsafePointer:
			if val.IsNil() {
				break
			}
			
			// For maps, we need to get the pointer to the map header
			if val.Kind() == reflect.Map {
				// This gets the underlying map header pointer in a way that's safe for the garbage collector
				// The map header is guaranteed to be stable for the lifetime of the map
				ptr = uintptr((*[2]uintptr)(unsafe.Pointer(&value.Value))[1])
			} else {
				// For other reference types, we can use the built-in Pointer() method
				ptr = val.Pointer()
			}
			
			if info, exists := visited[ptr]; exists {
				// If we've seen this reference before, show where it was first referenced
				// If the path is the same as the current path, it's a self-reference
				if info.Path == path || info.Path == "" {
					return "[Circular *1]"
				}
				// For nested objects, show the reference and the path where it was first seen
				return fmt.Sprintf("[Circular *%d: %s]", info.ID, info.Path)
			}

			refID := len(visited) + 1
			// Store the current path with this reference
			visited[ptr] = visitedInfo{
				ID:   refID,
				Path: path,
			}
			defer delete(visited, ptr)
		}
	}

	output := ""
	switch value.Type {
	case shared.String:
		str := value.Value.(string)
		if internal {
			// When internal is true, we're being called from within the interpreter
			// and need to add quotes around strings for proper representation
			output += fmt.Sprintf("\"%s\"", str)
		} else {
			// When internal is false, this is for user output (like print statements)
			// so we output the string without quotes
			output += str
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
		output = ""
	case shared.Object:
		output = formatObject(value.Value.(map[string]*shared.RuntimeValue), visited, indentLevel, path)
	case shared.Array:
		output = formatArray(value.Value.([]shared.RuntimeValue), visited, indentLevel, path)
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
