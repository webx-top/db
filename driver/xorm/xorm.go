package xorm

import (
	"errors"
	"fmt"
	"log"
	"time"

	. "github.com/coscms/xorm"
	. "github.com/webx-top/dbx/driver"
)

type Condition struct {
	SQL  string
	Args []interface{}
}

func New() *XORM {
	return &XORM{}
}

type XORM struct {
	*Engine
	defaultDB string
}

func (m *XORM) Connect(connectInfo interface{}, defaultDB string, timeout ...time.Duration) (err error) {
	m.defaultDB = defaultDB
	var engine, dsnString string = `mysql`, ``

	if dsn, ok := connectInfo.(string); ok {
		dsnString = dsn
	} else if info, ok := connectInfo.([]string); ok {
		switch len(info) {
		case 2:
			engine = info[0]
			dsnString = info[1]
		case 1:
			dsnString = info[0]
		}
	} else {
		return fmt.Errorf(`param error: The connectInfo data type is not supported: %T`, connectInfo)
	}

	m.Engine, err = NewEngine(engine, dsnString)
	if err != nil {
		log.Println("The database connection failed:", err)
	}
	err = m.Engine.Ping()
	if err != nil {
		log.Println("The database ping failed:", err)
	}
	m.Engine.CloseLog()
	m.Engine.OpenLog("base")
	return nil
}

func (m *XORM) GenSession(sel Selecter) *Session {
	cond := sel.Condition().Build().(*Condition)
	sess := m.Table(sel.Table()).Where(cond.SQL, cond.Args...)
	groups := sel.Group()
	if len(groups) > 0 {
		for _, group := range groups {
			sess = sess.GroupBy(group)
		}
	}

	cond = sel.Having().Build().(*Condition)
	sess = sess.Having(cond.SQL)
	return sess
}

func (m *XORM) All(sel Selecter, result interface{}, args ...string) error {
	skip, limit := sel.Limit()
	sess := m.GenSession(sel).Limit(limit, skip)
	sort := sel.Sort()
	if len(sort) > 0 {
		for _, sortField := range sort {
			if sortField[0] == '-' {
				sess = sess.Desc(sortField[1:])
			} else {
				sess = sess.Asc(sortField)
			}
		}
	}
	return sess.Find(result)
}

func (m *XORM) Count(sel Selecter, bean interface{}, args ...string) (int, error) {
	n, err := m.GenSession(sel).Count(bean)
	return int(n), err
}

func (m *XORM) One(sel Selecter, result interface{}, args ...string) error {
	sort := sel.Sort()
	sess := m.GenSession(sel).Limit(1, 0)
	if len(sort) > 0 {
		for _, sortField := range sort {
			if sortField[0] == '-' {
				sess = sess.Desc(sortField[1:])
			} else {
				sess = sess.Asc(sortField)
			}
		}
	}
	has, err := sess.Get(result)
	if err != nil {
		return err
	}
	if !has {
		return ErrNotExist
	}
	return nil
}

func (m *XORM) IsExists(err error) bool {
	return err == nil || err != ErrNotExist
}

func (m *XORM) Delete(collection string, condition CondBuilder, args ...string) error {
	cond := condition.Build().(*Condition)
	_, err := m.Table(collection).Where(cond.SQL, cond.Args...).Delete(nil)
	return err
}

func (m *XORM) Update(collection string, values H, condition CondBuilder, args ...string) error {
	return m.Collection(collection, args...).Update(condition.Build(), bson.M(values))
}

func (m *XORM) Insert(collection string, values H, args ...string) error {
	return m.Collection(collection, args...).Insert(bson.M(values))
}

func (m *XORM) Upsert(collection string, values H, condition CondBuilder, args ...string) (int, error) {
	info, err := m.Collection(collection, args...).Upsert(condition.Build(), bson.M(values))
	if err != nil {
		return 0, err
	}
	return info.Updated, err
}

func (m *XORM) AddIndex(collection string, keyInfo interface{}, args ...string) error {
	if key, ok := keyInfo.([]string); ok {
		return m.Collection(collection, args...).EnsureIndexKey(key...)
	}
	if key, ok := keyInfo.(string); ok {
		return m.Collection(collection, args...).EnsureIndexKey(key)
	}
	return nil
}

func (m *XORM) DropIndex(collection string, keyInfo interface{}, args ...string) error {
	if key, ok := keyInfo.([]string); ok {
		return m.Collection(collection, args...).DropIndex(key...)
	}
	if key, ok := keyInfo.(string); ok {
		return m.Collection(collection, args...).DropIndex(key)
	}
	return nil
}

func (m *XORM) DropTable(collection string, args ...string) error {
	return m.Collection(collection, args...).DropCollection()
}
