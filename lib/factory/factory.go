// added by swh@admpub.com
package factory

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
)

func New() *Factory {
	return &Factory{
		databases: make([]sqlbuilder.Database, 0),
	}
}

type Factory struct {
	databases []sqlbuilder.Database
}

func (f *Factory) AddDB(databases ...sqlbuilder.Database) *Factory {
	f.databases = append(f.databases, databases...)
	return f
}

func (f *Factory) SetDB(index int, database sqlbuilder.Database) *Factory {
	if len(f.databases) > index {
		f.databases[index] = database
	}
	return f
}

func (f *Factory) DB(index int) sqlbuilder.Database {
	if len(f.databases) > index {
		return f.databases[index]
	}
	if index == 0 {
		panic(`Not connected to any database`)
	}
	return f.DB(0)
}

func (f *Factory) Collection(collection string, args ...int) db.Collection {
	var index int
	if len(args) > 0 {
		index = args[0]
	}
	return f.DB(index).Collection(collection)
}

func (f *Factory) Find(collection string, args ...interface{}) db.Result {
	return f.Collection(collection).Find(args...)
}

func (f *Factory) FindDB(index int, collection string, args ...interface{}) db.Result {
	return f.Collection(collection, index).Find(args...)
}
