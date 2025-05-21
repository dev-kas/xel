package globals

import (
	"fmt"
	"os"
	"path/filepath"
	"xel/helpers"
	"xel/modules"
	xShared "xel/shared"

	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/evaluator"
	"github.com/dev-kas/virtlang-go/v3/parser"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

var importCache = map[string]*shared.RuntimeValue{}
var resolvingImports = map[string]bool{}

// Syntax:
/*
  - import("./path/to/file")
  - import("scope:package")
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

	modEnv := environment.NewEnvironment(&xShared.XelRootEnv)

	// Check if its a module name or a local file
	if libname[0] == '.' {
		// It's a local file
		libpath = filepath.Join(libname, libpath)
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
	} else {
		// It's a module name

		// First we'll check if it has this as a dependency in its manifest
		processRuntimeVal, err := env.LookupVar("proc")
		if err != nil {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("Failed to lookup proc: %v", err),
			}
		}

		if processRuntimeVal.Type != shared.Object {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("proc must be an object, but was %s (%s)", helpers.Stringify(*processRuntimeVal, false), shared.Stringify(processRuntimeVal.Type)),
			}
		}

		manifestRuntimeVal := processRuntimeVal.Value.(map[string]*shared.RuntimeValue)["manifest"]
		if manifestRuntimeVal == nil {
			return nil, &errors.RuntimeError{
				Message: "proc.manifest is not defined",
			}
		}

		if manifestRuntimeVal.Type != shared.Object {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("proc.manifest must be an object, but was %s (%s)", helpers.Stringify(*manifestRuntimeVal, false), shared.Stringify(manifestRuntimeVal.Type)),
			}
		}

		manifest := manifestRuntimeVal.Value.(map[string]*shared.RuntimeValue)
		dependencies := manifest["deps"].Value.(map[string]*shared.RuntimeValue)

		// Check if module is in the dependency list
		if _, exists := dependencies[libname]; !exists {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("module `%s` is not a dependency of this project", libname),
			}
		}

		constraint := dependencies[libname].Value.(string)
		constraint = constraint[1 : len(constraint)-1]

		pkgManifestPath, pkgManifest, resolutionErr := helpers.ResolveModule(libname, constraint)
		if resolutionErr != nil {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("Failed to resolve module `%s`: %v", libname, resolutionErr),
			}
		}

		manifestConverted := map[string]*shared.RuntimeValue{}
		nameVal := values.MK_STRING(fmt.Sprintf("\"%s\"", pkgManifest.Name))
		manifestConverted["name"] = &nameVal
		descVal := values.MK_STRING(fmt.Sprintf("\"%s\"", pkgManifest.Description))
		manifestConverted["description"] = &descVal
		versionVal := values.MK_STRING(fmt.Sprintf("\"%s\"", pkgManifest.Version))
		manifestConverted["version"] = &versionVal
		xelVal := values.MK_STRING(fmt.Sprintf("\"%s\"", pkgManifest.Xel))
		manifestConverted["xel"] = &xelVal
		engineVal := values.MK_STRING(fmt.Sprintf("\"%s\"", pkgManifest.Engine))
		manifestConverted["engine"] = &engineVal
		mainVal := values.MK_STRING(fmt.Sprintf("\"%s\"", pkgManifest.Main))
		manifestConverted["main"] = &mainVal
		authorVal := values.MK_STRING(fmt.Sprintf("\"%s\"", pkgManifest.Author))
		manifestConverted["author"] = &authorVal
		licenseVal := values.MK_STRING(fmt.Sprintf("\"%s\"", pkgManifest.License))
		manifestConverted["license"] = &licenseVal
		// Convert deps map separately
		depsObj := values.MK_OBJECT(map[string]*shared.RuntimeValue{})
		manifestConverted["deps"] = &depsObj
		depsMap := make(map[string]*shared.RuntimeValue)
		for k, v := range pkgManifest.Deps {
			depVal := values.MK_STRING(fmt.Sprintf("\"%s\"", v))
			depsMap[k] = &depVal
		}
		(*manifestConverted["deps"]).Value = depsMap

		modifiedProc := helpers.DeepCopyObject(processRuntimeVal.Value.(map[string]*shared.RuntimeValue))
		manifestVal := values.MK_OBJECT(manifestConverted)
		modifiedProc["manifest"] = &manifestVal

		modEnv.DeclareVar("proc", values.MK_OBJECT(modifiedProc), true)

		pkgPath := filepath.Dir(pkgManifestPath)
		libpath = filepath.Join(pkgPath, pkgManifest.Main)
	}

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

	libExports, err := evaluateModule(libpath, &modEnv)
	if err != nil {
		return nil, err
	}

	return libExports, nil
})

func evaluateModule(libpath string, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	content, readErr := os.ReadFile(libpath)
	if readErr != nil {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Failed to read file %s: %v", libpath, readErr),
		}
	}

	p := parser.New(libpath)
	lib, parserError := p.ProduceAST(string(content))
	if parserError != nil {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Failed to parse file %s: %v", libpath, parserError),
		}
	}

	libScope := environment.NewEnvironment(env)
	libScope.DeclareVar("__filename__", values.MK_STRING(fmt.Sprintf("\"%s\"", libpath)), true)
	libScope.DeclareVar("__dirname__", values.MK_STRING(fmt.Sprintf("\"%s\"", filepath.Dir(libpath))), true)
	placeholderExports := values.MK_NIL()
	importCache[libpath] = &placeholderExports

	libExports, evaluatorError := evaluator.Evaluate(lib, &libScope, nil)
	if evaluatorError != nil {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Failed to evaluate file %s: %v", libpath, evaluatorError),
		}
	}
	importCache[libpath] = libExports

	return libExports, nil
}
