// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package kiss_test

import (
	"fmt"
	"github.com/sdming/kiss"
	"reflect"
	"testing"
)

type SimpleType struct {
	A_uint8   uint8
	A_uint16  uint16
	A_uint32  uint32
	A_uint64  uint64
	A_unit    uint
	A_int8    int8
	A_int16   int16
	A_int32   int32
	A_int64   int64
	A_int     int
	A_float32 float32
	A_float64 float64
	A_byte    byte
	A_rune    rune
	A_bool    bool
	A_string  string
	//A_interface interface{}
}

var strMap = map[string]string{
	"A_uint8":   "8",
	"A_uint16":  "16",
	"A_uint32":  "32",
	"A_uint64":  "64",
	"A_unit":    "1024",
	"A_int8":    "-8",
	"A_int16":   "-16",
	"A_int32":   "-32",
	"A_int64":   "-64",
	"A_int":     "-1204",
	"A_float32": "3.2",
	"A_float64": "0",
	"A_byte":    "127",
	"A_rune":    "-127",
	"A_bool":    "true",
	"A_string":  "string"}

var map2 = map[string]interface{}{
	"A_uint8":   int(8),
	"A_uint16":  int64(16),
	"A_uint32":  byte(32),
	"A_uint64":  int8(64),
	"A_unit":    int(1024),
	"A_int8":    byte(8),
	"A_int16":   uint(16),
	"A_int32":   uint(32),
	"A_int64":   rune(64),
	"A_int":     int64(1204),
	"A_float32": int(0),
	"A_float64": "0",
	"A_byte":    int64(127),
	"A_rune":    int8(-127),
	"A_bool":    "true",
	"A_string":  false}

func TestExtend(t *testing.T) {
	var t1 SimpleType
	var t2 map[string]string = strMap
	v1 := reflect.ValueOf(&t1).Elem()

	dest := kiss.StructValue(v1)
	var src kiss.GetFunc = func(name string) (interface{}, bool) {
		x, ok := t2[name]
		return x, ok
	}
	kiss.Extend(dest, src)

	for i := 0; i < v1.Type().NumField(); i++ {
		f := v1.Type().Field(i)
		except, ok := t2[f.Name]
		actual := fmt.Sprint(v1.FieldByName(f.Name).Interface())
		if !ok || except != actual {
			t.Errorf("extend field %s fail, except %s, actual %s", f.Name, except, actual)
		}
	}
}

func TestExtend2(t *testing.T) {
	var t1 SimpleType
	t2 := map2
	v1 := reflect.ValueOf(&t1).Elem()

	dest := kiss.StructValue(v1)
	var src kiss.GetFunc = func(name string) (interface{}, bool) {
		x, ok := t2[name]
		return x, ok
	}
	kiss.Extend(dest, src)

	for i := 0; i < v1.Type().NumField(); i++ {
		f := v1.Type().Field(i)
		except, ok := t2[f.Name]
		actual := fmt.Sprint(v1.FieldByName(f.Name).Interface())
		if !ok || fmt.Sprint(except) != actual {
			t.Errorf("extend field %s fail, except %s, actual %s", f.Name, fmt.Sprint(except), actual)
		}
	}
}

func TestExtendStruct(t *testing.T) {
	var t1 SimpleType
	var t2 map[string]string = strMap
	v1 := reflect.ValueOf(&t1).Elem()
	var src kiss.GetFunc = func(name string) (interface{}, bool) {
		x, ok := t2[name]
		return x, ok
	}
	kiss.ExtdStruct(v1, src)

	for i := 0; i < v1.Type().NumField(); i++ {
		f := v1.Type().Field(i)
		except, ok := t2[f.Name]
		actual := fmt.Sprint(v1.FieldByName(f.Name).Interface())
		if !ok || except != actual {
			t.Errorf("extend field %s fail, except %s, actual %s", f.Name, except, actual)
		}
	}
}

func TestExtendStruct2(t *testing.T) {
	var t1 SimpleType
	t2 := strMap
	v1 := reflect.ValueOf(&t1).Elem()
	var src kiss.GetFunc = func(name string) (interface{}, bool) {
		x, ok := t2[name]
		return x, ok
	}
	kiss.ExtdStruct(v1, src)

	for i := 0; i < v1.Type().NumField(); i++ {
		f := v1.Type().Field(i)
		except, ok := t2[f.Name]
		actual := fmt.Sprint(v1.FieldByName(f.Name).Interface())
		if !ok || except != actual {
			t.Errorf("extend field %s fail, except %s, actual %s", f.Name, except, actual)
		}
	}
}

func TestExtendStruct3(t *testing.T) {
	var t1 SimpleType
	var t2 SimpleType

	t2.A_bool = true
	t2.A_byte = byte(127)
	t2.A_float32 = 3.2
	t2.A_float64 = 6.4
	t2.A_int = -1
	t2.A_int16 = -16
	t2.A_int32 = -32
	t2.A_int64 = -64
	t2.A_int8 = -8
	t2.A_string = "string"
	t2.A_uint16 = 16
	t2.A_uint32 = 32
	t2.A_uint64 = 64
	t2.A_uint8 = 8
	t2.A_unit = 1

	v1 := reflect.ValueOf(&t1).Elem()
	v2 := reflect.ValueOf(&t2).Elem()

	var src kiss.GetFunc = func(name string) (interface{}, bool) {
		return v2.FieldByName(name).Interface(), true
	}
	kiss.ExtdStruct(v1, src)

	for i := 0; i < v1.Type().NumField(); i++ {
		f := v1.Type().Field(i)
		except := v2.FieldByName(f.Name).Interface()
		actual := v1.FieldByName(f.Name).Interface()
		if except != actual {
			t.Errorf("extend field %s fail, except %s, actual %s", f.Name, except, actual)
		}
	}

}

func TestParseStruct(t *testing.T) {
	var t1 SimpleType
	var t2 map[string]string = strMap
	v1 := reflect.ValueOf(&t1).Elem()
	var src kiss.StrGetFunc = func(name string) (string, bool) {
		x, ok := t2[name]
		return x, ok
	}
	kiss.ParseStruct(v1, src)

	for i := 0; i < v1.Type().NumField(); i++ {
		f := v1.Type().Field(i)
		except, ok := t2[f.Name]
		actual := fmt.Sprint(v1.FieldByName(f.Name).Interface())
		if !ok || except != actual {
			t.Errorf("extend field %s fail, except %s, actual %s", f.Name, except, actual)
		}
	}
}
