// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package gotype

import (
	"reflect"
)

type MethodInfo struct {
	Method reflect.Method

	// Type
	Type reflect.Type

	// Value  
	Func reflect.Value

	// NumIn is number of input paramter
	NumIn int

	// NumOut is number of output paramter
	NumOut int

	// Ins is input parameter type
	In []reflect.Type

	// Outs is output parameter type
	Out []reflect.Type
}

func (m *MethodInfo) String() string {
	return m.Type.String()
}

func GetMethodInfo(method reflect.Method) *MethodInfo {

	typ := method.Type
	info := &MethodInfo{
		Method: method,
		Type:   typ,
		Func:   method.Func,
		NumIn:  typ.NumIn(),
		NumOut: typ.NumOut(),
		In:     make([]reflect.Type, typ.NumIn()),
		Out:    make([]reflect.Type, typ.NumOut()),
	}

	for j := 0; j < typ.NumIn(); j++ {
		info.In[j] = typ.In(j)
	}
	for j := 0; j < typ.NumOut(); j++ {
		info.Out[j] = typ.Out(j)
	}

	return info
}

func GetMethodInfoByValue(fnValue reflect.Value) *MethodInfo {

	typ := fnValue.Type()
	info := &MethodInfo{
		Type:   typ,
		Func:   fnValue,
		NumIn:  typ.NumIn(),
		NumOut: typ.NumOut(),
		In:     make([]reflect.Type, typ.NumIn()),
		Out:    make([]reflect.Type, typ.NumOut()),
	}

	for j := 0; j < typ.NumIn(); j++ {
		info.In[j] = typ.In(j)
	}
	for j := 0; j < typ.NumOut(); j++ {
		info.Out[j] = typ.Out(j)
	}

	return info
}
