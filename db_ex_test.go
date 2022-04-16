package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompoundsDelete(t *testing.T) {
	c := NewCompounds()
	c.AddKV(`name_a`, `value_a`)
	c.AddKV(`name_b`, `value_b`)
	c.AddKV(`name_c`, `value_c`)
	c.AddKV(`name_d`, `value_d`)
	c.AddKV(`name_e`, `value_e`)
	c.AddKV(`name_f`, `value_f`)
	c.Delete(`name_a`, `name_b`, `name_c`)
	toMap := func() map[string]interface{} {
		data := map[string]interface{}{}
		for _, _v := range c.Slice() {
			cd := _v.(Cond)
			for k, v := range cd {
				data[fmt.Sprint(k)] = v
			}
		}
		return data
	}
	assert.Equal(t, map[string]interface{}{
		"name_d": "value_d",
		"name_e": "value_e",
		"name_f": "value_f",
	}, toMap())
	c.Delete(`name_c`, `name_d`, `name_e`, `name_f`)
	assert.Equal(t, map[string]interface{}{}, toMap())
	assert.Equal(t, 0, c.Size())

	c.AddKV(`name_e`, `value_e`)
	assert.Equal(t, 1, c.Size())
	c.AddKV(`name_f`, `value_f`)
	assert.Equal(t, 2, c.Size())
	assert.Equal(t, map[string]interface{}{
		`name_e`: `value_e`,
		"name_f": "value_f",
	}, toMap())
	c.Delete(`name_f`)
	assert.Equal(t, 1, c.Size())
	assert.Equal(t, map[string]interface{}{
		`name_e`: `value_e`,
	}, toMap())
}
