package old

/*
#include "v8_wrap.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "runtime"
import "reflect"
import "time"

// The superclass of all JavaScript values and objects.
//
type Value struct {
	engine       *Engine
	self         unsafe.Pointer
	isType       int
	notType      int
	fieldOwnerId int64
}

func newValue(engine *Engine, self unsafe.Pointer) *Value {
	if self == nil {
		return nil
	}

	result := &Value{
		engine: engine,
		self:   self,
	}

	runtime.SetFinalizer(result, func(v *Value) {
		if traceDispose {
			println("v8.Value.Dispose()", v.self)
		}
		C.V8_DisposeValue(v.self)
	})

	return result
}

func (e *Engine) Undefined() *Value {
	if e._undefined == nil {
		e._undefined = newValue(e, C.V8_Undefined(e.self))
	}
	return e._undefined
}

func (e *Engine) Null() *Value {
	if e._null == nil {
		e._null = newValue(e, C.V8_Null(e.self))
	}
	return e._null
}

func (e *Engine) True() *Value {
	if e._true == nil {
		e._true = newValue(e, C.V8_True(e.self))
	}
	return e._true
}

func (e *Engine) False() *Value {
	if e._false == nil {
		e._false = newValue(e, C.V8_False(e.self))
	}
	return e._false
}

func (e *Engine) NewBoolean(value bool) *Value {
	if value {
		return e.True()
	}
	return e.False()
}

func (e *Engine) NewNumber(value float64) *Value {
	return newValue(e, C.V8_NewNumber(
		e.self, C.double(value),
	))
}

func (e *Engine) NewDate(value time.Time) *Value {
	return newValue(e, C.V8_NewDate(
		e.self, C.double(value.Unix()*1000),
	))
}

func (e *Engine) NewInteger(value int64) *Value {
	return newValue(e, C.V8_NewNumber(
		e.self, C.double(value),
	))
}

func (e *Engine) NewString(value string) *Value {
	valPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&value)).Data)
	return newValue(e, C.V8_NewString(
		e.self, (*C.char)(valPtr), C.int(len(value)),
	))
}

func (e *Engine) NewRangeError(message string) *Value {
	msgPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&message)).Data)
	return newValue(e, C.V8_Exception_RangeError(e.self, (*C.char)(msgPtr), C.int(len(message))))
}

func (e *Engine) NewReferenceError(message string) *Value {
	msgPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&message)).Data)
	return newValue(e, C.V8_Exception_ReferenceError(e.self, (*C.char)(msgPtr), C.int(len(message))))
}

func (e *Engine) NewSyntaxError(message string) *Value {
	msgPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&message)).Data)
	return newValue(e, C.V8_Exception_SyntaxError(e.self, (*C.char)(msgPtr), C.int(len(message))))
}

func (e *Engine) NewTypeError(message string) *Value {
	msgPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&message)).Data)
	return newValue(e, C.V8_Exception_TypeError(e.self, (*C.char)(msgPtr), C.int(len(message))))
}

func (e *Engine) NewError(message string) *Value {
	msgPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&message)).Data)
	return newValue(e, C.V8_Exception_Error(e.self, (*C.char)(msgPtr), C.int(len(message))))
}

func (v *Value) ToBoolean() bool {
	return C.V8_Value_ToBoolean(v.self) == 1
}

func (v *Value) ToNumber() float64 {
	return float64(C.V8_Value_ToNumber(v.self))
}

func (v *Value) ToInteger() int64 {
	return int64(C.V8_Value_ToInteger(v.self))
}

func (v *Value) ToUint32() uint32 {
	return uint32(C.V8_Value_ToUint32(v.self))
}

func (v *Value) ToInt32() int32 {
	return int32(C.V8_Value_ToInt32(v.self))
}

func (v *Value) ToString() string {
	cstring := C.V8_Value_ToString(v.self)
	gostring := C.GoString(cstring)
	C.free(unsafe.Pointer(cstring))
	return gostring
}

func (v *Value) ToTime() time.Time {
	return time.Unix(0, int64(C.V8_Value_ToInteger(v.self))*1e6)
}

func (v *Value) ToObject() *Object {
	if v == nil {
		return nil
	}
	return &Object{v, nil, nil}
}

func (v *Value) ToArray() *Array {
	if v == nil {
		return nil
	}
	return &Array{&Object{v, nil, nil}}
}

func (v *Value) ToRegExp() *RegExp {
	if v == nil {
		return nil
	}
	return &RegExp{&Object{v, nil, nil}, "", false, RF_None, false}
}

func (v *Value) ToFunction() *Function {
	if v == nil {
		return nil
	}
	return &Function{&Object{v, nil, nil}, nil, nil}
}

func (v *Value) ToExternal() *External {
	if v == nil {
		return nil
	}
	return &External{v, nil}
}

func (v *Value) String() string {
	return v.ToString()
}

// Keep the Object alive when it refence by JS
func (v *Value) setOwner(self interface{}) {
	// the object reference by engine
	if v.fieldOwnerId == 0 {
		v.engine.fieldOwnerId += 1
		v.fieldOwnerId = v.engine.fieldOwnerId
		v.engine.fieldOwners[v.fieldOwnerId] = self
		C.V8_Value_SetFieldOwnerInfo(v.self, unsafe.Pointer(v.engine), C.int64_t(v.fieldOwnerId))
	}
}

const (
	isUndefined     = 1 << iota
	isNull          = 1 << iota
	isTrue          = 1 << iota
	isFalse         = 1 << iota
	isString        = 1 << iota
	isFunction      = 1 << iota
	isArray         = 1 << iota
	isObject        = 1 << iota
	isBoolean       = 1 << iota
	isNumber        = 1 << iota
	isExternal      = 1 << iota
	isInt32         = 1 << iota
	isUint32        = 1 << iota
	isDate          = 1 << iota
	isBooleanObject = 1 << iota
	isNumberObject  = 1 << iota
	isStringObject  = 1 << iota
	isNativeError   = 1 << iota
	isRegExp        = 1 << iota
)

func (v *Value) checkJsType(typeCode int, check func(unsafe.Pointer) bool) bool {
	if v.isType&typeCode == typeCode {
		return true
	}

	if v.notType&typeCode == typeCode {
		return false
	}

	if check(v.self) {
		v.isType |= typeCode
		return true
	} else {
		v.notType |= typeCode
		return false
	}
}

func (v *Value) IsUndefined() bool {
	return v.checkJsType(isUndefined, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsUndefined(self) == 1
	})
}

func (v *Value) IsNull() bool {
	return v.checkJsType(isNull, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsNull(self) == 1
	})
}

func (v *Value) IsTrue() bool {
	return v.checkJsType(isTrue, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsTrue(self) == 1
	})
}

func (v *Value) IsFalse() bool {
	return v.checkJsType(isFalse, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsFalse(self) == 1
	})
}

func (v *Value) IsString() bool {
	return v.checkJsType(isString, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsString(self) == 1
	})
}

func (v *Value) IsFunction() bool {
	return v.checkJsType(isFunction, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsFunction(self) == 1
	})
}

func (v *Value) IsArray() bool {
	return v.checkJsType(isArray, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsArray(self) == 1
	})
}

func (v *Value) IsObject() bool {
	return v.checkJsType(isObject, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsObject(self) == 1
	})
}

func (v *Value) IsBoolean() bool {
	return v.checkJsType(isBoolean, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsBoolean(self) == 1
	})
}

func (v *Value) IsNumber() bool {
	return v.checkJsType(isNumber, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsNumber(self) == 1
	})
}

func (v *Value) IsExternal() bool {
	return v.checkJsType(isExternal, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsExternal(self) == 1
	})
}

func (v *Value) IsInt32() bool {
	return v.checkJsType(isInt32, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsInt32(self) == 1
	})
}

func (v *Value) IsUint32() bool {
	return v.checkJsType(isUint32, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsUint32(self) == 1
	})
}

func (v *Value) IsDate() bool {
	return v.checkJsType(isDate, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsDate(self) == 1
	})
}

func (v *Value) IsBooleanObject() bool {
	return v.checkJsType(isBooleanObject, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsBooleanObject(self) == 1
	})
}

func (v *Value) IsNumberObject() bool {
	return v.checkJsType(isNumberObject, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsNumberObject(self) == 1
	})
}

func (v *Value) IsStringObject() bool {
	return v.checkJsType(isStringObject, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsStringObject(self) == 1
	})
}

func (v *Value) IsNativeError() bool {
	return v.checkJsType(isNativeError, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsNativeError(self) == 1
	})
}

func (v *Value) IsRegExp() bool {
	return v.checkJsType(isRegExp, func(self unsafe.Pointer) bool {
		return C.V8_Value_IsRegExp(self) == 1
	})
}
