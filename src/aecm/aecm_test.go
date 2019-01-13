package main

import (
	"aecm/qnsql"
	"testing"
)

func Test_MysqlCheck(t *testing.T) {
	if !qnsql.DatabaseCheck("root", "") {
		t.Fatal("DatabaseCheck failed.")
	}
}

func Test_OpenDatabase(t *testing.T) {
	db := qnsql.OpenDatabase("root", "")
	if db == nil {
		t.Fatal("OpenDatabase failed.")
	}
	qnsql.CloseDatabase(db)
}
