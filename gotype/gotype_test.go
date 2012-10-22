// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package gotype_test

import (
	"fmt"
	"github.com/sdming/kiss/gotype"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type T1 struct {
	T1_x int
	T1_y int
}

type T2 struct {
	T2_x *int
	T2_y *int
}

type T3 struct {
	T1
	*T2
}

type Alias_T1 T1
type AliasInt int
type AliasBool bool

type SimpleType struct {
	A_uint8    uint8
	A_uint16   uint16
	A_uint32   uint32
	A_uint64   uint64
	A_unit     uint
	A_uint8_p  *uint8
	A_uint16_p *uint16
	A_uint32_p *uint32
	A_uint64_p *uint64
	A_unit_p   *uint

	A_int8    int8
	A_int16   int16
	A_int32   int32
	A_int64   int64
	A_int     int
	A_int8_p  *int8
	A_int16_p *int16
	A_int32_p *int32
	A_int64_p *int64
	A_int_p   *int

	A_float32   float32
	A_float64   float64
	A_float32_p *float32
	A_float64_p *float64

	A_byte         byte
	A_rune         rune
	A_bool         bool
	A_alias_bool   AliasBool
	A_alias_int    AliasInt
	A_string       string
	A_byte_p       *byte
	A_rune_p       *rune
	A_bool_p       *bool
	A_alias_bool_p *AliasBool
	A_alias_int_p  *AliasInt
	A_string_p     *string

	A_interface interface{}
}

func newSimpleType() SimpleType {
	var t SimpleType
	setPointer(reflect.ValueOf(&t))
	return t
}

type TypeStruct struct {
	SimpleType

	A_complex64    complex64
	A_complex128   complex128
	A_complex64_p  *complex64
	A_complex128_p *complex128

	A_unintptr uintptr

	A_t1         T1
	A_t1_p       *T1
	A_alias_t1   Alias_T1
	A_alias_t1_p *Alias_T1

	A_array_byte      [32]byte
	A_array_byte_p    *[32]byte
	A_array_bytep     [32]*byte
	A_array_t1        [32]T1
	A_array_t1_p      *[32]T1
	A_array_t1p       [32]*T1
	A_array_int_2     [2][2]int
	A_array_float64_3 [3][3][3]int
	A_array_interface [3]interface{}

	A_slice_byte      []byte
	A_slice_bytep     []*byte
	A_slice_t1        []T1
	A_slice_t1p       []*T1
	A_slice_interface []interface{}

	A_map_string_int       map[string]int
	A_map_intp_struct      map[*int]struct{ x, y float64 }
	A_map_string_interface map[string]interface{}

	A_chan_t1      chan T1
	A_chan_float64 chan<- float64
	A_chan_int     <-chan int

	_      int // pading
	A_func func(x, y int) (bool, error)
}

func (t TypeStruct) A_F1() {
}

func (t *TypeStruct) A_F2_P(x int) bool {
	return x > 0
}

func (t TypeStruct) A_F3(a, _ int, z float32) bool {
	return a > 0 && z > 0
}

func (t TypeStruct) A_F4(s string, values ...int) bool {
	return true
}

func F1() {
}

func F2(x int) bool {
	return x > 0
}

func F3(a, _ int, z float32) bool {
	return a > 0 && z > 0
}

func F4(s string, values ...int) bool {
	return true
}

type KindAssert func(k reflect.Kind) bool

var (
	debug bool = false
)

func methodNameN(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown method"
	}
	return f.Name()
}

func setPointer(v reflect.Value) {
	v = v.Elem()

	for i := 0; i < v.Type().NumField(); i++ {
		name := v.Type().Field(i).Name
		if strings.Contains(name, "_p") {
			p := v.FieldByName(name)
			f := v.FieldByName(strings.Replace(name, "_p", "", 1))
			if f.IsValid() && p.IsValid() && f.CanAddr() && p.CanSet() {
				p.Set(f.Addr())
			}
		}
	}
}

func newTypeStruct() TypeStruct {
	var t TypeStruct
	setPointer(reflect.ValueOf(&t))
	return t
}

func test(t *testing.T, b bool, msg string) {
	if !b {
		name := methodNameN(2)
		t.Errorf("%s fail : %s \n", name, msg)
	}
}

func testFieldsKind(t *testing.T, v reflect.Value, fields []string, fn KindAssert) {
	for _, field := range fields {
		fv := v.FieldByName(field)
		if !fv.IsValid() {
			continue
		}
		k := gotype.UnderlyingKind(fv)
		if debug {
			fmt.Println(field, fv)
		}

		if ok := fn(k); !ok {
			t.Errorf("field %s (kind %s) testing fail - %#v ", field, k, fn)
		}
	}
}

func testValueKind(t *testing.T, v reflect.Value, f KindAssert) {
	k := gotype.UnderlyingKind(v)
	if ok := f(k); !ok {
		t.Errorf("value %v (kind %v) test fail %#v ", v, k, f)
	}
}

func TestIsStruct(t *testing.T) {
	var t1 T1
	var at1 Alias_T1
	var t1_p *T1 = new(T1)
	var at1_p *Alias_T1 = &Alias_T1{}

	testValueKind(t, reflect.ValueOf(t1), gotype.IsStruct)
	testValueKind(t, reflect.ValueOf(at1), gotype.IsStruct)
	testValueKind(t, reflect.ValueOf(t1_p), gotype.IsStruct)
	testValueKind(t, reflect.ValueOf(at1_p), gotype.IsStruct)
}

func TestIsBool(t *testing.T) {
	data := newSimpleType()
	var fields = []string{"A_bool", "A_alias_bool", "A_bool_p", "A_alias_bool_p"}
	testFieldsKind(t, reflect.ValueOf(data), fields, gotype.IsBool)
}

func TestIsCollect(t *testing.T) {
	data := newTypeStruct()
	var fields = []string{
		"A_array_byte",
		"A_array_byte_p",
		"A_array_bytep",
		"A_array_t1",
		"A_array_t1_p",
		"A_array_t1p",
		"A_array_int_2",
		"A_array_float64_3",
		"A_array_interface",
		"A_slice_byte",
		"A_slice_bytep",
		"A_slice_t1",
		"A_slice_t1p",
		"A_slice_interface",
		"A_map_string_int",
		"A_map_intp_struct",
		"A_map_string_interface"}

	testFieldsKind(t, reflect.ValueOf(data), fields, gotype.IsCollect)
}

func TestIsNumeric(t *testing.T) {
	var fields = []string{
		"A_uint8",
		"A_uint16",
		"A_uint32",
		"A_uint64",
		"A_unit",
		"A_uint8_p",
		"A_uint16_p",
		"A_uint32_p",
		"A_uint64_p",
		"A_unit_p",
		"A_int8",
		"A_int16",
		"A_int32",
		"A_int64",
		"A_int",
		"A_int8_p",
		"A_int16_p",
		"A_int32_p",
		"A_int64_p",
		"A_int_p",
		"A_float32",
		"A_float64",
		"A_float32_p",
		"A_float64_p",
		"A_byte",
		"A_rune",
		"A_alias_int",
		"A_byte_p",
		"A_rune_p",
		"A_alias_int_p",
		"A_interface"}

	data := newSimpleType()
	data.A_interface = 101
	testFieldsKind(t, reflect.ValueOf(data), fields, gotype.IsNumeric)
}

func TestIsInt(t *testing.T) {
	var fields = []string{
		"A_int8",
		"A_int16",
		"A_int32",
		"A_int64",
		"A_int",
		"A_int8_p",
		"A_int16_p",
		"A_int32_p",
		"A_int64_p",
		"A_int_p",
		"A_rune",
		"A_rune_p",
		"A_alias_int",
		"A_alias_int_p",
		"A_interface"}

	data := newSimpleType()
	data.A_interface = int64(101)
	testFieldsKind(t, reflect.ValueOf(data), fields, gotype.IsInt)
}

func TestIsUint(t *testing.T) {
	var fields = []string{
		"A_uint8",
		"A_uint16",
		"A_uint32",
		"A_uint64",
		"A_unit",
		"A_uint8_p",
		"A_uint16_p",
		"A_uint32_p",
		"A_uint64_p",
		"A_unit_p",
		"A_interface"}

	data := newSimpleType()
	data.A_interface = uint64(101)
	testFieldsKind(t, reflect.ValueOf(data), fields, gotype.IsUint)
}

func TestIsFloat(t *testing.T) {
	var fields = []string{
		"A_float32",
		"A_float64",
		"A_float32_p",
		"A_float64_p",
		"A_interface"}

	data := newSimpleType()
	data.A_interface = float32(3.14)
	testFieldsKind(t, reflect.ValueOf(data), fields, gotype.IsFloat)
}

func TestIsString(t *testing.T) {
	var fields = []string{
		"A_string",
		"A_string_p",
		"A_interface"}

	data := newSimpleType()
	data.A_interface = "string"
	testFieldsKind(t, reflect.ValueOf(data), fields, gotype.IsString)
}

func TestIsSimple(t *testing.T) {
	data := newSimpleType()
	data.A_interface = "interface"
	fields := gotype.Fields(reflect.ValueOf(data).Type())
	testFieldsKind(t, reflect.ValueOf(data), fields, gotype.IsSimple)
}

func TestSetAddr(t *testing.T) {
	var data SimpleType
	v := reflect.ValueOf(&data).Elem()

	for i := 0; i < v.Type().NumField(); i++ {
		name := v.Type().Field(i).Name
		if strings.Contains(name, "_p") {
			p := v.FieldByName(name)
			f := v.FieldByName(strings.Replace(name, "_p", "", 1))
			gotype.SetAddr(p, f)
			if !f.IsValid() || (f.Kind() == reflect.Interface && f.IsNil()) {
				t.Errorf("set addr of %s to %s  fail ", strings.Replace(name, "_p", "", 1), name)
			}
		}
	}

}

// func TestFieldByName(t *testing.T) {
// 	var t1 TypeStruct
// 	v1 := reflect.ValueOf(t1)
// 	v2 := reflect.ValueOf(&t1)

// 	for i := 0; i < v1.Type().NumField(); i++ {
// 		name := v1.Type().Field(i).Name
// 		if f, ok := gotype.FieldByName(v1, name); !ok {
// 			t.Errorf("field %s FieldByName exist testing fail, return %s", name, f)
// 		}
// 		if f, ok := gotype.FieldByName(v2, name); !ok {
// 			t.Errorf("field %s FieldByName exist testing fail, return %s", name, f)
// 		}
// 	}

// }

// func TestFieldByNameFold(t *testing.T) {
// 	var t1 TypeStruct
// 	v1 := reflect.ValueOf(t1)
// 	v2 := reflect.ValueOf(&t1)

// 	for i := 0; i < v1.Type().NumField(); i++ {
// 		name := strings.ToLower(v1.Type().Field(i).Name)
// 		if f, ok := gotype.FieldByNameFold(v1, name); !ok {
// 			t.Errorf("field %s TestFieldByNameFold exist testing fail, return %s", name, f)
// 		}
// 		if f, ok := gotype.FieldByNameFold(v2, name); !ok {
// 			t.Errorf("field %s TestFieldByNameFold exist testing fail, return %s", name, f)
// 		}
// 	}
// }

// func TestMethodByName(t *testing.T) {

// 	var methods = []string{"A_F1",
// 		//"A_F2",
// 		"A_F3"}

// 	var t1 TypeStruct
// 	v := reflect.ValueOf(t1)

// 	for _, method := range methods {
// 		if m, ok := gotype.MethodByName(v, method); !ok {
// 			t.Errorf("method %s MethodByName exist testing fail", method, m)
// 		}
// 	}
// }

// func TestMethodByName2(t *testing.T) {

// 	var methods = []string{"A_F1",
// 		"A_F2",
// 		"A_F3"}

// 	var t1 TypeStruct
// 	v := reflect.ValueOf(&t1)

// 	for _, method := range methods {
// 		if m, ok := gotype.MethodByName(v, method); !ok {
// 			t.Errorf("method %s MethodByName exist testing fail", method, m)
// 		}
// 	}
// }

// func TestMethodByNameFold(t *testing.T) {

// 	var methods = []string{"a_F1",
// 		"A_f2",
// 		"a_f3"}

// 	var t1 TypeStruct
// 	v := reflect.ValueOf(&t1)

// 	for _, method := range methods {
// 		if m, ok := gotype.MethodByNameFold(v, method); !ok {
// 			t.Errorf("method %s TestMethodByNameFold exist testing fail", method, m)
// 		}
// 	}

/*
basic data type 
"A_uint8",
"A_uint16",
"A_uint32",
"A_uint64",
"A_unit",
"A_uint8_p",
"A_uint16_p",
"A_uint32_p",
"A_uint64_p",
"A_unit_p",
"A_int8",
"A_int16",
"A_int32",
"A_int64",
"A_int",
"A_int8_p",
"A_int16_p",
"A_int32_p",
"A_int64_p",
"A_int_p",

"A_float32",
"A_float64",
"A_float32_p",
"A_float64_p",

"A_byte",
"A_rune",
"A_byte_p",
"A_rune_p",
"A_alias_int",
"A_alias_int_p",

"A_bool",
"A_alias_bool",
"A_bool_p",
"A_alias_bool_p",

"A_string",		
"A_string_p",

"A_interface"
*/
