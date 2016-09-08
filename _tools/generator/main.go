package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/db/mysql"
)

var structTemplate = `//Generated by webx-top/db
package %[1]s

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	%[2]s
)

type %[3]s struct {
	trans	*factory.Transaction
	
%[4]s
}

func (this *%[3]s) SetTrans(trans *factory.Transaction) *%[3]s {
	this.trans = trans
	return this
}

func (this *%[3]s) Param() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetTrans(this.trans).SetCollection("%[5]s")
}

func (this *%[3]s) Get(mw func(db.Result) db.Result) error {
	return this.Param().SetRecv(this).SetMiddleware(mw).One()
}

func (this *%[3]s) List(mw func(db.Result) db.Result, page, size int) ([]*%[3]s, func() int64, error) {
	r := []*%[3]s{}
	counter, err := this.Param().SetPage(page).SetSize(size).SetRecv(&r).SetMiddleware(mw).List()
	return r, counter, err
}

func (this *%[3]s) ListByOffset(mw func(db.Result) db.Result, offset, size int) ([]*%[3]s, func() int64, error) {
	r := []*%[3]s{}
	counter, err := this.Param().SetOffset(offset).SetSize(size).SetRecv(&r).SetMiddleware(mw).List()
	return r, counter, err
}

func (this *%[3]s) Add(args ...*%[3]s) (interface{}, error) {
	var data = this
	if len(args)>0 {
		data = args[0]
	}
	return this.Param().SetSend(data).Insert()
}

func (this *%[3]s) Edit(mw func(db.Result) db.Result, args ...*%[3]s) error {
	var data = this
	if len(args)>0 {
		data = args[0]
	}
	return this.Param().SetSend(data).SetMiddleware(mw).Update()
}

func (this *%[3]s) Delete(mw func(db.Result) db.Result) error {
	return this.Param().SetMiddleware(mw).Delete()
}

`

var (
	user      *string
	pass      *string
	host      *string
	engine    *string
	database  *string
	targetDir *string
	prefix    *string
	pkgName   *string
	schema    *string
)

func main() {
	user = flag.String(`u`, `root`, `-u user`)
	pass = flag.String(`p`, ``, `-p password`)
	host = flag.String(`h`, `localhost`, `-p host`)
	engine = flag.String(`e`, `mysql`, `-e engine`)
	database = flag.String(`d`, `blog`, `-d database`)
	targetDir = flag.String(`o`, `dbschema`, `-o targetDir`)
	prefix = flag.String(`pre`, ``, `-pre prefix`)
	pkgName = flag.String(`pkg`, `dbschema`, `-pkg packageName`)
	schema = flag.String(`schema`, `public`, `-schema schemaName`)
	flag.Parse()
	var sess sqlbuilder.Database
	var err error
	switch *engine {
	case `mymysql`, `mysql`:
		fallthrough
	default:
		settings := mysql.ConnectionURL{
			Host:     *host,
			Database: *database,
			User:     *user,
			Password: *pass,
		}
		sess, err = mysql.Open(settings)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()
	tables, err := sess.Collections()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(`Found tables: %v`, tables)
	err = os.MkdirAll(*targetDir, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	allFields := map[string]map[string]factory.FieldInfo{}
	hasPrefix := len(*prefix) > 0
	for _, tableName := range tables {
		fieldMaxLength, fieldsInfo := GetTableInfo(*engine, sess, tableName)
		structName := TableToStructName(tableName, *prefix)
		imports := ``
		fieldBlock := ``
		maxLen := strconv.Itoa(fieldMaxLength / 2)
		fieldNames := map[string]factory.FieldInfo{}
		for key, field := range fieldsInfo {
			if key > 0 {
				fieldBlock += "\n"
			}
			fieldInfo := factory.FieldInfo{Options: []string{}}
			p := strings.Index(field["Type"], `(`)
			fieldInfo.DataType = DataType(field["Type"])
			if p > -1 {
				//fieldInfo.DataType = field["Type"][0:p]
				pr := strings.Index(field["Type"], `)`)
				if pr > -1 {
					opts := field["Type"][p+1 : pr]
					if len(opts) > 0 {
						if opts[0] == '\'' {
							for _, opt := range strings.Split(opts, `,`) {
								fieldInfo.Options = append(fieldInfo.Options, strings.Trim(opt, `'`))
							}
						} else if strings.Contains(opts, `,`) {
							opts := strings.Split(opts, `,`)
							switch len(opts) {
							case 2:
								fieldInfo.Min, err = strconv.Atoi(opts[0])
								if err != nil {
									panic(err)
								}
								fieldInfo.Max, err = strconv.Atoi(opts[1])
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
			fieldP := fmt.Sprintf(`%-`+maxLen+`s`, TableToStructName(field["Field"], ``))
			typeP := fmt.Sprintf(`%-8s`, fieldInfo.DataType)
			dbTag := field["Field"]
			bsonTag := field["Field"]
			if field["Key"] == "PRI" && field["Extra"] == "auto_increment" {
				dbTag += ",omitempty,pk"
				bsonTag += ",omitempty"
				fieldInfo.PrimaryKey = true
				fieldInfo.AutoIncrement = true
			} else if field["Key"] == "PRI" {
				fieldInfo.PrimaryKey = true
			}
			fieldInfo.Comment = field["Comment"]
			fieldInfo.DefaultValue = field["Default"]
			fieldBlock += "\t" + fieldP + "\t" + typeP + "\t"
			fieldBlock += "`db:\"" + dbTag + "\" bson:\"" + bsonTag + "\" comment:\"" + field["Comment"] + "\" json:\"" + field["Field"] + "\" xml:\"" + field["Field"] + "\"`"

			fieldNames[field["Field"]] = fieldInfo
		}
		noPrefixTableName := tableName
		if hasPrefix {
			noPrefixTableName = strings.TrimPrefix(tableName, *prefix)
		}
		content := fmt.Sprintf(structTemplate, *pkgName, imports, structName, fieldBlock, noPrefixTableName)

		saveAs := filepath.Join(*targetDir, structName) + `.go`
		file, err := os.Create(saveAs)
		if err == nil {
			_, err = file.WriteString(content)
		}
		if err != nil {
			log.Println(err)
		} else {
			log.Println(`Generated struct:`, structName)
		}

		allFields[noPrefixTableName] = fieldNames
	}

	content := `//Generated by webx-top/db
package %[1]s

import (
	"github.com/webx-top/db/lib/factory"
)

func init(){
	%[2]s
}

`
	dataContent := strings.Replace(fmt.Sprintf(`factory.Fields=%#v`+"\n", allFields), `map[string]factory.FieldInfo`, `map[string]*factory.FieldInfo`, -1)
	dataContent = strings.Replace(dataContent, `:factory.FieldInfo`, `:&factory.FieldInfo`, -1)
	content = fmt.Sprintf(content, *pkgName, dataContent)
	saveAs := filepath.Join(*targetDir, `init`) + `.go`
	file, err := os.Create(saveAs)
	if err == nil {
		_, err = file.WriteString(content)
	}
	if err != nil {
		log.Println(err)
	} else {
		log.Println(`Generated init.go`)
	}

	log.Println(`End.`)
}

func TableToStructName(tableName string, prefix string) string {
	if len(prefix) > 0 {
		tableName = strings.TrimPrefix(tableName, prefix)
	}
	tableName = strings.Title(tableName)
	return camleCase(tableName)
}

func DataType(dbDataType string) string {
	switch {
	case strings.HasPrefix(dbDataType, `int`):
		return `int`
	case strings.HasPrefix(dbDataType, `bigint`):
		return `int64`
	case strings.HasPrefix(dbDataType, `decimal`):
		return `float64`
	case strings.HasPrefix(dbDataType, `float`):
		return `float32`
	case strings.HasPrefix(dbDataType, `double`):
		return `float64`

	//postgreSQL
	case strings.HasPrefix(dbDataType, `boolean`):
		return `bool`
	case strings.HasPrefix(dbDataType, `oid`):
		return `int64`

	default:
		return `string`
	}
}

func camleCase(s string) string {
	vs := []rune(s)
	underline := rune('_')
	isUnderline := false
	vals := []rune{}
	for _, v := range vs {
		if v == underline {
			isUnderline = true
			continue
		}
		if isUnderline {
			v = unicode.ToUpper(v)
		}
		isUnderline = false
		vals = append(vals, v)
	}
	return string(vals)
}

func GetTableInfo(engine string, d sqlbuilder.Database, tableName string) (int, []map[string]string) {
	switch engine {
	case "mymysql", "mysql":
		fallthrough
	default:
		return getMySQLTableInfo(d, tableName)
	}
}
