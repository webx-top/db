package mongo

import (
	"fmt"
	"log"
	"time"

	. "github.com/webx-top/dbx/driver"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func New() *Mongo {
	return &Mongo{}
}

type Mongo struct {
	*mgo.Session
	defaultDB string
}

func (m *Mongo) Connect(connectInfo interface{}, defaultDB string, timeout ...time.Duration) (err error) {
	m.defaultDB = defaultDB
	if dsn, ok := connectInfo.(string); ok {
		m.Session, err = mgo.Dial(dsn)
	} else if info, ok := connectInfo.(*mgo.DialInfo); ok {
		m.Session, err = mgo.DialWithInfo(info)
	} else {
		return fmt.Errorf(`param error: The connectInfo data type is not supported: %T`, connectInfo)
	}
	if err != nil {
		return fmt.Errorf(`Can not connect to MongoDB server: %v`, err)
	}
	if len(timeout) > 0 {
		m.Session.SetSocketTimeout(timeout[0])
	} else {
		m.Session.SetSocketTimeout(time.Hour * 10)
	}
	return nil
}

func (m *Mongo) Query(collection string, condition interface{}, args ...string) *mgo.Query {
	return m.Collection(collection, args...).Find(condition)
}

func (m *Mongo) Collection(collection string, args ...string) *mgo.Collection {
	dbName := m.defaultDB
	if len(args) > 0 {
		dbName = args[0]
	}
	err := m.Session.Ping()
	if err != nil {
		log.Println("Lost connection to MongoDB!")
		m.Session.Refresh()
		err = m.Session.Ping()
		if err == nil {
			log.Println("Reconnect to MongoDB successful.")
		}
	}
	return m.Session.DB(dbName).C(collection)
}

func (m *Mongo) All(sel Selecter, result interface{}, args ...string) error {
	skip, limit := sel.Limit()
	sort := sel.Sort()
	sess := m.Query(sel.Table(), sel.Condition().Build(), args...).Limit(limit).Skip(skip)
	if len(sort) > 0 {
		sess.Sort(sort...)
	}
	return sess.All(result)
}

func (m *Mongo) Count(sel Selecter, _ interface{}, args ...string) (int, error) {
	skip, limit := sel.Limit()
	sort := sel.Sort()
	sess := m.Query(sel.Table(), sel.Condition().Build(), args...).Limit(limit).Skip(skip)
	if len(sort) > 0 {
		sess.Sort(sort...)
	}
	return sess.Count()
}

func (m *Mongo) One(sel Selecter, result interface{}, args ...string) error {
	skip, _ := sel.Limit()
	sort := sel.Sort()
	sess := m.Query(sel.Table(), sel.Condition(), args...).Limit(1).Skip(skip)
	if len(sort) > 0 {
		sess.Sort(sort...)
	}
	return sess.One(result)
}

func (m *Mongo) IsExists(err error) bool {
	return err == nil || err != mgo.ErrNotFound
}

func (m *Mongo) Delete(collection string, condition CondBuilder, args ...string) error {
	return m.Collection(collection, args...).Remove(condition.Build())
}

func (m *Mongo) Update(collection string, values H, condition CondBuilder, args ...string) error {
	return m.Collection(collection, args...).Update(condition.Build(), bson.M(values))
}

func (m *Mongo) Insert(collection string, values H, args ...string) error {
	return m.Collection(collection, args...).Insert(bson.M(values))
}

func (m *Mongo) Upsert(collection string, values H, condition CondBuilder, args ...string) (int, error) {
	info, err := m.Collection(collection, args...).Upsert(condition.Build(), bson.M(values))
	if err != nil {
		return 0, err
	}
	return info.Updated, err
}

func (m *Mongo) AddIndex(collection string, keyInfo interface{}, args ...string) error {
	if key, ok := keyInfo.([]string); ok {
		return m.Collection(collection, args...).EnsureIndexKey(key...)
	}
	if key, ok := keyInfo.(string); ok {
		return m.Collection(collection, args...).EnsureIndexKey(key)
	}
	return nil
}

func (m *Mongo) DropIndex(collection string, keyInfo interface{}, args ...string) error {
	if key, ok := keyInfo.([]string); ok {
		return m.Collection(collection, args...).DropIndex(key...)
	}
	if key, ok := keyInfo.(string); ok {
		return m.Collection(collection, args...).DropIndex(key)
	}
	return nil
}

func (m *Mongo) DropTable(collection string, args ...string) error {
	return m.Collection(collection, args...).DropCollection()
}

/*
TODO
func (m *Mongo) AddTable(collection string, args ...string) error {
	var collectionInf *mgo.CollectionInfo
	return m.Collection(collection, args...).Create(collectionInf)
}
*/
