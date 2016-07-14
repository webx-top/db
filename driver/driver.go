package driver

import (
	"github.com/admpub/core"
)

type DBDriver struct {
	Type    core.DbType
	Driver  func() core.Driver
	Dialect func() core.Dialect
}

func (d *DBDriver) IsZero() bool {
	return d.Dialect == nil && d.Dialect == nil
}

var defaultDBDriver = &DBDriver{
	Type:    `None`,
	Driver:  nil,
	Dialect: nil,
}

var dbDrivers = make(map[string]*DBDriver)

func RegDBDriver(name string, driver *DBDriver) {
	dbDrivers[name] = driver
}

func GetDBDriver(name string) *DBDriver {
	if dr, ok := dbDrivers[name]; ok {
		return dr
	}
	return defaultDBDriver
}

func DelDBDriver(name string) {
	if _, ok := dbDrivers[name]; ok {
		delete(dbDrivers, name)
	}
}

func ResetDBDriver() {
	dbDrivers = make(map[string]*DBDriver)
}

func RegDBDrivers() {
	for driverName, v := range dbDrivers {
		if driver := core.QueryDriver(driverName); driver == nil {
			core.RegisterDriver(driverName, v.Driver())
			core.RegisterDialect(v.Type, v.Dialect)
		}
	}
}
