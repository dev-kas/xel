package time

import (
	"time"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

// Parses a time string in given format and returns the time since the Unix Epoch in milliseconds
var parse = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{Message: "parse() takes 2 arguments"}
	}

	if args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "parse() expects a string as first argument"}
	}

	if args[1].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "parse() expects a string as second argument"}
	}

	timeInput, err := time.Parse(args[1].Value.(string), args[0].Value.(string))
	if err != nil {
		return nil, &errors.RuntimeError{Message: "parse() failed to parse time: " + err.Error()}
	}

	timeParsed := timeInput.UnixNano() / 1e6
	retVal := values.MK_NUMBER(float64(timeParsed))
	return &retVal, nil
})
