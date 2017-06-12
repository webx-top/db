package sqlbuilder

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"reflect"
	"strings"

	"upper.io/db.v3"
	"upper.io/db.v3/internal/immutable"
)

var (
	errZeroPageSize        = errors.New("Illegal page size (cannot be zero)")
	errMissingCursorColumn = errors.New("Missing cursor column")
)

type paginatorQuery struct {
	sel Selector

	cursorColumn       string
	cursorValue        interface{}
	cursorCond         db.Cond
	cursorReverseOrder bool

	pageSize   uint
	pageNumber uint
}

func newPaginator(sel Selector, pageSize uint) Paginator {
	pag := &paginator{}
	return pag.frame(func(pq *paginatorQuery) error {
		if pageSize < 1 {
			return errZeroPageSize
		}
		pq.pageSize = pageSize
		pq.sel = sel.Limit(int(pq.pageSize))
		return nil
	}).Page(0)
}

func (pq *paginatorQuery) count() (uint, error) {
	var count int
	row, err := pq.sel.(*selector).setColumns(db.Raw("count(1)")).QueryRow()

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return uint(count), nil
}

type paginator struct {
	fn   func(*paginatorQuery) error
	prev *paginator
}

var _ = immutable.Immutable(&paginator{})

func (pag *paginator) frame(fn func(*paginatorQuery) error) *paginator {
	return &paginator{prev: pag, fn: fn}
}

func (pag *paginator) Page(pageNumber uint) Paginator {
	return pag.frame(func(pq *paginatorQuery) error {
		pq.pageNumber = pageNumber
		pq.sel = pq.sel.Offset(int(pq.pageSize * pq.pageNumber))
		return nil
	})
}

func (pag *paginator) Cursor(column string) Paginator {
	return pag.frame(func(pq *paginatorQuery) error {
		pq.cursorColumn = column
		pq.cursorValue = nil
		return nil
	})
}

func (pag *paginator) NextPage(cursorValue interface{}) Paginator {
	return pag.frame(func(pq *paginatorQuery) error {
		if pq.cursorValue != nil && pq.cursorColumn == "" {
			return errMissingCursorColumn
		}
		pq.cursorValue = cursorValue
		pq.cursorReverseOrder = false
		pq.cursorCond = db.Cond{
			pq.cursorColumn + " >": cursorValue,
		}
		return nil
	})
}

func (pag *paginator) PrevPage(cursorValue interface{}) Paginator {
	return pag.frame(func(pq *paginatorQuery) error {
		if pq.cursorValue != nil && pq.cursorColumn == "" {
			return errMissingCursorColumn
		}
		pq.cursorValue = cursorValue
		pq.cursorReverseOrder = true
		pq.cursorCond = db.Cond{
			pq.cursorColumn + " <": cursorValue,
		}
		return nil
	})
}

func (pag *paginator) TotalPages() (uint, error) {
	pq, err := pag.build()
	if err != nil {
		return 0, err
	}

	count, err := pq.count()
	if err != nil {
		return 0, err
	}

	if pq.pageSize == 0 {
		return 0, errZeroPageSize
	}

	pages := uint(math.Ceil(float64(count) / float64(pq.pageSize)))

	return pages, nil
}

func (pag *paginator) All(dest interface{}) error {
	pq, err := pag.buildWithCursor()
	if err != nil {
		return err
	}
	err = pq.sel.All(dest)
	if err != nil {
		return err
	}
	if pq.cursorReverseOrder {
		val := reflect.ValueOf(dest)
		size := val.Elem().Len()
		swap := reflect.Swapper(val.Elem().Interface())
		for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
			swap(i, j)
		}
	}
	return nil
}

func (pag *paginator) One(dest interface{}) error {
	pq, err := pag.buildWithCursor()
	if err != nil {
		return err
	}
	return pq.sel.One(dest)
}

func (pag *paginator) Iterator() Iterator {
	pq, err := pag.buildWithCursor()
	if err != nil {
		return &iterator{nil, err}
	}
	return pq.sel.Iterator()
}

func (pag *paginator) IteratorContext(ctx context.Context) Iterator {
	pq, err := pag.buildWithCursor()
	if err != nil {
		return &iterator{nil, err}
	}
	return pq.sel.IteratorContext(ctx)
}

func (pag *paginator) String() string {
	pq, err := pag.buildWithCursor()
	if err != nil {
		panic(err.Error())
	}
	return pq.sel.String()
}

func (pag *paginator) Arguments() []interface{} {
	pq, err := pag.buildWithCursor()
	if err != nil {
		return nil
	}
	return pq.sel.Arguments()
}

func (pag *paginator) Query() (*sql.Rows, error) {
	pq, err := pag.buildWithCursor()
	if err != nil {
		return nil, err
	}
	return pq.sel.Query()
}

func (pag *paginator) QueryContext(ctx context.Context) (*sql.Rows, error) {
	pq, err := pag.buildWithCursor()
	if err != nil {
		return nil, err
	}
	return pq.sel.QueryContext(ctx)
}

func (pag *paginator) QueryRow() (*sql.Row, error) {
	pq, err := pag.buildWithCursor()
	if err != nil {
		return nil, err
	}
	return pq.sel.QueryRow()
}

func (pag *paginator) QueryRowContext(ctx context.Context) (*sql.Row, error) {
	pq, err := pag.buildWithCursor()
	if err != nil {
		return nil, err
	}
	return pq.sel.QueryRowContext(ctx)
}

func (pag *paginator) Prepare() (*sql.Stmt, error) {
	pq, err := pag.buildWithCursor()
	if err != nil {
		return nil, err
	}
	return pq.sel.Prepare()
}

func (pag *paginator) PrepareContext(ctx context.Context) (*sql.Stmt, error) {
	pq, err := pag.buildWithCursor()
	if err != nil {
		return nil, err
	}
	return pq.sel.PrepareContext(ctx)
}

func (pag *paginator) TotalItems() (uint, error) {
	pq, err := pag.build()
	if err != nil {
		return 0, err
	}
	return pq.count()
}

func (pag *paginator) build() (*paginatorQuery, error) {
	pq, err := immutable.FastForward(pag)
	if err != nil {
		return nil, err
	}
	return pq.(*paginatorQuery), nil
}

func (pag *paginator) buildWithCursor() (*paginatorQuery, error) {
	pq, err := immutable.FastForward(pag)
	if err != nil {
		return nil, err
	}

	pqq := pq.(*paginatorQuery)
	if pqq.cursorColumn != "" {
		orderBy := pqq.cursorColumn
		if pqq.cursorReverseOrder {
			if strings.HasPrefix(orderBy, "-") {
				orderBy = orderBy[1:]
			} else {
				orderBy = "-" + orderBy
			}
		}
		pqq.sel = pqq.sel.OrderBy(orderBy)
	}

	if pqq.cursorCond != nil {
		pqq.sel = pqq.sel.Where(pqq.cursorCond).Offset(0)
	}

	return pqq, nil
}

func (pag *paginator) Prev() immutable.Immutable {
	if pag == nil {
		return nil
	}
	return pag.prev
}

func (pag *paginator) Fn(in interface{}) error {
	if pag.fn == nil {
		return nil
	}
	return pag.fn(in.(*paginatorQuery))
}

func (pag *paginator) Base() interface{} {
	return &paginatorQuery{}
}
