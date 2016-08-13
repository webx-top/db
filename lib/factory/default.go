package factory

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
)

var DefaultFactory = New()

func Default() *Factory {
	return DefaultFactory
}

func Debug() bool {
	return DefaultFactory.Debug()
}

func SetDebug(on bool) *Factory {
	DefaultFactory.SetDebug(on)
	return DefaultFactory
}

func SetCacher(cc Cacher) *Factory {
	DefaultFactory.SetCacher(cc)
	return DefaultFactory
}

func AddDB(databases ...sqlbuilder.Database) *Factory {
	DefaultFactory.AddDB(databases...)
	return DefaultFactory
}

func AddSlaveDB(databases ...sqlbuilder.Database) *Factory {
	DefaultFactory.AddSlaveDB(databases...)
	return DefaultFactory
}

func SetCluster(index int, cluster *Cluster) *Factory {
	DefaultFactory.SetCluster(index, cluster)
	return DefaultFactory
}

func AddCluster(clusters ...*Cluster) *Factory {
	DefaultFactory.AddCluster(clusters...)
	return DefaultFactory
}

func GetCluster(index int) *Cluster {
	return DefaultFactory.GetCluster(index)
}

func Collection(collection string, args ...int) db.Collection {
	return DefaultFactory.Collection(collection, args...)
}

func Find(collection string, args ...interface{}) db.Result {
	return DefaultFactory.Find(collection, args...)
}

func FindR(collection string, args ...interface{}) db.Result {
	return DefaultFactory.FindR(collection, args...)
}

func FindDB(index int, collection string, args ...interface{}) db.Result {
	return DefaultFactory.FindDB(index, collection, args...)
}

func FindDBR(index int, collection string, args ...interface{}) db.Result {
	return DefaultFactory.FindDBR(index, collection, args...)
}

func CloseAll() {
	DefaultFactory.CloseAll()
}

// ================================
// API
// ================================

// Read ==========================
func All(param *Param) error {
	return DefaultFactory.All(param)
}

func List(param *Param) (func() int64, error) {
	return DefaultFactory.List(param)
}

func One(param *Param) error {
	return DefaultFactory.One(param)
}

func Count(param *Param) (int64, error) {
	return DefaultFactory.Count(param)
}

// Write ==========================

func Insert(param *Param) (interface{}, error) {
	return DefaultFactory.Insert(param)
}

func Update(param *Param) error {
	return DefaultFactory.Update(param)
}

func Delete(param *Param) error {
	return DefaultFactory.Delete(param)
}
