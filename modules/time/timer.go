package time

import (
	xShared "xel/shared"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/evaluator"
	"github.com/dev-kas/virtlang-go/v3/parser"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

// Cache, or otherwise stack overflow
var timerGenerator_ *shared.RuntimeValue

func timerGenerator() *shared.RuntimeValue {
	if timerGenerator_ != nil {
		return timerGenerator_
	}

	nilVal := values.MK_NIL()
	timerGenerator_ = &nilVal

	code := `
	const time = import("xel:time")
	class Timer {
		private startedAt
		private elapsed_ = 0
		private running = false

		public constructor() {}

		public start() {
			if (running) {
				return nil
			}
			startedAt = time.now()
			running = true
		}

		public stop() {
			if (running) {
				elapsed_ = elapsed_ + (time.now() - startedAt)
				running = false
			}
		}

		public reset() {
			elapsed_ = 0
			startedAt = time.now()
		}

		public elapsed() {
			if (running) {
				return elapsed_ + (time.now() - startedAt)
			}
			return elapsed_
		}
	}

	// Export the class
	Timer
	`

	p := parser.New("xel-internal:/time/timer")

	ast, perr := p.ProduceAST(code)
	if perr != nil {
		xShared.ColorPalette.Warning.Println("Failed to produce AST for timer generator: " + perr.Error())
		return &nilVal
	}

	env := environment.NewEnvironment(&xShared.XelRootEnv)
	output, eerr := evaluator.Evaluate(ast, &env, nil)
	if eerr != nil {
		xShared.ColorPalette.Warning.Println("Failed to evaluate timer generator: " + eerr.Error())
		return &nilVal
	}

	timerGenerator_ = output
	return output
}
