// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package gotype_test

import (
	"github.com/sdming/kiss/gotype"
	"reflect"
	"testing"
)

func TestToInt(t *testing.T) {
	data := newSimpleType()
	data.A_string = "-1"
	data.A_string_p = &data.A_string
	data.A_interface = "-1"

	v := reflect.ValueOf(data)
	for name, f := range gotype.ToMap(v) {
		_, err := gotype.ToInt(f)
		if err != nil {
			t.Errorf("field %s (%v) convert to int fail, %s", name, f.Interface(), err)
		}
	}
}

func TestToUint(t *testing.T) {
	data := newSimpleType()
	data.A_string = "11"
	data.A_string_p = &data.A_string
	data.A_interface = "11"

	v := reflect.ValueOf(data)
	for name, f := range gotype.ToMap(v) {
		_, err := gotype.ToUint(f)
		if err != nil {
			t.Errorf("field %s (%v) convert to uint fail, %s", name, f.Interface(), err)
		}
	}
}

func TestToFloat(t *testing.T) {
	data := newSimpleType()
	data.A_string = "3.14"
	data.A_string_p = &data.A_string
	data.A_interface = "3.14"

	v := reflect.ValueOf(&data)
	for name, f := range gotype.ToMap(v) {
		_, err := gotype.ToFloat(f)
		if err != nil {
			t.Errorf("field %s (%v) convert to float fail, %s", name, f.Interface(), err)
		}
	}
}

func TestToBool(t *testing.T) {
	data := newSimpleType()
	data.A_string = "true"
	data.A_string_p = &data.A_string
	data.A_interface = "true"

	v := reflect.ValueOf(&data)
	for name, f := range gotype.ToMap(v) {
		_, err := gotype.ToBool(f)
		if err != nil {
			t.Errorf("field %s (%v) convert to bool fail, %s", name, f.Interface(), err)
		}
	}
}

func TestToString(t *testing.T) {
	data := newSimpleType()
	data.A_string = "string"
	data.A_string_p = &data.A_string
	data.A_interface = "string"

	v := reflect.ValueOf(data)
	for name, f := range gotype.ToMap(v) {
		t.Logf("test to string  %s %s %s '\n", name, f, f.Kind())
		_, err := gotype.ToString(f)
		if err != nil {
			t.Errorf("field %s (%v) convert to string fail, %s", name, f.Interface(), err)
		}
	}

}

func TestToMap(t *testing.T) {
	data := newSimpleType()
	data.A_string = "string"
	data.A_string_p = &data.A_string
	data.A_interface = "string"

	v := reflect.ValueOf(data)
	m := gotype.ToMap(v)
	for i := 0; i < v.Type().NumField(); i++ {
		name := v.Type().Field(i).Name
		fv := v.FieldByName(name)
		actual, ok := m[name]
		if !ok {
			t.Errorf("field %s (%v) does not exists", name, fv)
		}
		if actual != fv {
			t.Errorf("convert %s fail, expect %s, actual %v", name, fv, actual)
		}
	}
}

func testFromStr(t *testing.T, str string, expect interface{}) {
	v, e := gotype.Atov(str, reflect.TypeOf(expect))
	if e != nil {
		t.Errorf("convert %s to type %v error %s", str, reflect.TypeOf(expect).Name(), e)
	}
	if v.Interface() != expect {
		t.Errorf("convert %s to type %v fail, expect %s, actual %v",
			str, reflect.TypeOf(expect).Name(), expect, v)
	}
}

func TestAtok(t *testing.T) {
	testFromStr(t, "true", true)
	testFromStr(t, "-1", int(-1))
	testFromStr(t, "-8", int8(-8))
	testFromStr(t, "-16", int16(-16))
	testFromStr(t, "-32", int32(-32))
	testFromStr(t, "-64", int64(-64))
	testFromStr(t, "1", uint(1))
	testFromStr(t, "8", uint8(8))
	testFromStr(t, "16", uint16(16))
	testFromStr(t, "32", uint32(32))
	testFromStr(t, "64", int64(64))
	testFromStr(t, "3.2", float32(3.2))
	testFromStr(t, "-6.4", float64(-6.4))
	testFromStr(t, "string", "string")
}
