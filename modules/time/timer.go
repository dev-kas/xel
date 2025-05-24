package time

import (
	"time"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var timer = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	var startedAt time.Time
	var elapsed time.Duration
	var running bool

	var nilVal = values.MK_NIL()

	start_impl := values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
		if !running {
			startedAt = time.Now()
			running = true
		}
		return &nilVal, nil
	})
	stop_impl := values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
		if running {
			elapsed += time.Since(startedAt)
			running = false
		}
		return &nilVal, nil
	})
	elapsed_impl := values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
		total := elapsed
		if running {
			total += time.Since(startedAt)
		}
		retVal := values.MK_NUMBER(float64(total.Nanoseconds() / 1e6))
		return &retVal, nil
	})
	retVal := values.MK_OBJECT(map[string]*shared.RuntimeValue{
		"start": &start_impl,
		"stop": &stop_impl,
		"elapsed": &elapsed_impl,
	})
	return &retVal, nil
})
