package main

import (
	"fmt"
	"github.com/sdming/kiss"
	"github.com/sdming/kiss/gotype"
	"reflect"
)

func main() {
	compare()
	convert()
	parse()
	extend()
}

type User struct {
	Name   string
	Age    uint8
	Single bool
}

func extend() {
	fmt.Println("extend struct")

	var u1 User
	var u2 User

	u2.Name = "tom"
	u2.Age = 11
	u2.Single = true

	v1 := reflect.ValueOf(&u1).Elem()
	v2 := reflect.ValueOf(&u2).Elem()

	var src kiss.GetFunc = func(name string) (interface{}, bool) {
		return v2.FieldByName(name).Interface(), true
	}
	kiss.ExtdStruct(v1, src)

	fmt.Println("u1", u1)

}

func parse() {
	fmt.Println("parse struct")

	var u1 User
	var u2 = map[string]string{
		"Name":   "tom",
		"Age":    "11",
		"Single": "true"}

	v1 := reflect.ValueOf(&u1).Elem()
	var src kiss.StrGetFunc = func(name string) (string, bool) {
		x, ok := u2[name]
		return x, ok
	}
	kiss.ParseStruct(v1, src)

	fmt.Println(u1)
}

func convert() {
	fmt.Println("convert...")

	var data map[string]interface{} = map[string]interface{}{
		"true":   true,
		"-1":     int(-1),
		"-8":     int8(-8),
		"-16":    int16(-16),
		"-32":    int32(-32),
		"-64":    int64(-64),
		"1":      uint(1),
		"8":      uint8(8),
		"16":     uint16(16),
		"32":     uint32(32),
		"64":     int64(64),
		"3.2":    float32(3.2),
		"-6.4":   float64(-6.4),
		"string": "string"}

	for s, t := range data {
		typ := reflect.ValueOf(t).Type()
		v, err := gotype.Atov(s, typ)
		if err == nil {
			fmt.Printf("convert string %s to type %v %v \n", s, typ, v.Interface() == t)
		} else {
			fmt.Printf("convert string %s to type %v error %v \n", s, typ, err)
		}
	}

}

func compare() {
	fmt.Println("compare...")

	var v1 []interface{} = []interface{}{
		int(16),
		int(10)}

	var v2 []interface{} = []interface{}{
		uint(16),
		float32(10)}

	for i, v := range v1 {

		fmt.Printf("%d\t (%v)%v = (%v)%v ? \t%v \n",
			i, reflect.ValueOf(v).Kind(), v, reflect.ValueOf(v2[i]).Kind(), v2[i], gotype.Equal(v, v2[i]))
	}

}
