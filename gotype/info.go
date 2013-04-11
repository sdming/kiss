// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package gotype

import (
	"reflect"
	"strings"
)

type MethodInfo struct {
	// Name 
	Name string

	// NameLowe is name in lower
	NameLower string

	// PkgPath is type's package path
	PkgPath string

	// Type
	Type reflect.Type

	// Value  
	Value reflect.Value

	// NumIn is number of input paramter
	NumIn int

	// NumOut is number of output paramter
	NumOut int

	// Ins is input parameter type
	In []reflect.Type

	// Outs is output parameter type
	Out []reflect.Type

	//Index for Type.Method
	//Index

	//method reflect.Method
}

func (m *MethodInfo) String() string {
	return m.Type.String()
}

func GetMethodInfo(v reflect.Value) *MethodInfo {

	typ := v.Type()
	info := &MethodInfo{
		Name:      typ.Name(),
		NameLower: strings.ToLower(typ.Name()),
		PkgPath:   typ.PkgPath(),
		Type:      typ,
		Value:     v,
		NumIn:     typ.NumIn(),
		NumOut:    typ.NumOut(),
		In:        make([]reflect.Type, typ.NumIn()),
		Out:       make([]reflect.Type, typ.NumOut()),
	}

	for j := 0; j < typ.NumIn(); j++ {
		info.In[j] = typ.In(j)
	}
	for j := 0; j < typ.NumOut(); j++ {
		info.Out[j] = typ.Out(j)
	}

	return info
}

// // StructInfo is 
// type StructInfo struct {
// 	// Name 
// 	Name string

// 	// NameLowe is name in lower
// 	NameLower string

// 	// PkgPath is type's package path
// 	PkgPath string

// 	// Type
// 	Type reflect.Type

// 	// Value
// 	//Value reflect.Value

// 	// NumField is number of field
// 	NumField int

// 	// Fields is 
// 	Fields []reflect.StructField
// }

// // FieldInfo is 
// type FieldInfo struct {
// 	// Name 
// 	Name string

// 	// NameLowe is name in lower
// 	NameLower string

// 	// Type
// 	Type reflect.Type

// 	// Kind
// 	//Kind reflect.Kind
// }

// func GetStructInfo(typ reflect.Type) *StructInfo {
// 	utype := UnderlyingType(typ)

// 	info := &StructInfo{
// 		Name:      typ.Name(),
// 		NameLower: strings.ToLower(typ.Name()),
// 		PkgPath:   typ.PkgPath(),
// 		Type:      typ,
// 		NumField:  utype.NumField(),
// 	}

// 	n := utype.NumField()
// 	for i := 0; i < n; i++ {
// 		ftype := utype.Field(i)
// 		finfo := FieldInfo{
// 			Name:      ftype.Name(),
// 			NameLower: strings.ToLower(ftype.Name()),
// 			Type:      ftype,
// 		}
// 	}

// 	return info
// }
