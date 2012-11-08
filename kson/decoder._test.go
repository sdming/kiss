// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package kson_test

import (
	"github.com/sdming/kiss/kson"
	"github.com/sdming/kiss/ktest"
	"reflect"
	"testing"
)

func TestParsePoco(t *testing.T) {
	var data string = `
	{
		Name:	value 		
		Int:	-1024		
		Float:	6.4			
		Bool:	true		
		Date: 	2012-12-21 
		String:	abcdefghijklmnopqrstuvwxyz/字符串 #:[]{} 
		Quote:	"[0,1,2,3,4,5,6,7,8,9]"  

		Json: 	` + "`" + `
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
			` + "`" + `		
		Xml: "
				<root>
					<!-- a node -->
					<text>
						I'll be back
					</text>
				</root>
				"
		Empty:	
	}
	`

	n, err := kson.Parse([]byte(data))
	if err != nil {
		t.Error("poco parse error", err)
		return
	}

	expect := Poco{
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
				</root>`,
		Empty: "",
	}

	rv := reflect.ValueOf(expect)
	typ := rv.Type()
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		name := f.Name
		fv := rv.FieldByName(name)

		child, ok := n.Child(name)
		if !ok {
			t.Errorf("can not find field %s of struct", name)
			continue
		}
		actual, err := child.String()
		if err != nil {
			t.Errorf("get string value of filed %s error", name)
			continue
		}
		ktest.EqualAsString(t, name, fv.Interface(), actual)
	}

	ktest.EqualAsString(t, "ChildBool", true, n.ChildBool("Bool"))
	ktest.EqualAsString(t, "ChildInt", -1024, n.ChildInt("Int"))
	ktest.EqualAsString(t, "ChildFloat", 6.4, n.ChildFloat("Float"))
	ktest.EqualAsString(t, "ChildString", "[0,1,2,3,4,5,6,7,8,9]", n.ChildString("Quote"))

	var newPoco Poco
	err = n.Value(&newPoco)
	if err != nil {
		t.Error("poco unmarshal error", err)
		return
	}

	pv := reflect.ValueOf(newPoco)
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		name := f.Name
		ktest.EqualAsString(t, "poco unmarshal "+name, rv.FieldByName(name).Interface(), pv.FieldByName(name).Interface())
	}
}

func TestParseList(t *testing.T) {
	list := `
	[
		line one
		"[line two]"
		"

		line three

		"
	]
	`

	n, err := kson.Parse([]byte(list))
	if err != nil {
		t.Error("list parse error", err)
		return
	}

	expect := []string{
		"line one",
		"[line two]",
		`

		line three

		`,
	}

	for i, v := range expect {
		actual := string(n.List[i].Literal)
		ktest.Equal(t, "list parse", v, actual)
	}

	newList := []string{}
	err = n.Value(&newList)
	if err != nil {
		t.Error("list unmarshal error", err)
		return
	}
	for i, v := range expect {
		actual := newList[i]
		ktest.Equal(t, "list unmarshal", v, actual)
	}
}

func TestParseHash(t *testing.T) {
	hash := `
	{				
		int:	-1024	
		float:	6.4		
		bool:	true	
		string:	string	
		text: 	"
				I'm not a great programmer,
				I'm a pretty good programmer with great habits.
				"
	} 
	`
	n, err := kson.Parse([]byte(hash))

	if err != nil {
		t.Error("hash parse error", err)
		return
	}

	expect := map[string]string{
		"int":    "-1024",
		"float":  "6.4",
		"bool":   "true",
		"string": "string",
		"text": `
				I'm not a great programmer,
				I'm a pretty good programmer with great habits.
				`,
	}

	for key, v := range expect {
		child, ok := n.Child(key)
		if !ok {
			t.Errorf("can not find key %s of hash", key)
			continue
		}
		actual, err := child.String()
		if err != nil {
			t.Errorf("get string value of child %s error", key)
			continue
		}

		ktest.Equal(t, "hash parse", v, actual)

	}

	newMap := map[string]string{}
	err = n.Value(newMap)
	if err != nil {
		t.Error("hash unmarshal error")
		return
	}
	for key, v := range newMap {
		ktest.Equal(t, "hask unmarshal "+key, v, expect[key])
	}

}

func TestParseConfig(t *testing.T) {
	config := defaultConfigString

	n, err := kson.Parse([]byte(config))
	t.Log(n)
	if err != nil {
		t.Error("config parse", err)
		return
	}

	ktest.Equal(t, "Log_Level", "debug", n.ChildString("Log_Level"))
	ktest.Equal(t, "Listen", 8000, n.ChildInt("Listen"))
	ktest.Equal(t, "role-user", "user", n.MustChild("Roles").List[0].ChildString("Name"))
	ktest.Equal(t, "role-user-allow", "/user", n.MustChild("Roles").List[0].MustChild("Allow").List[0].Literal)
	ktest.Equal(t, "dblog-host", "127.0.0.1", n.MustChild("Db_Log").ChildString("Host"))
	ktest.Equal(t, "env-auth", "http://auth.io", n.MustChild("Env").Hash["auth"].Literal)

	var newConfig Config
	err = kson.Unmarshal([]byte(config), &newConfig)
	if err != nil {
		t.Error("config unmarshal error", err)
		return
	}

	ktest.Equal(t, "Log_Level", "debug", newConfig.Log_Level)
	ktest.Equal(t, "Listen", 8000, newConfig.Listen)
	ktest.Equal(t, "role-user", "user", newConfig.Roles[0].Name)
	ktest.Equal(t, "role-user-allow", "/user", newConfig.Roles[0].Allow[0])
	ktest.Equal(t, "dblog-host", "127.0.0.1", newConfig.Db_Log.Host)
	ktest.Equal(t, "env-auth", "http://auth.io", newConfig.Env["auth"])
}

func TestUnmarshalStruct(t *testing.T) {
	data := `	
		 {
			T1_bool:true
			T1_int:10			
			T1_map:{
				a:{
					A_array_T3:[
						{
							T3_int:2031
							T3_string:T3_string_2031
						}		
						{
							T3_int:2032
							T3_string:T3_string_2032
						}		
					]			
					T2_array_uint8:[
						2001
						2002
						2003
					]
					T2_float:2000
				}
		
				b:{
					A_array_T3:[
						{
							T3_int:2131
							T3_string:T3_string_2131
						}
		
						{
							T3_int:2132
							T3_string:T3_string_2132
						}
		
					]					
					T2_array_uint8:[
						2101
						2102
						2103
					]
					T2_float:2100
				}		
			}
		
			T1_t3: {
					T3_int:31
					T3_string:T3_string_31
				}
			T1_string:T1_string
		}
	`
	var v T1
	err := kson.Unmarshal([]byte(data), &v)
	if err != nil {
		t.Error("Unmarshal struct", err)
		return
	}

	ktest.Equal(t, "T1_bool", true, v.T1_bool)
	ktest.Equal(t, "T1_int", 10, v.T1_int)
	ktest.Equal(t, "T1_string", "T1_string", v.T1_string)
	ktest.Equal(t, "T1_t3.T3_int", 31, v.T1_t3.T3_int)
	ktest.Equal(t, "T1_t3.T3_string", "T3_string_31", v.T1_t3.T3_string)
	ktest.Equal(t, "T1_map_a.T2_float", float32(2000), v.T1_map["a"].T2_float)
	ktest.Equal(t, "T1_map_a.T2_array_uint8", [3]uint{2001, 2002, 2003}, v.T1_map["a"].T2_array_uint8)
	ktest.Equal(t, "T1_map_a.A_array_T3.T3_int", 2031, v.T1_map["a"].A_array_T3[0].T3_int)
	ktest.Equal(t, "T1_map_a.A_array_T3.T3_string", "T3_string_2031", v.T1_map["a"].A_array_T3[0].T3_string)

}

// func performance() {

// 	t := newConfig()
// 	var b []byte
// 	var err error

// 	count := 10000

// 	start := time.Now()
// 	for i := 0; i < count; i++ {
// 		b, err = json.Marshal(t)
// 		if err != nil {
// 			fmt.Println("json.Marshal error:", err, b)
// 			return
// 		}
// 	}
// 	fmt.Println("json.Marshal", time.Since(start))

// 	start = time.Now()
// 	for i := 0; i < count; i++ {
// 		var p Config
// 		err = json.Unmarshal(b, &p)
// 		if err != nil {
// 			fmt.Println("json.Unmarshal error:", err)
// 			return
// 		}
// 	}

// 	fmt.Println("json.Unmarshal", time.Since(start))

// 	start = time.Now()
// 	for i := 0; i < count; i++ {
// 		b, err = kson.Marshal(t)
// 		if err != nil {
// 			fmt.Println("kson.Marshal error", err, b)
// 			break
// 		}
// 	}

// 	fmt.Println("kson.Marshal", time.Since(start))

// 	start = time.Now()
// 	for i := 0; i < count; i++ {
// 		var p Config
// 		err = kson.Unmarshal(b, &p)
// 		if err != nil {
// 			fmt.Println("kson.Unmarshal error", err)
// 			break
// 		}
// 	}

// 	fmt.Println("kson.Unmarshal", time.Since(start))
// }
