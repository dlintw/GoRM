package gorm

import (
	"testing"
)

type Person struct {
	Id int
	Name string
	Age int
}

func TestOpenDB(t *testing.T) {
	db, err := OpenDB("test.db")
	if err != nil || db == nil {
		t.Errorf("opening db file should not have failed")
	}
	err = db.Close()
	if err != nil {
		t.Errorf("closing db connection should not have failed either")
	}
}

func TestGetSingle(t *testing.T) {
	// return // gotta test scanning map into struct first
	
	db, err := OpenDB("test.db")
	defer db.Close()
	
	bob := Person{}
	
	err = db.Get(&bob, "name = ?", "bob")
	
	if err != nil {
		t.Error(err)
	}
	
	if bob.Name != "bob" || bob.Age != 24 || bob.Id != 2 {
		t.Errorf("bob was not filled out properly [%v]", bob)
	}
}
