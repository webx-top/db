package factory

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
)

type Transaction struct {
	sqlbuilder.Tx
	*Cluster
	*Factory
}

func (t *Transaction) Backend(param *Param) sqlbuilder.Backend {
	if t.Tx != nil {
		return t.Tx
	}
	if t.Cluster == nil {
		t.Cluster = t.Factory.Cluster(param.Index)
	}
	if param.ReadOrWrite == R {
		return t.Cluster.R()
	}
	return t.Cluster.W()
}

func (t *Transaction) Result(param *Param) db.Result {
	return t.C(param).Find(param.Args...)
}

func (t *Transaction) C(param *Param) db.Collection {
	return t.Backend(param).Collection(t.Cluster.Table(param.Collection))
}

// ================================
// API
// ================================

// Read ==========================

func (t *Transaction) SelectAll(param *Param) error {
	selector := t.Select(param)
	if param.Size > 0 {
		selector = selector.Limit(param.Size).Offset(param.Offset())
	}
	if param.SelectorMiddleware != nil {
		selector = param.SelectorMiddleware(selector)
	}
	return selector.All(param.ResultData)
}

func (t *Transaction) SelectOne(param *Param) error {
	selector := t.Select(param).Limit(1)
	if param.SelectorMiddleware != nil {
		selector = param.SelectorMiddleware(selector)
	}
	return selector.One(param.ResultData)
}

func (t *Transaction) Select(param *Param) sqlbuilder.Selector {
	selector := t.Backend(param).Select(param.Cols...).From(t.Table(param.Collection))
	if param.Joins == nil {
		return selector
	}
	for _, join := range param.Joins {
		coll := t.Table(join.Collection)
		if len(join.Alias) > 0 {
			coll += ` AS ` + join.Alias
		}
		switch strings.ToUpper(join.Type) {
		case "LEFT":
			selector = selector.LeftJoin(coll)
		case "RIGHT":
			selector = selector.RightJoin(coll)
		case "CROSS":
			selector = selector.CrossJoin(coll)
		case "INNER":
			selector = selector.FullJoin(coll)
		default:
			selector = selector.FullJoin(coll)
		}
		if len(join.Condition) > 0 {
			selector = selector.On(join.Condition)
		}
	}
	return selector
}

func (t *Transaction) All(param *Param) error {
	if param.Lifetime > 0 && t.Factory.cacher != nil {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.Lifetime)
	}
	res := t.Result(param)
	if param.Size > 0 {
		res = res.Limit(param.Size).Offset(param.Offset())
	}
	if param.Middleware != nil {
		res = param.Middleware(res)
	}
	return res.All(param.ResultData)
}

func (t *Transaction) List(param *Param) (func() int64, error) {

	if param.Lifetime > 0 && t.Factory.cacher != nil {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return func() int64 {
					return param.Total
				}, nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.Lifetime)
	}

	var res db.Result
	if param.Middleware == nil {
		param.CountFunc = func() int64 {
			if param.Total <= 0 {
				res := t.Result(param)
				count, _ := res.Count()
				param.Total = int64(count)
			}
			return param.Total
		}
		res = t.Result(param).Limit(param.Size).Offset(param.Offset())
	} else {
		param.CountFunc = func() int64 {
			if param.Total <= 0 {
				res := param.Middleware(t.Result(param)).OrderBy()
				count, _ := res.Count()
				param.Total = int64(count)
			}
			return param.Total
		}
		res = param.Middleware(t.Result(param).Limit(param.Size).Offset(param.Offset()))
	}
	return param.CountFunc, res.All(param.ResultData)
}

func (t *Transaction) One(param *Param) error {

	if param.Lifetime > 0 && t.Factory.cacher != nil {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.Lifetime)
	}

	res := t.Result(param)
	if param.Middleware != nil {
		res = param.Middleware(res)
	}
	return res.One(param.ResultData)
}

func (t *Transaction) Count(param *Param) (int64, error) {

	if param.Lifetime > 0 && t.Factory.cacher != nil {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return param.Total, nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.Lifetime)
	}

	var cnt uint64
	var err error

	res := t.Result(param)
	if param.Middleware != nil {
		res = param.Middleware(res)
	}
	cnt, err = res.Count()
	param.Total = int64(cnt)
	return param.Total, err
}

// Write ==========================

func (t *Transaction) Insert(param *Param) (interface{}, error) {
	param.ReadOrWrite = W
	return t.C(param).Insert(param.SaveData)
}

func (t *Transaction) Update(param *Param) error {
	param.ReadOrWrite = W
	res := t.Result(param)
	if param.Middleware != nil {
		res = param.Middleware(res)
	}
	return res.Update(param.SaveData)
}

func (t *Transaction) Delete(param *Param) error {
	param.ReadOrWrite = W
	res := t.Result(param)
	if param.Middleware != nil {
		res = param.Middleware(res)
	}
	return res.Delete()
}
