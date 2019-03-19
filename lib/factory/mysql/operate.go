package mysql

import (
	"regexp"

	"github.com/webx-top/db/lib/factory"
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
