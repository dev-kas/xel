package globals

import (
	"fmt"
	"os"
	"path/filepath"
	"xel/helpers"
	"xel/modules"
	xShared "xel/shared"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/evaluator"
	"github.com/dev-kas/virtlang-go/v2/parser"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

var importCache = map[string]*shared.RuntimeValue{}
var resolvingImports = map[string]bool{}

// Syntax:
/*
  - import("./path/to/file")
  - import("@scope:package")
  - import("package")
*/
var Import = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 1 {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("`import` takes one and only one argument, but given %d", len(args)),
		}
	}

	if args[0].Type != shared.String {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("`import` only accepts string as an argument, but was given %s (%s)", helpers.Stringify(args[0], false), shared.Stringify(args[0].Type)),
		}
	}

	libname := args[0].Value.(string)
	libname = libname[1 : len(libname)-1]

	if loader, isNative := modules.GetNativeModuleLoader(libname); isNative {
		return loader()
	}

	libpath := fmt.Sprintf("%s.xel", libname)
	useCache := true

	dirnameRuntimeVal, err := env.LookupVar("__dirname__")
	if err != nil {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Failed to lookup __dirname__: %v", err),
		}
	}

	if dirnameRuntimeVal.Type != shared.String {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("__dirname__ must be a string, but was %s (%s)", helpers.Stringify(*dirnameRuntimeVal, false), shared.Stringify(dirnameRuntimeVal.Type)),
		}
	}

	dirname := dirnameRuntimeVal.Value.(string)
	dirname = dirname[1 : len(dirname)-1]
	libpath = filepath.Join(dirname, libpath)

	if resolvingImports[libpath] {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Circular import detected: %s", libpath),
		}
	}

	if val, exists := importCache[libpath]; exists && useCache {
		return val, nil
	}

	resolvingImports[libpath] = true
	defer delete(resolvingImports, libpath)

	content, readErr := os.ReadFile(libpath)
	if readErr != nil {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Failed to read file %s: %v", libpath, readErr),
		}
	}

	p := parser.New()
	lib, parserError := p.ProduceAST(string(content))
	if parserError != nil {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Failed to parse file %s: %v", libpath, parserError),
		}
	}

	libScope := environment.NewEnvironment(&xShared.XelRootEnv)
	libScope.DeclareVar("__filename__", values.MK_STRING(fmt.Sprintf("\"%s\"", libpath)), true)
	libScope.DeclareVar("__dirname__", values.MK_STRING(fmt.Sprintf("\"%s\"", filepath.Dir(libpath))), true)
	placeholderExports := values.MK_NIL()
	importCache[libpath] = &placeholderExports

	libExports, evaluatorError := evaluator.Evaluate(lib, &libScope)
	if evaluatorError != nil {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Failed to evaluate file %s: %v", libpath, evaluatorError),
		}
	}
	importCache[libpath] = libExports

	return libExports, nil
})
