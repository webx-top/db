package main

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/webx-top/db"
	"github.com/webx-top/db/_tools/test/settings"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
)

func main() {
	c := settings.Connect()
	ctx := context.Background()
	trans, err := factory.NewTx(ctx)
	if err != nil {
		panic(err)
	}
	param := factory.NewParam()
	_, err = trans.SQLBuilder(param).Update(`vhost`).Set(echo.H{`disabled`: `Y`}).Where(db.Cond{`id`: 1}).Exec()
	if err != nil {
		panic(err)
	}
	err = trans.Rollback()
	if err != nil {
		panic(err)
	}
	_ = c
}
