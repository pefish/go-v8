package go_v8

import (
	"fmt"
	"testing"
)

func TestNewV8Isolate(t *testing.T) {
	defer func() {
		fmt.Println(isolatesNum)
	}()

	v8Isolate := NewV8Isolate()
	defer v8Isolate.Close()

	//v8Isolate1 := NewV8Isolate()
	//defer v8Isolate1.Close()
	//
	//v8Isolate2 := NewV8Isolate()
	//defer v8Isolate2.Close()

	fmt.Println(isolatesNum)
}
