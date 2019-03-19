package mysql

import (
	"context"
	"database/sql"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo/param"
)

const (
	// SQLTableComment 查询表注释的SQL
	SQLTableComment = `SELECT TABLE_NAME,TABLE_COMMENT FROM information_schema.TABLES WHERE table_schema=? AND TABLE_NAME=?`
	// SQLColumnComment 查询列注释的SQL
	SQLColumnComment = "SELECT COLUMN_NAME as `field`, column_comment as `description`, DATA_TYPE as `type`, CHARACTER_MAXIMUM_LENGTH as `max_length`, CHARACTER_OCTET_LENGTH as `octet_length`, NUMERIC_PRECISION as `precision` FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema=? AND table_name=?"
	// SQLShowCreate 查询建表语句的SQL
	SQLShowCreate  = "SHOW CREATE TABLE "
	SQLTableExists = "SELECT COUNT(1) AS count FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA=? AND TABLE_NAME=?"
	SQLFieldExists = "SELECT COUNT(1) AS count FROM information_schema.columns WHERE table_name=? AND column_name=?"
)

// CreateTableSQL 查询建表语句
func CreateTableSQL(linkID int, dbName string, tableName string) (string, error) {
	ctx := context.Background()
	db := factory.NewParam(linkID).DB()
	stmt, err := db.PrepareContext(ctx, SQLShowCreate+"`"+dbName+"`.`"+tableName+"`")
	if err != nil {
		return ``, err
	}
	recvTableName := sql.NullString{}
	recvCreateTableSQL := sql.NullString{}
	err = stmt.QueryRowContext(ctx, dbName, tableName).Scan(&recvTableName, &recvCreateTableSQL)
	if err != nil {
		return ``, err
	}
	return recvCreateTableSQL.String, err
}

// TableComment 查询表注释
func TableComment(linkID int, dbName string, tableName string) (string, error) {
	ctx := context.Background()
	db := factory.NewParam(linkID).DB()
	stmt, err := db.PrepareContext(ctx, SQLTableComment)
	if err != nil {
		return ``, err
	}
	recvTableName := sql.NullString{}
	recvTableComment := sql.NullString{}
	err = stmt.QueryRowContext(ctx, dbName, tableName).Scan(&recvTableName, &recvTableComment)
	if err != nil {
		return ``, err
	}
	return recvTableComment.String, err
}

// ColumnComment 查询表中某些列的注释
func ColumnComment(linkID int, dbName string, tableName string, fieldNames ...string) (map[string]param.StringMap, error) {
	ctx := context.Background()
	db := factory.NewParam(linkID).DB()
	sqlStr := SQLColumnComment
	if len(fieldNames) > 0 {
		if len(fieldNames) == 1 {
			sqlStr += `AND COLUMN_NAME = '` + com.AddSlashes(fieldNames[0]) + `'`
		} else {
			for key, val := range fieldNames {
				fieldNames[key] = com.AddSlashes(val)
			}
			sqlStr += `AND COLUMN_NAME IN ('` + strings.Join(fieldNames, `','`) + `')`
		}
	}
	stmt, err := db.PrepareContext(ctx, sqlStr)
	if err != nil {
		return nil, err
	}
	results := map[string]param.StringMap{}
	rows, err := stmt.QueryContext(ctx, dbName, tableName)
	if err != nil {
		return nil, err
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	indexes := make(map[string]int, len(cols))
	for idx, col := range cols {
		indexes[col] = idx
	}
	for rows.Next() {
		recv := make([]interface{}, len(cols))
		for idx := range cols {
			recv[idx] = interface{}(&sql.NullString{})
		}
		err := rows.Scan(recv...)
		if err != nil {
			return results, err
		}
		result := param.StringMap{}
		for col, idx := range indexes {
			result[col] = param.String(recv[idx].(*sql.NullString).String)
		}
		results[result.String(`field`)] = result
	}
	return results, err
}
