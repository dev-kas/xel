#ifndef FFI_H
#define FFI_H

#if defined(_WIN32)
#include <windows.h>
#else
#include <dlfcn.h>
#endif

#include <stdbool.h>

typedef void* LibHandle;


typedef enum {
    TYPE_VOID,
    TYPE_INT,
    TYPE_FLOAT,
    TYPE_LONG,
    TYPE_DOUBLE,
    TYPE_STRING,
    TYPE_BOOL,
    TYPE_UNKNOWN
} Type;

typedef struct {
    void* ret_val;
    Type  ret_type;
} FFReturn;

LibHandle load_library(const char* path);
void* get_symbol(LibHandle handle, const char* symbol);
void close_library(LibHandle handle);

FFReturn call(void* fn, void** args, const char** argt, int argc);
void* call_free(void* fn, void* ptr);

#endif
