package gorm

import (
	"testing"
)

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
