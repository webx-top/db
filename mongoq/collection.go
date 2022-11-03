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
	"strings"
	"sync"

	"reflect"

	"github.com/webx-top/db"
	mgo "github.com/webx-top/qmgo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Collection represents a mongodb collection.
type Collection struct {
	parent     *Source
	collection *mgo.Collection
}

var (
	// idCache should be a struct if we're going to cache more than just
	// _id field here
	idCache      = make(map[reflect.Type]string)
	idCacheMutex sync.RWMutex
)

// Find creates a result set with the given conditions.
func (col *Collection) Find(terms ...interface{}) db.Result {
	fields := []string{"*"}

	conditions := col.compileQuery(terms...)

	res := &result{}
	res = res.frame(func(r *resultQuery) error {
		r.c = col
		r.conditions = conditions
		r.fields = fields
		return nil
	})

	return res
}

var comparisonOperators = map[db.ComparisonOperator]string{
	db.ComparisonOperatorEqual:    "$eq",
	db.ComparisonOperatorNotEqual: "$ne",

	db.ComparisonOperatorLessThan:    "$lt",
	db.ComparisonOperatorGreaterThan: "$gt",

	db.ComparisonOperatorLessThanOrEqualTo:    "$lte",
	db.ComparisonOperatorGreaterThanOrEqualTo: "$gte",

	db.ComparisonOperatorIn:    "$in",
	db.ComparisonOperatorNotIn: "$nin",
}

func compare(field string, cmp db.Comparison) (string, interface{}) {
	op := cmp.Operator()
	value := cmp.Value()

	switch op {
	case db.ComparisonOperatorEqual:
		return field, value
	case db.ComparisonOperatorBetween:
		values := value.([]interface{})
		return field, mgo.M{
			"$gte": values[0],
			"$lte": values[1],
		}
	case db.ComparisonOperatorNotBetween:
		values := value.([]interface{})
		return "$or", []mgo.M{
			{field: mgo.M{"$gt": values[1]}},
			{field: mgo.M{"$lt": values[0]}},
		}
	case db.ComparisonOperatorIs:
		if value == nil {
			return field, mgo.M{"$exists": false}
		}
		return field, mgo.M{"$eq": value}
	case db.ComparisonOperatorIsNot:
		if value == nil {
			return field, mgo.M{"$exists": true}
		}
		return field, mgo.M{"$ne": value}
	case db.ComparisonOperatorRegExp, db.ComparisonOperatorLike:
		return field, primitive.Regex{Pattern: value.(string), Options: ""}
	case db.ComparisonOperatorNotRegExp, db.ComparisonOperatorNotLike:
		return field, mgo.M{"$not": primitive.Regex{Pattern: value.(string), Options: ""}}
	}

	if cmpOp, ok := comparisonOperators[op]; ok {
		return field, mgo.M{
			cmpOp: value,
		}
	}

	panic(fmt.Sprintf("Unsupported operator %v", op))
}

// compileStatement transforms conditions into something *mgo.Session can
// understand.
func compileStatement(cond db.Cond) mgo.M {
	conds := mgo.M{}

	// Walking over conditions
	for fieldI, value := range cond {
		field := strings.TrimSpace(fmt.Sprintf("%v", fieldI))

		if cmp, ok := value.(db.Comparison); ok {
			k, v := compare(field, cmp)
			conds[k] = v
			continue
		}

		var op string
		chunks := strings.SplitN(field, ` `, 2)

		if len(chunks) > 1 {
			switch chunks[1] {
			case `IN`:
				op = `$in`
			case `NOT IN`:
				op = `$nin`
			case `>`:
				op = `$gt`
			case `<`:
				op = `$lt`
			case `<=`:
				op = `$lte`
			case `>=`:
				op = `$gte`
			case `!=`, `<>`:
				op = `$ne`
			default:
				op = chunks[1]
			}
		}
		field = chunks[0]

		if op == "" {
			conds[field] = value
		} else {

			if v, y := conds[field]; y {
				if bsonM, ok := v.(mgo.M); ok {
					bsonM[op] = value
					continue
				}
			}
			conds[field] = mgo.M{op: value}
		}
	}

	return conds
}

// compileConditions compiles terms into something *mgo.Session can
// understand.
func (col *Collection) compileConditions(term interface{}) interface{} {

	switch t := term.(type) {
	case []interface{}:
		values := []interface{}{}
		for i := range t {
			value := col.compileConditions(t[i])
			if value != nil {
				values = append(values, value)
			}
		}
		if len(values) > 0 {
			return values
		}
	case db.Cond:
		return compileStatement(t)
	case db.Compound:
		values := []interface{}{}

		for _, s := range t.Sentences() {
			values = append(values, col.compileConditions(s))
		}

		var op string
		switch t.Operator() {
		case db.OperatorOr:
			op = `$or`
		default:
			op = `$and`
		}

		return mgo.M{op: values}
	}
	return nil
}

// compileQuery compiles terms into something that *mgo.Session can
// understand.
func (col *Collection) compileQuery(terms ...interface{}) interface{} {
	var query interface{}

	compiled := col.compileConditions(terms)

	if compiled != nil {
		conditions := compiled.([]interface{})
		if len(conditions) == 1 {
			query = conditions[0]
		} else {
			// this should be correct.
			// query = map[string]interface{}{"$and": conditions}

			// attempt to workaround https://jira.mongodb.org/browse/SERVER-4572
			mapped := map[string]interface{}{}
			for _, v := range conditions {
				for kk := range v.(map[string]interface{}) {
					mapped[kk] = v.(map[string]interface{})[kk]
				}
			}

			query = mapped
		}
	} else {
		query = nil
	}

	return query
}

// Name returns the name of the table or tables that form the collection.
func (col *Collection) Name() string {
	return col.collection.GetCollectionName()
}

// Truncate deletes all rows from the table.
func (col *Collection) Truncate() error {
	err := col.collection.DropCollection(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (col *Collection) InsertReturning(item interface{}) error {
	return db.ErrUnsupported
}

func (col *Collection) UpdateReturning(item interface{}) error {
	return db.ErrUnsupported
}

// Insert inserts an item (map or struct) into the collection.
func (col *Collection) Insert(item interface{}) (interface{}, error) {
	var err error

	id := getID(item)

	ctx := context.Background()

	if col.parent.versionAtLeast(2, 6, 0, 0) {
		// this breaks MongoDb older than 2.6
		if _, err = col.collection.Upsert(ctx, mgo.M{"_id": id}, item); err != nil {
			return nil, err
		}
	} else {
		// Allocating a new ID.
		var result *mgo.InsertOneResult
		if result, err = col.collection.InsertOne(ctx, mgo.M{"_id": id}); err != nil {
			return nil, err
		}
		_ = result

		// Now append data the user wants to append.
		if err = col.collection.UpdateOne(ctx, mgo.M{"_id": id}, item); err != nil {
			// Cleanup allocated ID
			if err := col.collection.Remove(ctx, mgo.M{"_id": id}); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	return id, nil
}

// Exists returns true if the collection exists.
func (col *Collection) Exists() bool {
	ctx := context.Background()
	query := col.parent.database.Collection(`system.namespaces`).Find(ctx, map[string]string{
		`name`: fmt.Sprintf(`%s.%s`, col.parent.database.GetDatabaseName(), col.Name()),
	})
	count, _ := query.Count()
	return count > 0
}

// Fetches object _id or generates a new one if object doesn't have one or the one it has is invalid
func getID(item interface{}) interface{} {
	v := reflect.ValueOf(item) // convert interface to Value
	v = reflect.Indirect(v)    // convert pointers

	switch v.Kind() {
	case reflect.Map:
		if inItem, ok := item.(map[string]interface{}); ok {
			if id, ok := inItem["_id"]; ok {
				bsonID, ok := id.(primitive.ObjectID)
				if ok {
					return bsonID
				}
			}
		}
	case reflect.Struct:
		t := v.Type()

		idCacheMutex.RLock()
		fieldName, found := idCache[t]
		idCacheMutex.RUnlock()

		if !found {
			for n := 0; n < t.NumField(); n++ {
				field := t.Field(n)
				if len(field.PkgPath) > 0 {
					continue // Private field
				}

				tag := field.Tag.Get("bson")
				if len(tag) == 0 {
					tag = field.Tag.Get("db")
				}

				if len(tag) == 0 {
					continue
				}

				parts := strings.Split(tag, ",")

				if parts[0] == "_id" {
					fieldName = field.Name
					idCacheMutex.Lock()
					idCache[t] = fieldName
					idCacheMutex.Unlock()
					break
				}
			}
		}
		if len(fieldName) > 0 {
			if bsonID, ok := v.FieldByName(fieldName).Interface().(primitive.ObjectID); ok {
				if !bsonID.IsZero() {
					return bsonID
				}
			} else {
				return v.FieldByName(fieldName).Interface()
			}
		}
	}

	return mgo.NewObjectID()
}