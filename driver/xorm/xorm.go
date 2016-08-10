package xorm

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/admpub/mapstruct"
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

func (m *XORM) Delete(bean interface{}, condition CondBuilder, args ...string) error {
	cond := condition.Build().(*Condition)
	_, err := m.Table(bean).Where(cond.SQL, cond.Args...).Delete(bean)
	return err
}

func (m *XORM) Update(bean interface{}, values H, condition CondBuilder, args ...string) error {
	cond := condition.Build().(*Condition)
	if err := mapstruct.Map2Struct(values, bean); err != nil {
		return err
	}
	cols := []string{}
	for key := range values {
		cols = append(cols, key)
	}
	_, err := m.Table(bean).Where(cond.SQL, cond.Args...).MustCols(cols...).Update(bean)
	return err
}

func (m *XORM) Insert(bean interface{}, values H, args ...string) error {
	if err := mapstruct.Map2Struct(values, bean); err != nil {
		return err
	}
	_, err := m.Table(bean).Insert(bean)
	return err
}

func (m *XORM) Upsert(bean interface{}, values H, condition CondBuilder, args ...string) (int, error) {
	cond := condition.Build().(*Condition)
	v := reflect.ValueOf(bean)
	if !v.CanSet() {
		v = v.Elem()
	}
	ve := reflect.New(v.Type())
	newBean := ve.Interface()
	has, err := m.Where(cond.SQL, cond.Args...).Get(newBean)
	if err != nil {
		return 0, err
	}
	if err := mapstruct.Map2Struct(values, bean); err != nil {
		return 0, err
	}
	var affected int64
	if has {
		affected, err = m.Engine.Where(cond.SQL, cond.Args...).Update(bean)
	} else {
		err = m.Engine.Insert(bean)
	}
	return int(affected), err
}

func (m *XORM) AddIndex(bean interface{}, keyInfo interface{}, args ...string) error {
	return m.Engine.CreateIndexes(bean)
}

func (m *XORM) DropIndex(bean interface{}, keyInfo interface{}, args ...string) error {
	return m.Engine.DropIndexes(bean)
}

func (m *XORM) DropTable(bean interface{}, args ...string) error {
	return m.DropTables(bean)
}
