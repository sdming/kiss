// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package gotype_test

import (
	"github.com/sdming/kiss/gotype"
	"reflect"
	"testing"
)

func TestCompare(t *testing.T) {
	t1 := newSimpleType()
	t1.A_string = "string"
	t1.A_string_p = &t1.A_string
	t1.A_interface = int(0)

	t2 := newSimpleType()
	t2.A_string = "string"
	t2.A_string_p = &t2.A_string
	t2.A_interface = float64(0)

	for name1, f1 := range gotype.ToMap(reflect.ValueOf(t1)) {
		for name2, f2 := range gotype.ToMap(reflect.ValueOf(&t2)) {
			if gotype.CanCompareKind(gotype.Underlying(f1).Kind(), gotype.Underlying(f2).Kind()) {
				equal := gotype.Equal(f1.Interface(), f2.Interface())
				if !equal {
					t.Errorf("Test Equal fail for %s (%s) and %s (%s)", name1, f1, name2, f2)
				}

				greater := gotype.Greater(f1.Interface(), f2.Interface())
				if greater {
					t.Errorf("Test Greater fail for %s (%s) and %s (%s)", name1, f1, name2, f2)
				}

				less := gotype.Less(f1.Interface(), f2.Interface())
				if less {
					t.Errorf("Test Less fail for %s (%s) and %s (%s)", name1, f1, name2, f2)
				}
			}
		}
	}
}
