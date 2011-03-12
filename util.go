package gorm

import (
	"reflect"
	"strings"
	"strconv"
	"os"
	"fmt"
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

func pluralizeString(str string) string {
	if strings.HasSuffix(str, "y") {
		str = str[:len(str) - 1] + "ie"
	}
	return str + "s"
}

func escapeString(str string, args ...interface{}) (result string, err os.Error) {
	if qmarks := strings.Count(str, "?"); qmarks != len(args) {
		return "", os.NewError(fmt.Sprintf("Incorrect number of arguments: have %d want %d", len(args), qmarks))
	}
	
	for i := 0; i < len(args); i++ {
		arg := args[i]
		argstr := ""
		switch arg := arg.(type) {
		case string:
			argstr = "'" + string(arg) + "'"
		case int:
			argstr = strconv.Itoa(arg)
		}
		str = strings.Replace(str, "?", argstr, 1)
	}
	
	return str, nil
}
