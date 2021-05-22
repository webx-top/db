package main

import (
	"database/sql"
	"log"

	"github.com/webx-top/db/lib/sqlbuilder"
)

func getPostgreSQLTableInfo(d sqlbuilder.Database, tableName string) (int, []map[string]string) {
	rows, err := d.Query(`SELECT column_name, data_type, is_nullable, column_default,
    CASE WHEN p.contype = 'p' THEN true ELSE false END AS primarykey,
    CASE WHEN p.contype = 'u' THEN true ELSE false END AS uniquekey
FROM pg_attribute f
    JOIN pg_class c ON c.oid = f.attrelid JOIN pg_type t ON t.oid = f.atttypid
    LEFT JOIN pg_attrdef d ON d.adrelid = c.oid AND d.adnum = f.attnum
    LEFT JOIN pg_namespace n ON n.oid = c.relnamespace
    LEFT JOIN pg_constraint p ON p.conrelid = c.oid AND f.attnum = ANY (p.conkey)
    LEFT JOIN pg_class AS g ON p.confrelid = g.oid
    LEFT JOIN INFORMATION_SCHEMA.COLUMNS s ON s.column_name=f.attname AND c.relname=s.table_name
WHERE c.relkind = 'r'::char AND c.relname = '` + tableName + `' AND s.table_schema = '` + cfg.Schema + `' AND f.attnum > 0 ORDER BY f.attnum`)
	if err != nil {
		log.Println(err)
		panic(err.Error())
	}
	fieldsInfo := []map[string]string{}
	fieldMaxLength := 0
	for rows.Next() {

		var (
			colField      sql.NullString
			colType       sql.NullString
			colCollation  sql.NullString
			colNull       sql.NullString
			colKey        sql.NullString
			colDefault    sql.NullString
			colExtra      sql.NullString
			colPrivileges sql.NullString
			colComment    sql.NullString
		)
		var isPK, isUnique bool
		err = rows.Scan(&colField, &colType, &colNull, &colDefault, &isPK, &isUnique)
		if err != nil {
			log.Println(err)
		}
		if isPK {
			colKey.String = `PRI`
			colExtra.String = `auto_increment`
		}
		result := map[string]string{
			"Field":      colField.String,
			"Type":       colType.String,
			"Collation":  colCollation.String,
			"Null":       colNull.String,
			"Key":        colKey.String,
			"Default":    colDefault.String,
			"Extra":      colExtra.String,
			"Privileges": colPrivileges.String,
			"Comment":    colComment.String,
		}
		for _, v := range result {
			sz := len(v)
			if sz > fieldMaxLength {
				fieldMaxLength = sz
			}
		}
		fieldsInfo = append(fieldsInfo, result)
		//log.Printf(`%#v`+"\n", remap)
	}
	return fieldMaxLength, fieldsInfo
}
