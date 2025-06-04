package threads

import (
	"context"
	"runtime"
	"sync"
	"time"
	"xel/modules"

	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

type ThreadStatus string

const (
	RUNNING  ThreadStatus = "RUNNING"
	FINISHED ThreadStatus = "FINISHED"
	KILLED   ThreadStatus = "KILLED"
)

type Thread struct {
	ID          int
	mu          sync.Mutex
	Result      chan *shared.RuntimeValue
	Error       chan *errors.RuntimeError
	ReturnValue *shared.RuntimeValue
	StartedAt   time.Time
	FinishedAt  time.Time
	Cancel      context.CancelFunc
	Context     context.Context
	Status      ThreadStatus
}

var maxRunningThreads = max(256, runtime.NumCPU()*32)

var threadLimiter = make(chan struct{}, maxRunningThreads)
var threadsGroup sync.WaitGroup
var threads []*Thread
var threadsMutex sync.Mutex

func module() (*shared.RuntimeValue, *errors.RuntimeError) {
	mod := values.MK_OBJECT(map[string]*shared.RuntimeValue{
		"spawn":      &spawn,
		"waitForAll": &waitForAll,
		"killAll":    &killAll,
	})

	return &mod, nil
}

func init() {
	modules.RegisterNativeModule("xel:threads", module)
}
