package main

import (
	"log"
	"reflect"
	"fmt"
	"sdegutis/sqlite"
)

type Shirt struct {
	size int
}

type Person struct {
	Id int
	Name string
	Age int
	// Shirts []byte
}

func stuff(obj interface{}) {
	if obj, ok := reflect.Typeof(obj).(*reflect.StructType); ok {
		// fldname := "Name"
		fmt.Println(reflect.Typeof(obj).String())
		fmt.Println(obj.Size())
	}
	fmt.Println(reflect.Typeof(obj).String())
}

func main() {
	conn, err := sqlite.Open("test.db")
	defer conn.Close()
	
	if err != nil {
		log.Fatal(err)
	}
	
    s, err := conn.Prepare("select * from persons")
    if err != nil {
            log.Fatal(err)
    }
    defer s.Finalize()
    err = s.Exec()
    if err != nil {
        log.Fatal(err)
    }
	
	for s.Next() {
		// pete := Person{}
		
		results, err := s.ResultsAsMap()
		
		// fmt.Println(s.ScanStruct(&pete))
		// fmt.Println(pete)
		
		fmt.Println(results, err)
	}
}

// func (s *Stmt) ScanStruct(obj interface{}) os.Error {
// 	objPtr, ok := reflect.NewValue(obj).(*reflect.PtrValue)
// 	if !ok {
// 		return os.NewError(fmt.Sprintf("%v", ok))
// 	}
// 	dataStruct, ok := objPtr.Elem().(*reflect.StructValue)
// 	if !ok {
// 		return os.NewError(fmt.Sprintf("%v", ok))
// 	}
// 	
// 	n := int(C.sqlite3_column_count(s.stmt))
// 	
// 	for i := 0; i < n; i++ {
// 		n := C.sqlite3_column_bytes(s.stmt, C.int(i))
// 		p := C.sqlite3_column_blob(s.stmt, C.int(i))
// 		colname := C.GoString(C.sqlite3_column_name(s.stmt, C.int(i)))
// 		
// 		if p == nil && n > 0 {
// 			return os.NewError("got nil blob")
// 		}
// 		
// 		var v interface{}
// 		
// 		structField := dataStruct.FieldByName(ToTitle(colname))
// 		
// 		var data []byte
// 		if n > 0 {
// 			data = (*[1<<30]byte)(unsafe.Pointer(p))[0:n]
// 		}
// 		switch structField.Type().(type) {
// 		case *reflect.SliceType:
// 			v = data
// 		case *reflect.StringType:
// 			v = string(data)
// 		case *reflect.BoolType:
// 			v = string(data) == "1"
// 		case *reflect.IntType:
// 			x, err := strconv.Atoi(string(data))
// 			if err != nil {
// 				return os.NewError("arg " + strconv.Itoa(i) + " as int: " + err.String())
// 			}
// 			v = x
// 		case *reflect.FloatType:
// 			x, err := strconv.Atof64(string(data))
// 			if err != nil {
// 				return os.NewError("arg " + strconv.Itoa(i) + " as float64: " + err.String())
// 			}
// 			v = x
// 		default:
// 			return os.NewError("unsupported type in Scan: " + reflect.Typeof(v).String())
// 		}
// 		
// 		if structField.CanSet() {
// 			nv := reflect.NewValue(v)
// 			structField.SetValue(nv)
// 		}
// 	}
// 	return nil
// }
