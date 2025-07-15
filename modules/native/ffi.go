package native

/*
#cgo CFLAGS: -I.
#include "ffi.h"
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/dev-kas/xel/shared"
)

type Library struct {
	Handle C.LibHandle
}

func (lib *Library) Close() {
	C.close_library(lib.Handle)
}

func (lib *Library) Call(name string, args []any) (any, error) {
	fnName := C.CString(name)
	defer C.free(unsafe.Pointer(fnName))

	sym := C.get_symbol(lib.Handle, fnName)
	if sym == nil {
		return nil, fmt.Errorf("symbol not found: %s", name)
	}

	argc := len(args)

	cArgs := C.malloc(C.size_t(argc) * C.size_t(unsafe.Sizeof(uintptr(0))))
	defer C.free(cArgs)

	cArgsTypes := C.malloc(C.size_t(argc) * C.size_t(unsafe.Sizeof(uintptr(0))))
	defer C.free(cArgsTypes)

	for i := 0; i < argc; i++ {
		ptr := (**C.void)(unsafe.Pointer(uintptr(cArgs) + uintptr(i)*unsafe.Sizeof(uintptr(0))))
		typePtr := (**C.char)(unsafe.Pointer(uintptr(cArgsTypes) + uintptr(i)*unsafe.Sizeof(uintptr(0))))

		switch v := args[i].(type) {
		case float64:
			p := (*C.double)(C.malloc(8))
			*p = C.double(v)
			*ptr = (*C.void)(unsafe.Pointer(p))

			t := C.CString("float")
			defer C.free(unsafe.Pointer(t))
			*typePtr = t
		case string:
			cs := C.CString(v)
			*ptr = (*C.void)(unsafe.Pointer(cs))

			t := C.CString("string")
			defer C.free(unsafe.Pointer(t))
			*typePtr = t
		case bool:
			b := (*C.char)(C.malloc(1))
			if v {
				*b = 1
			} else {
				*b = 0
			}
			*ptr = (*C.void)(unsafe.Pointer(b))

			t := C.CString("bool")
			defer C.free(unsafe.Pointer(t))
			*typePtr = t
		case nil:
			*ptr = nil
			t := C.CString("void")
			defer C.free(unsafe.Pointer(t))
			*typePtr = t
		default:
			return nil, fmt.Errorf("unsupported arg type at index %d: %T", i, v)
		}
	}

	ret := C.call(sym, (*unsafe.Pointer)(cArgs), (**C.char)(cArgsTypes), C.int(argc))

	var result any
	var freerName string

	switch ret.ret_type {
	case C.TYPE_INT, C.TYPE_LONG:
		result = int(*(*C.long)(ret.ret_val))
		freerName = "free_int"
	case C.TYPE_FLOAT, C.TYPE_DOUBLE:
		result = float64(*(*C.double)(ret.ret_val))
		freerName = "free_float"
	case C.TYPE_STRING:
		result = C.GoString((*C.char)(ret.ret_val))
		freerName = "free_string"
	case C.TYPE_BOOL:
		result = *(*bool)(ret.ret_val)
		freerName = "free_bool"
	case C.TYPE_VOID:
		result = nil
	default:
		return nil, fmt.Errorf("unsupported return type: %d", ret.ret_type)
	}

	if freerName != "" {
		freeSym := C.get_symbol(lib.Handle, C.CString(freerName))
		if freeSym != nil {
			C.call_free(freeSym, ret.ret_val)
		} else {
			C.free(ret.ret_val)
			shared.ColorPalette.Warning.Printf("Warning: symbol `%s` not found in library, memory leak possible\n", freerName)
		}
	}

	return result, nil
}

func loadLibrary(path string) (*Library, error) {
	libPath := C.CString(path)
	defer C.free(unsafe.Pointer(libPath))

	libHandle := C.load_library(libPath)
	if libHandle == nil {
		return nil, fmt.Errorf("failed to load library: %s", path)
	}
	return &Library{Handle: libHandle}, nil
}
