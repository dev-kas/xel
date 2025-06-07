package classes

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var properties = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{
			Message: "properties() takes exactly one argument",
		}
	}

	if args[0].Type != shared.Class {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("properties() takes a class as an argument, but got %s", shared.Stringify(args[0].Type)),
		}
	}

	class := args[0].Value.(values.ClassValue)
	properties := make([]shared.RuntimeValue, 0)

	for _, stmt := range class.Body {
		if stmt.GetType() == ast.ClassPropertyNode && stmt.(*ast.ClassProperty).IsPublic {
			name := stmt.(*ast.ClassProperty).Name
			properties = append(properties, values.MK_STRING(name))
		}
	}

	res := values.MK_ARRAY(properties)
	return &res, nil
})
