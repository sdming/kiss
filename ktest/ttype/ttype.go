// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package ttype

type AliasInt int
type AliasUint uint
type AliasBool bool
type AliasFloat float32
type AliasString string

type SimpleType struct {
	A_uint8  uint8
	A_uint16 uint16
	A_uint32 uint32
	A_uint64 uint64
	A_uint   uint

	a_uint8_p  uint8
	a_uint16_p uint16
	a_uint32_p uint32
	a_uint64_p uint64
	a_uint_p   uint

	A_uint8_p  *uint8
	A_uint16_p *uint16
	A_uint32_p *uint32
	A_uint64_p *uint64
	A_uint_p   *uint

	A_int8  int8
	A_int16 int16
	A_int32 int32
	A_int64 int64
	A_int   int

	a_int8_p  int8
	a_int16_p int16
	a_int32_p int32
	a_int64_p int64
	a_int_p   int

	A_int8_p  *int8
	A_int16_p *int16
	A_int32_p *int32
	A_int64_p *int64
	A_int_p   *int

	A_float32   float32
	A_float64   float64
	a_float32_p float32
	a_float64_p float64
	A_float32_p *float32
	A_float64_p *float64

	A_byte   byte
	A_rune   rune
	A_bool   bool
	A_string string

	a_byte_p   byte
	a_rune_p   rune
	a_bool_p   bool
	a_string_p string

	A_byte_p   *byte
	A_rune_p   *rune
	A_bool_p   *bool
	A_string_p *string

	A_alias_bool   AliasBool
	A_alias_int    AliasInt
	A_alias_uint   AliasUint
	A_alias_float  AliasFloat
	A_alias_string AliasString

	a_alias_bool_p   AliasBool
	a_alias_int_p    AliasInt
	a_alias_uint_p   AliasUint
	a_alias_float_p  AliasFloat
	a_alias_string_p AliasString

	A_alias_bool_p   *AliasBool
	A_alias_int_p    *AliasInt
	A_alias_uint_p   *AliasUint
	A_alias_float_p  *AliasFloat
	A_alias_string_p *AliasString
}

func NewSimpleType() SimpleType {
	t := SimpleType{
		A_uint8:    8,
		A_uint16:   16,
		A_uint32:   32,
		A_uint64:   64,
		A_uint:     1024,
		a_uint8_p:  8,
		a_uint16_p: 16,
		a_uint32_p: 32,
		a_uint64_p: 64,
		a_uint_p:   1024,

		A_int8:    -8,
		A_int16:   -16,
		A_int32:   -32,
		A_int64:   -64,
		A_int:     -1024,
		a_int8_p:  -8,
		a_int16_p: -16,
		a_int32_p: -32,
		a_int64_p: -64,
		a_int_p:   -1024,

		A_float32:   3.2,
		A_float64:   6.4,
		a_float32_p: 3.2,
		a_float64_p: 6.4,

		A_byte:     127,
		A_rune:     0x767d,
		A_bool:     true,
		A_string:   "string",
		a_byte_p:   127,
		a_rune_p:   0x767d,
		a_bool_p:   true,
		a_string_p: "string",

		A_alias_bool:   true,
		A_alias_int:    -1024,
		A_alias_uint:   1024,
		A_alias_float:  6.4,
		A_alias_string: "string",

		a_alias_bool_p:   true,
		a_alias_int_p:    -1024,
		a_alias_uint_p:   1024,
		a_alias_float_p:  6.4,
		a_alias_string_p: "string",
	}

	t.A_uint8_p = &t.a_uint8_p
	t.A_uint16_p = &t.a_uint16_p
	t.A_uint32_p = &t.a_uint32_p
	t.A_uint64_p = &t.a_uint64_p
	t.A_uint_p = &t.a_uint_p

	t.A_int8_p = &t.a_int8_p
	t.A_int16_p = &t.a_int16_p
	t.A_int32_p = &t.a_int32_p
	t.A_int64_p = &t.a_int64_p
	t.A_int_p = &t.a_int_p

	t.A_float32_p = &t.a_float32_p
	t.A_float64_p = &t.a_float64_p

	t.A_byte_p = &t.a_byte_p
	t.A_rune_p = &t.a_rune_p
	t.A_bool_p = &t.a_bool_p
	t.A_string_p = &t.a_string_p

	t.A_alias_bool_p = &t.a_alias_bool_p
	t.A_alias_int_p = &t.a_alias_int_p
	t.A_alias_uint_p = &t.a_alias_uint_p
	t.A_alias_float_p = &t.a_alias_float_p
	t.A_alias_string_p = &t.a_alias_string_p

	return t
}
