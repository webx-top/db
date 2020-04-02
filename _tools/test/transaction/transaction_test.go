package transaction_test

import (
	"context"
	"testing"

	"github.com/admpub/null"
	"github.com/webx-top/db"
	"github.com/webx-top/db/_tools/test/settings"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/testing/test"
)

var (
	id   = 1
	cond = db.Cond{`id`: id}
)

func verifyResult(t *testing.T, result string) null.StringMap {
	recv := null.StringMap{}
	err := factory.NewParam().SetCollection(`vhost`).SetRecv(&recv).SetArgs(cond).One()
	if err != nil {
		panic(err)
	}
	//com.Dump(recv)
	test.Eq(t, result, recv.String(`disabled`))
	return recv
}

func TestTransaction(t *testing.T) {
	c := settings.Connect()
	//reset
	_, err := factory.NewParam().SQLBuilder().Update(`vhost`).Set(echo.H{`disabled`: `N`}).Where(cond).Exec()
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
	_, err = param.SQLBuilder().Update(`vhost`).Set(echo.H{`disabled`: `Y`}).Where(cond).Exec()
	if err != nil {
		panic(err)
	}
	err = trans.Rollback()
	if err != nil {
		panic(err)
	}
	verifyResult(t, `N`)

	param2 := factory.NewParam()
	err = param2.SetTrans(trans).SetCollection(`vhost`).SetArgs(cond).SetSend(echo.H{`disabled`: `Y`}).Update()
	if err != nil {
		panic(err)
	}
	verifyResult(t, `Y`)

	//reset
	err = param2.SetTrans(trans).SetCollection(`vhost`).SetArgs(cond).SetSend(echo.H{`disabled`: `N`}).Update()
	if err != nil {
		panic(err)
	}
	verifyResult(t, `N`)

	_ = c
}

func TestTransaction2(t *testing.T) {
	c := settings.Connect()
	//reset
	_, err := factory.NewParam().SQLBuilder().Update(`vhost`).Set(echo.H{`disabled`: `N`}).Where(cond).Exec()
	if err != nil {
		panic(err)
	}
	verifyResult(t, `N`)
	ctx := context.Background()
	tx, err := factory.NewTx(ctx)
	if err != nil {
		panic(err)
	}
	p := factory.NewParam().SetTrans(tx)
	err = p.SetCollection(`vhost`).SetArgs(cond).SetSend(echo.H{`disabled`: `Y`}).Update()
	if err != nil {
		panic(err)
	}
	err = tx.Rollback()
	if err != nil {
		panic(err)
	}
	recv := verifyResult(t, `N`)

	ta := echo.NewTransaction(p)
	err = ta.Begin(ctx)
	if err != nil {
		panic(err)
	}
	recv[`id`] = null.String{`1000`, true}
	_, err = factory.NewParam().TransFrom(p).SetCollection(`vhost`).SetSend(recv).Insert()
	if err != nil {
		panic(err)
	}
	err = ta.Rollback(ctx)
	if err != nil {
		panic(err)
	}
	n, err := factory.NewParam().SetCollection(`vhost`).SetRecv(&recv).SetArgs(db.Cond{`id`: 1000}).Count()
	if err != nil {
		panic(err)
	}
	test.Eq(t, int64(0), n)
	_ = c
}
