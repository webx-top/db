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

type Data struct {
	User    *User
	Profile *Profile
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
}
