// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package kson

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sdming/kiss/gotype"
	"reflect"
	"runtime"
	"strings"
)

type encoder struct {
	bytes.Buffer
	deep int
}

func (e *encoder) indentOuter() {
	if e.deep > 0 {
		e.deep--
	}
}

func (e *encoder) indentInner() {
	e.deep++
}

func (e *encoder) indent() {
	for i := 0; i < e.deep; i++ {
		e.WriteString(indent)
	}
}

func (e *encoder) visitArray(v reflect.Value) {
	//fmt.Fprintln(e, "[")
	e.WriteByte('[')
	e.WriteByte('\n')
	e.indentInner()

	n := v.Len()
	for i := 0; i < n; i++ {
		e.indent()
		e.visitReflectValue(v.Index(i))
		//fmt.Fprintln(e, "")
		e.WriteByte('\n')
	}

	e.indentOuter()
	e.indent()
	//fmt.Fprintln(e, "]")
	e.WriteByte(']')
	e.WriteByte('\n')

	return
}

func stringNeedQuote(s string) (b bool, quote string) {
	if s == "" {
		return false, ""
	}
	start := s[0]
	if start == '[' || start == '{' || start == '`' || start == '"' || start == '\t' || start == ' ' ||
		strings.ContainsAny(s, "\r\n") {
		b = true

		if strings.Contains(s, "\"") {
			quote = "`"
		} else {
			quote = "\""
		}
		return
	}

	return false, ""
}

// TODO: adjust indent algorithm
func (e *encoder) visitReflectValue(v reflect.Value) {

	if !v.IsValid() {
		e.WriteString("")
	}

	//v = gotype.Underlying(v)
	kind := v.Kind()
	switch kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		e.WriteString(gotype.Value(v).Format())
	case reflect.String:
		s := v.String()
		if b, quote := stringNeedQuote(s); b {
			e.WriteString(quote)
			e.WriteString(s)
			e.WriteString(quote)
			//e.WriteString("\n")
		} else {
			e.WriteString(s)
		}
	case reflect.Struct:
		//fmt.Fprintln(e, "{")
		e.WriteByte('{')
		e.WriteByte('\n')
		e.indentInner()

		typ := v.Type()
		count := typ.NumField()
		for i := 0; i < count; i++ {
			f := typ.Field(i)

			if f.PkgPath != "" {
				continue
			}

			e.indent()
			//fmt.Fprint(e, f.Name)
			e.WriteString(f.Name)
			//fmt.Fprint(e, ":") 
			e.WriteByte(':')
			e.visitReflectValue(v.Field(i))
			//fmt.Fprintln(e, "")
			e.WriteByte('\n')
		}

		e.indentOuter()
		e.indent()
		//fmt.Fprintln(e, "}")
		e.WriteByte('}')
		e.WriteByte('\n')
	case reflect.Map:
		if !gotype.IsSimple(v.Type().Key().Kind()) {
			return
		}

		if v.IsNil() {
			//fmt.Fprint(e, "") //fmt.Fprint(e, "null")
			e.WriteByte(' ')
			break
		}

		//fmt.Fprintln(e, "{")
		e.WriteByte('{')
		e.WriteByte('\n')
		e.indentInner()

		keys := v.MapKeys()
		for _, k := range keys {
			e.indent()
			//fmt.Fprint(e, k)
			e.WriteString(k.String())
			//fmt.Fprint(e, ":") 
			e.WriteByte(':')
			e.visitReflectValue(v.MapIndex(k))
			//fmt.Fprintln(e, "") 
			e.WriteByte('\n')
		}
		e.indentOuter()
		e.indent()
		//fmt.Fprintln(e, "}")
		e.WriteByte('}')
		e.WriteByte('\n')
	case reflect.Slice:

		if v.IsNil() {
			//fmt.Fprint(e, "") //fmt.Fprint(e, "null")
			e.WriteByte(' ')
			break
		}
		e.visitArray(v)
	case reflect.Array:
		e.visitArray(v)
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() || !v.IsValid() {
			//fmt.Fprint(e, "") //fmt.Fprint(e, "null")
			e.WriteByte(' ')
			return
		}
		e.visitReflectValue(v.Elem())
	default:
		fmt.Fprint(e, v.Interface())
		//return errors.New("Unsupported type " + v.Type().String())
	}
	return
}

// Marshal returns the kson encoding of v.
func Marshal(a interface{}) (data []byte, err error) {

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = errors.New(fmt.Sprint(r))
			}
		}
	}()

	encoder := &encoder{}
	encoder.visitReflectValue(reflect.ValueOf(a))
	if err != nil {
		return nil, err
	}
	return encoder.Bytes(), nil

}
