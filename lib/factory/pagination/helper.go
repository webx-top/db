package pagination

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

type ListQuerier interface {
	Lister
	Context() echo.Context
}

type ListOffsetQuerier interface {
	OffsetLister
	Context() echo.Context
}

func QueryListPage(q ListQuerier, cond *db.Compounds, mw func(r db.Result) db.Result, paginationVarSuffix ...string) error {
	return QueryListPageAs(q, nil, cond, mw, paginationVarSuffix...)
}

func QueryListPageAs(q ListQuerier, recv interface{}, cond *db.Compounds, mw func(r db.Result) db.Result, paginationVarSuffix ...string) error {
	_, err := NewLister(q, recv, mw, cond.And()).Paging(q.Context(), paginationVarSuffix...)
	return err
}

func ListPage(q ListQuerier, cond *db.Compounds, sorts ...interface{}) error {
	length := len(sorts)
	if length > 0 {
		if mw, ok := sorts[0].(func(r db.Result) db.Result); ok {
			if length > 1 {
				if paginationVarSuffix, ok := sorts[1].(string); ok {
					return QueryListPage(q, cond, mw, paginationVarSuffix)
				}
			}
			return QueryListPage(q, cond, mw)
		}
	}
	return QueryListPage(q, cond, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	})
}

func ListPageAs(q ListQuerier, recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	length := len(sorts)
	if length > 0 {
		if mw, ok := sorts[0].(func(r db.Result) db.Result); ok {
			if length > 1 {
				if paginationVarSuffix, ok := sorts[1].(string); ok {
					return QueryListPageAs(q, recv, cond, mw, paginationVarSuffix)
				}
			}
			return QueryListPageAs(q, recv, cond, mw)
		}
	}
	return QueryListPageAs(q, recv, cond, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	})
}

func QueryListPageByOffset(q ListOffsetQuerier, cond *db.Compounds, mw func(r db.Result) db.Result, paginationVarSuffix ...string) error {
	return QueryListPageByOffsetAs(q, nil, cond, mw, paginationVarSuffix...)
}

func QueryListPageByOffsetAs(q ListOffsetQuerier, recv interface{}, cond *db.Compounds, mw func(r db.Result) db.Result, paginationVarSuffix ...string) error {
	_, err := NewOffsetLister(q, recv, mw, cond.And()).Paging(q.Context(), paginationVarSuffix...)
	return err
}

func ListPageByOffset(q ListOffsetQuerier, cond *db.Compounds, sorts ...interface{}) error {
	length := len(sorts)
	if length > 0 {
		if mw, ok := sorts[0].(func(r db.Result) db.Result); ok {
			if length > 1 {
				if paginationVarSuffix, ok := sorts[1].(string); ok {
					return QueryListPageByOffset(q, cond, mw, paginationVarSuffix)
				}
			}
			return QueryListPageByOffset(q, cond, mw)
		}
	}
	return QueryListPageByOffset(q, cond, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	})
}

func ListPageByOffsetAs(q ListOffsetQuerier, recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	length := len(sorts)
	if length > 0 {
		if mw, ok := sorts[0].(func(r db.Result) db.Result); ok {
			if length > 1 {
				if paginationVarSuffix, ok := sorts[1].(string); ok {
					return QueryListPageByOffsetAs(q, recv, cond, mw, paginationVarSuffix)
				}
			}
			return QueryListPageByOffsetAs(q, recv, cond, mw)
		}
	}
	return QueryListPageByOffsetAs(q, recv, cond, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	})
}
