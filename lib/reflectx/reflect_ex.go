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
func (f StructMap) Find(fieldPath string, isStructField bool) (tree *FieldInfo, exists bool) {
	tree = f.Tree
	for _, field := range strings.Split(fieldPath, `.`) {
		if len(field) == 0 {
			return nil, false
		}
		if isStructField {
			field = strings.Title(field)
		}
		var found bool
		for _, fieldInfo := range tree.Children {
			if fieldInfo == nil {
				continue
			}
			var equal bool
			if isStructField {
				equal = fieldInfo.Field.Name == field
			} else {
				equal = fieldInfo.Name == field
			}
			if equal {
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

func joinTableFieldPath(fieldInfo *FieldInfo, aliasOptionName string, tableFieldPath string) string {
	alias, _ := fieldInfo.Options[aliasOptionName]
	if len(alias) == 0 {
		if len(fieldInfo.Name) > 0 {
			alias = fieldInfo.Name
		} else {
			alias = fieldInfo.Field.Name
		}
	}
	tableFieldPath += `.` + alias
	return tableFieldPath
}

// FindTableField Find("User.Profile")
func (f StructMap) FindTableField(fieldPath string, isStructField bool, aliasOptionNames ...string) (tableFieldPath string, exists bool) {
	tree := f.Tree
	var aliasOptionName string
	if len(aliasOptionNames) > 0 {
		aliasOptionName = aliasOptionNames[0]
	}
	if len(aliasOptionName) == 0 {
		aliasOptionName = `alias`
	}
	for _, field := range strings.Split(fieldPath, `.`) {
		if len(field) == 0 {
			return strings.TrimPrefix(tableFieldPath, `.`), false
		}
		if isStructField {
			field = strings.Title(field)
		}
		var found bool
		for _, fieldInfo := range tree.Children {
			if fieldInfo == nil {
				continue
			}
			var equal bool
			if isStructField {
				equal = fieldInfo.Field.Name == field
			} else {
				equal = fieldInfo.Name == field
			}
			if equal {
				tableFieldPath = joinTableFieldPath(fieldInfo, aliasOptionName, tableFieldPath)
				tree = fieldInfo
				found = true
				break
			}
		}
		if !found {
			return strings.TrimPrefix(tableFieldPath, `.`), false
		}
		exists = true
	}
	tableFieldPath = strings.TrimPrefix(tableFieldPath, `.`)
	return
}
