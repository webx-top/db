package sqlbuilder

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/db"
)

type A struct {
	ID    uint   `db:"id"`
	Group string `db:"group"`
}

type C struct {
	AID  uint   `db:"aid"`
	Name string `db:"name"`
	A    *A     `db:"-,relation=ID:AID|gtZero,dbconn=link2,columns=colName&colName2,where=group:testc"`
}

type B struct {
	AID  uint   `db:"aid"`
	Name string `db:"name"`
	A    *A     `db:"-,relation=ID:AID|gtZero,dbconn=link3,columns=colName&colName2,where=group:testb"`
}

func TestRelationCache(t *testing.T) {
	c := &C{
		AID:  10,
		Name: `10`,
	}
	rValue := reflect.ValueOf(c)
	rType := reflect.TypeOf(c)
	typeMap := mapper.TypeMap(rType)
	options, ok := typeMap.Options[`relation`]
	if !ok {
		panic(`not ok`)
	}
	// for key, val := range typeMap.Options {
	// 	fmt.Printf("%v ==> %+v\n", key, val[0])
	// }
	assert.Equal(t, 4, len(typeMap.Options))
	for _, fieldInfo := range options {
		r, err := parseRelationData(fieldInfo)
		assert.NoError(t, err)
		fmt.Println(r)
		buildCondPrepare(fieldInfo, db.Cond{}, rValue, -1)
		sel := &selector{}
		hasMustCol := true
		buildSelector(fieldInfo, sel, `aid`, &hasMustCol, nil)
		r = getRelationCache(fieldInfo.Field)
		fmt.Println(r)
		assert.Equal(t, `link2`, fieldInfo.Options[DBConnTagName])
		assert.Equal(t, &[]*kv{
			{
				k: `group`,
				v: `testc`,
			},
		}, r.where)
	}

	b := &B{
		AID:  10,
		Name: `10`,
	}
	rValue = reflect.ValueOf(b)
	rType = reflect.TypeOf(b)
	typeMap = mapper.TypeMap(rType)
	options, ok = typeMap.Options[`relation`]
	if !ok {
		panic(`not ok`)
	}
	// for key, val := range typeMap.Options {
	// 	fmt.Printf("%v ==> %+v\n", key, val[0])
	// }
	assert.Equal(t, 4, len(typeMap.Options))
	for _, fieldInfo := range options {
		r, err := parseRelationData(fieldInfo)
		assert.NoError(t, err)
		fmt.Println(r)
		buildCondPrepare(fieldInfo, db.Cond{}, rValue, -1)
		sel := &selector{}
		hasMustCol := true
		buildSelector(fieldInfo, sel, `aid`, &hasMustCol, nil)
		r = getRelationCache(fieldInfo.Field)
		fmt.Println(r)
		assert.Equal(t, `link3`, fieldInfo.Options[DBConnTagName])
		assert.Equal(t, &[]*kv{
			{
				k: `group`,
				v: `testb`,
			},
		}, r.where)
	}
}
