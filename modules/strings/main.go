package strings

import (
	"github.com/dev-kas/xel/modules"

	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func module() (*shared.RuntimeValue, *errors.RuntimeError) {
	mod := values.MK_OBJECT(map[string]*shared.RuntimeValue{
		"charAt":      &charAt,
		"charCodeAt":  &charCodeAt,
		"includes":    &includes,
		"startsWith":  &startsWith,
		"endsWith":    &endsWith,
		"indexOf":     &indexOf,
		"lastIndexOf": &lastIndexOf,
		"concat":      &concat,
		"slice":       &slice,
		"substring":   &substring,
		"substr":      &substr,
		"lower":       &lower,
		"upper":       &upper,
		"trim":        &trim,
		"trimStart":   &trimStart,
		"trimEnd":     &trimEnd,
		"padStart":    &padStart,
		"padEnd":      &padEnd,
		"repeat":      &repeat,
		"replace":     &replace,
		"replaceAll":  &replaceAll,
		"split":       &split,
		"toArray":     &toArray,
		"format":      &format,
		"equal":       &equal,
	})

	return &mod, nil
}

func init() {
	modules.RegisterNativeModule("xel:strings", module)
}
