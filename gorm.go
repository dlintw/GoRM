package main

import (
	"fmt"
	"reflect"
	"strings"
	"gosqlite.googlecode.com/hg/sqlite"
)

type PersonStruct struct {
	age int
	name string
}

func getTypeName(obj interface{}) (typestr string, isPtr bool) {
	typ := reflect.Typeof(obj)
	typestr = typ.String()
	
	isPtr = strings.HasPrefix(typestr, "*")
	
	lastDotIndex := strings.LastIndex(typestr, ".")
	if lastDotIndex != -1 {
		typestr = typestr[lastDotIndex+1:]
	}
	
	return
}

func snakeCasedName(name string) string {
	return ""
}

func main() {
	fmt.Println(sqlite.Version())
	return
	
	steven1 := new(PersonStruct)
	var steven2 PersonStruct
	fmt.Println(getTypeName(steven1))
	fmt.Println(getTypeName(steven2))
}
