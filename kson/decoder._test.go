// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package kson_test

import (
	"fmt"
	"github.com/sdming/kiss/kson"
	"reflect"
	"strings"
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
		if strings.TrimSpace(actual) != strings.TrimSpace(fmt.Sprint(fv.Interface())) {
			t.Errorf("poco unmarshalerror, field:[%s]\n actual:%s \n expect:%s \n", name, actual, fmt.Sprint(fv.Interface()))
		}
	}

	if n.ChildBool("Bool") != true {
		t.Errorf("parse bool fail")
	}

	if n.ChildInt("Int") != -1024 {
		t.Errorf("parse int fail")
	}

	if n.ChildFloat("Float") != 6.4 {
		t.Errorf("parse float fail")
	}

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

		if strings.TrimSpace(fmt.Sprint(rv.FieldByName(name).Interface())) !=
			strings.TrimSpace(fmt.Sprint(pv.FieldByName(name).Interface())) {
			t.Errorf("poco unmarshal error, field [%s]", name)
		}
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
		if strings.TrimSpace(actual) != strings.TrimSpace(v) {
			t.Errorf("list parse error, index:[%d]\n actual:%s \n expect:%s \n", i, actual, v)
		}
	}

	newList := []string{}
	err = n.Value(&newList)
	if err != nil {
		t.Error("list unmarshal error", err)
		return
	}
	for i, v := range expect {
		actual := newList[i]
		if strings.TrimSpace(actual) != strings.TrimSpace(v) {
			t.Errorf("lis unmarshal error, index:[%d]\n actual:%s \n expect:%s \n", i, actual, v)
		}
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
		if strings.TrimSpace(actual) != strings.TrimSpace(v) {
			t.Errorf("hash parse error, key:[%s]\n actual:%s \n expect:%s \n", key, actual, v)
		}
	}

	newMap := map[string]string{}
	err = n.Value(newMap)
	if err != nil {
		t.Error("hash unmarshal error")
		return
	}
	for key, v := range newMap {
		if strings.TrimSpace(expect[key]) != strings.TrimSpace(v) {
			t.Errorf("hask unmarshal error, key:[%s]\n actual:%s \n expect:%s \n", key, v, expect[key])
		}
	}

}

func TestParseConfig(t *testing.T) {
	config := `
	{	
		Log_Level:	debug
		Listen:		8000

		Roles: [
			{
				Name:	user
				Allow:	[
					/user		
					/order
				]
			} 
			{
				Name:	*				
				Deny: 	[
					/user
					/order
				]
			} 
		]

		Db_Log:	{
			Driver:		mysql			
			Host: 		127.0.0.1
			User:		user
			Password:	password
			Database:	log
		}

		Env:	{
			auth:		http://auth.io]
			browser:	ie, chrome, firefox, safari
			key:
		}
	}	
	`
	n, err := kson.Parse([]byte(config))
	t.Log(n)
	if err != nil {
		t.Error("config parse", err)
		return
	}

}
