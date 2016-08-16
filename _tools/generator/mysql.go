package main

import (
	"database/sql"
	"log"

	"github.com/webx-top/db/lib/sqlbuilder"
)

func getMySQLTables(d sqlbuilder.Database) []string {
	rows, err := d.Query(`SHOW TABLES`)
	if err != nil {
		log.Fatal(err)
	}
	tables := []string{}
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			log.Println(err)
			continue
		}
		tables = append(tables, tableName)
	}
	return tables
}

func getMySQLTableInfo(d sqlbuilder.Database, tableName string) (int, []map[string]string) {
	rows, err := d.Query("SHOW FULL COLUMNS FROM `" + tableName + "`")
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

		err = rows.Scan(&colField, &colType, &colCollation, &colNull, &colKey, &colDefault, &colExtra, &colPrivileges, &colComment)
		if err != nil {
			log.Println(err)
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
