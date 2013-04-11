// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package kson_test

import (
	"encoding/json"
	"github.com/sdming/kiss/kson"
	"testing"
)

/*
	cpuprofile := "cpu.pprof"
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			fmt.Println("create cpu pprof file error", err)
			return
		}
		fmt.Println("write pprof file", cpuprofile)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}


memprofile := "mem.pprof"
	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			fmt.Println(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

*/

type T1 struct {
	T1_string string
	T1_int    int
	T1_bool   bool
	T1_map    map[string]*T2
	T1_t3     *T3
}

type T2 struct {
	T2_float       float32
	T2_array_uint8 [3]uint
	A_array_T3     [2]T3
}

type T3 struct {
	T3_int    int
	T3_string string
}

type Config struct {
	Log_Level string
	Listen    uint
	Roles     []Role
	Db_Log    Db
	Env       map[string]string
}

type Role struct {
	Name  string
	Allow []string
	Deny  []string
}

type Db struct {
	Driver   string
	Host     string
	User     string
	Password string
	Database string
}

type Poco struct {
	Name   string
	Int    int
	Float  float32
	Bool   bool
	Date   string
	String string
	Quote  string
	Json   string
	Xml    string
	Empty  string
}

var defaultConfig *Config = &Config{
	Log_Level: "debug",
	Listen:    8000,
	Roles: []Role{
		Role{
			Name: "user",
			Allow: []string{
				"/user",
				"/order"},
		},
		Role{
			Name: "*",
			Deny: []string{
				"/user",
				"/order"},
		},
	},
	Db_Log: Db{
		Driver:   "mysql",
		Host:     "127.0.0.1",
		User:     "user",
		Password: "Password",
		Database: "log",
	},
	Env: map[string]string{
		"auth":    `http://auth.io`,
		"browser": "ie, chrome, firefox, safari",
	},
}

var defaultConfigString string = `
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
			auth:		http://auth.io
			browser:	ie, chrome, firefox, safari
			key:
		}
	}	
	`

func BenchmarkMarshal(b *testing.B) {

	config := defaultConfig

	for i := 0; i < b.N; i++ {
		_, err := kson.Marshal(config)
		if err != nil {
			b.Error(err)
		}
	}

}

func BenchmarkUnmarshal(b *testing.B) {
	var data string = defaultConfigString
	var newConfig Config
	for i := 0; i < b.N; i++ {
		err := kson.Unmarshal([]byte(data), &newConfig)
		if err != nil {
			b.Error("config unmarshal error", err)
		}
	}

}

func BenchmarkJsonMarshal(b *testing.B) {

	config := defaultConfig
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(config)
		if err != nil {
			b.Error(err)
		}
	}

}

var defaultJsonConfingString = `{"Log_Level":"debug","Listen":8000,"Roles":[{"Name":"user","Allow":["/user","/order"],"Deny":null},{"Name":"*","Allow":null,"Deny":["/user","/order"]}],"Db_Log":{"Driver":"mysql","Host":"127.0.0.1","User":"user","Password":"Password","Database":"log"},"Env":{"auth":"http://auth.io","browser":"ie, chrome, firefox, safari"}}`

func BenchmarkJsonUnmarshal(b *testing.B) {
	var data string = defaultJsonConfingString
	var newConfig Config
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal([]byte(data), &newConfig)
		if err != nil {
			b.Error("config unmarshal error", err)
		}
	}
}
