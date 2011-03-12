package gorm

import (
	"os"
	"sdegutis/sqlite"
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
