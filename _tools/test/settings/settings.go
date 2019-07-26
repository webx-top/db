package settings

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/mysql"
)

var settings = mysql.ConnectionURL{
	Host:     "localhost",
	Database: "nging",
	User:     "root",
	Password: "root",
}

func Connect() db.Database {
	db.DefaultSettings.SetLogging(true)
	c, err := mysql.Open(settings)
	if err != nil {
		panic(err)
	}
	factory.AddDB(c)
	return c
}
