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

	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/db/mysql"
)

var structTemplate = `package %s
%s
type %s struct {
%s
}
`

func main() {
	user := flag.String(`u`, `root`, `-u user`)
	pass := flag.String(`p`, `root`, `-p password`)
	host := flag.String(`h`, `localhost`, `-p host`)
	engine := flag.String(`e`, `mysql`, `-e engine`)
	database := flag.String(`d`, `blog`, `-d database`)
	targetDir := flag.String(`o`, `model`, `-o targetDir`)
	prefix := flag.String(`pre`, `webx_`, `-pre prefix`)
	flag.Parse()
	var sess sqlbuilder.Database
	var err error
	switch *engine {
	case `mysql`:
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
	rows, err := sess.Query(`SHOW TABLES`)
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
	log.Printf(`Found tables: %v`, tables)
	err = os.MkdirAll(*targetDir, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	for _, tableName := range tables {
		rows, err := sess.Query("SHOW COLUMNS FROM `" + tableName + "`")
		if err != nil {
			log.Println(err)
			continue
		}
		fields, err := rows.Columns()
		if err != nil {
			log.Println(err)
			continue
		}
		var result = make([]interface{}, len(fields))
		for k, _ := range fields {
			v := ``
			result[k] = &v
		}
		fieldsInfo := []map[string]string{}
		fieldMaxLength := 0
		for rows.Next() {

			remap := map[string]string{}
			err = rows.Scan(result...)
			if err != nil {
				//log.Println(err)
			}

			for k, v := range result {
				remap[fields[k]] = *(v.(*string))
				sz := len(remap[fields[k]])
				if sz > fieldMaxLength {
					fieldMaxLength = sz
				}
			}
			fieldsInfo = append(fieldsInfo, remap)
			//log.Printf(`%#v`+"\n", remap)
		}
		structName := TableToStructName(tableName, *prefix)
		pkgName := `model`
		imports := ``
		fieldBlock := ``
		maxLen := strconv.Itoa(fieldMaxLength / 2)
		for key, field := range fieldsInfo {
			if key > 0 {
				fieldBlock += "\n"
			}
			fieldP := fmt.Sprintf(`%-`+maxLen+`s`, TableToStructName(field["Field"], ``))
			typeP := fmt.Sprintf(`%-8s`, DataType(field["Type"]))
			fieldBlock += "\t" + fieldP + "\t" + typeP + "\t`db:\"" + field["Field"] + "\"`"
		}
		content := fmt.Sprintf(structTemplate, pkgName, imports, structName, fieldBlock)

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
	case strings.HasPrefix(dbDataType, `bigint`):
		return `int64`
	case strings.HasPrefix(dbDataType, `decimal`):
		return `float64`
	case strings.HasPrefix(dbDataType, `float`):
		return `float64`
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
