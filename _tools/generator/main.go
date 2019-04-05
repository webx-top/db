package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/confl"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/db/mysql"
	//"github.com/webx-top/com"
)

func main() {
	parseFlag()
	var err error
	if len(configFile) > 0 {
		_, err = confl.DecodeFile(configFile, cfg)
		if err != nil {
			log.Fatal(err)
		}
	}
	//com.Dump(cfg)
	var sess sqlbuilder.Database
	switch cfg.Engine {
	case `mymysql`, `mysql`:
		fallthrough
	default:
		settings := mysql.ConnectionURL{
			Host:     cfg.Host,
			Database: cfg.Database,
			User:     cfg.Username,
			Password: cfg.Password,
		}
		sess, err = mysql.Open(settings)
	}
	if err != nil {
		log.Fatal(err)
	}
	cfg.Check()
	defer sess.Close()
	tables, err := sess.Collections()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(`Found tables: %v`, tables)
	err = os.MkdirAll(cfg.SchemaConfig.SaveDir, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	allFields := map[string]map[string]factory.FieldInfo{}
	hasPrefix := len(cfg.Prefix) > 0
	for _, tableName := range tables {
		structName := TableToStructName(tableName, cfg.Prefix)
		imports := ``
		goFields, fields, fieldNames := GetTableFields(cfg.Engine, sess, tableName)
		fieldBlock := strings.Join(goFields, "\n")
		noPrefixTableName := tableName
		if hasPrefix {
			noPrefixTableName = strings.TrimPrefix(tableName, cfg.Prefix)
		}
		var resets, asMap string
		for key, fieldName := range fieldNames {
			f := fields[fieldName]
			if key > 0 {
				resets += "\n"
				asMap += "\n"
			}
			resets += "	this." + f.GoName + " = " + ZeroValue(f.GoType)
			asMap += `	r["` + f.GoName + `"] = this.` + f.GoName
		}
		replaceMap := *replaces
		replaceMap["packageName"] = cfg.SchemaConfig.PackageName
		replaceMap["structName"] = structName
		replaceMap["attributes"] = fieldBlock
		replaceMap["reset"] = resets
		replaceMap["asMap"] = asMap
		replaceMap["tableName"] = noPrefixTableName
		replaceMap["beforeInsert"] = ""
		replaceMap["beforeUpdate"] = ""
		replaceMap["beforeDelete"] = ""
		replaceMap["afterInsert"] = ""

		importTime := false
		if cfg.AutoTimeFields != nil {
			_fieldNames, ok := cfg.AutoTimeFields.Insert[`*`]
			if !ok {
				_fieldNames, ok = cfg.AutoTimeFields.Insert[tableName]
			}
			if ok && len(_fieldNames) > 0 {
				beforeInsert := ``
				newLine := ``
				for _, _fieldName := range _fieldNames {
					fieldInf, ok := fields[_fieldName]
					if !ok {
						continue
					}
					switch fieldInf.GoType {
					case `uint`, `int`, `uint32`, `int32`, `int64`, `uint64`:
						beforeInsert += newLine + `this.` + fieldInf.GoName + ` = ` + fieldInf.GoType + `(time.Now().Unix())`
						newLine = "\n\t"
						importTime = true
					case `string`:
						//TODO
					}
				}
				afterInsert := ``
				newLine2 := ``
				newTab2 := ``
				for _, fieldInf := range fields {
					if fieldInf.AutoIncrement && fieldInf.PrimaryKey {
						beforeInsert += newLine + `this.` + fieldInf.GoName + ` = 0`
						newLine = "\n\t"
						var extData string
						if fieldInf.GoType != `int64` {
							extData = ` else if v, y := pk.(int64); y {
` + newTab2 + `			this.` + fieldInf.GoName + ` = ` + fieldInf.GoType + `(v)
` + newTab2 + `		}`
						}
						afterInsert += newLine2 + `if err == nil && pk != nil {
` + newTab2 + `		if v, y := pk.(` + fieldInf.GoType + `); y {
` + newTab2 + `			this.` + fieldInf.GoName + ` = v
` + newTab2 + `		}` + extData + `
` + newTab2 + `	}`
						newLine2 = "\n\t"
						newTab2 = "\t"
						break
					}
				}
				replaceMap["afterInsert"] = afterInsert
				replaceMap["beforeInsert"] = beforeInsert
			}
			_fieldNames, ok = cfg.AutoTimeFields.Update[`*`]
			if !ok {
				_fieldNames, ok = cfg.AutoTimeFields.Update[tableName]
			}
			if ok && len(_fieldNames) > 0 {
				beforeUpdate := ``
				newLine := ``
				for _, _fieldName := range _fieldNames {
					fieldInf, ok := fields[_fieldName]
					if !ok {
						continue
					}
					switch fieldInf.GoType {
					case `uint`, `int`, `uint32`, `int32`, `int64`, `uint64`:
						beforeUpdate += newLine + `this.` + fieldInf.GoName + ` = ` + fieldInf.GoType + `(time.Now().Unix())`
						newLine = "\n\t"
						importTime = true
					case `string`:
						//TODO
					}
				}
				replaceMap["beforeUpdate"] = beforeUpdate
			}
		}
		if importTime {
			imports += "\n\t" + `"time"`
		}

		replaceMap["imports"] = imports

		content := structTemplate
		for tag, val := range replaceMap {
			content = strings.Replace(content, `{{`+tag+`}}`, val, -1)
		}

		saveAs := filepath.Join(cfg.SchemaConfig.SaveDir, structName) + `.go`
		file, err := os.Create(saveAs)
		if err == nil {
			_, err = file.WriteString(content)
		}
		if err != nil {
			log.Println(err)
		} else {
			log.Println(`Generated schema struct:`, structName)
		}

		if len(cfg.ModelConfig.PackageName) > 0 && len(cfg.ModelConfig.SaveDir) > 0 {
			os.MkdirAll(cfg.ModelConfig.SaveDir, 0777)
			modelFile := filepath.Join(cfg.ModelConfig.SaveDir, structName) + `.go`
			_, err := os.Stat(modelFile)
			if err != nil && os.IsNotExist(err) {
				file, err := os.Create(modelFile)
				if err == nil {
					mr := *modelReplaces
					baseName := `Base`
					if len(cfg.ModelConfig.BaseName) > 0 {
						baseName = cfg.ModelConfig.BaseName
					}
					mr["packageName"] = cfg.ModelConfig.PackageName
					mr["imports"] = ""
					mr["structName"] = structName
					mr["baseName"] = baseName
					mr["schemaPackagePath"] = cfg.SchemaConfig.ImportPath
					mr["schemaPackageName"] = cfg.SchemaConfig.PackageName
					content := modelTemplate
					for tag, val := range mr {
						content = strings.Replace(content, `{{`+tag+`}}`, val, -1)
					}
					_, err = file.WriteString(content)
				}
				if err != nil {
					log.Println(err)
				} else {
					log.Println(`Generated model struct:`, structName)
				}
			}
		}

		allFields[noPrefixTableName] = fields
	}

	content := initFileTemplate
	dataContent := strings.Replace(fmt.Sprintf(`factory.Fields=%#v`+"\n", allFields), `map[string]map[string]factory.FieldInfo`, `map[string]map[string]*factory.FieldInfo`, -1)
	dataContent = strings.Replace(dataContent, `map[string]factory.FieldInfo`, ``, -1)
	dataContent = strings.Replace(dataContent, `:factory.FieldInfo`, `:`, -1)
	content = strings.Replace(content, `{{packageName}}`, cfg.SchemaConfig.PackageName, -1)
	content = strings.Replace(content, `{{initCode}}`, dataContent, -1)
	saveAs := filepath.Join(cfg.SchemaConfig.SaveDir, `init`) + `.go`
	file, err := os.Create(saveAs)
	if err == nil {
		_, err = file.WriteString(content)
	}
	if err != nil {
		log.Println(err)
	} else {
		log.Println(`Generated init.go`)
	}
	if len(cfg.ModelConfig.PackageName) > 0 && len(cfg.ModelConfig.SaveDir) > 0 {
		structName := `Base`
		if len(cfg.ModelConfig.BaseName) > 0 {
			structName = cfg.ModelConfig.BaseName
		}
		os.MkdirAll(cfg.ModelConfig.SaveDir, 0777)
		modelFile := filepath.Join(cfg.ModelConfig.SaveDir, structName) + `.go`
		_, err := os.Stat(modelFile)
		if err != nil && os.IsNotExist(err) {
			file, err := os.Create(modelFile)
			if err == nil {
				mr := *modelReplaces
				mr["packageName"] = cfg.ModelConfig.PackageName
				mr["imports"] = ""
				mr["structName"] = structName
				mr["schemaPackagePath"] = cfg.SchemaConfig.ImportPath
				mr["schemaPackageName"] = cfg.SchemaConfig.PackageName
				content := modelBaseTemplate
				for tag, val := range mr {
					content = strings.Replace(content, `{{`+tag+`}}`, val, -1)
				}
				_, err = file.WriteString(content)
			}
			if err != nil {
				log.Println(err)
			} else {
				log.Println(`Generated model struct:`, structName)
			}
		}
	}

	log.Println(`End.`)
}
