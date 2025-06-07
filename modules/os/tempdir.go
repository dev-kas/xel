package os

import (
	"os"

	"github.com/dev-kas/virtlang-go/v4/values"
)

var tempdir = values.MK_STRING(os.TempDir())
