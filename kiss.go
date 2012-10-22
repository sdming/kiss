// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package kiss

import (
	"github.com/sdming/kiss/gotype"
	"reflect"
	"strconv"
)

// value getter
type Getter interface {
	// get field value by name
	Get(name string) (v reflect.Value, ok bool)
}

// string value getter
type StrGetter interface {
	// get field value by name
	Get(name string) (v string, ok bool)
}

// value setter
type Setter interface {
	// set field value by name
	Set(name string, value reflect.Value) (ok bool)

	// get all fields name
	Fields() []string
}

// GetFunc is a func for getting value by name
type GetFunc func(name string) (interface{}, bool)

// implements get 
func (f GetFunc) Get(name string) (reflect.Value, bool) {
	x, ok := f(name)
	return reflect.ValueOf(x), ok
}

// GetFunc is a func for getting string value by name
type StrGetFunc func(name string) (string, bool)

// implements get 
func (f StrGetFunc) Get(name string) (string, bool) {
	return f(name)
}

// like jqery extend, copy value from src to dest
func Extend(dest Setter, src Getter) {
	if dest == nil || src == nil {
		return
	}

	for _, name := range dest.Fields() {
		if value, ok := src.Get(name); ok {
			dest.Set(name, value)
		}
	}
}

// copy value from src to dest, the dest must be struct
func ExtdStruct(dest reflect.Value, src Getter) {
	if !dest.IsValid() || dest.Kind() != reflect.Struct || src == nil {
		return
	}

	typ := dest.Type()
	num := typ.NumField()
	for i := 0; i < num; i++ {
		field := typ.Field(i)
		x, ok := src.Get(field.Name)
		if !ok || !x.IsValid() {
			continue
		}

		v := gotype.Value(dest.Field(i))
		v.Set(x)
	}

}

// copy value from src to dest, struct must be struct
func ParseStruct(dest reflect.Value, src StrGetter) {
	if !dest.IsValid() || dest.Kind() != reflect.Struct || src == nil {
		return
	}

	typ := dest.Type()
	num := typ.NumField()
	for i := 0; i < num; i++ {
		field := typ.Field(i)
		str, ok := src.Get(field.Name)
		if !ok || str == "" {
			continue
		}
		value := dest.Field(i)

		switch value.Kind() {
		case reflect.Bool:
			b, err := strconv.ParseBool(str)
			if err == nil {
				value.SetBool(b)
			}
		case reflect.Int:
			i, err := strconv.Atoi(str)
			if err == nil {
				value.SetInt((int64(i)))
			}
		case reflect.Int8:
			i, err := strconv.ParseInt(str, 0, 8)
			if err == nil {
				value.SetInt(i)
			}
		case reflect.Int16:
			i, err := strconv.ParseInt(str, 0, 16)
			if err == nil {
				value.SetInt(i)
			}
		case reflect.Int32:
			i, err := strconv.ParseInt(str, 0, 32)
			if err == nil {
				value.SetInt(i)
			}
		case reflect.Int64:
			i, err := strconv.ParseInt(str, 0, 64)
			if err == nil {
				value.SetInt(i)
			}
		case reflect.Uint:
			i, err := strconv.ParseUint(str, 0, 0)
			if err == nil {
				value.SetUint(i)
			}
		case reflect.Uint8:
			i, err := strconv.ParseUint(str, 0, 8)
			if err == nil {
				value.SetUint(i)
			}
		case reflect.Uint16:
			i, err := strconv.ParseUint(str, 0, 16)
			if err == nil {
				value.SetUint(i)
			}
		case reflect.Uint32:
			i, err := strconv.ParseUint(str, 0, 32)
			if err == nil {
				value.SetUint(i)
			}
		case reflect.Uint64:
			i, err := strconv.ParseUint(str, 0, 64)
			if err == nil {
				value.SetUint(i)
			}
		case reflect.Float32:
			i, err := strconv.ParseFloat(str, 32)
			if err == nil {
				value.SetFloat(i)
			}
		case reflect.Float64:
			i, err := strconv.ParseFloat(str, 64)
			if err == nil {
				value.SetFloat(i)
			}
		case reflect.String:
			value.SetString(str)
		default:
		}
	}
}
