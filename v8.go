package go_v8

/*
#cgo CXXFLAGS: -std=c++11 -DV8_COMPRESS_POINTERS
#cgo pkg-config: v8.pc
#include "v8_wrap.h"
*/
import "C"
import (
	"sync"
	"sync/atomic"
)

// V8Isolate 接口
type IV8Isolate interface {
	Close()
	RunInNewContext(src string) error
}

var (
	initV8      sync.Once
	isolatesNum int64
)

type V8Isolate struct { // 对应v8中的Isolate
	isolate C.IsolatePtr
}

func NewV8Isolate() *V8Isolate {
	initV8.Do(func() {
		C.InitV8() // 初始化V8，只能初始化一次
	})
	v8 := &V8Isolate{
		isolate: C.NewIsolate(),
	}
	atomic.AddInt64(&isolatesNum, 1)
	return v8
}

func (v8Isolate *V8Isolate) Close() {
	C.CloseIsolate(v8Isolate.isolate)
	atomic.AddInt64(&isolatesNum, -1)
	if isolatesNum == 0 { // 没了isolate就关闭V8
		C.CloseV8()
	}
}

func (v8Isolate *V8Isolate) RunInNewContext(src string) error {
	return nil
}
