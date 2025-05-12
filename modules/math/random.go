package math

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

var random = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) > 2 {
		return nil, &errors.RuntimeError{
			Message: "random() takes at most two arguments",
		}
	}

	min := 0.0
	max := 1.0

	if len(args) >= 1 {
		if args[0].Type == shared.Nil {
			min = 0.0
		} else if args[0].Type != shared.Number {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("random() takes a number as the first argument, but got %s", shared.Stringify(args[0].Type)),
			}
		} else {
			min = args[0].Value.(float64)
		}
	}

	if len(args) == 2 {
		if args[1].Type == shared.Nil {
			max = 1.0
		} else if args[1].Type != shared.Number {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("random() takes a number as the second argument, but got %s", shared.Stringify(args[1].Type)),
			}
		} else {
			max = args[1].Value.(float64)
		}
	} else {
		max = min + 1
	}

	if min > max {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("random() expects min <= max, but got min=%f, max=%f", min, max),
		}
	}

	if min == max {
		result := values.MK_NUMBER(min)
		return &result, nil
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := min + rng.Float64()*(max-min)

	res := values.MK_NUMBER(result)
	return &res, nil
})
