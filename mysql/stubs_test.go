//go:build !generated
// +build !generated

package mysql

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/db/lib/sqlbuilder"
)

func mustOpen() sqlbuilder.Database {
	return nil
}

func _mustOpen() sqlbuilder.Database {
	sess, err := Open(settings)
	if err != nil {
		panic(err.Error())
	}
	return sess
}

func TestMain(m *testing.M) {
	m.Run()
	log.Fatal(`Tests use generated code and a custom database, please use "make test".`)
}

func TestJSON(t *testing.T) {
	sess := _mustOpen()
	defer sess.Close()

	type _MyType struct {
		ID                        int64  `db:"id,omitempty"`
		JSONArray                 string `db:"json_array"`
		AutoCustomJSONObjectArray string `db:"auto_custom_json_object_array"`
	}
	var actual _MyType
	err := sess.Collection("my_types").Find(1).One(&actual)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), actual.ID)
	assert.Equal(t, `[1, 2, 3, 4]`, actual.JSONArray)
	assert.Equal(t, `[{"name": "Hello", "value": 0}, {"name": "World", "value": 0}]`, actual.AutoCustomJSONObjectArray)

	var actual2 _MyType
	err = sess.Collection("my_types").Find(mysql.FindInJSON(`json_array`, 3, `[2]`)).One(&actual2)
	assert.NoError(t, err)
	assert.Equal(t, `[1, 2, 3, 4]`, actual2.JSONArray)

	var actual3 _MyType
	err = sess.Collection("my_types").Find(mysql.FindInJSON(`auto_custom_json_object_array`, `World`, `*`, `name`)).One(&actual3)
	assert.NoError(t, err)
	assert.Equal(t, `[1, 2, 3, 4]`, actual3.JSONArray)

	var actual4 _MyType
	err = sess.Collection("my_types").Find(mysql.FindInJSON(`auto_custom_json_object_array`, `World`, `*`, `"name"`)).One(&actual4)
	assert.NoError(t, err)
	assert.Equal(t, `[1, 2, 3, 4]`, actual4.JSONArray)
	err = sess.Collection("my_types").Find(mysql.FindInJSON(`auto_custom_json_object_array`, `World`, `*`, `'name'`)).One(&actual4)
	assert.Error(t, err)
}
