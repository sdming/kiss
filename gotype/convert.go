// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package gotype

import (
	"fmt"
	"reflect"
	"strconv"
)

// convert to int64
func ToInt(input reflect.Value) (output int64, err error) {
	//defer checkError(err)
	if !input.IsValid() {
		return 0, newTypeErr(methodName(), "input is invalid", input)
	}

	input = Underlying(input)
	k := input.Kind()
	switch {
	case IsBool(k):
		if input.Bool() {
			return 1, nil
		}
		return 0, nil
	case IsInt(k):
		return input.Int(), nil
	case IsUint(k):
		return int64(input.Uint()), nil
	case IsFloat(k):
		return int64(input.Float()), nil
	case IsString(k):
		return strconv.ParseInt(input.String(), 0, 32)
	default:
	}
	return 0, newTypeErr(methodName(), "input can not convert to int", input)
}

// convert to uint64
func ToUint(input reflect.Value) (output uint64, err error) {
	//defer checkError(err)
	if !input.IsValid() {
		return 0, newTypeErr(methodName(), "input is invalid", input)
	}

	input = Underlying(input)
	k := input.Kind()
	switch {
	case IsBool(k):
		if input.Bool() {
			return 1, nil
		}
		return 0, nil
	case IsInt(k):
		return uint64(input.Int()), nil
	case IsUint(k):
		return input.Uint(), nil
	case IsFloat(k):
		return uint64(input.Float()), nil
	case IsString(k):
		return strconv.ParseUint(input.String(), 0, 64)
	default:
	}
	return 0, newTypeErr(methodName(), "input can not convert to uint", input)
}

// convert to float64
func ToFloat(input reflect.Value) (output float64, err error) {
	//defer checkError(err)
	if !input.IsValid() {
		return 0, newTypeErr(methodName(), "input is invalid", input)
	}

	input = Underlying(input)
	k := input.Kind()
	switch {
	case IsBool(k):
		if input.Bool() {
			return 1, nil
		}
		return 0, nil
	case IsInt(k):
		return float64(input.Int()), nil
	case IsUint(k):
		return float64(input.Uint()), nil
	case IsFloat(k):
		return input.Float(), nil
	case IsString(k):
		return strconv.ParseFloat(input.String(), 32)
	default:
	}
	return 0, newTypeErr(methodName(), "input can not convert to float", input)
}

// convert to bool
func ToBool(input reflect.Value) (output bool, err error) {
	//defer checkError(err)
	if !input.IsValid() {
		return false, newTypeErr(methodName(), "input is invalid", input)
	}

	input = Underlying(input)
	k := input.Kind()
	switch {
	case IsBool(k):
		return input.Bool(), nil
	case IsInt(k):
		return input.Int() > 0, nil
	case IsUint(k):
		return input.Uint() > 0, nil
	case IsFloat(k):
		return input.Float() > 0, nil
	case IsString(k):
		return strconv.ParseBool(input.String())
	default:
	}
	return false, newTypeErr(methodName(), "input can not convert to bool", input)
}

// convert to string
func ToString(input reflect.Value) (output string, err error) {
	//defer checkError(err)
	if !input.IsValid() {
		return "", nil
	}

	input = Underlying(input)
	switch input.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(input.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(input.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(input.Uint(), 10), nil
	case reflect.String:
		return input.String(), nil
	default:
	}

	return fmt.Sprint(input.Interface()), nil
}

// convert string to reflect.Value accoring to type( use base 0 to parse numeric )
func Atov(str string, typ reflect.Type) (output reflect.Value, err error) {
	if str == "" {
		return reflect.Zero(typ), nil
	}

	if IsSimple(typ.Kind()) {
		return Atok(str, typ.Kind())
	}

	if typ.Kind() == reflect.Ptr {
		if x, e := Atok(str, typ.Elem().Kind()); e != nil {
			output = x.Addr()
			return
		}
	}

	return reflect.Zero(typ), newTypeErr(methodName(), "can not convert string to "+typ.String(), nil)
}

// convert string to reflect.Value acoring to kind( use base 0 to parse numeric )
func Atok(str string, k reflect.Kind) (output reflect.Value, err error) {
	switch k {
	case reflect.Bool:
		b, e := strconv.ParseBool(str)
		return reflect.ValueOf(b), e
	case reflect.Int:
		i, e := strconv.Atoi(str)
		return reflect.ValueOf(i), e
	case reflect.Int8:
		i, e := strconv.ParseInt(str, 0, 8)
		return reflect.ValueOf(int8(i)), e
	case reflect.Int16:
		i, e := strconv.ParseInt(str, 0, 16)
		return reflect.ValueOf(int16(i)), e
	case reflect.Int32:
		i, e := strconv.ParseInt(str, 0, 32)
		return reflect.ValueOf(int32(i)), e
	case reflect.Int64:
		i, e := strconv.ParseInt(str, 0, 64)
		return reflect.ValueOf(i), e
	case reflect.Uint:
		i, e := strconv.ParseUint(str, 0, 0)
		return reflect.ValueOf(uint(i)), e
	case reflect.Uint8:
		i, err := strconv.ParseUint(str, 0, 8)
		return reflect.ValueOf(uint8(i)), err
	case reflect.Uint16:
		i, err := strconv.ParseUint(str, 0, 16)
		return reflect.ValueOf(uint16(i)), err
	case reflect.Uint32:
		i, err := strconv.ParseUint(str, 0, 32)
		return reflect.ValueOf(uint32(i)), err
	case reflect.Uint64:
		i, err := strconv.ParseUint(str, 0, 64)
		return reflect.ValueOf(i), err
	case reflect.Float32:
		i, err := strconv.ParseFloat(str, 32)
		return reflect.ValueOf(float32(i)), err
	case reflect.Float64:
		i, err := strconv.ParseFloat(str, 64)
		return reflect.ValueOf(i), err
	case reflect.String:
		return reflect.ValueOf(str), nil
	default:
	}
	err = newTypeErr(methodName(), "can not convert string to "+k.String(), nil)
	return
}

// convert struct to map[string]reflect.Value
func ToMap(input reflect.Value) map[string]reflect.Value {
	input = Underlying(input)
	if input.Kind() != reflect.Struct {
		panic(fmt.Sprintf("can not convert %v to map", input.Kind()))
	}

	typ := input.Type()
	count := typ.NumField()
	output := make(map[string]reflect.Value, count)

	for i := 0; i < count; i++ {
		f := typ.Field(i)
		if f.PkgPath != "" {
			continue
		}

		v := input.FieldByName(f.Name)
		if v.IsValid() {
			output[f.Name] = v
		}
	}
	return output
}

// convert bool to int, true=1; false=0
func BToi(b bool) int {
	if b {
		return 1
	}
	return 0
}
