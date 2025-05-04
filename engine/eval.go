package engine

import (
	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/evaluator"
	"github.com/dev-kas/VirtLang-Go/parser"
	"github.com/dev-kas/VirtLang-Go/shared"
)

func Eval(src string, globalizer func(*environment.Environment)) (*shared.RuntimeValue, error) {
	// create a new parser
	p := parser.New()

	// parse the source code
	program, parser_error := p.ProduceAST(src)
	if parser_error != nil {
		return nil, parser_error
	}

	// create a new environment
	env := environment.NewEnvironment(nil)

	// globalize the environment
	if globalizer != nil {
		globalizer(&env)
	}

	// evaluate the program
	result, eval_error := evaluator.Evaluate(program, &env)
	if eval_error != nil {
		return nil, eval_error
	}
	return result, nil
}
