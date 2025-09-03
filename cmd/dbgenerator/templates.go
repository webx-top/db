package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type structField struct {
	field    string
	typ      string
	dbTag    string
	bsonTag  string
	jsonTag  string
	xmlTag   string
	otherTag string
	comment  string
}

func (f *structField) String() string {
	otherTag := f.otherTag
	if len(otherTag) > 0 {
		otherTag = ` ` + otherTag
	}
	return fmt.Sprintf(memberTemplate, f.field, f.typ, f.dbTag, f.bsonTag, f.comment, f.jsonTag, f.xmlTag, otherTag)
}

var structFuncs = map[string]string{
	`Trans`:              `Trans`,
	`Use`:                `Use`,
	`SetContext`:         `SetContext`,
	`EventON`:            `EventON`,
	`EventOFF`:           `EventOFF`,
	`Context`:            `Context`,
	`SetConnID`:          `SetConnID`,
	`ConnID`:             `ConnID`,
	`SetNamer`:           `SetNamer`,
	`Namer`:              `Namer`,
	`Param`:              `Param`,
	`New`:                `New`,
	`AsKV`:               `AsKV`,
	`XObjects`:           `XObjects`,
	`Objects`:            `Objects`,
	`NewObjects`:         `NewObjects`,
	`InitObjects`:        `InitObjects`,
	`NewParam`:           `NewParam`,
	`SetParam`:           `SetParam`,
	`Get`:                `Get`,
	`Set`:                `Set`,
	`List`:               `List`,
	`ListByOffset`:       `ListByOffset`,
	`Insert`:             `Insert`,
	`Update`:             `Update`,
	`Updatex`:            `Updatex`,
	`UpdateByFields`:     `UpdateByFields`,
	`UpdatexByFields`:    `UpdatexByFields`,
	`Upsert`:             `Upsert`,
	`Delete`:             `Delete`,
	`Deletex`:            `Deletex`,
	`Count`:              `Count`,
	`Exists`:             `Exists`,
	`GroupBy`:            `GroupBy`,
	`KeyBy`:              `KeyBy`,
	`Setter`:             `Setter`,
	`UpdateField`:        `UpdateField`,
	`UpdateFields`:       `UpdateFields`,
	`UpdatexField`:       `UpdateField`,
	`UpdatexFields`:      `UpdateFields`,
	`UpdateValues`:       `UpdateValues`,
	`Reset`:              `Reset`,
	`AsMap`:              `AsMap`,
	`AsRow`:              `AsRow`,
	`FromRow`:            `FromRow`,
	`GetField`:           `GetField`,
	`GetAllFieldNames`:   `GetAllFieldNames`,
	`HasField`:           `HasField`,
	`ListPage`:           `ListPage`,
	`ListPageAs`:         `ListPageAs`,
	`ListPageByOffset`:   `ListPage`,
	`ListPageByOffsetAs`: `ListPageAs`,
	`BatchValidate`:      `BatchValidate`,
	`Validate`:           `Validate`,
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
	SchemaPackagePath string
	SchemaPackageName string
}

var memberTemplate = "\t%v\t%v\t`db:%q bson:%q comment:%q json:%q xml:%q`%s"

var (
	structTemplateObj *template.Template
	//modelBaseTemplateObj *template.Template
	modelTemplateObj *template.Template
	initTemplateObj  *template.Template
	once             sync.Once
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
		b, err = io.ReadAll(fs)
	} else {
		b, err = os.ReadFile(tempfile)
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
	default:
		panic("unknown template name: " + name)
	}
	return w.Bytes(), err
}
