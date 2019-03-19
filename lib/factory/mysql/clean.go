package mysql

import "github.com/webx-top/echo"

// Clean 清理掉表中不存在的字段
func Clean(linkID int, dbName string, tableName string, data echo.H, excludeFields ...string) (echo.H, error) {
	columns, err := ColumnComment(linkID, dbName, tableName)
	if err != nil {
		return data, err
	}
	if len(excludeFields) > 0 {
		for _, field := range excludeFields {
			if data.Has(field) {
				delete(data, field)
			}
		}
	}
	for field := range data {
		if _, ok := columns[field]; !ok {
			delete(data, field)
		}
	}
	return data, err
}
