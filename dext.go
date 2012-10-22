// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package kiss

import (
	"errors"
	"fmt"
	"github.com/sdming/kiss/gotype"
	"reflect"
)

// StructValue is value of struct 
type StructValue reflect.Value

// MapValue is value of map[string]interface{}
type MapValue reflect.Value

func mustStruct(v reflect.Value) {
	if v = gotype.Underlying(v); v.Kind() != reflect.Struct {
		panic(errors.New(fmt.Sprintf("the value is not a struct, it's %s", v.Kind())))
	}
}

func (s StructValue) value() reflect.Value {
	v := reflect.Value(s)
	mustStruct(v)
	return v
}

// Fields() get fields name slice of a struct
func (s StructValue) Fields() []string {
	return gotype.Fields(s.value().Type())
}

// Get return field value by name
func (s StructValue) Get(name string) (output reflect.Value, ok bool) {
	if fv := s.value().FieldByName(name); fv.IsValid() {
		return fv, true
	}
	return
}

// Set can set field value by name 
func (s StructValue) Set(name string, value reflect.Value) (ok bool) {
	fv := s.value().FieldByName(name)

	if !fv.IsValid() || !fv.CanSet() || !value.IsValid() {
		return false
	}
	gotype.Value(fv).Set(value)
	return true
}
