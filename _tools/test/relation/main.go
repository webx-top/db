package main

import (
	"fmt"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/null"
	_ "github.com/go-sql-driver/mysql"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/db/mysql"
	"github.com/webx-top/echo"
)

var settings = mysql.ConnectionURL{
	Host:     "localhost",
	Database: "nging",
	User:     "root",
	Password: "root",
}

type GroupAndVHosts struct {
	*dbschema.VhostGroup
	Vhosts []*dbschema.Vhost `db:"-,relation=id:group_id"` //relation=<外键>:<vhost的主键>
}

type GroupAndVHost struct {
	*dbschema.VhostGroup
	Vhost *dbschema.Vhost `db:"-,relation=id:group_id"` //relation=<外键>:<vhost的主键>
}

type JoinData struct {
	Group *dbschema.VhostGroup
	Vhost *dbschema.Vhost
}

func testJoin(c db.Database) {
	jrow := &JoinData{}
	c.Collection(`vhost_group`).Find().Callback(func(sel sqlbuilder.Selector) sqlbuilder.Selector {
		return sel.From(`vhost_group g`).Join(`vhost v`).On(`v.group_id=g.id`)
	}).One(jrow)
	echo.Dump(jrow)
}

func testOne2one(c db.Database) {
	row := &GroupAndVHost{}
	err := c.(sqlbuilder.Database).SelectFrom(`vhost_group`).Relation(`Vhost`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
		return sel.OrderBy(`-id`)
	}).One(row)
	if err != nil {
		panic(err)
	}
	echo.Dump(row)
}

func testMultiOne2one(c db.Database) {
	rows := []*GroupAndVHost{}
	//Relation 是可选的，用于增加额外条件
	err := c.(sqlbuilder.Database).SelectFrom(`vhost_group`).Relation(`Vhost`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
		return sel.OrderBy(`-id`)
	}).All(&rows)
	if err != nil {
		panic(err)
	}
	echo.Dump(rows)
}

func testMultiOne2many(c db.Database) {
	rows2 := []*GroupAndVHosts{}
	err := c.Collection(`vhost_group`).Find().Relation(`Vhosts`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
		return sel.OrderBy(`id`) //.ForceIndex(`group_id`)
	}).All(&rows2)
	if err != nil {
		panic(err)
	}
	echo.Dump(rows2)

}

func testMap(c db.Database) {
	row2 := null.StringMap{}
	err := c.(sqlbuilder.Database).SelectFrom(`vhost_group`).One(&row2)
	if err != nil {
		panic(err)
	}
	echo.Dump(row2)

	rows3 := null.StringMapSlice{}
	err = c.(sqlbuilder.Database).SelectFrom(`vhost_group`).Limit(2).All(&rows3)
	if err != nil {
		panic(err)
	}
	echo.Dump(rows3)
}

func main() {
	db.DefaultSettings.SetLogging(true)
	c, err := mysql.Open(settings)
	if err != nil {
		panic(err)
	}
	factory.AddDB(c)

	testJoin(c)
	return
	testOne2one(c)

	fmt.Println(`===========================================`)
	testMultiOne2one(c)

	fmt.Println(`===========================================`)
	testMultiOne2many(c)
	//验证map方式是否正常==============================
	testMap(c)
}
