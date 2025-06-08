package object

import (
	"xel/modules"

	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func module() (*shared.RuntimeValue, *errors.RuntimeError) {
	mod := values.MK_OBJECT(map[string]*shared.RuntimeValue{
		"create":  &create,
		"set":     &set,
		"get":     &get,
		"keys":    &keys,
		"values":  &values_,
		"delete":  &delete_,
		"has":     &has,
		"merge":   &merge,
		"clone":   &clone,
		"pick":    &pick,
		"omit":    &omit,
		"entries": &entries,
		"isEmpty": &isEmpty,
		"equal":   &equal,
	})

	return &mod, nil
}

func init() {
	modules.RegisterNativeModule("xel:object", module)
}
