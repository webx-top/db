// added by swh@admpub.com
package factory

import (
	"encoding/gob"
	"fmt"
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
)

func init() {
	gob.Register(&Param{})
}

type Join struct {
	Collection string
	Alias      string
	Condition  string
	Type       string
}

type Param struct {
	factory            *Factory
	Index              int //数据库对象元素所在的索引位置
	ReadOrWrite        int
	Collection         string //集合名或表名称
	Middleware         func(db.Result) db.Result
	SelectorMiddleware func(sqlbuilder.Selector) sqlbuilder.Selector
	TxMiddleware       func(*Transaction) error
	CountFunc          func() int64
	ResultData         interface{}   //查询后保存的结果
	Args               []interface{} //Find方法的条件参数
	Cols               []interface{} //使用Selector要查询的列
	Joins              []*Join
	SaveData           interface{} //增加和更改数据时要保存到数据库中的数据
	Offset             int
	Page               int           //页码
	Size               int           //每页数据量
	Total              int64         //数据表中符合条件的数据行数
	MaxAge             time.Duration //缓存有效时间（单位：秒），为0时代表临时关闭缓存，为-1时代表删除缓存
	trans              *Transaction
	cachedKey          string
	setter             *Setting
}

func NewParam(args ...*Factory) *Param {
	p := &Param{
		factory: DefaultFactory,
		Args:    make([]interface{}, 0),
		Cols:    make([]interface{}, 0),
		Joins:   make([]*Join, 0),
		Page:    1,
		Offset:  -1,
	}
	p.init(args...)
	return p
}

func (p *Param) init(args ...*Factory) *Param {
	if len(args) > 0 {
		p.factory = args[0]
	}
	//p.setter = &Setting{Param: p}
	return p
}

func (p *Param) Setter() *Setting {
	if p.setter == nil {
		p.setter = &Setting{Param: p}
	}
	return p.setter
}

func (p *Param) SetIndex(index int) *Param {
	p.Index = index
	return p
}

func (p *Param) SelectLink(index int) *Param {
	p.Index = index
	return p
}

func (p *Param) CachedKey() string {
	if len(p.cachedKey) == 0 {
		p.cachedKey = fmt.Sprintf(`%v-%v-%v-%v-%v-%v`, p.Index, p.Collection, p.Args, p.Offset, p.Page, p.Size)
	}
	return p.cachedKey
}

func (p *Param) SetCache(maxAge time.Duration, key ...string) *Param {
	p.MaxAge = maxAge
	if len(key) > 0 {
		p.cachedKey = key[0]
	}
	return p
}

func (p *Param) SetCachedKey(key string) *Param {
	p.cachedKey = key
	return p
}

func (p *Param) SetJoin(joins ...*Join) *Param {
	p.Joins = joins
	return p
}

func (p *Param) SetTx(tx sqlbuilder.Tx) *Param {
	p.trans = &Transaction{
		Tx:      tx,
		Factory: p.factory,
	}
	return p
}

func (p *Param) SetTrans(trans *Transaction) *Param {
	p.trans = trans
	return p
}

func (p *Param) SetRead() *Param {
	p.ReadOrWrite = R
	return p
}

func (p *Param) SetWrite() *Param {
	p.ReadOrWrite = W
	return p
}

func (p *Param) AddJoin(joinType string, collection string, alias string, condition string) *Param {
	p.Joins = append(p.Joins, &Join{
		Collection: collection,
		Alias:      alias,
		Condition:  condition,
		Type:       joinType,
	})
	return p
}

func (p *Param) SetCollection(collection string) *Param {
	p.Collection = collection
	return p
}

func (p *Param) SetMiddleware(middleware func(db.Result) db.Result) *Param {
	p.Middleware = middleware
	return p
}

func (p *Param) SetSelectorMiddleware(middleware func(sqlbuilder.Selector) sqlbuilder.Selector) *Param {
	p.SelectorMiddleware = middleware
	return p
}

// SetMW is SetMiddleware's alias.
func (p *Param) SetMW(middleware func(db.Result) db.Result) *Param {
	p.SetMiddleware(middleware)
	return p
}

func (p *Param) SetTxMiddleware(middleware func(*Transaction) error) *Param {
	p.TxMiddleware = middleware
	return p
}

func (p *Param) SetTxMW(middleware func(*Transaction) error) *Param {
	p.SetTxMiddleware(middleware)
	return p
}

// SetSelMW is SetSelectorMiddleware's alias.
func (p *Param) SetSelMW(middleware func(sqlbuilder.Selector) sqlbuilder.Selector) *Param {
	p.SetSelectorMiddleware(middleware)
	return p
}

func (p *Param) SetRecv(result interface{}) *Param {
	p.ResultData = result
	return p
}

func (p *Param) SetArgs(args ...interface{}) *Param {
	p.Args = args
	return p
}

func (p *Param) AddArgs(args ...interface{}) *Param {
	p.Args = append(p.Args, args...)
	return p
}

func (p *Param) SetCols(args ...interface{}) *Param {
	p.Cols = args
	return p
}

func (p *Param) AddCols(args ...interface{}) *Param {
	p.Cols = append(p.Cols, args...)
	return p
}

func (p *Param) SetSend(save interface{}) *Param {
	p.SaveData = save
	return p
}

func (p *Param) SetPage(n int) *Param {
	if n < 1 {
		p.Page = 1
	} else {
		p.Page = n
	}
	return p
}

func (p *Param) SetOffset(offset int) *Param {
	p.Offset = offset
	return p
}

func (p *Param) SetSize(size int) *Param {
	p.Size = size
	return p
}

func (p *Param) SetTotal(total int64) *Param {
	p.Total = total
	return p
}

func (p *Param) Trans() *Transaction {
	return p.trans
}

func (p *Param) TransTo(param *Param) *Param {
	param.trans = p.trans
	return p
}

func (p *Param) TransFrom(param *Param) *Param {
	p.trans = param.trans
	return p
}

func (p *Param) Begin() *Param {
	p.trans = p.MustTx()
	return p
}

func (p *Param) End(err error) error {
	if p.trans == nil || p.trans.Tx == nil {
		return nil
	}
	defer p.SetTrans(nil)
	if err == nil {
		return p.trans.Commit()
	}
	return p.trans.Rollback()
}

func (p *Param) GetOffset() int {
	if p.Offset > -1 {
		return p.Offset
	}
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.Size
}

func (p *Param) NewTx() (*Transaction, error) {
	return p.factory.NewTx(p.Index)
}

func (p *Param) Tx() (*Transaction, error) {
	if p.trans != nil {
		return p.trans, nil
	}
	var err error
	p.trans, err = p.NewTx()
	return p.trans, err
}

func (p *Param) MustTx() *Transaction {
	trans, err := p.Tx()
	if err != nil {
		panic(err.Error())
	}
	return trans
}

func (p *Param) T() *Transaction {
	if p.trans != nil {
		return p.trans
	}
	return p.factory.Transaction
}

func (p *Param) Result() db.Result {
	return p.T().Result(p)
}

// Read ==========================

func (p *Param) SelectAll() error {
	return p.T().SelectAll(p)
}

func (p *Param) SelectOne() error {
	return p.T().SelectOne(p)
}

func (p *Param) Select() sqlbuilder.Selector {
	return p.T().Select(p)
}

func (p *Param) All() error {
	return p.T().All(p)
}

func (p *Param) List() (func() int64, error) {
	return p.T().List(p)
}

func (p *Param) One() error {
	return p.T().One(p)
}

func (p *Param) Count() (int64, error) {
	return p.T().Count(p)
}

// Write ==========================

func (p *Param) Insert() (interface{}, error) {
	return p.T().Insert(p)
}

func (p *Param) Update() error {
	return p.T().Update(p)
}

func (p *Param) Delete() error {
	return p.T().Delete(p)
}
