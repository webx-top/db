// added by swh@admpub.com
package factory

import (
	"errors"

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

type Factory struct {
	databases []*Cluster
	cacher    Cacher
}

func (f *Factory) Debug() bool {
	return db.Debug
}

func (f *Factory) SetDebug(on bool) *Factory {
	db.Debug = on
	return f
}

func (f *Factory) SetCacher(cc Cacher) *Factory {
	f.cacher = cc
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
func (f *Factory) All(param *Param) error {
	if param.Lifetime > 0 && f.cacher != nil {
		data, err := f.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = f
				return nil
			}
		}
		defer f.cacher.Put(param.CachedKey(), param, param.Lifetime)
	}
	if param.Middleware == nil {
		return f.FindDBR(param.Index, param.Collection, param.Args...).All(param.Result)
	}
	return param.Middleware(f.FindDBR(param.Index, param.Collection, param.Args...)).All(param.Result)
}

func (f *Factory) PageList(param *Param) (func() int64, error) {

	if param.Lifetime > 0 && f.cacher != nil {
		data, err := f.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = f
				return func() int64 {
					return param.Total
				}, nil
			}
		}
		defer f.cacher.Put(param.CachedKey(), param, param.Lifetime)
	}

	if param.Middleware == nil {
		param.CountFunc = func() int64 {
			if param.Total <= 0 {
				count, _ := f.FindDBR(param.Index, param.Collection, param.Args...).Count()
				param.Total = int64(count)
			}
			return param.Total
		}
		return param.CountFunc, f.FindDBR(param.Index, param.Collection, param.Args...).Limit(param.Size).Offset(param.Offset()).All(param.Result)
	}
	param.CountFunc = func() int64 {
		if param.Total <= 0 {
			count, _ := param.Middleware(f.FindDBR(param.Index, param.Collection, param.Args...)).Count()
			param.Total = int64(count)
		}
		return param.Total
	}
	return param.CountFunc, param.Middleware(f.FindDBR(param.Index, param.Collection, param.Args...).Limit(param.Size).Offset(param.Offset())).All(param.Result)
}

func (f *Factory) One(param *Param) error {

	if param.Lifetime > 0 && f.cacher != nil {
		data, err := f.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = f
				return nil
			}
		}
		defer f.cacher.Put(param.CachedKey(), param, param.Lifetime)
	}

	if param.Middleware == nil {
		return f.FindDBR(param.Index, param.Collection, param.Args...).One(param.Result)
	}
	return param.Middleware(f.FindDBR(param.Index, param.Collection, param.Args...)).One(param.Result)
}

func (f *Factory) Count(param *Param) (int64, error) {

	if param.Lifetime > 0 && f.cacher != nil {
		data, err := f.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = f
				return param.Total, nil
			}
		}
		defer f.cacher.Put(param.CachedKey(), param, param.Lifetime)
	}

	var cnt uint64
	var err error
	if param.Middleware == nil {
		cnt, err = f.FindDBR(param.Index, param.Collection, param.Args...).Count()
	} else {
		cnt, err = param.Middleware(f.FindDBR(param.Index, param.Collection, param.Args...)).Count()
	}
	param.Total = int64(cnt)
	return param.Total, err
}

// Write ==========================

func (f *Factory) Insert(param *Param) (interface{}, error) {
	return f.Collection(param.Collection, param.Index, W).Insert(param.SaveData)
}

func (f *Factory) Update(param *Param) error {
	if param.Middleware == nil {
		return f.FindDB(param.Index, param.Collection, param.Args...).Update(param.SaveData)
	}
	return param.Middleware(f.FindDB(param.Index, param.Collection, param.Args...)).Update(param.SaveData)
}

func (f *Factory) Delete(param *Param) error {
	if param.Middleware == nil {
		return f.FindDB(param.Index, param.Collection, param.Args...).Delete()
	}
	return param.Middleware(f.FindDB(param.Index, param.Collection, param.Args...)).Delete()
}
