package time

import (
	"xel/modules"

	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func module() (*shared.RuntimeValue, *errors.RuntimeError) {
	mod := values.MK_OBJECT(map[string]*shared.RuntimeValue{
		"now":    &now,
		"sleep":  &sleep,
		"format": &format,
		"parse":  &parse,
		"timer":  &timer,
	})

	return &mod, nil
}

func init() {
	modules.RegisterNativeModule("xel:time", module)
}
