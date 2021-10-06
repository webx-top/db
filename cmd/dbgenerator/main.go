package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/admpub/confl"
	"github.com/admpub/regexp2"

	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/db/mysql"
	"github.com/webx-top/echo"
)

//go:generate go get github.com/admpub/bindata/v3/...
//go:generate go-bindata -fs -o bindata_assetfs.go template/...

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
	hashids := cfg.FieldHashID()
	defer sess.Close()
	tables, err := sess.Collections()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(`Found tables: %v`, tables)
	if !cfg.NotGenerated {
		err = os.MkdirAll(cfg.SchemaConfig.SaveDir, os.ModePerm)
		if err != nil {
			panic(err.Error())
		}
	}
	allFields := map[string]map[string]factory.FieldInfo{}
	modelInstancers := map[string]string{}
	hasPrefix := len(cfg.Prefix) > 0
	hasIngore := len(cfg.Ignore) > 0
	hasMatch := len(cfg.Match) > 0
	validTables := []string{}
	columns := map[string][]string{}

	for _, tableName := range tables {
		if hasIngore {
			matched, err := regexp2.MustCompile(cfg.Ignore, 0).MatchString(tableName)
			if err != nil {
				panic(err)
			}
			if matched {
				log.Println(`Ignore the table:`, tableName)
				continue
			}
		}
		if hasMatch {
			matched, err := regexp2.MustCompile(cfg.Match, 0).MatchString(tableName)
			if err != nil {
				panic(err)
			}
			if !matched {
				log.Println(`Ignore the table:`, tableName)
				continue
			}
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
		modelInstancers[structName] = `factory.NewMI("` + tableName + `",func(connID int) factory.Model { return &` + structName + `{base:*((&factory.Base{}).SetConnID(connID))} },"` + com.AddSlashes(structComment, '"') + `")`
		var typeFields []string
		if idf, ok := hashids[tableName]; ok {
			typeFields = append(typeFields, idf)
		}
		goFields, fields, fieldNames := GetTableFields(cfg.Engine, sess, tableName, map[string][]string{`hashids`: typeFields})
		noPrefixTableName := tableName
		if hasPrefix {
			noPrefixTableName = strings.TrimPrefix(tableName, cfg.Prefix)
		}
		columns[noPrefixTableName] = fieldNames
		hasHashids := false
		// template data
		tplSchema := &tempateDbschemaData{}
		for _, fieldName := range fieldNames {
			if !hasHashids && com.InSlice(fieldName, typeFields) {
				hasHashids = true
			}
			f := fields[fieldName]
			goTypeName := f.GoType
			if goTypeName == `[]byte` {
				goTypeName = `Bytes`
			} else {
				goTypeName = strings.Title(goTypeName)
			}
			var extPrefix string
			var extSuffix string
			if len(f.MyType) > 0 {
				extPrefix = f.MyType + `(`
				extSuffix = `)`
			}
			tplSchema.Resets.AddItem(echo.NewKV(f.GoName, ZeroValue(f.GoType)))
			tplSchema.TableAndStructFields.AddItem(echo.NewKV(f.Name, f.GoName).SetHKV(`dataType`, goTypeName).SetHKV(`convertStart`, extPrefix).SetHKV(`convertEnd`, extSuffix))
		}
		tplSchema.PackageName = cfg.SchemaConfig.PackageName
		tplSchema.StructName = structName
		tplSchema.StructComment = structComment
		tplSchema.StructAttributes = goFields
		tplSchema.TableName = noPrefixTableName

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
					convt := fieldInf.GoType
					if len(fieldInf.MyType) > 0 {
						convt = fieldInf.MyType
					}
					switch fieldInf.GoType {
					case `uint`, `int`, `uint32`, `int32`, `int64`, `uint64`:
						beforeInsert += newLine + `a.` + fieldInf.GoName + ` = ` + convt + `(time.Now().Unix())`
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
						beforeInsert += newLine + `a.` + fieldInf.GoName + ` = 0`
						newLine = "\n\t"
						var extData string
						var extPrefix string
						var extSuffix string
						if len(fieldInf.MyType) > 0 {
							extPrefix = fieldInf.MyType + `(`
							extSuffix = `)`
						}
						if fieldInf.GoType != `int64` {
							extData = ` else if v, y := pk.(int64); y {
` + newTab2 + `			a.` + fieldInf.GoName + ` = ` + extPrefix + fieldInf.GoType + `(v)` + extSuffix + `
` + newTab2 + `		}`
						}
						afterInsert += newLine2 + `if err == nil && pk != nil {
` + newTab2 + `		if v, y := pk.(` + fieldInf.GoType + `); y {
` + newTab2 + `			a.` + fieldInf.GoName + ` = ` + extPrefix + `v` + extSuffix + `
` + newTab2 + `		}` + extData + `
` + newTab2 + `	}`
						newLine2 = "\n\t"
						newTab2 = "\t"
						break
					}
				}
				tplSchema.AfterInsert = afterInsert
				tplSchema.BeforeInsert = beforeInsert
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
					convt := fieldInf.GoType
					if len(fieldInf.MyType) > 0 {
						convt = fieldInf.MyType
					}
					switch fieldInf.GoType {
					case `uint`, `int`, `uint32`, `int32`, `int64`, `uint64`:
						beforeUpdate += newLine + `a.` + fieldInf.GoName + ` = ` + convt + `(time.Now().Unix())`
						setUpdatedAt += newLine + `kvset["` + _fieldName + `"] = ` + convt + `(time.Now().Unix())`
						newLine = "\n\t"
						importTime = true
					case `string`:
						//TODO
					}
				}
				tplSchema.BeforeUpdate = beforeUpdate
				tplSchema.SetUpdatedAt = setUpdatedAt
			}
		}
		if importTime {
			tplSchema.Imports = append(tplSchema.Imports, `time`)
		}
		if hasHashids {
			tplSchema.Imports = append(tplSchema.Imports, `github.com/admpub/hashseq`)
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
					setDefault := newLine + `if len(a.` + fieldInf.GoName + `) == 0 { a.` + fieldInf.GoName + ` = "` + fieldInf.DefaultValue + `" }`
					beforeUpdate += setDefault
					beforeInsert += setDefault
					setUpdatedAt += newLine + `if val, ok := kvset["` + _fieldName + `"]; ok && val != nil { if v, ok := val.(string); ok && len(v) == 0 { kvset["` + _fieldName + `"] = "` + fieldInf.DefaultValue + `" } }`
					newLine = "\n\t"
				}
			}
		}
		if len(beforeInsert) > 0 {
			tplSchema.BeforeInsert += newLine + beforeInsert
		}
		if len(beforeUpdate) > 0 {
			tplSchema.BeforeUpdate += newLine + beforeUpdate
		}
		if len(setUpdatedAt) > 0 {
			tplSchema.SetUpdatedAt = newLine + setUpdatedAt
		}
		content, err := Template(`dbschema`, tplSchema)
		if err != nil {
			panic(err)
		}
		saveAs := filepath.Join(cfg.SchemaConfig.SaveDir, structName) + `.go`
		file, err := os.Create(saveAs)
		if err == nil {
			_, err = file.Write(content)
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
				var file *os.File
				file, err = os.Create(modelFile)
				if err == nil {
					tplModel := &tempateModelData{
						BaseName: `Base`,
					}
					if len(cfg.ModelConfig.BaseName) > 0 {
						tplModel.BaseName = cfg.ModelConfig.BaseName
					}
					tplModel.PackageName = cfg.ModelConfig.PackageName
					tplModel.StructName = structName
					tplModel.SchemaPackagePath = cfg.SchemaConfig.ImportPath
					tplModel.SchemaPackageName = cfg.SchemaConfig.PackageName
					var content []byte
					content, err = Template(`model`, tplModel)
					if err != nil {
						panic(err)
					}
					_, err = file.Write(content)
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
	tplSchemaInit := &tempateDbschemaInitData{
		DefaultDBKey: factory.DefaultDBKey,
	}
	tplSchemaInit.Prefix = cfg.Prefix
	tplSchemaInit.DBKey = cfg.DBKey
	tplSchemaInit.PackageName = cfg.SchemaConfig.PackageName
	dataContent := strings.Replace(fmt.Sprintf(`DBI.Fields.Register(%#v)`+"\n", allFields), `map[string]map[string]factory.FieldInfo`, `map[string]map[string]*factory.FieldInfo`, -1)
	dataContent = strings.Replace(dataContent, `map[string]factory.FieldInfo`, ``, -1)
	dataContent = strings.Replace(dataContent, `:factory.FieldInfo`, `:`, -1)
	dataContent += "\n\t" + fmt.Sprintf(`DBI.Columns=%#v`, columns) + "\n"
	dataContent += "\n\tDBI.Models.Register(factory.ModelInstancers{"
	var structNames []string
	for structName := range modelInstancers {
		structNames = append(structNames, structName)
	}
	sort.Strings(structNames)
	for _, structName := range structNames {
		modelInstancer := modelInstancers[structName]
		dataContent += "`" + structName + "`:" + modelInstancer + `,`
	}
	dataContent += "})\n"
	tplSchemaInit.InitCode = dataContent
	content, err := Template(`dbschema_init`, tplSchemaInit)
	if err != nil {
		panic(err)
	}
	saveAs := filepath.Join(cfg.SchemaConfig.SaveDir, `init`) + `.go`
	file, err := os.Create(saveAs)
	if err == nil {
		_, err = file.Write(content)
	}
	if err != nil {
		log.Println(err)
	} else {
		log.Println(`Generated init.go`)
	}
	if len(cfg.ModelConfig.PackageName) > 0 && len(cfg.ModelConfig.SaveDir) > 0 {
		tplModelBase := &tempateModelBaseData{
			StructName: `Base`,
		}
		if len(cfg.ModelConfig.BaseName) > 0 {
			tplModelBase.StructName = cfg.ModelConfig.BaseName
		}
		os.MkdirAll(cfg.ModelConfig.SaveDir, 0777)
		modelFile := filepath.Join(cfg.ModelConfig.SaveDir, tplModelBase.StructName) + `.go`
		_, err := os.Stat(modelFile)
		if err != nil && os.IsNotExist(err) {
			file, err := os.Create(modelFile)
			if err == nil {
				tplModelBase.PackageName = cfg.ModelConfig.PackageName
				tplModelBase.SchemaPackagePath = cfg.SchemaConfig.ImportPath
				tplModelBase.SchemaPackageName = cfg.SchemaConfig.PackageName
				content, err = Template(`model_base`, tplModelBase)
				if err != nil {
					panic(err)
				}
				_, err = file.Write(content)
			}
			if err != nil {
				log.Println(err)
			} else {
				log.Println(`Generated model struct:`, tplModelBase.StructName)
			}
		}
	}
	Format(cfg.SchemaConfig.SaveDir)
	if cfg.SchemaConfig.SaveDir != cfg.ModelConfig.SaveDir {
		Format(cfg.ModelConfig.SaveDir)
	}
	execBackupCommand(cfg, validTables)

	log.Println(`End.`)
}
