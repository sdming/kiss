// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package kson_test

import (
	"github.com/sdming/kiss/kson"
	"reflect"
	"testing"
)

func testMarshalSimpleType(t *testing.T, expect string, a interface{}) {
	b, err := kson.Marshal(a)
	if err != nil {
		t.Errorf("marshal fail: %v", a)
		return
	}

	s := string(b)
	t.Log(s)
	if s != expect {
		t.Errorf("Marshal fail: %v; expect %s; actual %s \n", reflect.TypeOf(a), expect, s)
	}
}

func TestMarshalSlice(t *testing.T) {
	var data []interface{} = []interface{}{
		true,
		1024,
		-1,
		[]interface{}{
			true,
			1024,
			"hello",
			-1},
		[3]string{
			"one",
			"two",
			"three"},
	}

	b, err := kson.Marshal(data)
	t.Log(string(b))
	if err != nil {
		t.Error(err)
	}

}

func TestMarshalMap(t *testing.T) {
	data := map[string]interface{}{
		"bool": true,
		"uint": 1024,
		"int":  -1,
		"slice": []interface{}{
			true,
			1024,
			"hello",
			-1},
		"array": [3]string{
			"one",
			"two",
			"three"},
		"map": map[string]string{
			"a": "a_1",
			"b": "b_1",
		},
		"string": "hello",
	}

	b, err := kson.Marshal(data)
	t.Log(string(b))
	if err != nil {
		t.Error(err)
	}

}

func TestMarshalStruct(t *testing.T) {
	var t1 T1
	t1.T1_map = make(map[string]*T2, 2)
	t1.T1_map["a"] = &T2{}
	t1.T1_map["b"] = &T2{}

	b, err := kson.Marshal(t1)
	t.Log(string(b))
	if err != nil {
		t.Error(err)
	}

}

func TestMarshalSimpleType(t *testing.T) {

	testMarshalSimpleType(t, "true", bool(true))
	testMarshalSimpleType(t, "-1", int(-1))
	testMarshalSimpleType(t, "-8", int8(-8))
	testMarshalSimpleType(t, "-16", int16(-16))
	testMarshalSimpleType(t, "-32", int32(-32))
	testMarshalSimpleType(t, "-64", int64(-64))
	testMarshalSimpleType(t, "1", uint(1))
	testMarshalSimpleType(t, "8", uint8(8))
	testMarshalSimpleType(t, "16", uint16(16))
	testMarshalSimpleType(t, "32", uint32(32))
	testMarshalSimpleType(t, "64", int64(64))
	testMarshalSimpleType(t, "3.2", float32(3.2))
	testMarshalSimpleType(t, "-6.4", float64(-6.4))
	testMarshalSimpleType(t, "string", string("string"))

}

func TestMarshalConfig(t *testing.T) {
	config := defaultConfig
	b, err := kson.Marshal(config)
	t.Log(string(b))
	if err != nil {
		t.Error(err)
	}
}

func TestMarshalPoco(t *testing.T) {
	p := Poco{
		Name:   "value",
		Int:    -1024,
		Float:  6.4,
		Bool:   true,
		Date:   "2012-12-21",
		String: `abcdefghijklmnopqrstuvwxyz/字符串 #:[]{}`,
		Quote:  "[0,1,2,3,4,5,6,7,8,9]",
		Json: `
			var father = {
			    "Name": "John",
			    "Age": 32,
			    "Children": [
			        {
			            "Name": "Richard",
			            "Age": 7
			        },
			        {
			            "Name": "Susan",
			            "Age": 4
			        }
			    ]
			};
		`,
		Xml: `
			<root>
			<!-- a node -->
				<text>
					I'll be back
				</text>
			</root>
		`,
		Empty: "",
	}

	b, err := kson.Marshal(p)
	t.Log(string(b))
	if err != nil {
		t.Error(err)
	}

}

func TestMarshalList(t *testing.T) {
	list := []string{
		"line one",
		`[line two]`,
		`

		line three

	`,
	}

	b, err := kson.Marshal(list)
	t.Log(string(b))
	if err != nil {
		t.Error(err)
	}

}

func TestMarshalHash(t *testing.T) {
	hash := map[string]interface{}{
		"int":    1024,
		"bool":   true,
		"string": "string",
		"text": ` 
			I'm not a great programmer, 
			I'm a pretty good programmer with great habits.
		`,
	}

	b, err := kson.Marshal(hash)
	t.Log(string(b))
	if err != nil {
		t.Error(err)
	}

}
