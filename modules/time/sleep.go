package time

import (
	"time"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

// Suspends the execution of the current thread for a specified duration in milliseconds
var sleep = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 || args[0].Type != shared.Number {
		return nil, &errors.RuntimeError{Message: "sleep expects a number"}
	}
	time.Sleep(time.Duration(args[0].Value.(float64)) * time.Millisecond)
	nilVal := values.MK_NIL()
	return &nilVal, nil
})
