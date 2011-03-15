package gorm

import (
	"reflect"
	"strings"
	"strconv"
	"os"
)

func getTypeName(obj interface{}) (typestr string) {
	typ := reflect.Typeof(obj)
	typestr = typ.String()

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
		str = str[:len(str)-1] + "ie"
	}
	return str + "s"
}

func scanMapIntoStruct(obj reflect.Value, objMap map[string][]byte) os.Error {
	dataStruct, ok := reflect.Indirect(obj).(*reflect.StructValue)
	if !ok {
		return os.NewError("expected a pointer to a struct")
	}

	for key, data := range objMap {
		structField := dataStruct.FieldByName(titleCasedName(key))
		if !structField.CanSet() {
			continue
		}

		var v interface{}

		switch structField.Type().(type) {
		case *reflect.SliceType:
			v = data
		case *reflect.StringType:
			v = string(data)
		case *reflect.BoolType:
			v = string(data) == "1"
		case *reflect.IntType:
			x, err := strconv.Atoi(string(data))
			if err != nil {
				return os.NewError("arg " + key + " as int: " + err.String())
			}
			v = x
		case *reflect.FloatType:
			x, err := strconv.Atof64(string(data))
			if err != nil {
				return os.NewError("arg " + key + " as float64: " + err.String())
			}
			v = x
		default:
			return os.NewError("unsupported type in Scan: " + reflect.Typeof(v).String())
		}

		structField.SetValue(reflect.NewValue(v))
	}

	return nil
}

func scanStructIntoMap(obj reflect.Value) (map[string]interface{}, os.Error) {
	dataStruct, ok := reflect.Indirect(obj).(*reflect.StructValue)
	if !ok {
		return nil, os.NewError("expected a pointer to a struct")
	}

	dataStructType := dataStruct.Type().(*reflect.StructType)

	mapped := make(map[string]interface{})

	for i := 0; i < dataStructType.NumField(); i++ {
		field := dataStructType.Field(i)
		fieldName := field.Name

		mapKey := snakeCasedName(fieldName)
		value := dataStruct.FieldByName(fieldName).Interface()

		mapped[mapKey] = value
	}

	return mapped, nil
}
