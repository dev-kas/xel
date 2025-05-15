package array

import (
	"xel/modules"

	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

func module() (*shared.RuntimeValue, *errors.RuntimeError) {
	mod := values.MK_OBJECT(map[string]*shared.RuntimeValue{
		"push":        &push,
		"pop":         &pop,
		"shift":       &shift,
		"unshift":     &unshift,
		"slice":       &slice,
		"splice":      &splice,
		"fill":        &fill,
		"reverse":     &reverse,
		"sort":        &sort,
		"map":         &map_,
		"filter":      &filter,
		"forEach":     &forEach,
		"reduce":      &reduce,
		"reduceRight": &reduceRight,
		"includes":    &includes,
		"indexOf":     &indexOf,
		"lastIndexOf": &lastIndexOf,
		"find":        &find,
		"findIndex":   &findIndex,
		"every":       &every,
		"some":        &some,
		"join":        &join,
		"concat":      &concat,
		"from":        &from,
		"of":          &of,
		"create":      &create,
	})

	return &mod, nil
}

func init() {
	modules.RegisterNativeModule("xel:array", module)
}
