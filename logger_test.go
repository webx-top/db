package db

import (
	"testing"
)

func TestBuildSQL(t *testing.T) {
	query := "select * from `a` where a=? and b=?"
	expected := "select * from `a` where a='1' and b='c'"
	actual := BuildSQL(query, 1, `c`)
	if expected != actual {
		t.Fatal(actual,"!=",expected)
	}
}