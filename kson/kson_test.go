// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package kson_test

import (
	_ "github.com/sdming/kiss/kson"
)

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
