package gorm

import (
	"os"
	"sdegutis/sqlite"
	"fmt"
	"reflect"
	"log"
)

type Conn struct {
	conn *sqlite.Conn
}

func (c *Conn) Close() os.Error {
	return c.conn.Close()
}

func OpenDB(filename string) (*Conn, os.Error) {
	conn, err := sqlite.Open("test.db")
	return &Conn{conn: conn}, err
}

func (c *Conn) Get(rowStruct interface{}, condition string, args ...interface{}) os.Error {
	tname, _ := getTypeName(rowStruct)
	tname = snakeCasedName(tname)
	tableName := pluralizeString(tname)
	
	condition, err := escapeString(condition, args...)
	if err != nil {
		return err
	}
	
	condition = fmt.Sprintf("where %v", condition)
	
    s, err := c.conn.Prepare(fmt.Sprintf("select * from %v %v", tableName, condition))
    if err != nil {
            log.Fatal(err)
    }
    defer s.Finalize()
    err = s.Exec()
    if err != nil {
        log.Fatal(err)
    }
	
	if s.Next() {
		results, err := s.ResultsAsMap()
		if err != nil {
			return err
		}
		scanMapIntoStruct(reflect.NewValue(rowStruct), results)
	} else {
		return os.NewError("did not find any results")
	}
	
	if s.Next() {
		return os.NewError("more than one row matched")
	}
	
	return nil
}

func (c *Conn) GetAll(rowsSlicePtr interface{}, condition string, args ...interface{}) os.Error {
	rowsPtrValue, _ := reflect.NewValue(rowsSlicePtr).(*reflect.PtrValue)
	rowsPtrType, ok := reflect.Typeof(rowsSlicePtr).(*reflect.PtrType)
	if !ok {
		return os.NewError("needs a *pointer* to a slice")
	}
	
	sliceValue, _ := rowsPtrValue.Elem().(*reflect.SliceValue)
	sliceType, ok := rowsPtrType.Elem().(*reflect.SliceType)
	if !ok {
		log.Fatalf("%p", sliceType)
		return os.NewError("needs a pointer to a *slice*")
	}
	
	sliceElementType := sliceType.Elem()
	
	tname, _ := getTypeName(rowsSlicePtr)
	tname = snakeCasedName(tname)
	tableName := pluralizeString(tname)
	
	condition, err := escapeString(condition, args...)
	if err != nil {
		return err
	}
	
	condition = fmt.Sprintf("where %v", condition)
	
    s, err := c.conn.Prepare(fmt.Sprintf("select * from %v %v", tableName, condition))
    if err != nil {
		log.Fatal(err)
    }
    defer s.Finalize()
    err = s.Exec()
    if err != nil {
        log.Fatal(err)
    }

	for s.Next() {
		newValue := reflect.MakeZero(sliceElementType)
		
		results, err := s.ResultsAsMap()
		if err != nil {
			return err
		}
		
		scanMapIntoStruct(newValue.Addr(), results)
		
		sliceValue.SetValue(reflect.Append(sliceValue, newValue))
	}
	
	return nil
}
