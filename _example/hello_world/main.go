package main

/*
#cgo CXXFLAGS: -std=c++11 -DV8_COMPRESS_POINTERS
#cgo pkg-config: v8.pc
#cgo LDFLAGS: -pthread
#include <stdlib.h>
#include "hello_world.h"
*/
import "C"
import (
	"runtime"
	time "time"
	_ "unsafe"
)

func main() {

	runtime.LockOSThread()
	C.v8_init()
	C.CHelloworld()

	// C.CHelloworld()

	time.Sleep(3 * time.Second)
}
