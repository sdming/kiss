// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package kson

import (
	"errors"
	"github.com/sdming/kiss/gotype"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"fmt"
)

const (
	indent   string = "\t"
	eof      byte   = byte(0)
	capacity int    = 8
)

const (
	NodeNone = iota
	NodeLiteral
	NodeHash
	NodeList
)

func nameOfNodeType(typ int) string {
	switch typ {
	case NodeLiteral:
		return "literal"
	case NodeHash:
		return "hash"
	case NodeList:
		return "list"
	case NodeNone:
		return "none"
	}

	panic(&InvalidNodeTypeError{typ})
}

// Node 
type Node struct {
	Type    int
	Literal string
	List    []*Node
	Hash    map[string]*Node
}

// type LiteralNode []byte

// type ListNode []*Node

// type HashNode map[string]*Node

type InvalidNodeTypeError struct {
	NodeType int
}

func (e *InvalidNodeTypeError) Error() string {
	return "invalid node type: " + nameOfNodeType(e.NodeType)
}

type NodeNotExistsError struct {
	Name string
}

func (e *NodeNotExistsError) Error() string {
	return e.Name + " is not exists"
}

// Child return child node by name, ok is false if name doesn't exist
func (n *Node) Child(name string) (child *Node, ok bool) {
	if n.Type == NodeHash {
		child, ok = n.Hash[name]
		return
	}
	return
}

// Child return child node by name, ok is false if name doesn't exist
func (n *Node) ChildFold(name string) (child *Node, ok bool) {
	if n.Type != NodeHash {
		return
	}
	child, ok = n.Hash[name]
	if ok {
		return
	}

	for key, x := range n.Hash {
		if strings.EqualFold(name, key) {
			return x, true
		}
	}
	return
}

// ChildInt return child value as int64
func (n *Node) ChildInt(name string) int64 {
	if n, ok := n.Child(name); ok {
		if i, err := n.Int(); err == nil {
			return i
		} else {
			panic(err)
		}
	}
	panic(&NodeNotExistsError{name})
}

// ChildUint return child value as uint64
func (n *Node) ChildUint(name string) uint64 {
	if n, ok := n.Child(name); ok {
		if i, err := n.Uint(); err == nil {
			return i
		} else {
			panic(err)
		}
	}
	panic(&NodeNotExistsError{name})
}

// ChildFloat return child value as float64
func (n *Node) ChildFloat(name string) float64 {
	if n, ok := n.Child(name); ok {
		if i, err := n.Float(); err == nil {
			return i
		} else {
			panic(err)
		}
	}
	panic(&NodeNotExistsError{name})
}

// ChildBool return child value as bool
func (n *Node) ChildBool(name string) bool {
	if n, ok := n.Child(name); ok {
		if i, err := n.Bool(); err == nil {
			return i
		} else {
			panic(err)
		}
	}
	panic(&NodeNotExistsError{name})
}

// ChildString return child value as string
func (n *Node) ChildString(name string) string {
	if n, ok := n.Child(name); ok {
		if i, err := n.String(); err == nil {
			return i
		} else {
			panic(err)
		}
	}
	panic(&NodeNotExistsError{name})
}

// ChildIntOrDefault return child value as int64, return defaultValue if child doesn't exist
func (n *Node) ChildIntOrDefault(name string, defaultValue int64) int64 {
	if n, ok := n.Child(name); ok {
		if i, err := n.Int(); err == nil {
			return i
		}
	}
	return defaultValue
}

// ChildUintOrDefault return child value as uint64, return defaultValue if child doesn't exist
func (n *Node) ChildUintOrDefault(name string, defaultValue uint64) uint64 {
	if n, ok := n.Child(name); ok {
		if i, err := n.Uint(); err == nil {
			return i
		}
	}
	return defaultValue
}

// ChildFloatOrDefault return child value as float64, return defaultValue if child doesn't exist
func (n *Node) ChildFloatOrDefault(name string, defaultValue float64) float64 {
	if n, ok := n.Child(name); ok {
		if i, err := n.Float(); err == nil {
			return i
		}
	}
	return defaultValue
}

// ChildBoolOrDefault return child value as bool, return defaultValue if child doesn't exist
func (n *Node) ChildBoolOrDefault(name string, defaultValue bool) bool {
	if n, ok := n.Child(name); ok {
		if i, err := n.Bool(); err == nil {
			return i
		}
	}
	return defaultValue
}

// ChildStringOrDefault return child value as string, return defaultValue if child doesn't exist
func (n *Node) ChildStringOrDefault(name string, defaultValue string) string {
	if n, ok := n.Child(name); ok {
		if i, err := n.String(); err == nil {
			return i
		}
	}
	return defaultValue
}

// Int returns n's underlying value, as an int64.
func (n *Node) Int() (i int64, err error) {
	if n.Type != NodeLiteral {
		err = &InvalidNodeTypeError{n.Type}
		return
	}
	return strconv.ParseInt(n.Literal, 0, 64)
}

// Uint returns n's underlying value, as an uint64.
func (n *Node) Uint() (i uint64, err error) {
	if n.Type != NodeLiteral {
		err = &InvalidNodeTypeError{n.Type}
		return
	}
	return strconv.ParseUint(n.Literal, 0, 64)
}

// Float returns n's underlying value, as an float64.
func (n *Node) Float() (f float64, err error) {
	if n.Type != NodeLiteral {
		err = &InvalidNodeTypeError{n.Type}
		return
	}
	return strconv.ParseFloat(n.Literal, 64)
}

// Bool returns n's underlying value, as an bool.
func (n *Node) Bool() (b bool, err error) {
	if n.Type != NodeLiteral {
		err = &InvalidNodeTypeError{n.Type}
		return
	}
	return strconv.ParseBool(n.Literal)
}

// String returns n's underlying value, as an string.
func (n *Node) String() (s string, err error) {
	if n.Type != NodeLiteral {
		err = &InvalidNodeTypeError{n.Type}
		return
	}
	return n.Literal, nil
}

// Value unmarshal data to the value pointed to by a.
func (n *Node) Value(a interface{}) (err error) {

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

	v := reflect.ValueOf(a)
	n.set(v)
	return
}

func (n *Node) setArray(v reflect.Value) {

	kind := v.Kind()
	if n.Type != NodeList || (kind != reflect.Slice && kind != reflect.Array) {
		return
	}

	typ := v.Type()
	if n.List == nil {
		if v.CanSet() {
			v.Set(reflect.Zero(typ))
		}
		return
	}

	l := len(n.List)
	if kind == reflect.Slice {
		if v.CanSet() {
			v.Set(reflect.MakeSlice(v.Type(), l, l))
		} else {
			return
		}
	}

	vl := v.Len()
	elemType := typ.Elem()
	elemKind := elemType.Kind()
	simple := gotype.IsSimple(elemKind)

	for i, x := range n.List {
		if i < vl { // capacity of array maybe less of i
			if simple && x.Type == NodeLiteral {
				gotype.Value(v.Index(i)).Parse(x.Literal)
			} else {
				x.set(v.Index(i))
			}
		}
	}

	if kind == reflect.Array && l < vl {
		for i := l; i < vl; i++ {
			z := reflect.Zero(v.Type().Elem()) //reset
			v.Index(i).Set(z)
		}
	}
}

func (n *Node) setMap(v reflect.Value) {

	kind := v.Kind()
	if n.Type != NodeHash || kind != reflect.Map {
		return
	}

	typ := v.Type()
	if n.Hash == nil {
		if v.CanSet() {
			v.Set(reflect.Zero(typ))
		}
		return
	}

	if typ.Key() != gotype.TypeString { // only support string key
		return
	}

	if v.IsNil() {
		if v.CanSet() {
			v.Set(reflect.MakeMap(typ))
		} else {
			return
		}
	}

	elemType := typ.Elem()
	elemKind := elemType.Kind()
	simple := gotype.IsSimple(elemKind)
	for name, x := range n.Hash {
		if simple && x.Type == NodeLiteral {
			mapElem, err := gotype.Atok(x.Literal, elemKind)
			if err == nil {
				v.SetMapIndex(reflect.ValueOf(name), mapElem)
			}
		} else {
			mapElem := reflect.New(elemType)
			x.set(mapElem)
		}
	}
}

func (n *Node) setObject(v reflect.Value) {

	kind := v.Kind()
	if n.Type != NodeHash || kind != reflect.Struct {
		return
	}

	typ := v.Type()
	if n.Hash == nil {
		if v.CanSet() {
			v.Set(reflect.Zero(typ))
			return
		}
	}

	numFiled := typ.NumField()
	for i := 0; i < numFiled; i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue
		}

		name := field.Name
		filedNode, ok := n.Child(name)
		if !ok {
			filedNode, ok = n.ChildFold(name)
		}
		if !ok {
			continue
		}

		fv := v.Field(i)
		if !fv.CanSet() {
			continue
		}

		kind := field.Type.Kind()
		if filedNode.Type == NodeLiteral && gotype.IsSimple(kind) {
			gotype.Value(fv).Parse(filedNode.Literal)
		} else {
			filedNode.set(fv)
		}
		// } else if filedNode.Type == NodeLiteral && kind == reflect.Ptr && gotype.IsSimple(field.Type.Elem().Kind()) {
		// 	if fv.IsNil() {
		// 		fv.Set(reflect.New(field.Type.Elem()))
		// 	}
		// 	gotype.Value(reflect.Indirect(fv)).Parse(filedNode.Literal)
		// } 
	}

}

func (n *Node) set(v reflect.Value) {

	// if !v.CanSet() || !v.IsValid() {
	// 	return
	// }

	kind := v.Kind()
	switch {
	case gotype.IsSimple(kind):
		if n.Literal != "" {
			gotype.Value(v).Parse(n.Literal)
		}
	case kind == reflect.Invalid || kind == reflect.Uintptr || kind == reflect.UnsafePointer || kind == reflect.Func || kind == reflect.Chan || kind == reflect.Complex64 || kind == reflect.Complex128:
		//TODO: unsupport
	case kind == reflect.Array:
		n.setArray(v)
	case kind == reflect.Slice:
		n.setArray(v)
	case kind == reflect.Map:
		n.setMap(v)
	case kind == reflect.Struct:
		n.setObject(v)
	case kind == reflect.Interface:
		if n.Type == NodeLiteral {
			gotype.Value(v).Parse(n.Literal)
		}
	case kind == reflect.Ptr:
		if v.IsNil() && v.CanSet() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		n.set(v.Elem())
	default:
		//TODO:
	}
}
