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
			if firstTime == true {
				firstTime = false
			} else {
				newstr = append(newstr, '_')
			}
			chr -= ('A' - 'a')
		}
		newstr = append(newstr, chr)
	}
	
	return string(newstr)
}

func titleCasedName(name string) string {
	newstr := make([]int, 0)
	upNextChar := true
	
	for _, chr := range name {
		switch {
		case upNextChar:
			upNextChar = false
			chr -= ('a' - 'A')
		case chr == '_':
			upNextChar = true
			continue
		}
		
		newstr = append(newstr, chr)
	}
	
	return string(newstr)
}
