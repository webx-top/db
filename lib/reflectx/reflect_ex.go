package reflectx

import (
	"reflect"
	"strings"
)

// StructMap returns a mapping of field strings to int slices representing
// the traversal down the struct to reach the field.
func (m *Mapper) StructMap(bean interface{}) *StructMap {
	return m.TypeMap(reflect.ValueOf(bean).Type())
}

// Find Find("user.profile")
func (f StructMap) Find(structFieldPath string) (tree *FieldInfo, exists bool) {
	tree = f.Tree
	for _, field := range strings.Split(structFieldPath, `.`) {
		field = strings.Title(field)
		var found bool
		for _, fieldInfo := range tree.Children {
			if fieldInfo == nil {
				continue
			}
			if fieldInfo.Field.Name == field {
				tree = fieldInfo
				found = true
				break
			}
		}
		if !found {
			return nil, false
		}
		exists = true
	}
	return
}

// FindTableField Find("User.Profile")
func (f StructMap) FindTableField(structFieldPath string, aliasOptionNames ...string) (tableFieldPath string, exists bool) {
	tree := f.Tree
	var aliasOptionName string
	if len(aliasOptionNames) > 0 {
		aliasOptionName = aliasOptionNames[0]
	}
	if len(aliasOptionName) == 0 {
		aliasOptionName = `alias`
	}
	for index, field := range strings.Split(structFieldPath, `.`) {
		field = strings.Title(field)
		var found bool
		for _, fieldInfo := range tree.Children {
			if fieldInfo == nil {
				continue
			}
			if fieldInfo.Field.Name == field {
				alias, _ := fieldInfo.Options[aliasOptionName]
				if index > 0 {
					tableFieldPath += `.`
				}
				if len(alias) == 0 {
					if len(fieldInfo.Name) > 0 {
						alias = fieldInfo.Name
					} else {
						alias = fieldInfo.Field.Name
					}
				}
				tableFieldPath += alias
				tree = fieldInfo
				found = true
				break
			}
		}
		if !found {
			return tableFieldPath, false
		}
		exists = true
	}
	return
}
