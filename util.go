package gorm

import (
	"reflect"
	"strings"
)

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
	newstr := make([]int, 0)
	firstTime := true
	
	for _, chr := range name {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if !firstTime {
				newstr = append(newstr, '_')
			}
			chr -= ('A' - 'a')
		}
		newstr = append(newstr, chr)
		firstTime = false
	}
	
	return string(newstr)
}
