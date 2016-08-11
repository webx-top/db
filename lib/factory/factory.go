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
		databases: make([]*Cluster, 0),
	}
}

func NewCluster() *Cluster {
	return &Cluster{
		masters: []sqlbuilder.Database{},
		slaves:  []sqlbuilder.Database{},
	}
}

type Cluster struct {
	masters []sqlbuilder.Database
	slaves  []sqlbuilder.Database
}

func (c *Cluster) W() sqlbuilder.Database {
	length := len(c.masters)
	if length == 0 {
		panic(`Not connected to any database`)
	}
	if length > 1 {
		return c.masters[rand.Intn(length-1)]
	}
	return c.masters[0]
}

func (c *Cluster) R() sqlbuilder.Database {
	length := len(c.slaves)
	if length == 0 {
		return c.W()
	}
	if length > 1 {
		return c.slaves[rand.Intn(length-1)]
	}
	return c.slaves[0]
}

func (c *Cluster) AddW(databases ...sqlbuilder.Database) {
	c.masters = append(c.masters, databases...)
}

func (c *Cluster) AddR(databases ...sqlbuilder.Database) {
	c.slaves = append(c.slaves, databases...)
}

func (c *Cluster) CloseAll() {
	for _, database := range c.masters {
		if err := database.Close(); err != nil {
			log.Println(err.Error())
		}
	}
	for _, database := range c.slaves {
		if err := database.Close(); err != nil {
			log.Println(err.Error())
		}
	}
}

type Factory struct {
	databases []*Cluster
}

func (f *Factory) AddDB(databases ...sqlbuilder.Database) *Factory {
	if len(f.databases) > 0 {
		f.databases[0].AddW(databases...)
	} else {
		c := NewCluster()
		c.AddW(databases...)
		f.databases = append(f.databases, c)
	}
	return f
}

func (f *Factory) AddSlaveDB(databases ...sqlbuilder.Database) *Factory {
	if len(f.databases) > 0 {
		f.databases[0].AddR(databases...)
	} else {
		c := NewCluster()
		c.AddR(databases...)
		f.databases = append(f.databases, c)
	}
	return f
}

func (f *Factory) SetDB(index int, database sqlbuilder.Database, args ...bool) *Factory {
	if len(f.databases) > index {
		var isMaster bool
		if len(args) > 0 {
			isMaster = args[0]
		}
		if isMaster {
			f.databases[index].AddW(database)
		} else {
			f.databases[index].AddR(database)
		}
	}
	return f
}

func (f *Factory) SetCluster(index int, cluster *Cluster) *Factory {
	if len(f.databases) > index {
		f.databases[index] = cluster
	}
	return f
}

func (f *Factory) DB(index int, args ...bool) sqlbuilder.Database {
	if len(f.databases) > index {
		var isMaster bool
		if len(args) > 0 {
			isMaster = args[0]
		}
		if isMaster {
			return f.databases[index].W()
		}
		return f.databases[index].R()
	}
	if index == 0 {
		panic(`Not connected to any database`)
	}
	return f.DB(0, args...)
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
	for _, cluster := range f.databases {
		cluster.CloseAll()
	}
}
