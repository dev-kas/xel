package helpers

import (
	"reflect"

	"github.com/dev-kas/virtlang-go/v2/shared"
)

// EqualRuntimeValues compares two VirtLang RuntimeValue objects for equality.
// It performs deep comparison based on the type and value of each RuntimeValue.
//
// Equality rules:
// - Nil: Two nil values are equal
// - Number: Two numbers are equal if their float64 values are equal
// - Boolean: Two booleans are equal if they have the same boolean value
// - String: Two strings are equal if they have the same string content
// - Object: Two objects are equal if they have the same properties with equal values
// - Array: Two arrays are equal if they have the same length and equal elements in the same order
// - Function/NativeFN: Functions are compared by reference (only equal if they are the same function)
func EqualRuntimeValues(a, b *shared.RuntimeValue) bool {
	// Handle nil pointers
	if a == nil || b == nil {
		return a == b
	}

	// Different types are never equal
	if a.Type != b.Type {
		return false
	}

	// Compare based on type
	switch a.Type {
	case shared.Nil:
		// All nil values are equal to each other
		return true

	case shared.Number:
		// Compare number values
		aNum, aOk := a.Value.(float64)
		bNum, bOk := b.Value.(float64)
		if !aOk || !bOk {
			return false
		}
		return aNum == bNum

	case shared.Boolean:
		// Compare boolean values
		aBool, aOk := a.Value.(bool)
		bBool, bOk := b.Value.(bool)
		if !aOk || !bOk {
			return false
		}
		return aBool == bBool

	case shared.String:
		// Compare string values
		aStr, aOk := a.Value.(string)
		bStr, bOk := b.Value.(string)
		if !aOk || !bOk {
			return false
		}
		return aStr == bStr

	case shared.Object:
		// Compare object properties
		aObj, aOk := a.Value.(map[string]shared.RuntimeValue)
		bObj, bOk := b.Value.(map[string]shared.RuntimeValue)
		if !aOk || !bOk {
			return false
		}

		// Check if they have the same number of properties
		if len(aObj) != len(bObj) {
			return false
		}

		// Check if all properties in a exist in b with equal values
		for key, aVal := range aObj {
			bVal, exists := bObj[key]
			if !exists {
				return false
			}

			// Recursively compare property values
			if !EqualRuntimeValues(&aVal, &bVal) {
				return false
			}
		}
		return true

	case shared.Array:
		// Compare array elements
		aArr, aOk := a.Value.([]shared.RuntimeValue)
		bArr, bOk := b.Value.([]shared.RuntimeValue)
		if !aOk || !bOk {
			return false
		}

		// Check if they have the same length
		if len(aArr) != len(bArr) {
			return false
		}

		// Check if all elements are equal
		for i := range aArr {
			if !EqualRuntimeValues(&aArr[i], &bArr[i]) {
				return false
			}
		}
		return true

	case shared.Function, shared.NativeFN:
		// For functions, compare by reference
		// Two functions are equal only if they are the same function
		return reflect.ValueOf(a.Value).Pointer() == reflect.ValueOf(b.Value).Pointer()

	default:
		// For unknown types, use reflect.DeepEqual as a fallback
		return reflect.DeepEqual(a.Value, b.Value)
	}
}
