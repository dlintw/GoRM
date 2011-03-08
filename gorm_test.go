package gorm

import (
	"testing"
)

func TestTypeName(t *testing.T) {
	type PersonStruct struct {
		age int
		name string
	}
	
	var pete PersonStruct
	tname, isptr := getTypeName(pete)
	if tname != "PersonStruct" {
		t.Errorf("Expected type %T to be PersonStruct, got %v\n", pete, tname)
	}
	if isptr != false {
		t.Errorf("Didn't expect %v to be a pointer\n", pete)
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
