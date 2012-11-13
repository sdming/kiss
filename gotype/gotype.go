// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package gotype

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

const (
	MaxUint8  uint8  = 255
	MaxUint16 uint16 = 65535
	MaxUint32 uint32 = 4294967295
	MaxUint64 uint64 = 18446744073709551615
	MaxInt8   int8   = 127
	MaxInt16  int16  = 32767
	MaxInt32  int32  = 2147483647
	MaxInt64  int64  = 9223372036854775807
	MinInt8   int8   = -128
	MinInt16  int16  = -32768
	MinInt32  int32  = -2147483648
	MinInt64  int64  = -9223372036854775808
)

var (
	TypeBool      = reflect.TypeOf(false)
	TypeInt       = reflect.TypeOf(int(0))
	TypeInt8      = reflect.TypeOf(int8(0))
	TypeInt16     = reflect.TypeOf(int16(0))
	TypeInt32     = reflect.TypeOf(int32(0))
	TypeInt64     = reflect.TypeOf(int64(0))
	TypeUint      = reflect.TypeOf(uint(0))
	TypeUint8     = reflect.TypeOf(uint8(0))
	TypeUint16    = reflect.TypeOf(uint16(0))
	TypeUint32    = reflect.TypeOf(uint32(0))
	TypeUint64    = reflect.TypeOf(uint64(0))
	TypeFloat32   = reflect.TypeOf(float32(0))
	TypeFloat64   = reflect.TypeOf(float64(0))
	TypeString    = reflect.TypeOf("")
	TypeByteSlice = reflect.TypeOf([]byte(nil))
)

// returns the name of the calling method
func methodName() string {
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown method"
	}
	return f.Name()
}

// returns the name of the calling method, Caller(N)
func methodNameN(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown method"
	}
	return f.Name()
}

func checkError(e error) {
	if x := recover(); x != nil {
		if e1, ok := x.(error); ok {
			e = e1
		}
		e = newTypeErr(methodNameN(2), "run time panic", x)
	}
}

// go type error
type TypeError struct {
	Func    string      // the failing function 
	Message string      // the error message	
	Inner   interface{} // reference value
}

// Error interface of TypeError
func (e *TypeError) Error() string {
	if e.Inner == nil {
		return fmt.Sprintf("gotype func %s error, %s", e.Func, e.Message)
	}

	if x, ok := e.Inner.(error); ok {
		return fmt.Sprintf("gotype func %s error, %s, inner error %v",
			e.Func, e.Message, x.Error())
	}

	return fmt.Sprintf("gotype func %s error, %s, reference %v", e.Func, e.Message, e.Inner)
}

func newTypeErr(fn string, msg string, x interface{}) *TypeError {
	e := &TypeError{Func: fn, Message: msg, Inner: x}
	return e
}

// is struct?
func IsStruct(kind reflect.Kind) bool {
	return kind == reflect.Struct
}

// is boolean?
func IsBool(kind reflect.Kind) bool {
	return kind == reflect.Bool
}

// is collect?
func IsCollect(kind reflect.Kind) bool {
	switch kind {
	case reflect.Array, reflect.Slice, reflect.Map:
		return true
	}
	return false
}

// is numeric type?
func IsNumeric(kind reflect.Kind) bool {
	return IsInt(kind) || IsUint(kind) || IsFloat(kind)
}

// is int type?
func IsInt(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	}
	return false
}

// is Uint type?
func IsUint(kind reflect.Kind) bool {
	switch kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}
	return false
}

// is float type?
func IsFloat(kind reflect.Kind) bool {
	switch kind {
	case reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

// is string?
func IsString(kind reflect.Kind) bool {
	return kind == reflect.String
}

// is simple type(boo, string, numeric)?
func IsSimple(kind reflect.Kind) bool {
	return IsBool(kind) || IsNumeric(kind) || IsString(kind)
}

// set a as addrss of b
func SetAddr(a reflect.Value, b reflect.Value) {
	if !a.IsValid() || !b.IsValid() || !a.CanSet() {
		return
	}
	if b.CanAddr() {
		a.Set(b.Addr())
	} else if b.Kind() == reflect.Ptr {
		a.Set(b)
	}
}

// can compare a & b ?
func CanCompareValue(a, b reflect.Value) bool {
	if !a.IsValid() || !b.IsValid() {
		return false
	}
	a, b = Underlying(a), Underlying(b)
	return CanCompareKind(a.Kind(), b.Kind())
}

// can compare a & b ?
func CanCompareKind(a, b reflect.Kind) bool {
	if IsNumeric(a) && IsNumeric(b) {
		return true
	}
	if IsString(a) && IsString(b) {
		return true
	}
	if IsBool(a) && IsBool(b) {
		return true
	}
	return false
}

// is a pointer?
func CanPointer(kind reflect.Kind) bool {
	switch kind {
	case reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Func:
		return true
	}
	return false
}

// has element type?
func CanTypeElem(kind reflect.Kind) bool {
	switch kind {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	}
	return false
}

// value can element?
func CanValueElem(kind reflect.Kind) bool {
	switch kind {
	case reflect.Ptr, reflect.Interface:
		return true
	}
	return false
}

// get underlying element type
func UnderlyingType(typ reflect.Type) reflect.Type {
	if CanTypeElem(typ.Kind()) {
		return typ.Elem()
	}
	return typ
}

// get kind of underlying value 
func UnderlyingKind(v reflect.Value) reflect.Kind {
	return Underlying(v).Kind()
}

// get underlying value of *struct
func Underlying(v reflect.Value) reflect.Value {
	if !v.IsValid() {
		return v
	}

	k := v.Kind()
	if k == reflect.Ptr || k == reflect.Interface {
		return v.Elem()
	}
	return v
}

// get public field by name
func FieldByName(v reflect.Value, name string) (field reflect.Value, ok bool) {
	field = Underlying(v).FieldByName(strings.Title(name))
	if field.IsValid() {
		return field, true
	}
	return
}

// get public field by name ignore case
func FieldByNameFold(v reflect.Value, name string) (field reflect.Value, ok bool) {
	field = Underlying(v).FieldByNameFunc(func(x string) bool {
		return strings.EqualFold(name, x)
	})
	if field.IsValid() {
		return field, true
	}
	return
}

// get public method by name
func MethodByName(v reflect.Value, name string) (m reflect.Value, ok bool) {
	m = v.MethodByName(strings.Title(name))
	if m.IsValid() && m.Kind() == reflect.Func {
		return m, true
	}
	return
}

// get public method by name ignore case
func MethodByNameFold(v reflect.Value, name string) (m reflect.Value, ok bool) {
	typ := v.Type()
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		if method.PkgPath != "" {
			continue
		}
		if strings.EqualFold(method.Name, name) {
			return MethodByName(v, method.Name)
		}
	}
	return
}

// get fields name
func Fields(typ reflect.Type) []string {
	if typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Interface {
		typ = typ.Elem()
	}

	count := typ.NumField()
	names := make([]string, 0, count)
	for i := 0; i < count; i++ {
		f := typ.Field(i)
		if f.PkgPath != "" {
			continue
		}
		names = append(names, f.Name)
	}
	return names
}

// //call method
// func SafeCall(fn reflect.Value, args []reflect.Value) (result []reflect.Value, err error) {
// 	defer func() {
// 		if x := recover(); x != nil {
// 			err = RaiseError("SafeCall", "Call "+ToS(fn), x)
// 		}
// 	}()
// 	return fn.Call(args), nil
// }
