package threads

import (
	"context"
	"sync"
	"time"
	"xel/modules"

	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

type ThreadStatus string

const (
	RUNNING  ThreadStatus = "RUNNING"
	FINISHED ThreadStatus = "FINISHED"
	KILLED   ThreadStatus = "KILLED"
)

type Thread struct {
	ID         int
	Result     chan *shared.RuntimeValue
	Error      chan *errors.RuntimeError
	StartedAt  time.Time
	FinishedAt time.Time
	Cancel     context.CancelFunc
	Context    context.Context
	Status     ThreadStatus
}

var threads []*Thread
var threadsMutex sync.Mutex

func module() (*shared.RuntimeValue, *errors.RuntimeError) {
	mod := values.MK_OBJECT(map[string]*shared.RuntimeValue{
		"spawn":      &spawn,
		"waitForAll": &waitForAll,
		// "killAll":    &killAll,
	})

	return &mod, nil
}

func init() {
	modules.RegisterNativeModule("xel:threads", module)
}
