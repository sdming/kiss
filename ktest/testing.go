// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package ktest

import (
	"fmt"
	"github.com/sdming/kiss/gotype"
	"strings"
	"testing"
)

// type KT struct {
// 	testing.T
// }

// func NewKT(t testing.T) *KT {
// 	return &KT{t}
// }

func Equal(t *testing.T, name string, expect, actual interface{}) {
	if !gotype.Equal(expect, actual) {
		t.Errorf("%s test equal fail: expect=%v; actual=%v", name, expect, actual)
	}
}

func EqualAsString(t *testing.T, name string, expect, actual interface{}) {
	if strings.TrimSpace(fmt.Sprint(actual)) != strings.TrimSpace(fmt.Sprint(expect)) {
		t.Errorf("%s test equal fail: expect=%v; actual=%v", name, expect, actual)
	}
}
