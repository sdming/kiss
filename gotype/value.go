// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package gotype

import (
	"fmt"
	"reflect"
)

// Value is a alias of reflect.Value
type Value reflect.Value

// get underlying value of *struct
func ValueOf(i interface{}) Value {
	return Value(reflect.ValueOf(i))
}

// return string
func (v Value) String() string {
	return fmt.Sprint(v.Value())
}

func (v Value) kindTest(fn func(reflect.Kind) bool) bool {
	fv := v.Value()
	if !fv.IsValid() {
		return false
	}
	return fn(Underlying(fv).Kind())
}

// is struct?
func (v Value) IsStruct() bool {
	return v.kindTest(IsStruct)
}

// is boolean?
func (v Value) IsBool() bool {
	return v.kindTest(IsBool)
}

// is boolean?
func (v Value) IsCollect() bool {
	return v.kindTest(IsCollect)
}

// is numeric type?
func (v Value) IsNumeric() bool {
	return v.kindTest(IsNumeric)
}

// is numeric type?
func (v Value) IsInt() bool {
	return v.kindTest(IsInt)
}

// is numeric type?
func (v Value) IsUint() bool {
	return v.kindTest(IsUint)
}

// is numeric type?
func (v Value) IsFloat() bool {
	return v.kindTest(IsFloat)
}

// is numeric type?
func (v Value) IsString() bool {
	return v.kindTest(IsString)
}

// is numeric type?
func (v Value) IsSimple() bool {
	return v.kindTest(IsSimple)
}

// origin reflect.Value
func (v Value) Value() reflect.Value {
	return reflect.Value(v)
}

// can compare with i ?
func (v Value) CanCompareTo(i reflect.Value) bool {
	fv := v.Value()
	if !fv.IsValid() || !i.IsValid() {
		return false
	}
	a, b := Underlying(fv), Underlying(i)
	return CanCompareKind(a.Kind(), b.Kind())
}

// can call Pointer()
func (v Value) CanPointer() bool {
	return v.kindTest(CanPointer)
}

// can call Elem()?
func (v Value) CanElem() bool {
	return v.kindTest(CanValueElem)
}

// return underlying type
func (v Value) UnderlyingType() reflect.Type {
	return UnderlyingType(v.Value().Type())
}

// return underlying kind
func (v Value) UnderlyingKind() reflect.Kind {
	return UnderlyingKind(v.Value())
}

// return underlying value
func (v Value) Underlying() reflect.Value {
	return Underlying(v.Value())
}

// return fields name
func (v Value) Fields() []string {
	return Fields(v.Value().Type())
}

func setBool(src, dest reflect.Value) {
	if dest.Kind() == reflect.Bool {
		src.SetBool(dest.Bool())
	} else {
		if x, err := ToBool(dest); err == nil {
			src.SetBool(x)
		}
	}
}

func setInt(src, dest reflect.Value) {
	if IsInt(dest.Kind()) {
		src.SetInt(dest.Int())
	} else {
		if x, err := ToInt(dest); err == nil {
			src.SetInt(x)
		}
	}
}

func setString(src, dest reflect.Value) {
	if dest.Kind() == reflect.String {
		src.SetString(dest.String())
	} else {
		if x, err := ToString(dest); err == nil {
			src.SetString(x)
		}
	}
}

func setUint(src, dest reflect.Value) {
	if IsUint(dest.Kind()) {
		src.SetUint(dest.Uint())
	} else {
		if x, err := ToUint(dest); err == nil {
			src.SetUint(x)
		}
	}
}

func setFloat(src, dest reflect.Value) {
	if IsFloat(dest.Kind()) {
		src.SetFloat(dest.Float())
	} else {
		if x, err := ToFloat(dest); err == nil {
			src.SetFloat(x)
		}
	}
}

// set value to x
func (v Value) Set(x reflect.Value) {

	inner := v.Value()
	if !inner.IsValid() || !inner.CanSet() || !inner.IsValid() {
		return
	}

	//fmt.Println("inner= ", UnderlyingKind(inner), inner.IsValid(), inner.CanSet(), inner.IsValid())
	switch UnderlyingKind(inner) {
	case reflect.Bool:
		setBool(inner, x)
	case reflect.Float32, reflect.Float64:
		setFloat(inner, x)
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		setInt(inner, x)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		setUint(inner, x)
	case reflect.String:
		setString(inner, x)
	default:
		inner.Set(x)
	}
}

// prase value from a string
func (v Value) FromStr(str string) {

	inner := v.Value()
	if !inner.IsValid() || !inner.CanSet() || !inner.IsValid() {
		return
	}

	kind := UnderlyingKind(inner)
	if IsSimple(kind) {
		if x, err := Atok(str, kind); err == nil {
			inner.Set(x)
		}
	} else if inner.Kind() == reflect.Ptr {
		if x, err := Atok(str, inner.Elem().Kind()); err == nil {
			inner.Set(x.Addr())
		}
	}
}
