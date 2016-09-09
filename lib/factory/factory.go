// added by swh@admpub.com
package factory

import (
	"errors"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
)

const (
	// R : read mode
	R = iota
	// W : write mode
	W
)

var (
	ErrNotFoundKey = errors.New(`not found the key`)
)

func New() *Factory {
	f := &Factory{
		databases: make([]*Cluster, 0),
	}
	f.Transaction = &Transaction{
		Factory: f,
	}
	return f
}

type Factory struct {
	*Transaction
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

func (f *Factory) Cacher() Cacher {
	return f.cacher
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

func (f *Factory) AddDBToCluster(index int, databases ...sqlbuilder.Database) *Factory {
	if len(f.databases) > index {
		f.databases[index].AddW(databases...)
	} else {
		c := NewCluster()
		c.AddW(databases...)
		f.databases = append(f.databases, c)
	}
	return f
}

func (f *Factory) AddSlaveDBToCluster(index int, databases ...sqlbuilder.Database) *Factory {
	if len(f.databases) > index {
		f.databases[index].AddR(databases...)
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

func (f *Factory) GetCluster(index int) *Cluster {
	return f.Cluster(index)
}

func (f *Factory) Tx(param *Param) error {
	if param.TxMiddleware == nil {
		return nil
	}
	c := f.Cluster(param.Index)
	trans := &Transaction{
		Cluster: c,
		Factory: f,
	}
	fn := func(tx sqlbuilder.Tx) error {
		trans.Tx = tx
		return param.TxMiddleware(trans)
	}
	return c.W().Tx(fn)
}

func (f *Factory) NewTx(args ...int) (trans *Transaction, err error) {
	var index int
	if len(args) > 0 {
		index = args[0]
	}
	c := f.Cluster(index)
	trans = &Transaction{
		Cluster: c,
		Factory: f,
	}
	trans.Tx, err = c.W().NewTx()
	return
}

func (f *Factory) CloseAll() {
	for _, cluster := range f.databases {
		cluster.CloseAll()
	}
}
