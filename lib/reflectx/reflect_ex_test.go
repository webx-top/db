package reflectx

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
)

type User struct {
	Name string `db:"name"`
	Id   uint   `db:"id,pk"`
	Age  uint
}

type Profile struct {
	Email  string `db:"email"`
	Mobile string `db:"mobile"`
}

type Group struct {
	Discount float64 `db:"discount"`
	Name     string  `db:"name"`
}

type Data struct {
	User    *User
	Profile *Profile `db:"profile"`
	Group   *Group   `db:",alias=g"`
}

func TestFind(t *testing.T) {
	mapper := NewMapper(`db`)
	data := &Data{}
	typeMap := mapper.StructMap(data)

	assert.Equal(t, `UserProfile.EMail`, com.PascalCaseWith(`user_profile.e_mail`, '.'))

	row := reflect.ValueOf(data)
	fv := mapper.FieldByName(row, `profile`)
	assert.True(t, reflect.Indirect(fv).IsZero())
	fv = mapper.FieldByName(row, `profile.email`)
	assert.True(t, reflect.Indirect(fv).IsZero())

	_, exists := typeMap.Find(`user.name`, true)
	assert.True(t, exists)

	_, exists = typeMap.Find(`user.age`, true)
	assert.True(t, exists)

	_, exists = typeMap.Find(`user.no`, true)
	assert.False(t, exists)

	fieldInfo, exists := typeMap.Find(`User.Id`, true)
	assert.True(t, exists)
	_, exists = fieldInfo.Options[`pk`]
	assert.True(t, exists)

	fieldPath, exists := typeMap.FindTableField(`User.Id`, true)
	assert.True(t, exists)
	assert.Equal(t, `User.id`, fieldPath)

	fieldPath, exists = typeMap.FindTableField(`Profile.Email`, true)
	assert.True(t, exists)
	assert.Equal(t, `profile.email`, fieldPath)

	fieldPath, exists = typeMap.FindTableField(`group.discount`, true)
	assert.True(t, exists)
	assert.Equal(t, `g.discount`, fieldPath)

}

func TestFindMap(t *testing.T) {
	mapper := NewMapper(`db`)
	data := &Data{}
	typeMap := mapper.StructMap(data)
	tableFields, pk := typeMap.FindTableFieldByMap(map[string]map[string]interface{}{
		`user`: {
			`id`: nil,
		},
		`profile`: {
			`mobile`: nil,
		},
	}, true)
	assert.Equal(t, `User.id`, pk[0])
	for field := range tableFields {
		fmt.Println(`table field:`, field)
	}
	assert.Equal(t, `Id`, tableFields[`User.id`].FieldInfo.Field.Name)
	assert.Equal(t, nil, tableFields[`User.id`].RawData)
	assert.Equal(t, []string{`user`, `id`}, tableFields[`User.id`].RawPath)
	assert.Equal(t, `Mobile`, tableFields[`profile.mobile`].FieldInfo.Field.Name)
}

func TestSlice(t *testing.T) {
	mapper := NewMapper(`db`)
	data := &[]*Data{}
	typeMap := mapper.StructMap(data)
	assert.Equal(t, `Data`, typeMap.Tree.Name)
	assert.Equal(t, []int(nil), typeMap.Tree.Index)
	assert.Nil(t, typeMap.Tree.Parent)
	assert.Equal(t, ``, typeMap.Tree.Path)
	assert.Equal(t, `Data`, typeMap.Tree.Zero.Type().Name())
}
