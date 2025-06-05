//go:build !cgo

package native

import (
	"fmt"
)

// Stub implementation for platforms without CGO
func loadLibrary(path string) (*Library, error) {
	return nil, fmt.Errorf("native module loading is not supported without CGO")
}

// Stub Library implementation
type Library struct{}

func (lib *Library) Close() {}

func (lib *Library) Call(name string, args []any) (any, error) {
	return nil, fmt.Errorf("native module calls are not supported without CGO")
}