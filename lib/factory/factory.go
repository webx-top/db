// added by swh@admpub.com
package factory

import (
	"log"
	"math/rand"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
)

func New() *Factory {
	return &Factory{
		databases: make([]sqlbuilder.Database, 0),
		masters:   make([]int, 0),
		slaves:    make([]int, 0),
	}
}

type Factory struct {
	databases []sqlbuilder.Database
	masters   []int
	slaves    []int
}

func (f *Factory) AddDB(databases ...sqlbuilder.Database) *Factory {
	f.masters = append(f.masters, len(f.databases))
	f.databases = append(f.databases, databases...)
	return f
}

func (f *Factory) AddSlaveDB(databases ...sqlbuilder.Database) *Factory {
	f.slaves = append(f.slaves, len(f.databases))
	f.databases = append(f.databases, databases...)
	return f
}

func (f *Factory) SetDB(index int, database sqlbuilder.Database, args ...bool) *Factory {
	if len(f.databases) > index {
		f.databases[index] = database
		if len(args) > 0 {
			isMaster := args[0]
			found := false
			for key, mindex := range f.masters {
				if mindex == index {
					if !isMaster {
						f.masters = append(f.masters[:key], f.masters[key+1:]...)
					}
					found = true
					break
				}
			}
			if !found && isMaster {
				f.masters = append(f.masters, index)
			}
			found = false
			for key, mindex := range f.slaves {
				if mindex == index {
					if isMaster {
						f.slaves = append(f.slaves[:key], f.slaves[key+1:]...)
					}
					found = true
					break
				}
			}
			if !found && !isMaster {
				f.slaves = append(f.slaves, index)
			}
		}
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

func (f *Factory) MasterDB() sqlbuilder.Database {
	length := len(f.masters)
	var index int
	if length > 1 {
		index = f.masters[rand.Intn(length-1)]
	}
	if len(f.databases) > index {
		return f.databases[index]
	}
	if index == 0 {
		panic(`Not connected to any database`)
	}
	return f.DB(0)
}

func (f *Factory) SlaveDB() sqlbuilder.Database {
	length := len(f.slaves)
	var index int
	if length > 1 {
		index = f.slaves[rand.Intn(length-1)]
	}
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

func (f *Factory) CloseAll() {
	for _, database := range f.databases {
		if err := database.Close(); err != nil {
			log.Println(err.Error())
		}
	}
}
