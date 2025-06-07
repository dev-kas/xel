package classes

import (
	"xel/modules"

	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func module() (*shared.RuntimeValue, *errors.RuntimeError) {
	mod := values.MK_OBJECT(map[string]*shared.RuntimeValue{
		"instanceOf": &instanceOf,
		"getClass":   &getClass,
		"methods":    &methods,
		"properties": &properties,
	})

	return &mod, nil
}

func init() {
	modules.RegisterNativeModule("xel:classes", module)
}
