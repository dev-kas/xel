package helpers

import (
	"github.com/dev-kas/virtlang-go/v2/shared"
)

// IsTruthy determines whether a VirtLang RuntimeValue should be considered
// "truthy" in boolean contexts (like if statements and while loops).
//
// Truthiness rules:
// - Boolean: true is truthy, false is falsy
// - Number: non-zero is truthy, zero is falsy
// - String: non-empty is truthy, empty is falsy
// - Nil: always falsy
// - Object: always truthy (even empty objects)
// - Array: always truthy (even empty arrays)
// - Function: always truthy
// - NativeFN: always truthy
// - ClassInstance: always truthy
// - Class: always truthy
// - Unknown: always truthy
func IsTruthy(value *shared.RuntimeValue) bool {
	if value == nil {
		return false
	}

	switch value.Type {
	case shared.Boolean:
		// Boolean values are truthy if they're true
		return value.Value.(bool)

	case shared.Number:
		// Numbers are truthy if they're non-zero
		num := value.Value.(float64)
		return num != 0

	case shared.String:
		// Strings are truthy if they're non-empty
		str := value.Value.(string)
		return len(str) > 0

	case shared.Nil:
		// Nil is always falsy
		return false

	case shared.Object, shared.Array, shared.Function, shared.NativeFN, shared.ClassInstance, shared.Class:
		// Objects, arrays, and functions are always truthy
		return true

	default:
		// For any unknown types, default to truthy
		return true
	}
}
