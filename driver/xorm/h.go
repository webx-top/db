package xorm

import (
	"strings"

	. "github.com/webx-top/dbx/driver"
)

type XH struct {
	*H
}

func (h *XH) Build() interface{} {
	c := &Condition{
		SQL:  ``,
		Args: []interface{}{},
	}
	t := ``
	for _, key := range h.Keys() {
		val := h.Get(key)
		if key[0] == '$' {
			key = key[1:]
			switch key {
			case "and":
				if vs, ok := val.([]*XH); ok {
					t2 := ``
					c.SQL += t + `(`
					for _, v := range vs {
						b := v.Build().(*Condition)
						c.SQL += t2 + b.SQL
						c.Args = append(c.Args, b.Args...)
						t2 = ` AND `
					}
					c.SQL += `)`
				}
			case "or":
				if vs, ok := val.([]*XH); ok {
					t2 := ``
					c.SQL += t + `(`
					for _, v := range vs {
						b := v.Build().(*Condition)
						c.SQL += t2 + b.SQL
						c.Args = append(c.Args, b.Args...)
						t2 = ` OR `
					}
					c.SQL += `)`
				}
			}
			continue
		}
		if vs, ok := val.(map[string][]interface{}); ok {
			isIN := false
			for k, v := range vs {
				isIN = k == `$in` || k == `$nin`
				if isIN {
					op := `in`
					if k == `$nin` {
						op = `NOT IN`
					}
					c.SQL += t + key + ` ` + op + ` (` + strings.Repeat(`?`, len(v)) + `)`
					c.Args = append(c.Args, v...)
					t = ` AND `
				}
				val = v
				break
			}
			if isIN {
				continue
			}
		}
		op := `=`
		if vs, ok := val.(map[string]string); ok {
			for k, v := range vs {
				op = operatorFlag(k)
				val = v
				break
			}
		} else if vs, ok := val.([]string); ok {
			switch len(vs) {
			case 2:
				op = operatorFlag(vs[0])
				val = vs[1]
			case 1:
				val = vs[0]
			}
		}
		c.SQL += t + key + op + `?`
		c.Args = append(c.Args, val)
		t = ` AND `
	}
	return c
}

func operatorFlag(ident string) string {
	op := `=`
	switch ident {
	case "$gt":
		op = `>`
	case "$lt":
		op = `<`
	case "$gte":
		op = `>=`
	case "$lte":
		op = `<=`
	case "$ne":
		op = `!=`
	}
	return op
}
