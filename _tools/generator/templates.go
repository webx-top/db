package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var memberTemplate = "\t%v\t%v\t`db:\"%v\" bson:\"%v\" comment:\"%v\" json:\"%v\" xml:\"%v\"`"
var replaces = &map[string]string{
	"packageName":  "",
	"imports":      "",
	"structName":   "",
	"attributes":   "",
	"reset":        "",
	"asMap":        "",
	"asRow":        "",
	"tableName":    "",
	"beforeInsert": "",
	"afterInsert":  "",
	"beforeUpdate": "",
	"setUpdatedAt": "",
	"beforeDelete": "",
}
var structFuncs = map[string]string{
	`Trans`:        `Trans`,
	`Use`:          `Use`,
	`Objects`:      `Objects`,
	`NewObjects`:   `NewObjects`,
	`NewParam`:     `NewParam`,
	`SetParam`:     `SetParam`,
	`Get`:          `Get`,
	`List`:         `List`,
	`ListByOffset`: `ListByOffset`,
	`Add`:          `Add`,
	`Edit`:         `Edit`,
	`Upsert`:       `Upsert`,
	`Delete`:       `Delete`,
	`Count`:        `Count`,
}
var structTemplate = `// @generated Do not edit this file, which is automatically generated by the generator.

package {{packageName}}

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	{{imports}}
)

type {{structName}} struct {
	param   *factory.Param
	trans	*factory.Transaction
	objects []*{{structName}}
	namer   func(string) string
	connID  int
	
{{attributes}}
}

func (this *{{structName}}) Trans() *factory.Transaction {
	return this.trans
}

func (this *{{structName}}) Use(trans *factory.Transaction) factory.Model {
	this.trans = trans
	return this
}

func (this *{{structName}}) SetConnID(connID int) factory.Model {
	this.connID = connID
	return this
}

func (this *{{structName}}) New(structName string, connID ...int) factory.Model {
	if len(connID) > 0 {
		return factory.NewModel(structName,connID[0]).Use(this.trans)
	}
	return factory.NewModel(structName,this.connID).Use(this.trans)
}

func (this *{{structName}}) Objects() []*{{structName}} {
	if this.objects == nil {
		return nil
	}
	return this.objects[:]
}

func (this *{{structName}}) NewObjects() *[]*{{structName}} {
	this.objects = []*{{structName}}{}
	return &this.objects
}

func (this *{{structName}}) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(this.connID).SetTrans(this.trans).SetCollection(this.Name_()).SetModel(this)
}

func (this *{{structName}}) SetNamer(namer func (string) string) factory.Model {
	this.namer = namer
	return this
}

func (this *{{structName}}) Name_() string {
	if this.namer != nil {
		return this.namer("{{tableName}}")
	}
	return factory.TableNamerGet("{{tableName}}")(this)
}

func (this *{{structName}}) FullName_(connID ...int) string {
	if len(connID) > 0 {
		return factory.DefaultFactory.Cluster(connID[0]).Table(this.Name_())
	}
	return factory.DefaultFactory.Cluster(this.connID).Table(this.Name_())
}

func (this *{{structName}}) SetParam(param *factory.Param) factory.Model {
	this.param = param
	return this
}

func (this *{{structName}}) Param() *factory.Param {
	if this.param == nil {
		return this.NewParam()
	}
	return this.param
}

func (this *{{structName}}) Get(mw func(db.Result) db.Result, args ...interface{}) error {
	return this.Param().SetArgs(args...).SetRecv(this).SetMiddleware(mw).One()
}

func (this *{{structName}}) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = this.NewObjects()
	}
	return this.Param().SetArgs(args...).SetPage(page).SetSize(size).SetRecv(recv).SetMiddleware(mw).List()
}

func (this *{{structName}}) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = this.NewObjects()
	}
	return this.Param().SetArgs(args...).SetOffset(offset).SetSize(size).SetRecv(recv).SetMiddleware(mw).List()
}

func (this *{{structName}}) Add() (pk interface{}, err error) {
	{{beforeInsert}}
	pk, err = this.Param().SetSend(this).Insert()
	{{afterInsert}}
	return
}

func (this *{{structName}}) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	{{beforeUpdate}}
	return this.Setter(mw, args...).SetSend(this).Update()
}

func (this *{{structName}}) Setter(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	return this.Param().SetArgs(args...).SetMiddleware(mw)
}

func (this *{{structName}}) SetField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) error {
	return this.SetFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (this *{{structName}}) SetFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) error {
	{{setUpdatedAt}}
	return this.Setter(mw, args...).SetSend(kvset).Update()
}

func (this *{{structName}}) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = this.Param().SetArgs(args...).SetSend(this).SetMiddleware(mw).Upsert(func(){
		{{beforeUpdate}}
	},func(){
		{{beforeInsert}}
	})
	{{afterInsert}}
	return 
}

func (this *{{structName}}) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	{{beforeDelete}}
	return this.Param().SetArgs(args...).SetMiddleware(mw).Delete()
}

func (this *{{structName}}) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return this.Param().SetArgs(args...).SetMiddleware(mw).Count()
}

func (this *{{structName}}) Reset() *{{structName}} {
{{reset}}
	return this
}

func (this *{{structName}}) AsMap() map[string]interface{} {
	r := map[string]interface{}{}
{{asMap}}
	return r
}

func (this *{{structName}}) AsRow() map[string]interface{} {
	r := map[string]interface{}{}
{{asRow}}
	return r
}

func (this *{{structName}}) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = this.AsRow()
	}
	return factory.BatchValidate("{{tableName}}", kvset)
}

func (this *{{structName}}) Validate(field string, value interface{}) error {
	return factory.Validate("{{tableName}}", field, value)
}

`

var modelReplaces = &map[string]string{
	"packageName":       "",
	"imports":           "",
	"structName":        "",
	"schemaPackagePath": "",
	"schemaPackageName": "",
	"baseName":          "",
}
var modelBaseTemplate = `package {{packageName}}

import (
	"github.com/webx-top/echo"
	{{imports}}
)

type {{structName}} struct {
	echo.Context
}
`

var modelTemplate = `package {{packageName}}

import (
	"{{schemaPackagePath}}"
	"github.com/webx-top/echo"
	{{imports}}
)
				
func New{{structName}}(ctx echo.Context) *{{structName}} {
	return &{{structName}}{
		{{structName}}: &{{schemaPackageName}}.{{structName}}{},
		{{baseName}}:   &{{baseName}}{Context: ctx},
	}
}

type {{structName}} struct {
	*{{schemaPackageName}}.{{structName}}
	*{{baseName}}
}

`

var initFileTemplate = `// @generated Do not edit this file, which is automatically generated by the generator.

package {{packageName}}

import (
	"github.com/webx-top/db/lib/factory"
)

func init(){
	{{initCode}}
}

`

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
	args := []string{
		"--default-character-set=" + cfg.Charset,
		"--single-transaction",
		"--opt",
		"-d", //加上此参数代表只导出表结构，不导出数据
		"-h" + host,
		"-P" + port,
		"-u" + cfg.Username,
		"-p" + cfg.Password,
		cfg.Database,
	}
	args = append(args, tables...)
	cmd := exec.Command("mysqldump", args...)
	fp, err := os.Create(cfg.Backup)
	if err != nil {
		log.Println(`Failed to backup:`, err)
	}
	defer fp.Close()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(`Failed to backup:`, err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(`Failed to backup:`, err)
	}
	if _, err := io.Copy(fp, stdout); err != nil {
		log.Fatal(`Failed to backup:`, err)
	}
	b, err := ioutil.ReadFile(cfg.Backup)
	if err != nil {
		log.Fatal(err)
	}
	b = cleanRegExp.ReplaceAll(b, []byte(` `))
	err = ioutil.WriteFile(cfg.Backup, b, 0666)
	if err != nil {
		log.Fatal(err)
	}
}
