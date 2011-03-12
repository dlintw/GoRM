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
	db, err := OpenDB("test.db")
	defer db.Close()
	
	var bob Person
	err = db.Get(&bob, "name = ?", "bob")
	
	if err != nil {
		t.Error(err)
	}
	
	if bob.Name != "bob" || bob.Age != 24 || bob.Id != 2 {
		t.Errorf("bob was not filled out properly [%v]", bob)
	}
}

func TestGetSingleById(t *testing.T) {
	db, err := OpenDB("test.db")
	defer db.Close()
	
	var bob Person
	err = db.Get(&bob, 2)
	
	if err != nil {
		t.Error(err)
	}

	if bob.Name != "bob" || bob.Age != 24 || bob.Id != 2 {
		t.Errorf("bob was not filled out properly [%v]", bob)
	}
}


func TestGetMultiple(t *testing.T) {
	db, err := OpenDB("test.db")
	defer db.Close()
	
	var peoples []Person
	err = db.GetAll(&peoples, "id > 0")
	
	if err != nil {
		t.Error(err)
	}
	
	if len(peoples) != 2 {
		t.Errorf("wrong number of people returned, should be 2, but got %d", len(peoples))
	}
	
	hasBob := false
	hasJohn := false
	
	for _, guy := range peoples {
		if guy.Name == "john" && guy.Id == 1 && guy.Age == 42 {
			hasJohn = true
		}
		if guy.Name == "bob" && guy.Id == 2 && guy.Age == 24 {
			hasBob = true
		}
	}
	
	if !hasBob || !hasJohn {
		t.Errorf("peoples was not filled out properly %v", peoples)
	}
}
