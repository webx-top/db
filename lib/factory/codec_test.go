package factory

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo/param"
)

type testMapp struct {
	Id         uint32
	Name       string
	Address    string
	MaxRetries uint
}

func (t *testMapp) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r[`Id`] = t.Id
		r[`Name`] = t.Name
		r[`Address`] = t.Address
		r[`MaxRetries`] = t.MaxRetries
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case `Id`:
			r[`Id`] = t.Id
		case `Name`:
			r[`Name`] = t.Name
		case `Address`:
			r[`Address`] = t.Address
		case `MaxRetries`:
			r[`MaxRetries`] = t.MaxRetries
		}
	}
	return r
}

func (t *testMapp) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r[`id`] = t.Id
		r[`name`] = t.Name
		r[`address`] = t.Address
		r[`max_retries`] = t.MaxRetries
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case `id`:
			r[`id`] = t.Id
		case `name`:
			r[`name`] = t.Name
		case `address`:
			r[`address`] = t.Address
		case `max_retries`:
			r[`max_retries`] = t.MaxRetries
		}
	}
	return r
}

func (t *testMapp) FromRow(row map[string]interface{}) {
	for field, v := range row {
		switch field {
		case `id`:
			t.Id = param.AsUint32(v)
		case `name`:
			t.Name = param.AsString(v)
		case `address`:
			t.Address = param.AsString(v)
		case `max_retries`:
			t.MaxRetries = param.AsUint(v)
		}
	}
}

func (t *testMapp) Set(key interface{}, value ...interface{}) {
	switch k := key.(type) {
	case map[string]interface{}:
		for kk, vv := range k {
			t.Set(kk, vv)
		}
	default:
		var (
			kk string
			vv interface{}
		)
		if k, y := key.(string); y {
			kk = k
		} else {
			kk = fmt.Sprint(key)
		}
		if len(value) > 0 {
			vv = value[0]
		}
		switch kk {
		case `Id`:
			t.Id = param.AsUint32(vv)
		case `Name`:
			t.Name = param.AsString(vv)
		case `Address`:
			t.Address = param.AsString(vv)
		case `MaxRetries`:
			t.MaxRetries = param.AsUint(vv)
		}
	}
}

func jsonEncode(data interface{}) string {
	b, _ := json.Marshal(data)
	return string(b)
}

func xmlEncode(data interface{}) string {
	b, _ := xml.Marshal(data)
	return string(b)
}

func TestCodecEncode(t *testing.T) {
	m := &testMapp{
		Id:         10000,
		Name:       `test`,
		Address:    `testAddress`,
		MaxRetries: 3,
	}
	c := NewCodec(m, KeystyleCamelCase)
	expected := map[string]interface{}{
		`id`:         uint32(10000),
		`name`:       `test`,
		`address`:    `testAddress`,
		`maxRetries`: uint(3),
	}
	actual := c.MakeMap()
	assert.Equal(t, expected, actual)
	assert.Equal(t, "{\"address\":\"testAddress\",\"id\":10000,\"maxRetries\":3,\"name\":\"test\"}", jsonEncode(c))
	c.SetColumns(`id`, `name`)
	assert.Equal(t, "{\"id\":10000,\"name\":\"test\"}", jsonEncode(c))
	assert.Equal(t, "<item><id>10000</id><name>test</name></item>", xmlEncode(c))

	c.SetColumns()
	c.SetKeystyle(KeystylePascalCase)
	expected = map[string]interface{}{
		`Id`:         uint32(10000),
		`Name`:       `test`,
		`Address`:    `testAddress`,
		`MaxRetries`: uint(3),
	}
	actual = c.MakeMap()
	assert.Equal(t, expected, actual)
	assert.Equal(t, "<Item><Address>testAddress</Address><Id>10000</Id><MaxRetries>3</MaxRetries><Name>test</Name></Item>", xmlEncode(c))

	c.SetKeystyle(KeystyleSnakeCase)
	expected = map[string]interface{}{
		`id`:          uint32(10000),
		`name`:        `test`,
		`address`:     `testAddress`,
		`max_retries`: uint(3),
	}
	actual = c.MakeMap()
	assert.Equal(t, expected, actual)
	assert.Equal(t, "{\"address\":\"testAddress\",\"id\":10000,\"max_retries\":3,\"name\":\"test\"}", jsonEncode(c))
	c.SetColumns(`id`, `name`, `max_retries`)
	assert.Equal(t, "{\"id\":10000,\"max_retries\":3,\"name\":\"test\"}", jsonEncode(c))
	assert.Equal(t, "<item><id>10000</id><max_retries>3</max_retries><name>test</name></item>", xmlEncode(c))

}

func TestCodecDecode(t *testing.T) {
	m := &testMapp{}
	c := NewCodec(m, KeystyleSnakeCase)
	err := xml.Unmarshal([]byte(`<item><id>10000</id><name>test</name><max_retries>3</max_retries></item>`), c)
	assert.NoError(t, err)
	assert.Equal(t, &testMapp{
		Id:   10000,
		Name: `test`,
		//Address:    `testAddress`,
		MaxRetries: 3,
	}, m)
}
