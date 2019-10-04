package factory_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/db/_tools/generator/dbschema"
	"github.com/webx-top/db/lib/factory"
)

func TestEvent(t *testing.T) {
	dbi := dbschema.DBI
	buf := new(bytes.Buffer)

	// 定义事件
	dbi.On(`creating`, func(m factory.Model) error {
		buf.WriteString(`creating.`)
		println(`creating.`)
		return nil
	}, `config`)
	dbi.On(`created`, func(m factory.Model) error {
		buf.WriteString(`created.`)
		println(`created.`)
		return nil
	}, `config`)

	// 调用事件
	m := &dbschema.Config{}
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
}
