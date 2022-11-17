package transaction_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/db/_tools/test/dbschema"
	"github.com/webx-top/db/_tools/test/settings"

	"github.com/admpub/log"
	"github.com/admpub/null"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/testing/test"
)

var (
	id   = 1
	cond = db.Cond{`id`: id}
	c    db.Database
)

func TestMain(m *testing.M) {
	defer log.Close()
	c = settings.Connect()
	m.Run()
	c.Close()
}

/*
func TestOneMap(t *testing.T) {
	recv := map[string]*null.String{}
	err := factory.NewParam().SetCollection(`nging_vhost`).SetRecv(&recv).SetArgs(cond).One()
	if err != nil {
		panic(err)
	}
	com.Dump(recv)
}
*/

func verifyResult(t *testing.T, result string) null.StringMap {
	recv := null.StringMap{}
	err := factory.NewParam().SetCollection(`nging_vhost`).SetRecv(&recv).SetArgs(cond).One()
	if err != nil {
		panic(err)
	}
	//com.Dump(recv)
	test.Eq(t, result, recv.String(`disabled`))
	return recv
}

func TestTransaction(t *testing.T) {
	//reset
	_, err := factory.NewParam().SQLBuilder().Update(`nging_vhost`).Set(echo.H{`disabled`: `N`}).Where(cond).Exec()
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
	_, err = param.SQLBuilder().Update(`nging_vhost`).Set(echo.H{`disabled`: `Y`}).Where(cond).Exec()
	if err != nil {
		panic(err)
	}
	err = param.Rollback(ctx)
	if err != nil {
		panic(err)
	}
	verifyResult(t, `N`)

	param2 := factory.NewParam()
	err = param2.SetTrans(trans).SetCollection(`nging_vhost`).SetArgs(cond).SetSend(echo.H{`disabled`: `Y`}).Update()
	if err != nil {
		panic(err)
	}
	verifyResult(t, `Y`)

	//reset
	err = param2.SetTrans(trans).SetCollection(`nging_vhost`).SetArgs(cond).SetSend(echo.H{`disabled`: `N`}).Update()
	if err != nil {
		panic(err)
	}
	verifyResult(t, `N`)
}

func TestTransaction2(t *testing.T) {
	//reset
	_, err := factory.NewParam().SQLBuilder().Update(`nging_vhost`).Set(echo.H{`disabled`: `N`}).Where(cond).Exec()
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
	err = p.SetCollection(`nging_vhost`).SetArgs(cond).SetSend(echo.H{`disabled`: `Y`}).Update()
	if err != nil {
		panic(err)
	}
	err = p.Rollback(ctx)
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
	_, err = factory.NewParam().TransFrom(p).SetCollection(`nging_vhost`).SetSend(recv).Insert()
	if err != nil {
		panic(err)
	}
	err = ta.Rollback(ctx)
	if err != nil {
		panic(err)
	}
	n, err := factory.NewParam().SetCollection(`nging_vhost`).SetRecv(&recv).SetArgs(db.Cond{`id`: 1000}).Count()
	if err != nil {
		panic(err)
	}
	test.Eq(t, int64(0), n)
}

// TestContext 测试用ctx初始化后再开启事物也能生效
func TestContext(t *testing.T) {
	ctx := defaults.NewMockContext()
	ctx.SetTransaction(echo.NewTransaction(factory.NewParam()))
	m := dbschema.NewNgingVhost(ctx)
	err := m.UpdateField(nil, `disabled`, `Y`, `id`, 1)
	assert.NoError(t, err)
	verifyResult(t, `Y`)

	ctx.Begin()
	err = m.UpdateField(nil, `disabled`, `N`, `id`, 1)
	assert.NoError(t, err)
	ctx.Rollback()
	verifyResult(t, `Y`)
}
