// added by swh@admpub.com
package factory

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
)

type Setting struct {
	*Param
}

func (s *Setting) CachedKey(key string) *Setting {
	s.Param.SetCachedKey(key)
	return s
}

func (s *Setting) Join(joins ...*Join) *Setting {
	s.Param.SetJoin(joins...)
	return s
}

func (s *Setting) Tx(tx sqlbuilder.Tx) *Setting {
	s.Param.SetTx(tx)
	return s
}

func (s *Setting) Trans(trans *Transaction) *Setting {
	s.Param.SetTrans(trans)
	return s
}

func (s *Setting) R() *Setting {
	s.Param.SetRead()
	return s
}

func (s *Setting) W() *Setting {
	s.Param.SetWrite()
	return s
}

func (s *Setting) AddJoin(joinType string, collection string, alias string, condition string) *Setting {
	s.Param.AddJoin(joinType, collection, alias, condition)
	return s
}

func (s *Setting) Collection(collection string) *Setting {
	s.Param.SetCollection(collection)
	return s
}

func (s *Setting) C(collection string) *Setting {
	s.Collection(collection)
	return s
}

func (s *Setting) Middleware(middleware func(db.Result) db.Result) *Setting {
	s.Param.SetMiddleware(middleware)
	return s
}

func (s *Setting) SelectorMiddleware(middleware func(sqlbuilder.Selector) sqlbuilder.Selector) *Setting {
	s.Param.SetSelectorMiddleware(middleware)
	return s
}

// SetMW is SetMiddleware's alias.
func (s *Setting) MW(middleware func(db.Result) db.Result) *Setting {
	s.Param.SetMW(middleware)
	return s
}

func (s *Setting) TxMiddleware(middleware func(*Transaction) error) *Setting {
	s.Param.SetTxMiddleware(middleware)
	return s
}

func (s *Setting) TxMW(middleware func(*Transaction) error) *Setting {
	s.Param.SetTxMW(middleware)
	return s
}

// SetSelMW is SetSelectorMiddleware's alias.
func (s *Setting) SelMW(middleware func(sqlbuilder.Selector) sqlbuilder.Selector) *Setting {
	s.Param.SetSelMW(middleware)
	return s
}

func (s *Setting) Recv(result interface{}) *Setting {
	s.Param.SetResult(result)
	return s
}

func (s *Setting) Args(args ...interface{}) *Setting {
	s.Param.SetArgs(args)
	return s
}

func (s *Setting) AddArgs(args ...interface{}) *Setting {
	s.Param.AddArgs(args...)
	return s
}

func (s *Setting) Cols(args ...interface{}) *Setting {
	s.Param.SetCols(args)
	return s
}

func (s *Setting) AddCols(args ...interface{}) *Setting {
	s.Param.AddCols(args...)
	return s
}

func (s *Setting) SaveAs(save interface{}) *Setting {
	s.Param.SetSave(save)
	return s
}

func (s *Setting) Page(n int) *Setting {
	s.Param.SetPage(n)
	return s
}

func (s *Setting) Offset(offset int) *Setting {
	s.Param.SetOffset(offset)
	return s
}

func (s *Setting) Size(size int) *Setting {
	s.Param.SetSize(size)
	return s
}

func (s *Setting) Total(total int64) *Setting {
	s.Param.SetTotal(total)
	return s
}

func (s *Setting) TransTo(param *Param) *Setting {
	s.Param.TransTo(param)
	return s
}

func (s *Setting) TransFrom(param *Param) *Setting {
	s.Param.TransFrom(param)
	return s
}

func (s *Setting) Begin() *Setting {
	s.Param.Begin()
	return s
}
