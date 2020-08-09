package ant

import (
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"github.com/kr/pretty"
	"log"
	"strconv"
	"strings"
)

/////////////// types converters /////////
func pbTypesToGoType(tp string) string {
	s := tp
	switch tp {
	case "int64", "sint64":
		s = "int64"
	case "double":
		s = "float64"
	case "float":
		s = "float32"
	case "int32", "sint32":
		s = "int32"
	case "uint32":
		s = "uint32"
	case "uint64":
		s = "uint64"
	case "fixed32":
		s = "uint32"
	case "fixed64":
		s = "uint64"
	case "sfixed32":
		s = "int32"
	case "sfixed64":
		s = "int64"

	case "bool":
		s = "bool"
	case "string":
		s = "string"
	case "bytes":
		s = "[]byte"
	}
	return s
}

func pbTypesToGoFlatTypes(tp string) string {
	s := tp
	switch tp {
	case "int64", "sint64", "int32",
		"sint32", "uint32", "uint64", "fixed32",
		"fixed64", "sfixed32", "sfixed64":
		s = "int"
	case "double":
		s = "float64"
	case "float":
		s = "float32"

	case "bool":
		s = "bool"
	case "string":
		s = "string"
	case "bytes":
		s = "[]byte"
	}
	return s
}

func pbTypesToJavaType(tp string) string {
	s := tp
	switch tp {
	case "int32", "sint32",
		"uint32", "fixed32",
		"sfixed32":
		s = "int"
	case "int64", "sint64",
		"uint64", "fixed64",
		"sfixed64":
		s = "long"
	case "double":
		s = "double"
	case "float":
		s = "float"

	case "bool":
		s = "boolean"
	case "string":
		s = "String"
	case "bytes":
		s = "byte[]" //or  "ByteString" -  PB default converts to this
	}
	return s
}

//////////////// Helpers ////////////////////
func NoErr(err error) {
	noErr(err)
}

func noErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func StrToInt(str string, defualt int) int {
	r64, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return defualt
	}
	return int(r64)
}

func IntToStr(num int) string {
	return strconv.Itoa(num)
}

func PertyPrint(a interface{}) {
	fmt.Printf("%# v \n", pretty.Formatter(a))
}

/////////////// Hashes /////////////

// Produces hash id for servernames in range of [0..2^31) - so Int in Java and most other langs can handle it (no negative number).
func Hash(str string) int {
	//sh1 := md5.Sum([]byte(str))
	sh1 := sha512.Sum512([]byte(str))
	b := sh1[0]
	// Avoid negative numbers
	b = b >> 1
	bytes := []byte{b, sh1[1], sh1[2], sh1[3]}

	res := binary.BigEndian.Uint32(bytes)
	//fmt.Println(res, int32(res))
	return int(res)
}

func StrToInt32Hash(string string) int32 {
	return int32(Hash(string))
}

func MyHash(string string) int {
	h := 15485862
	a := 7

	for i := 0; i < len(string); i++ {
		h = ((h * a) + int(int8(string[0]))) / 3
	}
	return h
}

func MyHash2(string string) int {
	s := strings.Split(string, ".")
	rs := IntToStr((Hash(s[0])%1000)+1) + "00" + IntToStr(Hash(s[1])%100000)
	//fmt.Println("+++++++"+rs)
	return StrToInt(rs, -1)
}
