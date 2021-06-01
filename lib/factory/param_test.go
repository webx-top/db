package factory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
)

// [SWH|+]
func TestMap(t *testing.T) {
	keys, values, e := sqlbuilder.Map([]interface{}{`user`, `admin`, `id`, 1, `email`, `swh@admpub.com`}, nil)
	assert.Equal(
		t,
		[]string{`user`, `id`, `email`},
		keys,
	)
	assert.Equal(
		t,
		[]interface{}{`admin`, 1, `swh@admpub.com`},
		values,
	)
	assert.Equal(
		t,
		nil,
		e,
	)
	keys, values, e = sqlbuilder.Map([]interface{}{`user`, `admin`, `id`, 1, `email`}, nil)
	assert.Equal(
		t,
		[]string{`user`, `id`, `email`},
		keys,
	)
	assert.Equal(
		t,
		[]interface{}{`admin`, 1, nil},
		values,
	)
	assert.Equal(
		t,
		nil,
		e,
	)
}

type Example struct {
	ReturnTo string `db:"return_to"`
	User     string `db:"user"`
	ID       int    `db:"id"`
}

func TestStructField(t *testing.T) {
	a := &Example{
		ReturnTo: `http://www.webx-top/`,
		User:     `admin`,
		ID:       1,
	}
	p := NewParam()
	p.UsingStructField(a, `ID`, `User`)
	assert.Equal(
		t,
		[]interface{}{`id`, 1, `user`, `admin`},
		p.SaveData.(*db.KeysValues).Slice(),
	)
}
