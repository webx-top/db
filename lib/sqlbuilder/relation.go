package sqlbuilder

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/admpub/errors"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
)

type BuilderChainFunc func(Selector) Selector

func (b *sqlBuilder) Relation(name string, fn BuilderChainFunc) SQLBuilder {
	if b.relationMap == nil {
		b.relationMap = make(map[string]BuilderChainFunc)
	}
	b.relationMap[name] = fn
	return b
}

func (b *sqlBuilder) RelationMap() map[string]BuilderChainFunc {
	return b.relationMap
}

func (sel *selector) Relation(name string, fn BuilderChainFunc) Selector {
	sel.SQLBuilder().Relation(name, fn)
	return sel
}

func eachField(t reflect.Type, fn func(field reflect.StructField, val string, name string, relations []string) error) error {
	for i := 0; i < t.NumField(); i++ {
		val := t.Field(i).Tag.Get("relation")
		name := t.Field(i).Name
		field := t.Field(i)

		if len(val) > 0 && val != "-" {
			relations := strings.SplitN(val, ",", 2)
			if len(relations) != 2 {
				return fmt.Errorf("relation tag error, length must 2,but get %v", relations)
			}

			err := fn(field, val, name, relations)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type Name_ interface {
	Name_() string
}

var (
	ErrUnableDetermineTableName = errors.New(`Unable to determine table name`)
	TableName                   = DefaultTableName
)

func DefaultTableName(data interface{}, retry ...bool) (string, error) {
	switch m := data.(type) {
	case Name_:
		return m.Name_(), nil
	default:
		if len(retry) > 0 && retry[0] {
			return ``, ErrUnableDetermineTableName
		}
	}
	value := reflect.ValueOf(data)
	if value.IsNil() {
		return ``, errors.WithMessagef(errors.New("model argument cannot be nil pointer passed"), `%T`, data)
	}
	tp := reflect.Indirect(value).Type()
	if tp.Kind() == reflect.Interface {
		tp = reflect.Indirect(value).Elem().Type()
	}

	if tp.Kind() != reflect.Slice {
		return ``, fmt.Errorf("model argument must slice, but get %T", data)
	}

	tpEl := tp.Elem()
	//Compatible with []*Struct or []Struct
	if tpEl.Kind() == reflect.Ptr {
		tpEl = tpEl.Elem()
	}
	//fmt.Printf("[TableName] %s ========>%[1]T, %[1]v\n", tpEl.Name(), reflect.New(tpEl).Interface())
	name, err := DefaultTableName(reflect.New(tpEl).Interface(), true)
	if err == ErrUnableDetermineTableName {
		name = com.SnakeCase(tpEl.Name())
		err = nil
	}
	return name, err
}

// RelationOne is get the associated relational data for a single piece of data
func RelationOne(builder SQLBuilder, data interface{}) error {
	refVal := reflect.Indirect(reflect.ValueOf(data))
	t := refVal.Type()

	return eachField(t, func(field reflect.StructField, val string, name string, relations []string) error {
		var foreignModel reflect.Value
		// if field type is slice then one-to-many ,eg: []*Struct
		if field.Type.Kind() == reflect.Slice {
			foreignModel = reflect.New(field.Type)
			table, err := TableName(foreignModel.Interface())
			if err != nil {
				return err
			}
			// batch get field values
			// Since the structure is slice, there is no need to new Value
			sel := builder.Select().From(table).Where(db.Cond{
				relations[1]: mapper.FieldByName(refVal, relations[0]).Interface(),
			})
			if chains := builder.RelationMap(); chains != nil {
				if chainFn, ok := chains[name]; ok {
					sel = chainFn(sel)
				}
			}
			err = sel.All(foreignModel.Interface())
			if err != nil && err != db.ErrNoMoreRows {
				return err
			}

			if reflect.Indirect(foreignModel).Len() == 0 {
				// If relation data is empty, must set empty slice
				// Otherwise, the JSON result will be null instead of []
				refVal.FieldByName(name).Set(reflect.MakeSlice(field.Type, 0, 0))
			} else {
				refVal.FieldByName(name).Set(foreignModel.Elem())
			}

		} else {
			// If field type is struct the one-to-one,eg: *Struct
			foreignModel = reflect.New(field.Type.Elem())
			table, err := TableName(foreignModel.Interface())
			if err != nil {
				return err
			}
			sel := builder.Select().From(table).Where(db.Cond{
				relations[1]: mapper.FieldByName(refVal, relations[0]).Interface(),
			})
			if chains := builder.RelationMap(); chains != nil {
				if chainFn, ok := chains[name]; ok {
					sel = chainFn(sel)
				}
			}
			err = sel.All(foreignModel.Interface())
			// If one-to-one NoRows is not an error that needs to be terminated
			if err != nil && err != db.ErrNoMoreRows {
				return err
			}

			if err == nil {
				refVal.FieldByName(name).Set(foreignModel)
			}
		}
		return nil
	})
}

// RelationAll is gets the associated relational data for multiple pieces of data
func RelationAll(builder SQLBuilder, data interface{}) error {
	refVal := reflect.Indirect(reflect.ValueOf(data))

	l := refVal.Len()

	if l == 0 {
		return nil
	}

	// get the struct field in slice
	t := reflect.Indirect(refVal.Index(0)).Type()

	return eachField(t, func(field reflect.StructField, val string, name string, relations []string) error {
		relVals := make([]interface{}, 0)
		relValsMap := make(map[interface{}]interface{}, 0)

		// get relation field values and unique
		for j := 0; j < l; j++ {
			v := mapper.FieldByName(refVal.Index(j), relations[0]).Interface()
			relValsMap[v] = nil
		}

		for k := range relValsMap {
			relVals = append(relVals, k)
		}

		var foreignModel reflect.Value
		// if field type is slice then one to many ,eg: []*Struct
		if field.Type.Kind() == reflect.Slice {
			foreignModel = reflect.New(field.Type)
			table, err := TableName(foreignModel.Interface())
			if err != nil {
				return err
			}
			// batch get field values
			// Since the structure is slice, there is no need to new Value
			sel := builder.Select().From(table).Where(db.Cond{
				relations[1]: db.In(relVals),
			})
			if chains := builder.RelationMap(); chains != nil {
				if chainFn, ok := chains[name]; ok {
					sel = chainFn(sel)
				}
			}
			err = sel.All(foreignModel.Interface())
			if err != nil && err != db.ErrNoMoreRows {
				return err
			}

			fmap := make(map[interface{}]reflect.Value)

			// Combine relation data as a one-to-many relation
			// For example, if there are multiple images under an article
			// we use the article ID to associate the images, map[1][]*Images
			for n := 0; n < reflect.Indirect(foreignModel).Len(); n++ {
				val := reflect.Indirect(foreignModel).Index(n)
				fid := mapper.FieldByName(val, relations[1])
				if _, has := fmap[fid.Interface()]; !has {
					fmap[fid.Interface()] = reflect.New(reflect.SliceOf(field.Type.Elem())).Elem()
				}
				fmap[fid.Interface()] = reflect.Append(fmap[fid.Interface()], val)
			}

			// Set the result to the model
			for j := 0; j < l; j++ {
				fid := mapper.FieldByName(refVal.Index(j), relations[0])
				if value, has := fmap[fid.Interface()]; has {
					reflect.Indirect(refVal.Index(j)).FieldByName(name).Set(value)
				} else {
					// If relation data is empty, must set empty slice
					// Otherwise, the JSON result will be null instead of []
					reflect.Indirect(refVal.Index(j)).FieldByName(name).Set(reflect.MakeSlice(field.Type, 0, 0))
				}
			}
		} else {
			// If field type is struct the one to one,eg: *Struct
			foreignModel = reflect.New(field.Type.Elem())

			// Batch get field values, but must new slice []*Struct
			fi := reflect.New(reflect.SliceOf(foreignModel.Type()))

			b := builder

			table, err := TableName(foreignModel.Interface())
			if err != nil {
				return err
			}
			sel := b.Select().From(table).Where(db.Cond{
				relations[1]: db.In(relVals),
			})
			if chains := builder.RelationMap(); chains != nil {
				if chainFn, ok := chains[name]; ok {
					sel = chainFn(sel)
				}
			}
			err = sel.All(foreignModel.Interface())
			if err != nil && err != db.ErrNoMoreRows {
				return err
			}

			// Combine relation data as a one-to-one relation
			fmap := make(map[interface{}]reflect.Value)
			for n := 0; n < reflect.Indirect(fi).Len(); n++ {
				val := reflect.Indirect(fi).Index(n)
				fid := mapper.FieldByName(val, relations[1])
				fmap[fid.Interface()] = val
			}

			// Set the result to the model
			for j := 0; j < l; j++ {
				fid := mapper.FieldByName(refVal.Index(j), relations[0])
				if value, has := fmap[fid.Interface()]; has {
					reflect.Indirect(refVal.Index(j)).FieldByName(name).Set(value)
				}
			}
		}

		return nil
	})
}
