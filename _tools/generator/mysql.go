package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/echo/param"
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

func getMySQLTableFields(db sqlbuilder.Database, tableName string) ([]string, map[string]factory.FieldInfo, []string) {

	fieldMaxLength, fieldsInfo := getMySQLTableInfo(db, tableName)
	goFields := []string{}
	fields := map[string]factory.FieldInfo{}
	fieldNames := make([]string, len(fieldsInfo))
	for key, field := range fieldsInfo {
		goField, fieldInfo := getMySQLFieldInfo(field, fieldMaxLength, fields)
		goFields = append(goFields, goField)
		fields[fieldInfo.Name] = fieldInfo
		fieldNames[key] = fieldInfo.Name
	}
	return goFields, fields, fieldNames
}

func getMySQLFieldInfo(field map[string]string, maxLength int, fields map[string]factory.FieldInfo) (string, factory.FieldInfo) {

	fieldInfo := factory.FieldInfo{Options: []string{}}
	p := strings.Index(field["Type"], `(`)
	fieldInfo.Name = field["Field"]
	if p > -1 {
		fieldInfo.DataType = field["Type"][0:p]
		pr := strings.Index(field["Type"], `)`)
		if pr > -1 {
			opts := field["Type"][p+1 : pr]
			var isNum bool
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
						isNum = true
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
			if isNum {
				numStr := strings.Repeat(`9`, fieldInfo.MaxSize)
				if fieldInfo.Precision > 0 {
					end := fieldInfo.MaxSize - fieldInfo.Precision //(4,2): 9999=>99.99
					numStr = numStr[:end] + `.` + numStr[end:]
				}
				fieldInfo.Max = param.AsFloat64(numStr)
				if fieldInfo.Unsigned {
					fieldInfo.Min = 0
				} else {
					fieldInfo.Min = fieldInfo.Max * -1
				}
			}
		}
	} else {
		if vs := strings.Split(field["Type"], ` `); len(vs) > 1 && vs[1] == `unsigned` {
			fieldInfo.Unsigned = true
		}
	}

	fieldInfo.GoType = DataType(&fieldInfo)
	fieldInfo.GoName = TableToStructName(fieldInfo.Name, ``)

	//避免和默认方法名冲突，对于已经存在方法名的字段，在其名称后加后缀“V+编号”
	if _, exists := structFuncs[fieldInfo.GoName]; exists {
		var suffix string
		for i := 0; ; i++ {
			if i > 0 {
				suffix = fmt.Sprintf(`V%d`, i)
			} else {
				suffix = `V`
			}
			exists = false
			for _, f := range fields {
				if f.GoName == fieldInfo.GoName+suffix {
					exists = true
					break
				}
			}
			if !exists {
				break
			}
		}
		fieldInfo.GoName += suffix
	}

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
	if rg, ok := mysqlNumericRange[fieldInfo.DataType]; ok {
		if fieldInfo.Unsigned {
			fieldInfo.Min = rg.Unsigned.Min
			fieldInfo.Max = rg.Unsigned.Max
		} else {
			fieldInfo.Min = rg.Signed.Min
			fieldInfo.Max = rg.Signed.Max
		}
		maxSize := len(fmt.Sprint(fieldInfo.Max))
		if maxSize > fieldInfo.MaxSize {
			fieldInfo.Max = com.Float64(strings.Repeat(`9`, fieldInfo.MaxSize))
			if fieldInfo.Min < 0 {
				fieldInfo.Min = fieldInfo.Max * -1
			}
		}
	}
	encodedField := fieldInfo.Name
	if len(cfg.EncFieldFormat) > 0 && cfg.EncFieldFormat != `table` {
		encodedField = fieldInfo.GoName
	}
	fieldBlock := fmt.Sprintf(memberTemplate, fieldP, typeP, dbTag, bsonTag, fieldInfo.Comment, encodedField, encodedField)
	return fieldBlock, fieldInfo
}

type NumericRange struct {
	Min float64
	Max float64
}
type NumericRanges struct {
	Unsigned *NumericRange
	Signed   *NumericRange
}

var mysqlNumericRange = map[string]*NumericRanges{
	`tinyint`: {
		Unsigned: &NumericRange{Min: 0, Max: 255},
		Signed:   &NumericRange{Min: -128, Max: 127},
	},
	`smallint`: {
		Unsigned: &NumericRange{Min: 0, Max: 65535},
		Signed:   &NumericRange{Min: -32768, Max: 32767},
	},
	`mediumint`: {
		Unsigned: &NumericRange{Min: 0, Max: 16777215},
		Signed:   &NumericRange{Min: -8388608, Max: 8388607},
	},
	`int`: {
		Unsigned: &NumericRange{Min: 0, Max: 4294967295},
		Signed:   &NumericRange{Min: -2147483648, Max: 2147483647},
	},
	`bigint`: {
		Unsigned: &NumericRange{Min: 0, Max: 18446744073709551615},
		Signed:   &NumericRange{Min: -9233372036854775808, Max: 9233372036854775807},
	},
}
