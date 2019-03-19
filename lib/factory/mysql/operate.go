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
