package main

import (
	"fmt"

	"github.com/admpub/nging/application/dbschema"
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
	//Prefix:   "t_",
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

func main() {
	db.DefaultSettings.SetLogging(true)
	c, err := mysql.Open(settings)
	if err != nil {
		panic(err)
	}
	factory.AddDB(c)
	row := &GroupAndVHost{}
	err = c.SelectFrom(`vhost_group`).Relation(`Vhost`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
		return sel.OrderBy(`-id`)
	}).One(row)
	if err != nil {
		panic(err)
	}
	echo.Dump(row)
	return

	fmt.Println(`===========================================`)

	rows := []*GroupAndVHost{}
	//Relation 是可选的，用于增加额外条件
	err = c.SelectFrom(`vhost_group`).Relation(`Vhost`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
		return sel.OrderBy(`-id`)
	}).All(&rows)
	if err != nil {
		panic(err)
	}
	echo.Dump(rows)

	fmt.Println(`===========================================`)

	rows2 := []*GroupAndVHosts{}
	err = c.Collection(`vhost_group`).Find().Relation(`Vhosts`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
		return sel.OrderBy(`id`) //.ForceIndex(`group_id`)
	}).All(&rows2)
	if err != nil {
		panic(err)
	}
	echo.Dump(rows2)
}
