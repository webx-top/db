package driver

import (
	"github.com/admpub/core"
)

type DatabaseDriver struct {
	Type    core.DbType
	Driver  func() core.Driver
	Dialect func() core.Dialect
}

func (d *DatabaseDriver) IsZero() bool {
	return d.Dialect == nil && d.Dialect == nil
}

var defaultDatabaseDriver = &DatabaseDriver{
	Type:    `None`,
	Driver:  nil,
	Dialect: nil,
}

var databaseDrivers = make(map[string]*DatabaseDriver)

func RegDatabaseDriver(name string, driver *DatabaseDriver) {
	databaseDrivers[name] = driver
}

func GetDatabaseDriver(name string) *DatabaseDriver {
	if dr, ok := databaseDrivers[name]; ok {
		return dr
	}
	return defaultDatabaseDriver
}

func DelDatabaseDriver(name string) {
	if _, ok := databaseDrivers[name]; ok {
		delete(databaseDrivers, name)
	}
}

func ResetDatabaseDriver() {
	databaseDrivers = make(map[string]*DatabaseDriver)
}

func RegDatabaseDrivers() {
	for driverName, v := range databaseDrivers {
		if driver := core.QueryDriver(driverName); driver == nil {
			core.RegisterDriver(driverName, v.Driver())
			core.RegisterDialect(v.Type, v.Dialect)
		}
	}
}
