package transaction_test

import (
	"context"
	"testing"

	"github.com/admpub/nging/application/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/db/_tools/test/settings"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/testing/test"
)

func verifyResult(t *testing.T, result string) {
	recv := &dbschema.Vhost{}
	err := factory.NewParam().SetCollection(`vhost`).SetRecv(recv).SetArgs(db.Cond{`id`: 1}).One()
	if err != nil {
		panic(err)
	}
	test.Eq(t, result, recv.Disabled)
}

func TestTransaction(t *testing.T) {
	c := settings.Connect()
	//reset
	_, err := factory.NewParam().SQLBuilder().Update(`vhost`).Set(echo.H{`disabled`: `N`}).Where(db.Cond{`id`: 1}).Exec()
	if err != nil {
		panic(err)
	}
	verifyResult(t, `N`)
	ctx := context.Background()
	trans, err := factory.NewTx(ctx)
	if err != nil {
		panic(err)
	}
	param := factory.NewParam().SetTrans(trans)
	_, err = param.SQLBuilder().Update(`vhost`).Set(echo.H{`disabled`: `Y`}).Where(db.Cond{`id`: 1}).Exec()
	if err != nil {
		panic(err)
	}
	err = trans.Rollback()
	if err != nil {
		panic(err)
	}
	verifyResult(t, `N`)

	param2 := factory.NewParam()
	err = param2.SetTrans(trans).SetCollection(`vhost`).SetArgs(`id`, 1).SetSend(echo.H{`disabled`: `Y`}).Update()
	if err != nil {
		panic(err)
	}
	verifyResult(t, `Y`)
	_ = c
}
