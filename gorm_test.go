package gorm

import (
	"testing"
	"io/ioutil"
	"bytes"
	"reflect"
)

type Person struct {
	Id   int
	Name string
	Age  int
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
	db, _ := OpenDB("test.db")
	defer db.Close()

	var bob Person
	err := db.Get(&bob, "name = ?", "bob")

	if err != nil {
		t.Error(err)
	}

	if bob.Name != "bob" || bob.Age != 24 || bob.Id != 2 {
		t.Errorf("bob was not filled out properly [%v]", bob)
	}
}

func TestGetSingleById(t *testing.T) {
	db, _ := OpenDB("test.db")
	defer db.Close()

	var bob Person
	err := db.Get(&bob, 2)

	if err != nil {
		t.Error(err)
	}

	if bob.Name != "bob" || bob.Age != 24 || bob.Id != 2 {
		t.Errorf("bob was not filled out properly [%v]", bob)
	}
}

func copyTemp(t *testing.T, path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("could not read supposedly 'copyable' file")
	}

	f, err := ioutil.TempFile("", "gorm-sqlite-prefix")
	if err != nil {
		t.Errorf("could not create tempfile for writing")
	}
	f.Write(data)
	fname := f.Name()
	f.Close()

	return fname
}

func TestCopyTemp(t *testing.T) {
	name := "test.db"
	tmpName := copyTemp(t, name)
	if name == tmpName {
		t.Errorf("copyTemp should have given a filename other than %q", name)
	}
	data1, err1 := ioutil.ReadFile(name)
	data2, err2 := ioutil.ReadFile(tmpName)
	if err1 != nil || err2 != nil || bytes.Compare(data1, data2) != 0 {
		t.Errorf("copyTemp did not copy the file correctly.")
	}
}

func TestInsertViaSave(t *testing.T) {
	db, _ := OpenDB(copyTemp(t, "test.db"))
	defer db.Close()

	var james Person
	james.Name = "james"
	james.Age = 29

	err := db.Save(&james)
	if err != nil {
		t.Error(err)
	}

	if james.Id == 0 {
		t.Error("james should have a new Id, but doesnt. sigh.")
	}

	var people []Person
	db.GetAll(&people, "")

	if len(people) != 3 {
		t.Errorf("the db should contain 3 people now, but it contains %d", len(people))
	}
}

func TestSave(t *testing.T) {
	db, _ := OpenDB(copyTemp(t, "test.db"))
	defer db.Close()

	newName := "Fred Jones"

	var bob Person
	db.Get(&bob, 2)

	bob.Name = newName
	db.Save(&bob)

	var fred Person
	err := db.Get(&fred, 2)

	if err != nil {
		t.Error(err)
	}

	if fred.Name != newName {
		t.Errorf("name should have been %q, got %q instead", newName, fred.Name)
	}
}

func TestGetMultiple(t *testing.T) {
	db, _ := OpenDB("test.db")
	defer db.Close()

	var peoples []Person
	err := db.GetAll(&peoples, "")

	if err != nil {
		t.Error(err)
	}

	if len(peoples) != 2 {
		t.Errorf("wrong number of people returned, should be 2, but got %d", len(peoples))
	}

	comparablePeoples := []Person{
		Person{Name: "john", Id: 1, Age: 42},
		Person{Name: "bob", Id: 2, Age: 24},
	}

	if !reflect.DeepEqual(peoples, comparablePeoples) {
		t.Errorf("peoples was not filled out properly %v", peoples)
	}
}

func TestGetMultipleWithoutCondition(t *testing.T) {
	db, _ := OpenDB("test.db")
	defer db.Close()

	var people1 []Person
	db.GetAll(&people1, "id > 0")

	var people2 []Person
	err := db.GetAll(&people2, "")
	if err != nil {
		t.Error(err)
	}

	if len(people1) != len(people2) {
		t.Errorf("wrong number of people returned, should be %d, but got %d", len(people1), len(people2))
	}

	if !reflect.DeepEqual(people1, people2) {
		t.Errorf("people2 was not filled out properly %v", people2)
	}
}
