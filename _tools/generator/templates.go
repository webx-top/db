package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type structField struct {
	field   string
	typ     string
	dbTag   string
	bsonTag string
	jsonTag string
	xmlTag  string
	comment string
}

func (f *structField) String() string {
	return fmt.Sprintf(memberTemplate, f.field, f.typ, f.dbTag, f.bsonTag, f.comment, f.jsonTag, f.xmlTag)
}

var structFuncs = map[string]string{
	`Trans`:         `Trans`,
	`Use`:           `Use`,
	`Objects`:       `Objects`,
	`NewObjects`:    `NewObjects`,
	`NewParam`:      `NewParam`,
	`SetParam`:      `SetParam`,
	`Get`:           `Get`,
	`List`:          `List`,
	`ListByOffset`:  `ListByOffset`,
	`Add`:           `Add`,
	`Edit`:          `Edit`,
	`Upsert`:        `Upsert`,
	`Delete`:        `Delete`,
	`Count`:         `Count`,
	`GroupBy`:       `GroupBy`,
	`KeyBy`:         `KeyBy`,
	`Setter`:        `Setter`,
	`SetField`:      `SetField`,
	`SetFields`:     `SetFields`,
	`Reset`:         `Reset`,
	`AsMap`:         `AsMap`,
	`AsRow`:         `AsRow`,
	`BatchValidate`: `BatchValidate`,
	`Validate`:      `Validate`,
}

type tempateDbschemaData struct {
	PackageName          string
	Imports              []string
	StructName           string
	StructComment        string
	StructAttributes     []string
	TableName            string
	BeforeInsert         string
	AfterInsert          string
	BeforeUpdate         string
	SetUpdatedAt         string
	BeforeDelete         string
	Resets               echo.KVList // v.K: struct field; v.V: zero value
	TableAndStructFields echo.KVList // v.K: table field;  v.V: struct field
}

type tempateDbschemaInitData struct {
	PackageName  string
	Prefix       string
	InitCode     string
	DBKey        string
	DefaultDBKey string
}

type tempateModelData struct {
	PackageName       string
	Imports           []string
	StructName        string
	BaseName          string
	SchemaPackagePath string
	SchemaPackageName string
}

type tempateModelBaseData struct {
	PackageName       string
	Imports           []string
	StructName        string
	SchemaPackagePath string
	SchemaPackageName string
}

var memberTemplate = "\t%v\t%v\t`db:%q bson:%q comment:%q json:%q xml:%q`"

var (
	structTemplateObj    *template.Template
	modelBaseTemplateObj *template.Template
	modelTemplateObj     *template.Template
	initTemplateObj      *template.Template
	once                 sync.Once
)

func getTemplate(name string) string {
	if len(cfg.TemplateDir) == 0 {
		cfg.TemplateDir = filepath.Join(os.Getenv("GOPATH"), `src/github.com/webx-top/db/_tools/generator/template`)
	}
	switch name {
	case `dbschema`:
		return filepath.Join(cfg.TemplateDir, `dbschema.gotpl`)
	case `dbschema_init`:
		return filepath.Join(cfg.TemplateDir, `dbschema_init.gotpl`)
	case `model`:
		return filepath.Join(cfg.TemplateDir, `model.gotpl`)
	case `model_base`:
		return filepath.Join(cfg.TemplateDir, `model_base.gotpl`)
	default:
		panic("unknown template name: " + name)
	}
}

func getTemplateContent(name string) (string, error) {
	tempfile := getTemplate(name)
	var b []byte
	var err error
	if !com.FileExists(tempfile) {
		var fs http.File
		fs, err = AssetFile().Open(`template/` + name + `.gotpl`)
		if err != nil {
			return "", err
		}
		defer fs.Close()
		b, err = ioutil.ReadAll(fs)
	} else {
		b, err = ioutil.ReadFile(tempfile)
	}
	if err != nil {
		return ``, err
	}
	return string(b), nil
}

func initTemplate() {
	tpl, err := getTemplateContent(`dbschema`)
	if err != nil {
		panic(err)
	}
	structTemplateObj, err = template.New(`dbschema`).Parse(tpl)
	if err != nil {
		panic(err)
	}

	tpl, err = getTemplateContent(`dbschema_init`)
	if err != nil {
		panic(err)
	}
	initTemplateObj, err = template.New(`dbschema_init`).Parse(tpl)
	if err != nil {
		panic(err)
	}

	tpl, err = getTemplateContent(`model`)
	if err != nil {
		panic(err)
	}
	modelTemplateObj, err = template.New(`model`).Parse(tpl)
	if err != nil {
		panic(err)
	}

	tpl, err = getTemplateContent(`model_base`)
	if err != nil {
		panic(err)
	}
	modelBaseTemplateObj, err = template.New(`model_base`).Parse(tpl)
	if err != nil {
		panic(err)
	}
}

func Template(name string, data interface{}) ([]byte, error) {
	once.Do(initTemplate)
	w := bytes.NewBuffer(nil)
	var err error
	switch name {
	case `dbschema`:
		err = structTemplateObj.ExecuteTemplate(w, name, data)
	case `dbschema_init`:
		err = initTemplateObj.ExecuteTemplate(w, name, data)
	case `model`:
		err = modelTemplateObj.ExecuteTemplate(w, name, data)
	case `model_base`:
		err = modelBaseTemplateObj.ExecuteTemplate(w, name, data)
	default:
		panic("unknown template name: " + name)
	}
	return w.Bytes(), err
}
