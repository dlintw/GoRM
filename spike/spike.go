package main

import (
	"log"
	"fmt"
	"gosqlite.googlecode.com/hg/sqlite"
)

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
		fields := make([]interface{}, 3)
		
		for i, _ := range fields {
			var s string
			fields[i] = &s
		}
		
		s.Scan(fields...)
		// fields is now a slice of pointers to strings
		
		for _, v := range fields {
			if v, ok := v.(*string); ok {
				fmt.Printf("%v\n", *v)
			}
		}
	}
}
