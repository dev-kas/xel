package os

import (
	"runtime"

	"github.com/dev-kas/virtlang-go/v4/values"
)

var platform = values.MK_STRING(runtime.GOOS)
