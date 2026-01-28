package sqlbuilder

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/webx-top/db"
)

var (
	relationCache      = make(map[string]*parsedRelation)
	relationCacheMutex sync.RWMutex
)

func getRelationCache(t reflect.StructField) *parsedRelation {
	relationCacheMutex.RLock()
	//println(`key:`, t.Type.String()+`#`+string(t.Tag))
	r := relationCache[t.Type.String()+`#`+string(t.Tag)]
	relationCacheMutex.RUnlock()
	return r
}

func setRelationCache(t reflect.StructField, r *parsedRelation) {
	relationCacheMutex.Lock()
	relationCache[t.Type.String()+`#`+string(t.Tag)] = r
	relationCacheMutex.Unlock()
}

type fieldScope string

const (
	fieldScopeAll   fieldScope = `all`
	fieldScopeFirst fieldScope = `first`
	fieldScopeLast  fieldScope = `last`
)

type kv struct {
	k string
	v interface{}

	// - cross field -
	field string
	scope fieldScope
}

func (s *kv) String() string {
	return fmt.Sprintf(`{k: %+v, v: %v, field: %v, scope: %v}`, s.k, s.v, s.field, s.scope)
}

func (s *kv) getValue(refVal reflect.Value, sliceLen int) interface{} {
	if len(s.field) == 0 {
		return s.v
	}
	if sliceLen < 0 {
		return mapper.FieldByName(refVal, s.field).Interface()
	}
	if sliceLen == 0 {
		return nil
	}
	v := s.v
	switch s.scope {
	case fieldScopeFirst:
		v = mapper.FieldByName(refVal.Index(0), s.field).Interface()
	case fieldScopeLast:
		v = mapper.FieldByName(refVal.Index(sliceLen-1), s.field).Interface()
	default:
		relValsMap := map[interface{}]struct{}{}
		for j := 0; j < sliceLen; j++ {
			v := mapper.FieldByName(refVal.Index(j), s.field).Interface()
			relValsMap[v] = struct{}{}
		}
		relVals := make([]interface{}, 0, len(relValsMap))
		for k := range relValsMap {
			relVals = append(relVals, k)
		}
		if len(relVals) == 0 {
			return nil
		}
		if len(relVals) == 1 {
			v = relVals[0]
		} else {
			v = db.In(relVals)
		}
	}
	return v
}

type colType struct {
	col    interface{}
	colStr string
	typ    string
}

func (s *colType) String() string {
	return fmt.Sprintf(`{col: %+v, colStr: %v, typ: %v}`,
		s.col, s.colStr, s.typ)
}

type selectorArgs struct {
	orderby []interface{}
	offset  int
	limit   int
	groupby []interface{}
	columns []*colType
}

func (s *selectorArgs) String() string {
	return fmt.Sprintf(`{orderby: %+v, offset: %v, limit: %v, groupby: %+v, columns: %+v}`,
		s.orderby, s.offset, s.limit, s.groupby, s.columns)
}

func NewParsedRelation(relations []string, pipes []Pipe) *parsedRelation {
	return &parsedRelation{
		relations: relations,
		pipes:     pipes,
	}
}

type parsedRelation struct {
	relations    []string
	pipes        []Pipe
	where        *[]*kv
	selectorArgs *selectorArgs
	mutex        sync.RWMutex
}

func (r *parsedRelation) String() string {
	var where interface{}
	if r.where != nil {
		where = *r.where
	}
	return fmt.Sprintf(`{relations: %+v, pipes: %v, where: %+v, selectorArgs: %v}`,
		r.relations, r.pipes, where, r.selectorArgs)
}

func (r *parsedRelation) setWhere(where *[]*kv) {
	r.mutex.Lock()
	r.where = where
	r.mutex.Unlock()
}

func (r *parsedRelation) Where() (where *[]*kv) {
	r.mutex.RLock()
	where = r.where
	r.mutex.RUnlock()
	return
}

func (r *parsedRelation) SelectorArgs() (selectorArgs *selectorArgs) {
	r.mutex.RLock()
	selectorArgs = r.selectorArgs
	r.mutex.RUnlock()
	return
}

func (r *parsedRelation) setSelectorArgs(selectorArgs *selectorArgs) {
	r.mutex.Lock()
	r.selectorArgs = selectorArgs
	r.mutex.Unlock()
}
