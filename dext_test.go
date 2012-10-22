// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package kiss_test

import (
	"github.com/sdming/kiss"
	"reflect"
	"testing"
)

func TestStructValue(t *testing.T) {
	var data SimpleType
	v := kiss.StructValue(reflect.ValueOf(data))

	typ := reflect.TypeOf(data)

	for i := 0; i < typ.NumField(); i++ {
		name := typ.Field(i).Name
		except := reflect.ValueOf(data).FieldByName(name).Interface()
		actual, ok := v.Get(name)
		if !ok {
			t.Errorf("get field %s %s", name, ok)
		} else if except != actual.Interface() {
			t.Errorf("get field %s fail, except %#v actual %#v", name, except, actual)
		}
	}
}
