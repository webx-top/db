package reflectx

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	_, exists := typeMap.Find(`user.name`)
	assert.True(t, exists)

	_, exists = typeMap.Find(`user.age`)
	assert.True(t, exists)

	_, exists = typeMap.Find(`user.no`)
	assert.False(t, exists)

	fieldInfo, exists := typeMap.Find(`User.Id`)
	assert.True(t, exists)
	_, exists = fieldInfo.Options[`pk`]
	assert.True(t, exists)

	fieldPath, exists := typeMap.FindTableField(`User.Id`)
	assert.True(t, exists)
	assert.Equal(t, `User.id`, fieldPath)

	fieldPath, exists = typeMap.FindTableField(`Profile.Email`)
	assert.True(t, exists)
	assert.Equal(t, `profile.email`, fieldPath)

	fieldPath, exists = typeMap.FindTableField(`group.discount`)
	assert.True(t, exists)
	assert.Equal(t, `g.discount`, fieldPath)
}
