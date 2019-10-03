package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/admpub/confl"
	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/db/mysql"
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
	modelInstancers := map[string]string{}
	hasPrefix := len(cfg.Prefix) > 0
	hasIngore := len(cfg.Ignore) > 0
	hasMatch := len(cfg.Match) > 0
	validTables := []string{}
	columns := map[string][]string{}
	for _, tableName := range tables {
		if hasIngore && regexp.MustCompile(cfg.Ignore).MatchString(tableName) {
			log.Println(`Ignore the table:`, tableName)
			continue
		}
		if hasMatch && !regexp.MustCompile(cfg.Match).MatchString(tableName) {
			log.Println(`Ignore the table:`, tableName)
			continue
		}
		validTables = append(validTables, tableName)
		if cfg.NotGenerated {
			continue
		}
		structName := TableToStructName(tableName, cfg.Prefix)
		structComment, err := GetTableComment(cfg.Engine, sess, tableName)
		if err != nil {
			panic(err)
		}
		modelInstancers[structName] = `factory.NewMI("` + tableName + `",func(connID int) factory.Model { return &` + structName + `{connID: connID} },"` + com.AddSlashes(structComment, '"') + `")`
		var imports string
		goFields, fields, fieldNames := GetTableFields(cfg.Engine, sess, tableName)
		fieldBlock := strings.Join(goFields, "\n")
		noPrefixTableName := tableName
		if hasPrefix {
			noPrefixTableName = strings.TrimPrefix(tableName, cfg.Prefix)
		}
		columns[noPrefixTableName] = fieldNames
		var resets, asMap, asRow, setCase, fromMapCase string
		for key, fieldName := range fieldNames {
			f := fields[fieldName]
			if key > 0 {
				resets += "\n"
				asMap += "\n"
				asRow += "\n"
				setCase += "\n"
				fromMapCase += "\n"
			}
			resets += "	this." + f.GoName + " = " + ZeroValue(f.GoType)
			asMap += `	r["` + f.GoName + `"] = this.` + f.GoName
			asRow += `	r["` + f.Name + `"] = this.` + f.GoName
			goTypeName := f.GoType
			if goTypeName == `byte[]` {
				goTypeName = `bytes`
			}
			setCase += `				case "` + f.GoName + `": this.` + f.GoName + ` = param.As` + strings.Title(goTypeName) + `(vv)`
			fromMapCase += `				case "` + f.Name + `": this.` + f.GoName + ` = param.As` + strings.Title(goTypeName) + `(value)`
		}
		replaceMap := *replaces
		replaceMap["packageName"] = cfg.SchemaConfig.PackageName
		replaceMap["structName"] = structName
		replaceMap["structComment"] = structComment
		replaceMap["attributes"] = fieldBlock
		replaceMap["reset"] = resets
		replaceMap["asMap"] = asMap
		replaceMap["asRow"] = asRow
		replaceMap["fromMapCase"] = fromMapCase
		replaceMap["setCase"] = setCase
		replaceMap["tableName"] = noPrefixTableName
		replaceMap["beforeInsert"] = ""
		replaceMap["beforeUpdate"] = ""
		replaceMap["setUpdatedAt"] = ""
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
				setUpdatedAt := ``
				newLine := ``
				for _, _fieldName := range _fieldNames {
					fieldInf, ok := fields[_fieldName]
					if !ok {
						continue
					}
					switch fieldInf.GoType {
					case `uint`, `int`, `uint32`, `int32`, `int64`, `uint64`:
						beforeUpdate += newLine + `this.` + fieldInf.GoName + ` = ` + fieldInf.GoType + `(time.Now().Unix())`
						setUpdatedAt += newLine + `kvset["` + _fieldName + `"] = ` + fieldInf.GoType + `(time.Now().Unix())`
						newLine = "\n\t"
						importTime = true
					case `string`:
						//TODO
					}
				}
				replaceMap["beforeUpdate"] = beforeUpdate
				replaceMap["setUpdatedAt"] = setUpdatedAt
			}
		}
		if importTime {
			imports += "\n\t" + `"time"`
		}
		beforeInsert := ``
		beforeUpdate := ``
		setUpdatedAt := ``
		newLine := ``
		for _, _fieldName := range fieldNames {
			fieldInf := fields[_fieldName]
			switch fieldInf.GoType {
			case `string`:
				if len(fieldInf.DefaultValue) > 0 {
					setDefault := newLine + `if len(this.` + fieldInf.GoName + `) == 0 { this.` + fieldInf.GoName + ` = "` + fieldInf.DefaultValue + `" }`
					beforeUpdate += setDefault
					beforeInsert += setDefault
					setUpdatedAt += newLine + `if val, ok := kvset["` + _fieldName + `"]; ok && val != nil { if v, ok := val.(string); ok && len(v) == 0 { kvset["` + _fieldName + `"] = "` + fieldInf.DefaultValue + `" } }`
					newLine = "\n\t"
				}
			}
		}
		if len(beforeInsert) > 0 {
			replaceMap["beforeInsert"] += newLine + beforeInsert
		}
		if len(beforeUpdate) > 0 {
			replaceMap["beforeUpdate"] += newLine + beforeUpdate
		}
		if len(setUpdatedAt) > 0 {
			replaceMap["setUpdatedAt"] = newLine + setUpdatedAt
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
	if cfg.NotGenerated {
		execBackupCommand(cfg, validTables)
		log.Println(`End.`)
		return
	}
	content := initFileTemplate
	content = strings.Replace(content, `{{prefix}}`, cfg.Prefix, -1)
	content = strings.Replace(content, `{{dbKey}}`, cfg.DBKey, -1)
	dataContent := strings.Replace(fmt.Sprintf(`DBI.Fields.Register(%#v)`+"\n", allFields), `map[string]map[string]factory.FieldInfo`, `map[string]map[string]*factory.FieldInfo`, -1)
	dataContent = strings.Replace(dataContent, `map[string]factory.FieldInfo`, ``, -1)
	dataContent = strings.Replace(dataContent, `:factory.FieldInfo`, `:`, -1)
	dataContent += "\n\t" + fmt.Sprintf(`DBI.Columns=%#v`, columns) + "\n"
	dataContent += "\n\tDBI.Models.Register(factory.ModelInstancers{"
	for structName, modelInstancer := range modelInstancers {
		dataContent += "`" + structName + "`:" + modelInstancer + `,`
	}
	dataContent += "})\n"
	content = strings.Replace(content, `{{packageName}}`, cfg.SchemaConfig.PackageName, -1)
	if cfg.DBKey != factory.DefaultDBKey {
		dataContent = `factory.DBIRegister(DBI,"` + cfg.DBKey + `")` + "\n\t" + dataContent
	} else {
		content = strings.Replace(content, `factory.NewDBI()`, `factory.DefaultDBI`, -1)
	}
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

	execBackupCommand(cfg, validTables)

	log.Println(`End.`)
}
