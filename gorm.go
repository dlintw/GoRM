package gorm

import (
	"os"
	"sdegutis/sqlite"
	"fmt"
	
	// testing
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
		scanMapIntoStruct(rowStruct, results)
	} else {
		return os.NewError("did not find any results")
	}
	
	if s.Next() {
		return os.NewError("more than one row matched")
	}
	
	return nil
}
