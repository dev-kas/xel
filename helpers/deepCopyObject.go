package helpers

import (
	"github.com/dev-kas/virtlang-go/v4/shared"
)

// DeepCopyObject creates a deep copy of a map of RuntimeValue objects
func DeepCopyObject(obj map[string]*shared.RuntimeValue) map[string]*shared.RuntimeValue {
	cloned := make(map[string]*shared.RuntimeValue)
	for k, v := range obj {
		cloned[k] = DeepCopyRuntimeValue(v)
	}
	return cloned
}

// DeepCopyRuntimeValue creates a deep copy of a single RuntimeValue
func DeepCopyRuntimeValue(rtv *shared.RuntimeValue) *shared.RuntimeValue {
	if rtv == nil {
		return nil
	}

	var clonedValue *shared.RuntimeValue
	switch rtv.Type {
	case shared.Nil:
		clonedValue = &shared.RuntimeValue{Type: shared.Nil, Value: nil}

	case shared.Number:
		if num, ok := rtv.Value.(float64); ok {
			clonedValue = &shared.RuntimeValue{Type: shared.Number, Value: num}
		}

	case shared.Boolean:
		if b, ok := rtv.Value.(bool); ok {
			clonedValue = &shared.RuntimeValue{Type: shared.Boolean, Value: b}
		}

	case shared.String:
		if str, ok := rtv.Value.(string); ok {
			clonedValue = &shared.RuntimeValue{Type: shared.String, Value: str}
		}

	case shared.Object:
		if obj, ok := rtv.Value.(map[string]*shared.RuntimeValue); ok {
			clonedValue = &shared.RuntimeValue{Type: shared.Object, Value: DeepCopyObject(obj)}
		}

	case shared.Array:
		if arr, ok := rtv.Value.([]shared.RuntimeValue); ok {
			clonedArr := make([]shared.RuntimeValue, len(arr))
			for i, val := range arr {
				clonedArr[i] = *DeepCopyRuntimeValue(&val)
			}
			clonedValue = &shared.RuntimeValue{Type: shared.Array, Value: clonedArr}
		}

	case shared.Function, shared.NativeFN:
		// For functions, we just copy the reference since they are not meant to be deep copied
		clonedValue = &shared.RuntimeValue{Type: rtv.Type, Value: rtv.Value}

	default:
		// For any other type, just copy the value directly
		clonedValue = &shared.RuntimeValue{Type: rtv.Type, Value: rtv.Value}
	}

	return clonedValue
}
