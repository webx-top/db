// Copyright (c) 2012-present The upper.io/db authors. All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package mongo

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"encoding/json"

	"github.com/webx-top/db"
	"github.com/webx-top/db/internal/immutable"
	mgo "github.com/webx-top/qmgo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type resultQuery struct {
	c *Collection

	fields     []string
	limit      int
	offset     int
	sort       []string
	conditions interface{}
	groupBy    []interface{}

	pageSize           uint
	pageNumber         uint
	cursorColumn       string
	cursorValue        interface{}
	cursorCond         db.Cond
	cursorReverseOrder bool
}

type result struct {
	iter  mgo.CursorI
	err   error
	errMu sync.Mutex

	fn   func(*resultQuery) error
	prev *result
}

var _ = immutable.Immutable(&result{})

func (res *result) frame(fn func(*resultQuery) error) *result {
	return &result{prev: res, fn: fn}
}

func (r *resultQuery) and(terms ...interface{}) error {
	if r.conditions == nil {
		return r.where(terms...)
	}

	r.conditions = map[string]interface{}{
		"$and": []interface{}{
			r.conditions,
			r.c.compileQuery(terms...),
		},
	}
	return nil
}

func (r *resultQuery) where(terms ...interface{}) error {
	r.conditions = r.c.compileQuery(terms...)
	return nil
}

func (res *result) And(terms ...interface{}) db.Result {
	return res.frame(func(r *resultQuery) error {
		return r.and(terms...)
	})
}

func (res *result) Where(terms ...interface{}) db.Result {
	return res.frame(func(r *resultQuery) error {
		return r.where(terms...)
	})
}

func (res *result) Paginate(pageSize uint) db.Result {
	return res.frame(func(r *resultQuery) error {
		r.pageSize = pageSize
		return nil
	})
}

func (res *result) Page(pageNumber uint) db.Result {
	return res.frame(func(r *resultQuery) error {
		r.pageNumber = pageNumber
		return nil
	})
}

func (res *result) Cursor(cursorColumn string) db.Result {
	return res.frame(func(r *resultQuery) error {
		r.cursorColumn = cursorColumn
		return nil
	})
}

func (res *result) NextPage(cursorValue interface{}) db.Result {
	return res.frame(func(r *resultQuery) error {
		r.cursorValue = cursorValue
		r.cursorReverseOrder = false
		r.cursorCond = db.Cond{
			r.cursorColumn: mgo.M{"$gt": cursorValue},
		}
		return nil
	})
}

func (res *result) PrevPage(cursorValue interface{}) db.Result {
	return res.frame(func(r *resultQuery) error {
		r.cursorValue = cursorValue
		r.cursorReverseOrder = true
		r.cursorCond = db.Cond{
			r.cursorColumn: mgo.M{"$lt": cursorValue},
		}
		return nil
	})
}

func (res *result) TotalEntries() (uint64, error) {
	return res.Count()
}

func (res *result) TotalPages() (uint, error) {
	count, err := res.Count()
	if err != nil {
		return 0, err
	}

	rq, err := res.build()
	if err != nil {
		return 0, err
	}

	if rq.pageSize < 1 {
		return 1, nil
	}

	total := uint(math.Ceil(float64(count) / float64(rq.pageSize)))
	return total, nil
}

// Limit determines the maximum limit of results to be returned.
func (res *result) Limit(n int) db.Result {
	return res.frame(func(r *resultQuery) error {
		r.limit = n
		return nil
	})
}

// Offset determines how many documents will be skipped before starting to grab
// results.
func (res *result) Offset(n int) db.Result {
	return res.frame(func(r *resultQuery) error {
		r.offset = n
		return nil
	})
}

// OrderBy determines sorting of results according to the provided names. Fields
// may be prefixed by - (minus) which means descending order, ascending order
// would be used otherwise.
func (res *result) OrderBy(fields ...interface{}) db.Result {
	return res.frame(func(r *resultQuery) error {
		ss := make([]string, len(fields))
		for i, field := range fields {
			ss[i] = fmt.Sprintf(`%v`, field)
		}
		r.sort = ss
		return nil
	})
}

// String satisfies fmt.Stringer
func (res *result) String() string {
	return ""
}

// Select marks the specific fields the user wants to retrieve.
func (res *result) Select(fields ...interface{}) db.Result {
	return res.frame(func(r *resultQuery) error {
		fieldslen := len(fields)
		r.fields = make([]string, 0, fieldslen)
		for i := 0; i < fieldslen; i++ {
			r.fields = append(r.fields, fmt.Sprintf(`%v`, fields[i]))
		}
		return nil
	})
}

// All dumps all results into a pointer to an slice of structs or maps.
func (res *result) All(dst interface{}) error {
	rq, err := res.build()
	if err != nil {
		return err
	}

	q, err := rq.query()
	if err != nil {
		return err
	}

	if rq.c.parent.LoggingEnabled() || rq.c.parent.LoggingElapsedMs() > 0 {
		defer func(start time.Time) {
			rq.c.parent.Log(&db.QueryStatus{
				Query: rq.debugQuery("Find.All"),
				Err:   err,
				Start: start,
				End:   time.Now(),
			})
		}(time.Now())
	}

	err = q.All(dst)
	if mgo.IsErrNoDocuments(err) {
		return db.ErrNoMoreRows
	}
	return err
}

// Group is used to group results that have the same value in the same column
// or columns.
func (res *result) Group(fields ...interface{}) db.Result {
	return res.frame(func(r *resultQuery) error {
		r.groupBy = fields
		return nil
	})
}

// One fetches only one result from the resultset.
func (res *result) One(dst interface{}) error {
	rq, err := res.build()
	if err != nil {
		return err
	}

	q, err := rq.query()
	if err != nil {
		return err
	}

	if rq.c.parent.LoggingEnabled() || rq.c.parent.LoggingElapsedMs() > 0 {
		defer func(start time.Time) {
			rq.c.parent.Log(&db.QueryStatus{
				Query: rq.debugQuery("Find.One"),
				Err:   err,
				Start: start,
				End:   time.Now(),
			})
		}(time.Now())
	}

	err = q.One(dst)
	if mgo.IsErrNoDocuments(err) {
		return db.ErrNoMoreRows
	}
	return err
}

func (res *result) Err() error {
	res.errMu.Lock()
	defer res.errMu.Unlock()

	return res.err
}

func (res *result) setErr(err error) {
	res.errMu.Lock()
	res.err = err
	res.errMu.Unlock()
}

func (res *result) Next(dst interface{}) bool {
	if res.iter == nil {
		rq, err := res.build()
		if err != nil {
			return false
		}

		q, err := rq.query()
		if err != nil {
			return false
		}

		if rq.c.parent.LoggingEnabled() || rq.c.parent.LoggingElapsedMs() > 0 {
			defer func(start time.Time) {
				rq.c.parent.Log(&db.QueryStatus{
					Query: rq.debugQuery("Find.Next"),
					Err:   err,
					Start: start,
					End:   time.Now(),
				})
			}(time.Now())
		}

		res.iter = q.Cursor()
	}

	if !res.iter.Next(dst) {
		res.setErr(res.iter.Err())
		return false
	}

	return true
}

// Delete remove the matching items from the collection.
func (res *result) Delete() error {
	rq, err := res.build()
	if err != nil {
		return err
	}

	if rq.c.parent.LoggingEnabled() || rq.c.parent.LoggingElapsedMs() > 0 {
		defer func(start time.Time) {
			rq.c.parent.Log(&db.QueryStatus{
				Query: rq.debugQuery("Remove"),
				Err:   err,
				Start: start,
				End:   time.Now(),
			})
		}(time.Now())
	}

	ctx := context.Background()

	_, err = rq.c.collection.RemoveAll(ctx, rq.conditions)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the result set.
func (r *result) Close() error {
	var err error
	if r.iter != nil {
		err = r.iter.Close()
		r.iter = nil
	}
	return err
}

// Update modified matching items from the collection with values of the given
// map or struct.
func (res *result) Update(src interface{}) (err error) {
	updateSet := map[string]interface{}{"$set": src}

	rq, err := res.build()
	if err != nil {
		return err
	}

	if rq.c.parent.LoggingEnabled() || rq.c.parent.LoggingElapsedMs() > 0 {
		defer func(start time.Time) {
			rq.c.parent.Log(&db.QueryStatus{
				Query: rq.debugQuery("Update"),
				Err:   err,
				Start: start,
				End:   time.Now(),
			})
		}(time.Now())
	}

	ctx := context.Background()
	_, err = rq.c.collection.UpdateAll(ctx, rq.conditions, updateSet)
	if err != nil {
		return err
	}
	return nil
}

func (res *result) build() (*resultQuery, error) {
	rqi, err := immutable.FastForward(res)
	if err != nil {
		return nil, err
	}

	rq := rqi.(*resultQuery)
	if !rq.cursorCond.Empty() {
		if err := rq.and(rq.cursorCond); err != nil {
			return nil, err
		}
	}

	if rq.cursorColumn != "" {
		if rq.cursorReverseOrder {
			rq.sort = append(rq.sort, "-"+rq.cursorColumn)
		} else {
			rq.sort = append(rq.sort, rq.cursorColumn)
		}
	}
	if rq.conditions == nil {
		rq.conditions = mgo.D{}
	}
	return rq, nil
}

// query executes a mgo query.
func (r *resultQuery) query() (mgo.QueryI, error) {
	if len(r.groupBy) > 0 {
		return nil, db.ErrUnsupported
	}

	ctx := context.Background()
	q := r.c.collection.Find(ctx, r.conditions)

	if r.pageSize > 0 {
		r.offset = int(r.pageSize * (r.pageNumber - 1))
		r.limit = int(r.pageSize)
	}

	if r.offset > 0 {
		q.Skip(int64(r.offset))
	}

	if r.limit > 0 {
		q.Limit(int64(r.limit))
	}

	if len(r.sort) > 0 {
		q.Sort(r.sort...)
	}

	selectedFields := mgo.M{}
	if len(r.fields) > 0 {
		for _, field := range r.fields {
			if field == `*` {
				break
			}
			selectedFields[field] = true
		}
	}

	if r.cursorReverseOrder {
		ids := make([]primitive.ObjectID, 0, r.limit)

		iter := q.Select(mgo.M{"_id": true}).Cursor()
		defer iter.Close()

		var item map[string]primitive.ObjectID
		for iter.Next(&item) {
			ids = append(ids, item["_id"])
		}

		r.conditions = mgo.M{"_id": mgo.M{"$in": ids}}

		q = r.c.collection.Find(ctx, r.conditions)
	}

	if len(selectedFields) > 0 {
		q.Select(selectedFields)
	}

	return q, nil
}

func (res *result) Exists() (bool, error) {
	total, err := res.Count()
	if err != nil {
		return false, err
	}
	if total > 0 {
		return true, nil
	}
	return false, nil
}

// Count counts matching elements.
func (res *result) Count() (total uint64, err error) {
	rq, err := res.build()
	if err != nil {
		return 0, err
	}

	if rq.c.parent.LoggingEnabled() || rq.c.parent.LoggingElapsedMs() > 0 {
		defer func(start time.Time) {
			rq.c.parent.Log(&db.QueryStatus{
				Query: rq.debugQuery("Find.Count"),
				Err:   err,
				Start: start,
				End:   time.Now(),
			})
		}(time.Now())
	}

	ctx := context.Background()
	q := rq.c.collection.Find(ctx, rq.conditions)

	var c int64
	c, err = q.Count()

	return uint64(c), err
}

func (res *result) Prev() immutable.Immutable {
	if res == nil {
		return nil
	}
	return res.prev
}

func (res *result) Fn(in interface{}) error {
	if res.fn == nil {
		return nil
	}
	return res.fn(in.(*resultQuery))
}

func (res *result) Base() interface{} {
	return &resultQuery{}
}

func (r *resultQuery) debugQuery(action string) string {
	query := fmt.Sprintf("db.%s.%s", r.c.collection.GetCollectionName(), action)

	if r.conditions != nil {
		query = fmt.Sprintf("%s.conds(%v)", query, r.conditions)
	}
	if r.limit > 0 {
		query = fmt.Sprintf("%s.limit(%d)", query, r.limit)
	}
	if r.offset > 0 {
		query = fmt.Sprintf("%s.offset(%d)", query, r.offset)
	}
	if len(r.fields) > 0 {
		selectedFields := mgo.M{}
		for _, field := range r.fields {
			if field == `*` {
				break
			}
			selectedFields[field] = true
		}
		if len(selectedFields) > 0 {
			query = fmt.Sprintf("%s.select(%v)", query, selectedFields)
		}
	}
	if len(r.groupBy) > 0 {
		escaped := make([]string, len(r.groupBy))
		for i := range r.groupBy {
			escaped[i] = string(mustJSON(r.groupBy[i]))
		}
		query = fmt.Sprintf("%s.groupBy(%v)", query, strings.Join(escaped, ", "))
	}
	if len(r.sort) > 0 {
		escaped := make([]string, len(r.sort))
		for i := range r.sort {
			escaped[i] = string(mustJSON(r.sort[i]))
		}
		query = fmt.Sprintf("%s.sort(%s)", query, strings.Join(escaped, ", "))
	}
	return query
}

func mustJSON(in interface{}) (out []byte) {
	out, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return out
}

func (res *result) Relation(name string, fn interface{}) db.Result {
	panic(db.ErrUnsupported)
}

func (res *result) ForceIndex(index string) db.Result {
	panic(db.ErrUnsupported)
}

func (res *result) Callback(m interface{}) db.Result {
	panic(db.ErrUnsupported)
}
