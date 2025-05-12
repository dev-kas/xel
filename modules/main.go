package modules

import (
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
)

type NativeModuleLoader func() (*shared.RuntimeValue, *errors.RuntimeError)

var NativeModuleRegistry = make(map[string]NativeModuleLoader)

func RegisterNativeModule(name string, loader NativeModuleLoader) {
	if _, exists := NativeModuleRegistry[name]; exists {
		panic("Native module " + name + " already exists")
	}
	NativeModuleRegistry[name] = loader
}

func GetNativeModuleLoader(name string) (NativeModuleLoader, bool) {
	loader, exists := NativeModuleRegistry[name]
	return loader, exists
}
