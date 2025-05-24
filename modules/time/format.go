package time

import (
	"fmt"
	"time"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

// Formats a time in milliseconds since the unix epoch in given format or ISO 8601 by default
var format = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) > 2 || len(args) < 1 {
		return nil, &errors.RuntimeError{Message: "format() takes 1 or 2 arguments"}
	}

	if args[0].Type != shared.Number {
		return nil, &errors.RuntimeError{Message: "format() expects a number as first argument"}
	}

	timeInput := time.Unix(0, int64(args[0].Value.(float64))*1e6)
	timeFormat := time.RFC3339 // ISO 8601: "2006-01-02T15:04:05Z07:00"

	if len(args) == 2 {
		if args[1].Type != shared.String {
			return nil, &errors.RuntimeError{Message: "format() expects a string as second argument"}
		}
		timeFormat = args[1].Value.(string)
		timeFormat = timeFormat[1 : len(timeFormat)-1]
	}

	formatted := timeInput.Format(timeFormat)
	retVal := values.MK_STRING(fmt.Sprintf("%q", formatted))
	return &retVal, nil
})
