package shared

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/dev-kas/virtlang-go/v4/debugger"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/fatih/color"
)

var RuntimeVersion string
var EngineVersion string

var XelRootEnv *environment.Environment = environment.NewEnvironment(nil)
var XelRootDebugger *debugger.Debugger = debugger.NewDebugger(XelRootEnv)

func XelDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	path := filepath.Join(homeDir, ".xel")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	return path
}

// Config represents the application configuration
type Config struct {
	DefaultTemplate     string   `json:"DefaultTemplate"`
	ModulePaths         []string `json:"ModulePaths"`
	PackageRegistryURI  string   `json:"PackageRegistryURI"`
	AllowInstallScripts bool     `json:"AllowInstallScripts"`
}

// ProjectManifest represents the metadata and configuration of a Xel project
//
// Fields:
// - Name: The name of the Xel project
// - Description: A brief description of the project
// - Version: The current version of the project
// - Xel: Required Xel runtime version (e.g., "^0.6.0")
// - Engine: Required VirtLang engine version (e.g., "^2.1.0")
// - Main: The main entry point file (relative to project root)
// - Deps: Project dependencies (key-value pairs of package names and versions)
// - Author: The author of the project
// - License: The license under which the project is distributed
//
// Example:
//
//	{
//	    "name": "my-project",
//	    "description": "A sample Xel project",
//	    "version": "1.0.0",
//	    "xel": "^0.6.0",
//	    "engine": "^2.1.0",
//	    "main": "src/main.xel",
//	    "deps": {
//	        "my-package": "^1.0.0"
//	    },
//	    "author": "John Doe",
//	    "license": "MIT"
//	}
type ProjectManifest struct {
	Name        string             `json:"name"`                 // Project name
	Description string             `json:"description"`          // Project description
	Version     string             `json:"version"`              // Project version
	Xel         *string            `json:"xel,omitempty"`        // Required Xel runtime version
	Engine      *string            `json:"engine,omitempty"`     // Required VirtLang engine version
	Main        string             `json:"main"`                 // Main entry point file
	Deps        *map[string]string `json:"deps,omitempty"`       // Project dependencies
	Author      string             `json:"author"`               // Project author
	License     string             `json:"license"`              // Project license
	Tags        []string           `json:"tags,omitempty"`       // Project tags
	Deprecated  *string            `json:"deprecated,omitempty"` // Project deprecation message
}

type Lockfile map[string]struct {
	Algorithm string `json:"algorithm"` // Algorithm used for hashing
	Hash      string `json:"hash"`      // Hash of the file
	URL       string `json:"url"`       // URL of the file
	Version   string `json:"version"`   // Version of the package
}

// XelConfig holds the application configuration
var XelConfig Config

// ColorPalette provides centralized color configuration
var ColorPalette = struct {
	// Standard colors
	Warning *color.Color
	Prompt  *color.Color
	Info    *color.Color
	Error   *color.Color

	// Specific color functions
	Welcome     func(format string, a ...interface{})
	Version     func(format string, a ...interface{}) string
	PromptStr   func(format string, a ...interface{}) string
	GrayMessage *color.Color
	ExitMessage *color.Color
}{
	// Basic colors
	Warning: color.New(color.FgYellow),
	Prompt:  color.New(color.FgBlue),
	Info:    color.New(color.FgCyan),
	Error:   color.New(color.FgRed),

	// Specific color functions
	Welcome:     color.Cyan,
	Version:     color.CyanString,
	PromptStr:   color.BlueString,
	GrayMessage: color.RGB(105, 105, 105),
	ExitMessage: color.New(color.FgHiRed),
}

func init() {
	xelDir := XelDir()
	templatesPath := filepath.Join(xelDir, "templates")
	if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
		err = os.MkdirAll(templatesPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// Path to the config file
	configPath := filepath.Join(xelDir, "config.json")

	// Check if config file exists, if not create it with default values
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		XelConfig = Config{
			DefaultTemplate:    "default",
			ModulePaths:        []string{filepath.Join(xelDir, "modules")},
			PackageRegistryURI: "https://pkg.xel.glitchiethedev.com/api/v1/",
			AllowInstallScripts: true,
		}

		configJSON, err := json.MarshalIndent(XelConfig, "", "  ")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(configPath, configJSON, 0644)
		if err != nil {
			panic(err)
		}
	} else {
		// Read existing config file
		configData, err := os.ReadFile(configPath)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(configData, &XelConfig)
		if err != nil {
			panic(err)
		}
	}
}
