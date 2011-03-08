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