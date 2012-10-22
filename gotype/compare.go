// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package gotype

import (
	"fmt"
	"reflect"
)

// type code int
// var NOT code = -128

func mustCanCompare(a, b reflect.Value) {
	if CanCompareValue(a, b) == false {
		panic(newTypeErr(methodName(), fmt.Sprintf("can not compare %s and %s", a.Type(), b.Type()), nil))
	}
}

// a == b?, just compare simple data type
func Equal(a interface{}, b interface{}) bool {

	if a == b {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	avalue, bvalue := Underlying(reflect.ValueOf(a)), Underlying(reflect.ValueOf(b))
	mustCanCompare(avalue, bvalue)
	akind, bkind := avalue.Kind(), bvalue.Kind()

	//fmt.Printf("equal %v == %v ?; value %s, %s; kind %s, %s \n", a, b, avalue, bvalue, akind, bkind)

	switch {
	case IsBool(akind) && IsBool(bkind):
		return avalue.Bool() == bvalue.Bool()
	case IsInt(akind) && IsInt(bkind):
		return avalue.Int() == bvalue.Int()
	case IsUint(akind) && IsUint(bkind):
		return avalue.Uint() == bvalue.Uint()
	case IsFloat(akind) && IsFloat(bkind):
		return avalue.Float() == bvalue.Float()
	case IsString(akind) && IsString(bkind):
		return avalue.String() == bvalue.String()
	case CanPointer(akind) && CanPointer(bkind):
		return avalue.Pointer() == bvalue.Pointer()
	case IsInt(akind) && IsNumeric(bkind):
		if i, err := ToInt(bvalue); err == nil {
			return avalue.Int() == i
		} else {
			return false
		}
	case IsUint(akind) && IsNumeric(bkind):
		if i, err := ToUint(bvalue); err == nil {
			return avalue.Uint() == i
		} else {
			return false
		}
	case IsFloat(akind) && IsNumeric(bkind):
		if i, err := ToFloat(bvalue); err == nil {
			return avalue.Float() == i
		} else {
			return false
		}
	}
	return a == b
}

// a > b?, just compare simple data type
func Greater(a interface{}, b interface{}) (result bool) {
	if a == b {
		return false
	}

	avalue, bvalue := Underlying(reflect.ValueOf(a)), Underlying(reflect.ValueOf(b))
	mustCanCompare(avalue, bvalue)
	akind, bkind := avalue.Kind(), bvalue.Kind()

	switch {
	case IsBool(akind) && IsBool(bkind):
		return BToi(avalue.Bool()) > BToi(bvalue.Bool())
	case IsInt(akind) && IsInt(bkind):
		return avalue.Int() > bvalue.Int()
	case IsUint(akind) && IsUint(bkind):
		return avalue.Uint() > bvalue.Uint()
	case IsFloat(akind) && IsFloat(bkind):
		return avalue.Float() > bvalue.Float()
	case IsString(akind) && IsString(bkind):
		return avalue.String() > bvalue.String()
	case IsInt(akind) && IsNumeric(bkind):
		if i, err := ToInt(bvalue); err == nil {
			return avalue.Int() > i
		} else {
			return false
		}
	case IsUint(akind) && IsNumeric(bkind):
		if i, err := ToUint(bvalue); err == nil {
			return avalue.Uint() > i
		} else {
			return false
		}
	case IsFloat(akind) && IsNumeric(bkind):
		if i, err := ToFloat(bvalue); err == nil {
			return avalue.Float() > i
		} else {
			return false
		}
	case IsNumeric(akind) && IsNumeric(bkind):
		afloat, aerr := ToFloat(avalue)
		bfloat, berr := ToFloat(bvalue)
		if aerr != nil && berr != nil {
			return afloat > bfloat
		}
	}
	panic(newTypeErr(methodName(), fmt.Sprintf("Greater is not defiend for %s and %s ",
		reflect.TypeOf(a), reflect.TypeOf(b)), nil))
}

// a < b?, just compare simple data type
func Less(a interface{}, b interface{}) bool {
	if a == b {
		return false
	}

	avalue, bvalue := Underlying(reflect.ValueOf(a)), Underlying(reflect.ValueOf(b))
	mustCanCompare(avalue, bvalue)

	akind, bkind := avalue.Kind(), bvalue.Kind()
	switch {
	case IsBool(akind) && IsBool(bkind):
		return BToi(avalue.Bool()) > BToi(bvalue.Bool())
	case IsInt(akind) && IsInt(bkind):
		return avalue.Int() < bvalue.Int()
	case IsUint(akind) && IsUint(bkind):
		return avalue.Uint() < bvalue.Uint()
	case IsFloat(akind) && IsFloat(bkind):
		return avalue.Float() < bvalue.Float()
	case IsString(akind) && IsString(bkind):
		return avalue.String() < bvalue.String()
	case IsInt(akind) && IsNumeric(bkind):
		if i, err := ToInt(bvalue); err == nil {
			return avalue.Int() < i
		} else {
			return false
		}
	case IsUint(akind) && IsNumeric(bkind):
		if i, err := ToUint(bvalue); err == nil {
			return avalue.Uint() < i
		} else {
			return false
		}
	case IsFloat(akind) && IsNumeric(bkind):
		if i, err := ToFloat(bvalue); err == nil {
			return avalue.Float() < i
		} else {
			return false
		}

	case IsNumeric(akind) && IsNumeric(bkind):
		afloat, aerr := ToFloat(avalue)
		bfloat, berr := ToFloat(bvalue)
		if aerr != nil && berr != nil {
			return afloat < bfloat
		}
	}
	panic(newTypeErr(methodName(), fmt.Sprintf("Less is not defiend for %s and %s ",
		reflect.TypeOf(a), reflect.TypeOf(b)), nil))
}

// compare a and b, return 0 (a==b), 1 (a > b), -1 (a<b) 
func Compare(a interface{}, b interface{}) int {
	if Equal(a, b) {
		return 0
	}
	if Greater(a, b) {
		return 1
	}
	if Less(a, b) {
		return -1
	}

	panic(newTypeErr(methodName(), fmt.Sprintf("can not compare %s and %s ",
		reflect.TypeOf(a), reflect.TypeOf(b)), nil))
}

// // slice is equal ?
// func SliceEqual(a interface{}, b interface{}) bool {
// 	panic("SliceEqual todo")
// }

// // Map is equal ?
// func MapEqual(a interface{}, b interface{}) bool {
// 	panic("MapEqual todo")
// }

// func Contains(a interface{}, b interface{}) bool {
// 	panic("Contains todo")
// }

// func ContainsAny(a interface{}, b interface{}) bool {
// 	panic("Contains todo")
// }
