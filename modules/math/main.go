package math

import (
	"xel/modules"

	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

func module() (*shared.RuntimeValue, *errors.RuntimeError) {
	mod := values.MK_OBJECT(map[string]*shared.RuntimeValue{
		"abs":      &abs,
		"round":    &round,
		"floor":    &floor,
		"ceil":     &ceil,
		"sign":     &sign,
		"sin":      &sin,
		"cos":      &cos,
		"tan":      &tan,
		"trunc":    &trunc,
		"sum":      &sum,
		"mean":     &mean,
		"median":   &median,
		"random":   &random,
		"degToRad": &degToRad,
		"radToDeg": &radToDeg,
		"atan":     &atan,
		"atan2":    &atan2,
		"acos":     &acos,
		"asin":     &asin,
		"exp":      &exp,
		"log":      &log,
		"log2":     &log2,
		"log10":    &log10,
		"pow":      &pow,
		"sqrt":     &sqrt,
		"cbrt":     &cbrt,
		"clamp":    &clamp,
		"max":      &max,
		"min":      &min,

		"PI": &pi,
		"E":  &e,
	})

	return &mod, nil
}

func init() {
	modules.RegisterNativeModule("xel:math", module)
}
