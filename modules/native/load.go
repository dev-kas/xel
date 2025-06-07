package native

import (
	"fmt"
	"path/filepath"
	xShared "xel/shared"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var load = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{
			Message: "load() takes exactly 1 argument",
		}
	}

	if args[0].Type != shared.String {
		return nil, &errors.RuntimeError{
			Message: "load() expects string path as argument",
		}
	}

	path := args[0].Value.(string)

	dirname := filepath.Dir(xShared.XelRootDebugger.CurrentFile)
	libpath := filepath.Join(dirname, path)

	lib, err := loadLibrary(libpath)
	if err != nil {
		return nil, &errors.RuntimeError{
			Message: err.Error(),
		}
	}

	nilVal := values.MK_NIL()

	callVal := values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
		if len(args) != 2 {
			return nil, &errors.RuntimeError{
				Message: "call() takes exactly 2 arguments",
			}
		}
		if args[0].Type != shared.String {
			return nil, &errors.RuntimeError{
				Message: "call() expects string name as first argument",
			}
		}
		name := args[0].Value.(string)
		var vlargs []shared.RuntimeValue
		if len(args) > 1 {
			if args[1].Type != shared.Array {
				return nil, &errors.RuntimeError{
					Message: "call() expects array of arguments as second argument",
				}
			}
			vlargs = args[1].Value.([]shared.RuntimeValue)
		}

		var ffargs []any
		for i, arg := range vlargs {
			switch arg.Type {
			case shared.String:
				ffargs = append(ffargs, arg.Value.(string))
			case shared.Number:
				ffargs = append(ffargs, arg.Value.(float64))
			case shared.Boolean:
				ffargs = append(ffargs, arg.Value.(bool))
			case shared.Nil:
				ffargs = append(ffargs, nil)
			default:
				return nil, &errors.RuntimeError{
					Message: "call() invalid argument type " + shared.Stringify(arg.Type) + " at index " + fmt.Sprintf("%d", i),
				}
			}
		}

		res, err := lib.Call(name, ffargs)
		if err != nil {
			return nil, &errors.RuntimeError{
				Message: err.Error(),
			}
		}

		var retVal shared.RuntimeValue

		switch v := res.(type) {
		case string:
			retVal = values.MK_STRING(v)
		case int:
			retVal = values.MK_NUMBER(float64(v))
		case float64:
			retVal = values.MK_NUMBER(v)
		case bool:
			retVal = values.MK_BOOL(v)
		case nil:
			retVal = nilVal
		default:
			return nil, &errors.RuntimeError{
				Message: "call() foreign function produced an invalid return value",
			}
		}

		return &retVal, nil
	})

	unloadVal := values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
		lib.Close()
		return &nilVal, nil
	})

	retVal := values.MK_OBJECT(map[string]*shared.RuntimeValue{
		"call":   &callVal,
		"unload": &unloadVal,
	})

	return &retVal, nil
})
