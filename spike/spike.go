package main

import (
	"log"
	"gosqlite.googlecode.com/hg/sqlite"
)

func main() {
	conn, err := sqlite.Open("test.go")
	
	if err != nil {
		log.Fatal(err)
	}
	
	println(conn)
	
	conn.Close()
}
