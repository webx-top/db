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
	Password: "root",
}

type GroupAndVHost struct {
	*dbschema.VhostGroup
	Vhosts []*dbschema.Vhost `db:"-" relation:"id,group_id"`
}

func main() {
	db.DefaultSettings.SetLogging(true)
	c, err := mysql.Open(settings)
	if err != nil {
		panic(err)
	}
	factory.AddDB(c)
	rows := []*GroupAndVHost{}
	err = c.SelectFrom(`vhost_group`).Relation(`Vhosts`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
		return sel.OrderBy(`-id`).Limit(3)
	}).All(&rows)
	if err != nil {
		panic(err)
	}
	echo.Dump(rows)

	fmt.Println(`===========================================`)

	rows2 := []*GroupAndVHost{}
	err = c.Collection(`vhost_group`).Find().Relation(`Vhosts`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
		return sel.OrderBy(`-id`).Limit(2)
	}).All(&rows2)
	if err != nil {
		panic(err)
	}
	echo.Dump(rows2)
}
