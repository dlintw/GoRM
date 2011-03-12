package gorm

import (
	"testing"
	"reflect"
	"os"
)

func TestTypeName(t *testing.T) {
	type PersonStruct struct {
		age int
		name string
	}
	
	var pete PersonStruct
	tname := getTypeName(pete)
	if tname != "PersonStruct" {
		t.Errorf("Expected type %T to be PersonStruct, got %v\n", pete, tname)
	}
}

func TestSnakeCasing(t *testing.T) {
	names := map[string]string {
		"ThisThat": "this_that",
		"WhatIAm": "what_i_am",
		"IAmNot": "i_am_not",
		"Shop": "shop",
	}
	for key, val := range names {
		if name := snakeCasedName(key); name != val {
			t.Errorf("Expected [%v] to translate to [%v], got [%v]\n", key, val, name)
		}
	}
}

func TestTitleCasing(t *testing.T) {
	names := map[string]string {
		"this_that": "ThisThat",
		"what_i_am": "WhatIAm",
		"i_am_not": "IAmNot",
		"shop": "Shop",
	}
	for key, val := range names {
		if name := titleCasedName(key); name != val {
			t.Errorf("Expected [%v] to translate to [%v], got [%v]\n", key, val, name)
		}
	}
}

func TestPluralizeString(t *testing.T) {
	names := map[string]string {
		"person": "persons",
		"yak": "yaks",
		"ghost": "ghosts",
		"party": "parties",
	}
	for key, val := range names {
		if name := pluralizeString(key); name != val {
			t.Errorf("Expected [%v] to translate to [%v], got [%v]\n", key, val, name)
		}
	}
}


func TestEscapeString(t *testing.T) {
	nameFuncs := map[func()(string, os.Error)]string {
		func()(string, os.Error) { return escapeString("where name = ?", "jack") }: "where name = 'jack'",
		func()(string, os.Error) { return escapeString("where age = ?", 42) }: "where age = 42",
		func()(string, os.Error) { return escapeString("where name = ? and age = ?", "jack", 42) }: "where name = 'jack' and age = 42",
	}
	
	for key, val := range nameFuncs {
		if str, err := key(); str != val {
			t.Errorf("Expected [%v] to translate to [%v], got [%v] with error [%v]\n", key, val, str, err)
		}
	}
	str, err := escapeString("where age = ?", 42, "jack")
	if str != "" || err == nil {
		t.Errorf("Expected incorrect argument count error, didn't get it.")
	}
	str, err = escapeString("where name = ? and age = ?", 42)
	if str != "" || err == nil {
		t.Errorf("Expected incorrect argument count error, didn't get it.")
	}
}

func TestScanStructIntoMap(t *testing.T) {
	var pete Person
	pete.Name = "bob"
	pete.Age = 32
	pete.Id = 7
	
	peteMap, err := scanStructIntoMap(reflect.NewValue(&pete))
	if err != nil {
		t.Error(err)
	}
	
	peteComparableMap := map[string]interface{} {
		"name": "bob",
		"age": 32,
		"id": 7,
	}
	
	if !reflect.DeepEqual(peteMap, peteComparableMap) {
		t.Errorf("pete's map was not filled out properly. have %v, want %v", peteMap, peteComparableMap)
	}
}

func TestScanMapIntoStruct(t *testing.T) {
	personMap := map[string][]byte {
		"name": []byte("bob"),
		"id": []byte("2"),
		"age": []byte("42"),
	}
	
	bob := Person{}
	
	err := scanMapIntoStruct(reflect.NewValue(&bob), personMap)
	if err != nil {
		t.Error(err)
	}
	if bob.Name != "bob" || bob.Age != 42 || bob.Id != 2 {
		t.Errorf("struct was not filled out right; got %v with error %v", bob, err)
	}
}
