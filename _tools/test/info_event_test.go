package main_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/webx-top/db/_tools/test/dbschema"

	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/factory/pagination"
	"github.com/webx-top/echo"
)

func TestObjectsSize(t *testing.T) {
	m := dbschema.NewNgingVhost(nil)
	m.SetObjects([]*dbschema.NgingVhost{
		{},
		{},
		{},
	})
	size := pagination.ObjectsSize(m)
	assert.Equal(t, 3, len(m.Objects()))
	assert.Equal(t, 3, size)
	size = pagination.ObjectsSize(&[]*dbschema.NgingVhost{
		{},
		{},
	})
	assert.Equal(t, 2, size)
}

func TestTransaction(t *testing.T) {
	var tr echo.Transaction = echo.NewTransaction(factory.NewParam())
	_, ok := tr.(*echo.BaseTransaction).Transaction.(*factory.Param)
	assert.True(t, ok)
}

func TestEvent(t *testing.T) {
	dbi := dbschema.DBI
	buf := new(bytes.Buffer)

	// 定义事件
	dbi.On(`creating`, func(m factory.Model, _ ...string) error {
		buf.WriteString(`creating.`)
		println(`creating.`)
		return nil
	}, `nging_vhost`)
	dbi.On(`created`, func(m factory.Model, _ ...string) error {
		buf.WriteString(`created.`)
		println(`created.`)
		return nil
	}, `nging_vhost`)

	// 调用事件
	m := &dbschema.NgingVhost{}
	dbi.Fire(`creating`, m, nil)
	assert.Equal(
		t,
		`creating.`,
		buf.String(),
	)
	dbi.Fire(`created`, m, nil)
	assert.Equal(
		t,
		`creating.created.`,
		buf.String(),
	)

	b, err := json.MarshalIndent(dbi.Events, ``, `  `)
	if err != nil {
		panic(err)
	}
	assert.Contains(t, string(b), `.TestEvent.func1`)
}
