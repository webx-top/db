package mongo

import (
	"fmt"
	"log"
	"reflect"
	"time"

	. "github.com/webx-top/dbx/driver"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ Driver = New()

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

func (m *Mongo) Count(sel Selecter, bean interface{}, args ...string) (int, error) {
	skip, limit := sel.Limit()
	sort := sel.Sort()
	collection := sel.Table()
	if len(collection) == 0 {
		collection = m.CollectionName(bean)
	}
	sess := m.Query(collection, sel.Condition().Build(), args...).Limit(limit).Skip(skip)
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

func (m *Mongo) Delete(bean interface{}, condition CondBuilder, args ...string) error {
	collection := m.CollectionName(bean)
	return m.Collection(collection, args...).Remove(condition.Build())
}

func (m *Mongo) Update(bean interface{}, values H, condition CondBuilder, args ...string) error {
	collection := m.CollectionName(bean)
	return m.Collection(collection, args...).Update(condition.Build(), bson.M(values.Map()))
}

func (m *Mongo) Insert(bean interface{}, values H, args ...string) error {
	collection := m.CollectionName(bean)
	return m.Collection(collection, args...).Insert(bson.M(values.Map()))
}

func (m *Mongo) Upsert(bean interface{}, values H, condition CondBuilder, args ...string) (int, error) {
	collection := m.CollectionName(bean)
	info, err := m.Collection(collection, args...).Upsert(condition.Build(), bson.M(values.Map()))
	if err != nil {
		return 0, err
	}
	return info.Updated, err
}

func (m *Mongo) AddIndex(bean interface{}, keyInfo interface{}, args ...string) error {
	collection := m.CollectionName(bean)
	if key, ok := keyInfo.([]string); ok {
		return m.Collection(collection, args...).EnsureIndexKey(key...)
	}
	if key, ok := keyInfo.(string); ok {
		return m.Collection(collection, args...).EnsureIndexKey(key)
	}
	return nil
}

func (m *Mongo) DropIndex(bean interface{}, keyInfo interface{}, args ...string) error {
	collection := m.CollectionName(bean)
	if key, ok := keyInfo.([]string); ok {
		return m.Collection(collection, args...).DropIndex(key...)
	}
	if key, ok := keyInfo.(string); ok {
		return m.Collection(collection, args...).DropIndex(key)
	}
	return nil
}

func (m *Mongo) DropTable(bean interface{}, args ...string) error {
	collection := m.CollectionName(bean)
	return m.Collection(collection, args...).DropCollection()
}

/*
TODO
func (m *Mongo) AddTable(collection string, args ...string) error {
	var collectionInf *mgo.CollectionInfo
	return m.Collection(collection, args...).Create(collectionInf)
}
*/

func (m *Mongo) CollectionName(bean interface{}) string {
	if v, ok := bean.(string); ok {
		return v
	}
	v := reflect.ValueOf(bean)
	if !v.CanSet() {
		v = v.Elem()
	}
	name := v.Type().Name()
	return snakeCasedName(name)
}

func snakeCasedName(name string) string {
	newstr := make([]rune, 0)
	for idx, chr := range name {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if idx > 0 {
				newstr = append(newstr, '_')
			}
			chr -= ('A' - 'a')
		}
		newstr = append(newstr, chr)
	}

	return string(newstr)
}
