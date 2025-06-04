// Package globals provides core functionality for module imports and dependency resolution
// in the Xel programming language. It handles both local file imports and remote package
// resolution, including dependency management and circular import detection.
package globals

import (
	// Standard library imports
	"fmt"
	"os"
	"path/filepath"

	// Internal imports
	"xel/helpers"
	"xel/modules"
	xShared "xel/shared"

	// External dependencies
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/evaluator"
	"github.com/dev-kas/virtlang-go/v4/parser"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

// importCache maintains a global cache of imported modules to prevent redundant
// file I/O and parsing operations. The cache uses the module's absolute path
// as the key and stores the evaluated module exports as the value.
// This significantly improves performance for frequently imported modules.
var importCache = map[string]*shared.RuntimeValue{}

// resolvingImports tracks modules that are currently being imported to detect
// and prevent circular dependencies. The map uses the module path as the key
// and a boolean flag to indicate if the import is in progress.
// This helps in providing meaningful error messages for circular dependencies.
var resolvingImports = map[string]bool{}

// Import is the native function implementation for the `import` statement in Xel.
// It handles module resolution, caching, and circular dependency detection.
//
// Supported import patterns:
//   - Local file: import("./path/to/file")
//   - Scoped package: import("@scope/package")
//   - Regular package: import("package")
//
// The function performs the following operations:
// 1. Validates the import statement syntax
// 2. Resolves the module path (local or remote)
// 3. Checks for circular dependencies
// 4. Returns cached module if available
// 5. Otherwise, loads and evaluates the module
var Import = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	// Ensure exactly one argument is provided to the import statement
	if len(args) != 1 {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("`import` requires exactly one argument, but received %d", len(args)),
		}
	}

	// The import specifier must be a string literal
	if args[0].Type != shared.String {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("`import` expects a string argument, but received %s (%s)",
				helpers.Stringify(args[0], false),
				shared.Stringify(args[0].Type)),
		}
	}

	// Extract the module name from the string literal
	libname := args[0].Value.(string)

	// Check if the requested module is a native module
	// Native modules are implemented in Go and provide core functionality
	if loader, isNative := modules.GetNativeModuleLoader(libname); isNative {
		return loader()
	}

	// Construct the default module path with .xel extension
	libpath := fmt.Sprintf("%s.xel", libname)

	// Flag to control module caching behavior
	// Can be modified to force reload in development
	useCache := true

	// Create a new environment for the module being imported
	// This inherits from the root Xel environment
	modEnv := environment.NewEnvironment(xShared.XelRootEnv)

	// Determine if this is a local file import or a package import
	// Local imports start with './' or '../'
	if libname[0] == '.' {
		// Handle local file import (relative to current file)
		// Get the directory of the current file from the debugger context
		dirname := filepath.Dir(xShared.XelRootDebugger.CurrentFile)

		// Resolve the full path to the target module
		libpath = filepath.Join(dirname, libpath)
	} else {
		// Handle package import (from node_modules or package registry)

		// For package imports, we need to verify the package is declared in the project's dependencies
		// and resolve its actual location from the package manager's cache

		// Look up the 'proc' variable which contains process-related information including the manifest
		processRuntimeVal, err := env.LookupVar("proc")
		if err != nil {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("Failed to lookup process information: %v", err),
			}
		}

		// Ensure proc is an object as expected
		if processRuntimeVal.Type != shared.Object {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("Expected 'proc' to be an object, but got %s (%s)",
					helpers.Stringify(*processRuntimeVal, false),
					shared.Stringify(processRuntimeVal.Type)),
			}
		}

		// Extract the manifest from the process object
		processObj := processRuntimeVal.Value.(map[string]*shared.RuntimeValue)
		manifestRuntimeVal, exists := processObj["manifest"]
		if !exists || manifestRuntimeVal == nil {
			return nil, &errors.RuntimeError{
				Message: "Project manifest is not available in process object",
			}
		}

		// Ensure the manifest is an object with the expected structure
		if manifestRuntimeVal.Type != shared.Object {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("Expected manifest to be an object, but got %s (%s)",
					helpers.Stringify(*manifestRuntimeVal, false),
					shared.Stringify(manifestRuntimeVal.Type)),
			}
		}

		// Extract the manifest object and its dependencies
		manifest := manifestRuntimeVal.Value.(map[string]*shared.RuntimeValue)

		// Get the dependencies map from the manifest
		depsVal, depsExists := manifest["deps"]
		if !depsExists || depsVal.Type != shared.Object {
			return nil, &errors.RuntimeError{
				Message: "Project manifest is missing or has invalid 'dependencies' section",
			}
		}
		dependencies := depsVal.Value.(map[string]*shared.RuntimeValue)

		// Verify the requested package is in the dependencies
		pkgConstraint, exists := dependencies[libname]
		if !exists {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("Package '%s' is not listed in the project's dependencies", libname),
			}
		}

		// Extract the version constraint
		constraint := pkgConstraint.Value.(string)

		// Resolve the module using the package manager to get its actual location
		pkgManifestPath, pkgManifest, resolutionErr := helpers.ResolveModule(libname, constraint)
		if resolutionErr != nil {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("Failed to resolve package '%s' (constraint: %s): %v",
					libname, constraint, resolutionErr),
			}
		}

		// Convert the resolved package manifest into a runtime-compatible format
		// This ensures the imported module has access to the correct manifest information
		manifestConverted := map[string]*shared.RuntimeValue{}

		// Helper to store string values in the manifest
		addStringField := func(key, value string) {
			val := values.MK_STRING(value)
			manifestConverted[key] = &val
		}

		// Convert each manifest field to a runtime string value
		addStringField("name", pkgManifest.Name)
		addStringField("description", pkgManifest.Description)
		addStringField("version", pkgManifest.Version)
		addStringField("xel", pkgManifest.Xel)
		addStringField("engine", pkgManifest.Engine)
		addStringField("main", pkgManifest.Main)
		addStringField("author", pkgManifest.Author)
		addStringField("license", pkgManifest.License)

		// Convert dependencies map to runtime format
		depsMap := make(map[string]*shared.RuntimeValue, len(pkgManifest.Deps))
		for depName, depConstraint := range pkgManifest.Deps {
			depVal := values.MK_STRING(depConstraint)
			depsMap[depName] = &depVal
		}
		depsObj := values.MK_OBJECT(depsMap)
		manifestConverted["deps"] = &depsObj

		// Create a deep copy of the process object to avoid modifying the original
		modifiedProc := helpers.DeepCopyObject(processRuntimeVal.Value.(map[string]*shared.RuntimeValue))

		// Create a new manifest value with the converted package manifest
		manifestVal := values.MK_OBJECT(manifestConverted)
		modifiedProc["manifest"] = &manifestVal

		// Update the module's environment with the modified process object
		// This ensures the imported module sees itself with its own manifest
		modEnv.DeclareVar("proc", values.MK_OBJECT(modifiedProc), true)

		// Resolve the full path to the package's main entry point
		pkgPath := filepath.Dir(pkgManifestPath)
		libpath = filepath.Join(pkgPath, pkgManifest.Main)
	}

	// Check for circular dependencies before proceeding with the import
	if resolvingImports[libpath] {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Circular import detected while importing '%s'", libpath),
		}
	}

	// Return cached module if available and caching is enabled
	if val, exists := importCache[libpath]; exists && useCache {
		return val, nil
	}

	// Mark this module as being imported to detect circular dependencies
	resolvingImports[libpath] = true
	// Ensure we clean up the resolvingImports map when done
	defer delete(resolvingImports, libpath)

	// Evaluate the module and get its exports
	libExports, err := evaluateModule(libpath, modEnv)
	if err != nil {
		return nil, err
	}

	return libExports, nil
})

// evaluateModule loads, parses, and evaluates a module file, returning its exports.
// It handles setting up the module environment and populating special variables.
//
// Parameters:
//   - libpath: The absolute path to the module file to evaluate
//   - env: The parent environment to inherit from
//
// Returns:
//   - The module's exports as a RuntimeValue
//   - An error if any step of the evaluation fails
func evaluateModule(libpath string, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	// Read the module source code
	content, readErr := os.ReadFile(libpath)
	if readErr != nil {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Failed to read module file '%s': %v", libpath, readErr),
		}
	}

	// Parse the source code into an AST
	p := parser.New(libpath)
	lib, parserError := p.ProduceAST(string(content))
	if parserError != nil {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Syntax error in '%s': %v", libpath, parserError),
		}
	}

	// Create a new scope for this module, inheriting from the parent
	libScope := environment.NewEnvironment(env)

	// Set up module-specific variables
	// Set up module-specific variables with proper string values
	filenameVal := values.MK_STRING(libpath)
	dirnameVal := values.MK_STRING(filepath.Dir(libpath))

	libScope.DeclareVar("__filename__", filenameVal, true) // constant
	libScope.DeclareVar("__dirname__", dirnameVal, true)   // constant

	// Add a placeholder to the cache to handle circular references
	placeholderExports := values.MK_NIL()
	importCache[libpath] = &placeholderExports

	// Evaluate the module's AST in the new scope
	libExports, evaluatorError := evaluator.Evaluate(
		lib,
		libScope,
		xShared.XelRootDebugger,
	)

	// Handle evaluation errors
	if evaluatorError != nil {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Runtime error in '%s': %v", libpath, evaluatorError),
		}
	}

	// Update the cache with the actual exports
	importCache[libpath] = libExports

	return libExports, nil
}
