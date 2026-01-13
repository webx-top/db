package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	mySQLUtil "github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/echo/param"
)

func getMySQLTableComment(d sqlbuilder.Database, tableName string) (string, error) {
	row, err := d.QueryRow(mySQLUtil.SQLTableComment, d.Name(), tableName)
	if err != nil {
		return ``, err
	}
	recvTableName := sql.NullString{}
	recvTableComment := sql.NullString{}
	err = row.Scan(&recvTableName, &recvTableComment)
	if err != nil {
		return ``, fmt.Errorf(`TableComment.Scan: %v`, err)
	}
	return recvTableComment.String, err
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
		} else if colExtra.String == `VIRTUAL GENERATED` || colExtra.String == `STORED GENERATED` {
			continue
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

func getMySQLTableFields(cfg *config, db sqlbuilder.Database, tableName string, typeMap map[string][]string) ([]string, map[string]factory.FieldInfo, []string) {

	fieldMaxLength, fieldsInfo := getMySQLTableInfo(db, tableName)
	goFields := []string{}
	fields := map[string]factory.FieldInfo{}
	fieldNames := make([]string, len(fieldsInfo))
	var sets []func(*structField, *factory.FieldInfo)
	if typeMap != nil {
		if typef, ok := typeMap[`hashids`]; ok && len(typef) > 0 {
			sets = append(sets, func(sf *structField, fi *factory.FieldInfo) {
				if com.InSlice(fi.Name, typef) {
					sf.typ = `hashseq.ID`
					fi.MyType = sf.typ
				}
			})
		}
	}
	if cfg.AutoTimeFields != nil {
		_fieldNames, _ := cfg.AutoTimeFields.GetUpdateFieldNames(tableName)
		_ifieldNames, _ := cfg.AutoTimeFields.GetInsertFieldNames(tableName)
		_fieldNames = append(_fieldNames, _ifieldNames...)
		if len(_fieldNames) > 0 {
			sets = append(sets, func(sf *structField, fi *factory.FieldInfo) {
				if com.InSlice(fi.Name, _fieldNames) {
					switch fi.GoType {
					case `uint`, `int`, `uint32`, `int32`, `int64`, `uint64`:
						sf.otherTag = `form_decoder:"time2unix" form_encoder:"unix2time"`
					case `string`:
						//TODO
					}
				}
			})
		}
	}
	for key, field := range fieldsInfo {
		goField, fieldInfo := getMySQLFieldInfo(cfg, field, fieldMaxLength, fields)
		for _, set := range sets {
			set(goField, &fieldInfo)
		}
		goFields = append(goFields, goField.String())
		fields[fieldInfo.Name] = fieldInfo
		fieldNames[key] = fieldInfo.Name
	}
	return goFields, fields, fieldNames
}

func getMySQLFieldInfo(cfg *config, field map[string]string, maxLength int, fields map[string]factory.FieldInfo) (*structField, factory.FieldInfo) {

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
				if fieldInfo.Precision > 0 {
					numStr := strings.Repeat(`9`, fieldInfo.MaxSize*255)
					end := fieldInfo.MaxSize - fieldInfo.Precision //(4,2): 9999=>99.99
					numStr = numStr[:end] + `.` + numStr[end:]
					fieldInfo.Max = param.AsFloat64(numStr)
				} else {
					fieldInfo.Max = float64(fieldInfo.MaxSize) * 255
				}
				if fieldInfo.Unsigned {
					fieldInfo.Min = 0
				} else {
					fieldInfo.Min = (fieldInfo.Max - 1) / 2 * -1
					fieldInfo.Max = fieldInfo.Min*-1 + 1
				}
			}
		}
	} else {
		vs := strings.Split(field["Type"], ` `)
		fieldInfo.DataType = vs[0]
		if len(vs) > 1 && vs[1] == `unsigned` {
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
	var otherTag string
	fieldInfo.Comment = field["Comment"]
	fieldInfo.DefaultValue = field["Default"]
	var magicTags []string
	if strings.HasPrefix(fieldInfo.Comment, "`") { // `magicTag1,magicTag2`
		p := strings.Index(fieldInfo.Comment[1:], "`")
		if p > -1 {
			magicTags = strings.Split(fieldInfo.Comment[1:p+1], `,`)
			fieldInfo.Comment = fieldInfo.Comment[p+2:]
		}
	}
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
			} else if len(magicTags) > 0 {
				for _, t := range magicTags {
					switch t {
					case `omitempty`:
						dbTag += ",omitempty"
						bsonTag += ",omitempty"
					case `pk`:
						dbTag += ",pk"
						fieldInfo.PrimaryKey = true
					case `timestamp`:
						if fieldInfo.IsInteger() {
							otherTag = `form_decoder:"time2unix" form_encoder:"unix2time"`
						}
					}
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
	jsonTag := fieldInfo.Name
	xmlTag := fieldInfo.Name
	if cfg.FieldEncodeType(`json`) != `table` {
		jsonTag = fieldInfo.GoName
	}
	if cfg.FieldEncodeType(`xml`) != `table` {
		xmlTag = fieldInfo.GoName
	}
	if cfg.FieldEncodeType(`bson`) != `table` {
		bsonTag = fieldInfo.GoName
	}
	if cfg.FieldEncodeType(`db`) != `table` {
		dbTag = fieldInfo.GoName
	}
	for _, t := range magicTags {
		parts := strings.SplitN(t, `:`, 2) // omit:json|xml or omit:encode
		var targets []string
		if len(parts) >= 2 {
			temp := map[string]struct{}{}
			for _, tg := range strings.Split(parts[1], `|`) {
				tg := strings.TrimSpace(tg)
				if len(tg) == 0 {
					continue
				}
				if _, ok := temp[tg]; ok {
					continue
				}
				temp[tg] = struct{}{}
				targets = append(targets, tg)
			}
		}
		switch parts[0] {
		case `omit`:
			for _, tg := range targets {
				switch tg {
				case `json`:
					jsonTag = `-`
				case `xml`:
					xmlTag = `-`
				case `encode`:
					xmlTag = `-`
					jsonTag = `-`
				default:
				}
			}
		case `-`:
			xmlTag = `-`
			jsonTag = `-`
		case `omitempty`:
			for _, tg := range targets {
				switch tg {
				case `json`:
					jsonTag += `,omitempty`
				case `xml`:
					xmlTag += `,omitempty`
				case `encode`:
					xmlTag += `,omitempty`
					jsonTag += `,omitempty`
				default:
				}
			}
		case `time`:
			if len(targets) > 0 && len(targets[0]) > 0 {
				switch targets[0] {
				case `unix`:
					if fieldInfo.IsInteger() {
						otherTag = `form_decoder:"time2unix" form_encoder:"unix2time"`
					}
				case `unixmilli`:
					if fieldInfo.IsInteger() {
						otherTag = `form_decoder:"time2unixmilli" form_encoder:"unixmilli2time"`
					}
				case `unixmicro`:
					if fieldInfo.IsInteger() {
						otherTag = `form_decoder:"time2unixmicro" form_encoder:"unixmicro2time"`
					}
				case `unixnano`: // 整数/小数/字符串
					if fieldInfo.IsNumeric() || fieldInfo.GoType == `string` {
						otherTag = `form_decoder:"time2unixnano" form_encoder:"unixnano2time"`
					}
				}
			}
		case `i18n`:
			fieldInfo.Multilingual = true
		}
	}
	fieldBlock := &structField{
		field:    fieldP,
		typ:      typeP,
		comment:  fieldInfo.Comment,
		dbTag:    dbTag,
		bsonTag:  bsonTag,
		jsonTag:  jsonTag,
		xmlTag:   xmlTag,
		otherTag: otherTag,
	}
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

/*
mysqldump 参数说明：
-d 			结构(--no-data:不导出任何数据，只导出数据库表结构)
-t 			数据(--no-create-info:只导出数据，而不添加CREATE TABLE 语句)
-n 			(--no-create-db:只导出数据，而不添加CREATE DATABASE 语句）
-R 			(--routines:导出存储过程以及自定义函数)
-E 			(--events:导出事件)
--triggers 	(默认导出触发器，使用--skip-triggers屏蔽导出)
-B 			(--databases:导出数据库列表，单个库时可省略）
--tables 	表列表（单个表时可省略）
*/

var cleanRegExp = regexp.MustCompile(` AUTO_INCREMENT=[0-9]*\s*`)

func execBackupCommand(cfg *config, tables []string) {
	if len(cfg.Backup) == 0 || len(tables) == 0 {
		return
	}
	log.Println(`Starting backup:`, tables)
	var port, host string
	if p := strings.LastIndex(cfg.Host, `:`); p > 0 {
		host = cfg.Host[0:p]
		port = cfg.Host[p+1:]
	} else {
		host = cfg.Host
	}
	if len(port) == 0 {
		port = `3306`
	}
	if len(host) == 0 {
		host = `127.0.0.1`
	}
	args := []string{
		"--default-character-set=" + cfg.Charset,
		"--single-transaction",
		"--set-gtid-purged=OFF",
		"--opt",
		"-d", //加上此参数代表只导出表结构，不导出数据
		"-h" + host,
		"-P" + port,
		"-u" + cfg.Username,
		"-p" + cfg.Password,
		cfg.Database,
	}
	var structFile, dataFile string
	files := strings.SplitN(cfg.Backup, `|`, 2)
	switch len(files) {
	case 2:
		dataFile = strings.TrimSpace(files[1])
		fallthrough
	case 1:
		structFile = strings.TrimSpace(files[0])
	}
	for index, saveFile := range []string{structFile, dataFile} {
		if len(saveFile) == 0 {
			continue
		}
		info := strings.SplitN(saveFile, ":", 2)
		saveFile = info[0]
		var saveTables []string
		if len(info) > 1 && len(info[1]) > 0 {
			saveTables = strings.Split(info[1], `,`)
		}
		if len(saveTables) == 0 {
			saveTables = append(saveTables, tables...)
		}
		cmdArgs := append([]string{}, args...)
		cmdArgs = append(cmdArgs, saveTables...)
		if index > 0 {
			cmdArgs[4] = `-t` //导出数据
		}
		var cmd *exec.Cmd
		if len(cfg.Container) > 0 {
			cmd = exec.Command("docker", append([]string{
				`exec`, cfg.Container, "mysqldump",
			}, cmdArgs...)...)
		} else {
			cmd = exec.Command("mysqldump", cmdArgs...)
		}
		cmd.Stderr = os.Stderr
		fp, err := os.Create(saveFile)
		if err != nil {
			log.Fatal(`Failed to backup:`, err)
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fp.Close()
			log.Fatal(`Failed to backup(StdoutPipe):`, err)
		}
		close := func() {
			fp.Close()
			stdout.Close()
		}
		if err := cmd.Start(); err != nil {
			close()
			log.Fatal(`Failed to backup(Start):`, err)
		}
		if _, err := io.Copy(fp, stdout); err != nil {
			close()
			log.Fatal(`Failed to backup(io.Copy):`, err)
		}
		if err := cmd.Wait(); err != nil {
			close()
			log.Fatal(`Failed to backup(Wait):`, err)
		}
		close()
		if index == 0 {
			b, err := os.ReadFile(saveFile)
			if err != nil {
				log.Fatal(err)
			}
			b = cleanRegExp.ReplaceAll(b, []byte(` `))
			err = os.WriteFile(saveFile, b, 0666)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
