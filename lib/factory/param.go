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
	CountFunc          func() int64
	ResultData         interface{}   //查询后保存的结果
	Args               []interface{} //Find方法的条件参数
	Cols               []interface{} //使用Selector要查询的列
	Joins              []*Join
	SaveData           interface{}   //增加和更改数据时要保存到数据库中的数据
	Page               int           //页码
	Size               int           //每页数据量
	Total              int64         //数据表中符合条件的数据行数
	Lifetime           time.Duration //缓存生存时间
	cachedKey          string
	offset             int
}

func NewParam(args ...*Factory) *Param {
	p := &Param{
		factory: DefaultFactory,
		Args:    make([]interface{}, 0),
		Cols:    make([]interface{}, 0),
		Joins:   make([]*Join, 0),
		Page:    1,
		offset:  -1,
	}
	if len(args) > 0 {
		p.factory = args[0]
	}
	return p
}

func (p *Param) SetIndex(index int) *Param {
	p.Index = index
	return p
}

func (p *Param) CachedKey() string {
	if len(p.cachedKey) == 0 {
		p.cachedKey = fmt.Sprintf(`%v-%v-%v-%v-%v`, p.Index, p.Collection, p.Args, p.Page, p.Size)
	}
	return p.cachedKey
}

func (p *Param) SetCachedKey(key string) *Param {
	p.cachedKey = key
	return p
}

func (p *Param) SetJoin(joins ...*Join) *Param {
	p.Joins = joins
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

// SetSelMW is SetSelectorMiddleware's alias.
func (p *Param) SetSelMW(middleware func(sqlbuilder.Selector) sqlbuilder.Selector) *Param {
	p.SetSelectorMiddleware(middleware)
	return p
}

func (p *Param) SetResult(result interface{}) *Param {
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

func (p *Param) SetSave(save interface{}) *Param {
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
	p.offset = offset
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

func (p *Param) Offset() int {
	if p.offset > -1 {
		return p.offset
	}
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.Size
}

func (p *Param) Result() db.Result {
	return p.factory.Result(p)
}

// Read ==========================

func (p *Param) SelectAll() error {
	return p.factory.SelectAll(p)
}

func (p *Param) SelectOne() error {
	return p.factory.SelectOne(p)
}

func (p *Param) Select() sqlbuilder.Selector {
	return p.factory.Select(p)
}

func (p *Param) All() error {
	return p.factory.All(p)
}

func (p *Param) List() (func() int64, error) {
	return p.factory.List(p)
}

func (p *Param) One() error {
	return p.factory.One(p)
}

func (p *Param) Count() (int64, error) {
	return p.factory.Count(p)
}

// Write ==========================

func (p *Param) Insert() (interface{}, error) {
	return p.factory.Insert(p)
}

func (p *Param) Update() error {
	return p.factory.Update(p)
}

func (p *Param) Delete() error {
	return p.factory.Delete(p)
}
