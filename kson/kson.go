// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package kson

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sdming/kiss/gotype"
	"io"
	"io/ioutil"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
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

var commentsReg *regexp.Regexp = regexp.MustCompile(`^[\s]*#`)

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
func (n *Node) MustChild(name string) *Node {
	if n.Type == NodeHash {
		child, ok := n.Hash[name]
		if ok {
			return child
		}
	}
	panic("child is not exists:" + name)
}

// Child return child node by name, ok is false if name doesn't exist
func (n *Node) Child(name string) (child *Node, ok bool) {
	if n.Type == NodeHash {
		child, ok = n.Hash[name]
		return
	}
	return
}

// Query return child node by path, just like query
func (n *Node) Query(path string) (child *Node, ok bool) {
	names := strings.Split(path, " ")
	if names == nil || len(names) == 0 {
		return child, false
	}

	current := n
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		if child, ok = current.Child(name); !ok {
			return
		} else {
			current = child
		}
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

// Slice returns []string.
func (n *Node) Slice() (data []string, err error) {
	if n.Type != NodeList {
		err = &InvalidNodeTypeError{n.Type}
		return
	}

	data = make([]string, 0, len(n.List))
	if n.List == nil {
		return
	}

	for _, child := range n.List {
		if child.Type == NodeLiteral {
			data = append(data, child.Literal)
		}
	}
	return data, nil
}

// Map returns map[string]string
func (n *Node) Map() (data map[string]string, err error) {
	if n.Type != NodeHash {
		err = &InvalidNodeTypeError{n.Type}
		return
	}

	data = make(map[string]string, len(n.Hash))
	if n.Hash == nil {
		return
	}

	for name, value := range n.Hash {
		if value.Type == NodeLiteral {
			data[name] = value.Literal
		}
	}
	return
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

	// fmt.Println("setmap", nameOfNodeType(n.Type), v.Type(), v.Kind())
	// fmt.Println(n.Dump())

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
			var mapElem reflect.Value
			if elemType.Kind() == reflect.Ptr {
				mapElem = reflect.New(elemType.Elem())
			} else {
				mapElem = reflect.New(elemType)
			}
			x.set(mapElem)
			v.SetMapIndex(reflect.ValueOf(name), mapElem)
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
	}

}

func (n *Node) set(v reflect.Value) {

	// if !v.CanSet() || !v.IsValid() {
	// 	return
	// }

	// fmt.Println("set==", nameOfNodeType(n.Type), v.Type(), v.Kind())
	// fmt.Println(n.Dump())

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

func (n *Node) dumpto(w *indentWriter) {
	switch n.Type {
	case NodeNone:
		return
	case NodeLiteral:
		//w.WriteIndent()
		w.WriteString(n.Literal)
	case NodeList:
		w.WriteString("[")
		w.WriteString("\n")
		w.Inner()
		for _, child := range n.List {
			w.WriteIndent()
			child.dumpto(w)
			w.WriteString("\n")
		}
		w.Outer()
		w.WriteIndent()
		w.WriteString("]")
	case NodeHash:
		w.WriteString("{")
		w.WriteString("\n")
		w.Inner()

		names := make([]string, 0, len(n.Hash))
		for name, _ := range n.Hash {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			child := n.Hash[name]
			w.WriteIndent()
			w.WriteString(name)
			w.WriteString(":")
			child.dumpto(w)
			w.WriteString("\n")
		}
		w.Outer()
		w.WriteIndent()
		w.WriteString("}")
	}
}

// Dump return dump of node as string
func (n *Node) Dump() string {
	w := &indentWriter{Indent: "\t"}
	n.dumpto(w)
	return string(w.Bytes())
}

type indentWriter struct {
	bytes.Buffer
	Deep   int
	Indent string
}

func (w *indentWriter) Outer() {
	if w.Deep > 0 {
		w.Deep--
	}
}

func (w *indentWriter) Inner() {
	w.Deep++
}

func (w *indentWriter) WriteIndent() {
	for i := 0; i < w.Deep; i++ {
		w.WriteString(w.Indent)
	}
}

// ParseFile parse a file, remove any line start with #(comments), add {\} auto at begin and end
func ParseFile(filename string) (node *Node, err error) {

	var f []byte
	if f, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	buff := bytes.NewBuffer(f)
	var data bytes.Buffer
	data.WriteString("{\n")

	for {
		var line []byte
		end := false

		if line, err = buff.ReadBytes('\n'); err != nil {
			if err != io.EOF {
				return
			} else {
				end = true
			}
		}

		if commentsReg.Match(line) {
			continue
		}
		if _, err = data.Write(line); err != nil {
			return
		}

		if end {
			break
		}
	}

	data.WriteString("\n}\n")
	node, err = Parse(data.Bytes())
	return
}
