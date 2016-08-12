// added by swh@admpub.com
package factory

import (
	"errors"
	"log"
	"math/rand"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
)

const (
	R = iota
	W
)

var ErrNotFoundKey = errors.New(`not found the key`)

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
	prefix  string
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

func (c *Cluster) Prefix() string {
	return c.prefix
}

func (c *Cluster) Table(tableName string) string {
	return c.prefix + tableName
}

func (c *Cluster) SetPrefix(prefix string) {
	c.prefix = prefix
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

func (c *Cluster) SetW(index int, database sqlbuilder.Database) error {
	if len(c.masters) > index {
		c.masters[index] = database
		return nil
	}
	return ErrNotFoundKey
}

func (c *Cluster) SetR(index int, database sqlbuilder.Database) error {
	if len(c.masters) > index {
		c.slaves[index] = database
		return nil
	}
	return ErrNotFoundKey
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

func (f *Factory) Debug() bool {
	return db.Debug
}

func (f *Factory) SetDebug(on bool) *Factory {
	db.Debug = on
	return f
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

func (f *Factory) SetCluster(index int, cluster *Cluster) *Factory {
	if len(f.databases) > index {
		f.databases[index] = cluster
	}
	return f
}

func (f *Factory) AddCluster(clusters ...*Cluster) *Factory {
	f.databases = append(f.databases, clusters...)
	return f
}

func (f *Factory) Cluster(index int) *Cluster {
	if len(f.databases) > index {
		return f.databases[index]
	}
	if index == 0 {
		panic(`Not connected to any database`)
	}
	return f.Cluster(0)
}

func (f *Factory) Collection(collection string, args ...int) db.Collection {
	var index int
	switch len(args) {
	case 2:
		index = args[0]
		if args[1] == R {
			c := f.Cluster(index)
			collection = c.Table(collection)
			return c.R().Collection(collection)
		}
	case 1:
		index = args[0]
	}
	c := f.Cluster(index)
	collection = c.Table(collection)
	return c.W().Collection(collection)
}

func (f *Factory) Find(collection string, args ...interface{}) db.Result {
	return f.Collection(collection).Find(args...)
}

func (f *Factory) FindR(collection string, args ...interface{}) db.Result {
	return f.Collection(collection, 0, R).Find(args...)
}

func (f *Factory) FindDB(index int, collection string, args ...interface{}) db.Result {
	return f.Collection(collection, index).Find(args...)
}

func (f *Factory) FindDBR(index int, collection string, args ...interface{}) db.Result {
	return f.Collection(collection, index, R).Find(args...)
}

func (f *Factory) CloseAll() {
	for _, cluster := range f.databases {
		cluster.CloseAll()
	}
}

// ================================
// API
// ================================

// Read ==========================

func (f *Factory) All(collection string, fn func(db.Result) db.Result, result interface{}, args ...interface{}) error {
	return f.AllFromDB(0, collection, fn, result, args...)
}

func (f *Factory) AllFromDB(index int, collection string, fn func(db.Result) db.Result, result interface{}, args ...interface{}) error {
	if fn == nil {
		return f.FindDBR(index, collection, args...).All(result)
	}
	return fn(f.FindDBR(index, collection, args...)).All(result)
}

func (f *Factory) One(collection string, fn func(db.Result) db.Result, result interface{}, args ...interface{}) error {
	return f.OneFromDB(0, collection, fn, result, args...)
}

func (f *Factory) OneFromDB(index int, collection string, fn func(db.Result) db.Result, result interface{}, args ...interface{}) error {
	if fn == nil {
		return f.FindDBR(index, collection, args...).One(result)
	}
	return fn(f.FindDBR(index, collection, args...)).One(result)
}

func (f *Factory) Count(collection string, fn func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return f.CountFromDB(0, collection, fn, args...)
}

func (f *Factory) CountFromDB(index int, collection string, fn func(db.Result) db.Result, args ...interface{}) (int64, error) {
	var cnt uint64
	var err error
	if fn == nil {
		cnt, err = f.FindDBR(index, collection, args...).Count()
	} else {
		cnt, err = fn(f.FindDBR(index, collection, args...)).Count()
	}

	return int64(cnt), err
}

// Write ==========================

func (f *Factory) Insert(collection string, data interface{}) (interface{}, error) {
	return f.InsertToDB(0, collection, data)
}

func (f *Factory) InsertToDB(index int, collection string, data interface{}) (interface{}, error) {
	return f.Collection(collection, index, W).Insert(data)
}

func (f *Factory) Update(collection string, fn func(db.Result) db.Result, data interface{}, args ...interface{}) error {
	return f.UpdateToDB(0, collection, fn, data, args...)
}

func (f *Factory) UpdateToDB(index int, collection string, fn func(db.Result) db.Result, data interface{}, args ...interface{}) error {
	if fn == nil {
		return f.FindDB(index, collection, args...).Update(data)
	}
	return fn(f.FindDB(index, collection, args...)).Update(data)
}

func (f *Factory) Delete(collection string, fn func(db.Result) db.Result, args ...interface{}) error {
	return f.DeleteFromDB(0, collection, fn, args...)
}

func (f *Factory) DeleteFromDB(index int, collection string, fn func(db.Result) db.Result, args ...interface{}) error {
	if fn == nil {
		return f.FindDB(index, collection, args...).Delete()
	}
	return fn(f.FindDB(index, collection, args...)).Delete()
}
