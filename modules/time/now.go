package time

import (
	"time"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

// Returns the time since the Unix Epoch in milliseconds
var now = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	timeNow := time.Now().UnixNano() / 1e6
	retVal := values.MK_NUMBER(float64(timeNow))
	return &retVal, nil
})
