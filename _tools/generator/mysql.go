package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/sqlbuilder"
)

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
		sz := len(colField.String)
		if sz > fieldMaxLength {
			fieldMaxLength = sz
		}
		fieldsInfo = append(fieldsInfo, result)
		//log.Printf(`%#v`+"\n", remap)
	}
	return fieldMaxLength, fieldsInfo
}

func getMySQLTableFields(db sqlbuilder.Database, tableName string) ([]string, map[string]factory.FieldInfo) {

	fieldMaxLength, fieldsInfo := getMySQLTableInfo(db, tableName)
	goFields := []string{}
	fields := map[string]factory.FieldInfo{}
	for _, field := range fieldsInfo {
		goField, fieldInfo := getMySQLFieldInfo(field, fieldMaxLength)
		goFields = append(goFields, goField)
		fields[fieldInfo.Name] = fieldInfo
	}
	return goFields, fields
}

func getMySQLFieldInfo(field map[string]string, maxLength int) (string, factory.FieldInfo) {

	fieldInfo := factory.FieldInfo{Options: []string{}}
	p := strings.Index(field["Type"], `(`)
	fieldInfo.Name = field["Field"]
	if p > -1 {
		fieldInfo.DataType = field["Type"][0:p]
		pr := strings.Index(field["Type"], `)`)
		if pr > -1 {
			opts := field["Type"][p+1 : pr]
			if len(opts) > 0 {
				var err error
				if opts[0] == '\'' {
					for _, opt := range strings.Split(opts, `,`) {
						fieldInfo.Options = append(fieldInfo.Options, strings.Trim(opt, `'`))
					}
				} else if strings.Contains(opts, `,`) {
					opts := strings.Split(opts, `,`)
					switch len(opts) {
					case 2:
						fieldInfo.MaxSize, err = strconv.Atoi(opts[0])
						if err != nil {
							panic(err)
						}
						fieldInfo.Precision, err = strconv.Atoi(opts[1])
						if err != nil {
							panic(err)
						}
					}
				} else {
					fieldInfo.MaxSize, err = strconv.Atoi(opts)
					if err != nil {
						panic(err)
					}
				}
			}
			if vs := strings.Split(field["Type"][pr:], ` `); len(vs) > 1 && vs[1] == `unsigned` {
				fieldInfo.Unsigned = true
			}
		}
	} else {
		if vs := strings.Split(field["Type"], ` `); len(vs) > 1 && vs[1] == `unsigned` {
			fieldInfo.Unsigned = true
		}
	}

	fieldInfo.GoType = DataType(&fieldInfo)
	fieldInfo.GoName = TableToStructName(fieldInfo.Name, ``)
	fieldP := fmt.Sprintf(`%-*s`, maxLength, fieldInfo.GoName)
	typeP := fmt.Sprintf(`%-8s`, fieldInfo.GoType)
	dbTag := fieldInfo.Name
	bsonTag := fieldInfo.Name
	fieldInfo.Comment = field["Comment"]
	fieldInfo.DefaultValue = field["Default"]
	if field["Key"] == "PRI" && field["Extra"] == "auto_increment" {
		dbTag += ",omitempty,pk"
		bsonTag += ",omitempty"
		fieldInfo.PrimaryKey = true
		fieldInfo.AutoIncrement = true
	} else {
		if field["Key"] == "PRI" {
			dbTag += ",pk"
			fieldInfo.PrimaryKey = true
		}
		if len(fieldInfo.Comment) > 0 {
			//支持注释内容为：`omitempty`我是注释内容
			if fieldInfo.Comment == "`omitempty`" {
				dbTag += ",omitempty"
				bsonTag += ",omitempty"
				fieldInfo.Comment = ""
			} else if strings.HasPrefix(fieldInfo.Comment, "`") {
				p := strings.Index(fieldInfo.Comment[1:], "`")
				if p > -1 {
					for _, t := range strings.Split(fieldInfo.Comment[1:p+1], `,`) {
						fmt.Println(t)
						switch t {
						case `omitempty`:
							dbTag += ",omitempty"
							bsonTag += ",omitempty"
						case `pk`:
							dbTag += ",pk"
							fieldInfo.PrimaryKey = true
						}
					}
					fieldInfo.Comment = fieldInfo.Comment[p+2:]
				}
			}
		}
	}
	fieldBlock := fmt.Sprintf(memberTemplate, fieldP, typeP, dbTag, bsonTag, fieldInfo.Comment, fieldInfo.Name, fieldInfo.Name)
	return fieldBlock, fieldInfo
}
