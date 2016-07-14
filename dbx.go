// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbx

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/admpub/core"
	"github.com/webx-top/dbx/config"
	"github.com/webx-top/dbx/driver"
)

const Version = "0.1.0"

func connect(conf *config.Config) (*core.DB, core.Dialect, error) {
	driverName := conf.Engine
	dataSourceName := conf.String()
	driver := core.QueryDriver(driverName)
	if driver == nil {
		return nil, nil, fmt.Errorf("Unsupported driver name: %v", driverName)
	}

	uri, err := driver.Parse(driverName, dataSourceName)
	if err != nil {
		return nil, nil, err
	}

	dialect := core.QueryDialect(uri.DbType)
	if dialect == nil {
		return nil, nil, fmt.Errorf("Unsupported dialect type: %v", uri.DbType)
	}

	db, err := core.Open(driverName, dataSourceName)
	if err != nil {
		return nil, nil, err
	}

	err = dialect.Init(db, uri, driverName, dataSourceName)
	if err != nil {
		return nil, nil, err
	}

	return db, dialect, err
}

/*
func RegsterDriver() bool {
	providedDrvsNDialects := map[string]struct {
		dbType     core.DbType
		getDriver  func() core.Driver
		getDialect func() core.Dialect
	}{
		"mssql":    {"mssql", func() core.Driver { return &odbcDriver{} }, func() core.Dialect { return &mssql{} }},
		"odbc":     {"mssql", func() core.Driver { return &odbcDriver{} }, func() core.Dialect { return &mssql{} }}, // !nashtsai! TODO change this when supporting MS Access
		"mysql":    {"mysql", func() core.Driver { return &mysqlDriver{} }, func() core.Dialect { return &mysql{} }},
		"mymysql":  {"mysql", func() core.Driver { return &mymysqlDriver{} }, func() core.Dialect { return &mysql{} }},
		"postgres": {"postgres", func() core.Driver { return &pqDriver{} }, func() core.Dialect { return &postgres{} }},
		"sqlite3":  {"sqlite3", func() core.Driver { return &sqlite3Driver{} }, func() core.Dialect { return &sqlite3{} }},
		"oci8":     {"oracle", func() core.Driver { return &oci8Driver{} }, func() core.Dialect { return &oracle{} }},
		"goracle":  {"oracle", func() core.Driver { return &goracleDriver{} }, func() core.Dialect { return &oracle{} }},
	}

	for driverName, v := range providedDrvsNDialects {
		if driver := core.QueryDriver(driverName); driver == nil {
			core.RegisterDriver(driverName, v.getDriver())
			core.RegisterDialect(v.dbType, v.getDialect)
		}
	}
	return true
}
*/

func close(engine *Engine) {
	engine.Close()
}

// NewEngine new a db manager according to the parameter. Currently support four
// drivers
func NewEngine(conf *config.Config) (*Engine, error) {
	driver.RegDBDrivers()

	db, dialect, err := connect(conf)
	if err != nil {
		return nil, err
	}

	engine := &Engine{}
	engine.Init(db, dialect)

	runtime.SetFinalizer(engine, close)

	return engine, nil
}

// AddSlaveDB .
func (engine *Engine) AddSlaveDB(conf *config.Config) error {
	db, dialect, err := connect(conf)
	if err != nil {
		return err
	}
	dialect.SetLogger(engine.TLogger.Base.Logger)
	engine.slaveDB = append(engine.slaveDB, &DB{DB: db, Dialect: dialect})
	return nil
}

// AddMasterDB .
func (engine *Engine) AddMasterDB(conf *config.Config) error {
	db, dialect, err := connect(conf)
	if err != nil {
		return err
	}
	dialect.SetLogger(engine.TLogger.Base.Logger)
	engine.masterDB = append(engine.masterDB, &DB{DB: db, Dialect: dialect})
	return engine
}

// Clone clone an engine
func (engine *Engine) Clone() (*Engine, error) {
	newEngine := *engine
	return &newEngine
}

func (engine *Engine) Init(db *core.DB, dialect core.Dialect) *Engine {

	logger := NewSimpleLogger(os.Stdout)
	logger.SetLevel(core.LOG_INFO)

	dialect.SetLogger(logger)

	engine.masterDB = []*DB{&DB{DB: db, Dialect: dialect}}
	engine.slaveDB = []*DB{}
	engine.dialect = dialect
	engine.Tables = make(map[reflect.Type]*core.Table)
	engine.mutex = &sync.RWMutex{}
	engine.DbxTagIdentifier = "dbx"
	engine.RelTagIdentifier = "rel"
	engine.TLogger = NewTLogger(logger)
	engine.TZLocation = time.Local

	engine.SetMapper(core.NewCacheMapper(new(core.SnakeMapper)))

	return engine
}
