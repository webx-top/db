package mysql

import (
	"context"
	"database/sql"
	"regexp"

	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo/param"
)

var (
	reTable = regexp.MustCompile("CREATE TABLE `[^`]+` \\(")
)

// CopyTableStruct 复制表结构到新表
func CopyTableStruct(srcLinkID int, srcDBName string, srcTableName string,
	destLinkID int, destDBName string, destTableName string) error {
	ddl, err := CreateTableSQL(srcLinkID, srcDBName, srcTableName)
	if err != nil {
		return err
	}
	db := factory.NewParam(destLinkID).DB()
	ddl = reTable.ReplaceAllString(ddl, "CREATE TABLE `"+destTableName+"` (")
	_, err = db.Exec(ddl)
	return err
}

// TableExists 查询表是否存在
func TableExists(linkID int, dbName string, tableName string) (bool, error) {
	ctx := context.Background()
	db := factory.NewParam(linkID).DB()
	stmt, err := db.PrepareContext(ctx, SQLTableExists)
	if err != nil {
		return false, err
	}
	recv := sql.NullString{}
	err = stmt.QueryRowContext(ctx, dbName, tableName).Scan(&recv)
	if err != nil {
		return false, err
	}
	return param.String(recv.String).Int64() > 0, err
}

// FieldExists 查询表字段是否存在
func FieldExists(linkID int, tableName string, fieldName string) (bool, error) {
	ctx := context.Background()
	db := factory.NewParam(linkID).DB()
	stmt, err := db.PrepareContext(ctx, SQLFieldExists)
	if err != nil {
		return false, err
	}
	recv := sql.NullString{}
	err = stmt.QueryRowContext(ctx, tableName, fieldName).Scan(&recv)
	if err != nil {
		return false, err
	}
	return param.String(recv.String).Int64() > 0, err
}

// MoveToNewTable 移动数据到新表
func MoveToNewTable(linkID int, srcTableName string, destTableName string, condition string, src2DestFieldMapping ...map[string]string) (int64, error) {
	sqlStr := "INSERT INTO `" + destTableName + "`"
	var srcFields, destFields string
	if len(src2DestFieldMapping) > 0 {
		var sep string
		for srcField, destField := range src2DestFieldMapping[0] {
			srcFields += sep + "`" + srcField + "`"
			destFields += sep + "`" + destField + "`"
			sep = ","
		}
		destFields = `(` + destFields + `)`
	}
	sqlStr += destFields
	if len(srcFields) < 1 {
		srcFields = `*`
	}
	sqlStr += ` SELECT ` + srcFields + ` FROM ` + "`" + srcTableName + "`"
	if len(condition) > 0 {
		sqlStr += ` WHERE ` + condition
	}
	db := factory.NewParam(linkID).DB()
	result, err := db.Exec(sqlStr)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return affected, err
	}
	if affected > 0 {
		sqlStr := "DELETE FROM `" + srcTableName + "`"
		if len(condition) > 0 {
			sqlStr += ` WHERE ` + condition
		}
		db := factory.NewParam(linkID).DB()
		result, err := db.Exec(sqlStr)
		if err != nil {
			return affected, err
		}
		affected, err = result.RowsAffected()
	}
	return affected, err
}
