package classes

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var methods = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{
			Message: "methods() takes exactly one argument",
		}
	}

	if args[0].Type != shared.Class {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("methods() takes a class as an argument, but got %s", shared.Stringify(args[0].Type)),
		}
	}

	class := args[0].Value.(values.ClassValue)
	methods := make([]shared.RuntimeValue, 0)

	for _, stmt := range class.Body {
		if stmt.GetType() == ast.ClassMethodNode && stmt.(*ast.ClassMethod).IsPublic {
			name := stmt.(*ast.ClassMethod).Name
			methods = append(methods, values.MK_STRING(name))
		}
	}

	res := values.MK_ARRAY(methods)
	return &res, nil
})
