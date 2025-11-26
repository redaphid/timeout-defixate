//go:build darwin

package main

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>

void doLockScreen() {
    void* handle = dlopen("/System/Library/PrivateFrameworks/login.framework/Versions/Current/login", RTLD_LAZY);
    if (handle) {
        void (*lock)() = dlsym(handle, "SACLockScreenImmediate");
        if (lock) {
            lock();
        }
        dlclose(handle);
    }
}
*/
import "C"

func lockScreen() {
	C.doLockScreen()
}
