#include "ffi.h"

LibHandle load_library(const char* path) {
#if defined(_WIN32)
    return LoadLibraryA(path);
#else
    return dlopen(path, RTLD_LAZY);
#endif
}

void* get_symbol(LibHandle handle, const char* symbol) {
#if defined(_WIN32)
    return GetProcAddress((HMODULE)handle, symbol);
#else
    return dlsym(handle, symbol);
#endif
}

void close_library(LibHandle handle) {
#if defined(_WIN32)
    FreeLibrary((HMODULE)handle);
#else
    dlclose(handle);
#endif
}

FFReturn call(void* fn, void** args, const char** argt, int argc) {
    typedef FFReturn (*trampoline_fn)(void**, const char**, int);
    trampoline_fn real_fn = (trampoline_fn)fn;
    return real_fn(args, argt, argc);
}

void* call_free(void* fn, void* ptr) {
    typedef void* (*trampoline_fn)(void*);
    trampoline_fn real_fn = (trampoline_fn)fn;
    return real_fn(ptr);
}

