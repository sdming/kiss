// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package gotype_test

import (
	//"fmt"
	"github.com/sdming/kiss/gotype"
	"reflect"
	"testing"
)

func testFieldsValue(t *testing.T, v reflect.Value, fields []string, fn func(v gotype.Value) bool) {
	for _, field := range fields {
		fv := v.FieldByName(field)
		k := gotype.Value(fv)
		if ok := fn(k); !ok {
			t.Errorf("field %s testing fail - %#v ", field, fn)
		}
	}
}

func TestValueIsStruct(t *testing.T) {
	var t1 T1
	var at1 Alias_T1
	var t1_p *T1 = new(T1)
	var at1_p *Alias_T1 = &Alias_T1{}

	test(t, gotype.ValueOf(t1).IsStruct(), "T1  ")
	test(t, gotype.ValueOf(at1).IsStruct(), "Alias_T1 ")
	test(t, gotype.ValueOf(t1_p).IsStruct(), "*T1  ")
	test(t, gotype.ValueOf(at1_p).IsStruct(), "*Alias_T1 ")
}

func TestValueIsBool(t *testing.T) {
	data := newSimpleType()
	var fields = []string{"A_bool", "A_alias_bool", "A_bool_p", "A_alias_bool_p"}
	testFieldsValue(t, reflect.ValueOf(data), fields,
		func(v gotype.Value) bool { return v.IsBool() })
}

func TestValueIsCollect(t *testing.T) {
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

	testFieldsValue(t, reflect.ValueOf(data), fields,
		func(v gotype.Value) bool { return v.IsCollect() })

}

func TestValueIsNumeric(t *testing.T) {
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
	testFieldsValue(t, reflect.ValueOf(data), fields,
		func(v gotype.Value) bool { return v.IsNumeric() })

}

func TestValueIsInt(t *testing.T) {
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
	testFieldsValue(t, reflect.ValueOf(data), fields,
		func(v gotype.Value) bool { return v.IsInt() })
}

func TestValueIsUint(t *testing.T) {
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
	testFieldsValue(t, reflect.ValueOf(data), fields,
		func(v gotype.Value) bool { return v.IsUint() })
}

func TestValueIsFloat(t *testing.T) {
	var fields = []string{
		"A_float32",
		"A_float64",
		"A_float32_p",
		"A_float64_p",
		"A_interface"}

	data := newSimpleType()
	data.A_interface = float32(3.14)
	testFieldsValue(t, reflect.ValueOf(data), fields,
		func(v gotype.Value) bool { return v.IsFloat() })
}

func TestValueIsString(t *testing.T) {
	var fields = []string{
		"A_string",
		"A_string_p",
		"A_interface"}

	data := newSimpleType()
	data.A_interface = "string"
	testFieldsValue(t, reflect.ValueOf(data), fields,
		func(v gotype.Value) bool { return v.IsString() })
}

func TestValueIsSimple(t *testing.T) {
	data := newSimpleType()
	data.A_interface = "interface"
	fields := gotype.Fields(reflect.ValueOf(data).Type())
	testFieldsValue(t, reflect.ValueOf(data), fields,
		func(v gotype.Value) bool { return v.IsSimple() })
}
