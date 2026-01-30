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

// QueryListPage 分页查询
// q: 查询器
// cond: 条件
// mw: 查询构建函数
// paginationVarSuffix: 分页变量后缀
func QueryListPage(q ListQuerier, cond *db.Compounds, mw func(r db.Result) db.Result, paginationVarSuffix ...string) error {
	return QueryListPageAs(q, nil, cond, mw, paginationVarSuffix...)
}

// QueryListPageAs 分页查询
// q: 查询器
// recv: 结果接收器
// cond: 条件
// mw: 查询构建函数
// paginationVarSuffix: 分页变量后缀
func QueryListPageAs(q ListQuerier, recv interface{}, cond *db.Compounds, mw func(r db.Result) db.Result, paginationVarSuffix ...string) error {
	_, err := NewLister(q, recv, mw, cond.And()).Paging(q.Context(), paginationVarSuffix...)
	return err
}

// ListPage 分页查询
// q: 查询器
// cond: 条件
// sorts: 排序方式 或 查询构建函数
//
// 例如：
//
//	ListPage(q, cond, func(r db.Result) db.Result {
//		return r.OrderBy("id", "name")
//	})
//
// 或者：
//
//	ListPage(q, cond, "id", "name")
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

// ListPageAs 分页查询
//
// 例如：
//
//	ListPageAs(q, recv, cond, func(r db.Result) db.Result {
//		return r.OrderBy("id", "name")
//	})
//
// 或者：
//
//	ListPageAs(q, recv, cond, "id", "name")
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

// QueryListPageByOffset 分页查询
// q: 查询器
// cond: 条件
// mw: 查询构建函数
// paginationVarSuffix: 分页变量后缀
func QueryListPageByOffset(q ListOffsetQuerier, cond *db.Compounds, mw func(r db.Result) db.Result, paginationVarSuffix ...string) error {
	return QueryListPageByOffsetAs(q, nil, cond, mw, paginationVarSuffix...)
}

// QueryListPageByOffsetAs 分页查询
// q: 查询器
// recv: 结果接收器
// cond: 条件
// mw: 查询构建函数
// paginationVarSuffix: 分页变量后缀
func QueryListPageByOffsetAs(q ListOffsetQuerier, recv interface{}, cond *db.Compounds, mw func(r db.Result) db.Result, paginationVarSuffix ...string) error {
	_, err := NewOffsetLister(q, recv, mw, cond.And()).Paging(q.Context(), paginationVarSuffix...)
	return err
}

// ListPageByOffset 分页查询
//
// q: 查询器
// cond: 条件
// sorts: 排序字段列表 或 查询构建函数
//
// 例如：
//
//	ListPageByoffset(q, cond, func(r db.Result) db.Result {
//		return r.OrderBy("id", "name")
//	})
//
// 或者：
//
//	ListPageByoffset(q, cond, "id", "name")
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

// ListPageByoffsetAs 分页查询
//
// q: 查询器
// recv: 结果接收器
// cond: 条件
// sorts: 排序字段列表 或 查询构建函数
//
// 例如：
//
//	ListPageByoffsetAs(q, recv, cond, func(r db.Result) db.Result {
//		return r.OrderBy("id", "name")
//	})
//
// 或者：
//
//	ListPageByoffsetAs(q, recv, cond, "id", "name")
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
