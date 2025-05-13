package shared

import "github.com/dev-kas/virtlang-go/v2/environment"

var RuntimeVersion string
var EngineVersion string

var XelRootEnv environment.Environment = environment.NewEnvironment(nil)
