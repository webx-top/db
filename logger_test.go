package db

import (
	"testing"
)

type myInt int

func TestBuildSQL(t *testing.T) {
	query := "select * from `a` where a=? and b=?"
	expected := "select * from `a` where a=1 and b='c'"
	actual := BuildSQL(query, 1, `c`)
	if expected != actual {
		t.Fatal(actual, "!=", expected)
	}
	myI := myInt(1)
	actual = BuildSQL(query, myI, `c`)
	if expected != actual {
		t.Fatal(actual, "!=", expected)
	}
	actual = BuildSQL(query, &myI, `c`)
	if expected != actual {
		t.Fatal(actual, "!=", expected)
	}
	expected = "select * from `a` where a='1' and b='c'"
	actual = BuildSQL(query, `1`, `c`)
	if expected != actual {
		t.Fatal(actual, "!=", expected)
	}
}
