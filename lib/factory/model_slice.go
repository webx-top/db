package factory

import (
	"fmt"

	"github.com/webx-top/echo/param"
)

type Slicex[T Model] []T

// Range applies the given function to each element of the slice.
// It will return an error if any of the calls to the given function
// return an error.
func (s Slicex[T]) Range(fn func(m Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

// RangeRaw applies the given function to each element of the slice.
// It returns the first error it encounters, if any.
// If no error is encountered, it returns nil.
func (s Slicex[T]) RangeRaw(fn func(m T) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

// GroupBy groups the slice by the given keyField.
// It returns a map of keyField values to slices of T.
// The keyField value is used as the key in the map,
// and the value is a slice of T containing all the rows
// that have the given keyField value.
func (s Slicex[T]) GroupBy(keyField string) map[string][]T {
	r := map[string][]T{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []T{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

// KeyBy returns a map of keyField values to rows of T.
// The keyField value is used as the key in the map,
// and the value is the row of T that has the given keyField value.
func (s Slicex[T]) KeyBy(keyField string) map[string]T {
	r := map[string]T{}
	for _, row := range s {
		dmap := row.AsMap(keyField)
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

// AsKV transforms the slice into a map of keyField values to valueField values.
// It returns a param.Store containing the transformed data.
// The keyField value is used as the key in the map,
// and the value is the valueField value of the row that has the given keyField value.
func (s Slicex[T]) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap(keyField, valueField)
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

// Transform transforms the slice into a slice of param.Store.
// It applies the given transformation function to each element of the slice.
// The transformation function takes a map of keyField values to valueField values
// and returns a param.Store containing the transformed data.
// The keyField value is used as the key in the map,
// and the value is the valueField value of the row that has the given keyField value.
func (s Slicex[T]) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

// FromList transforms the given data interface into a slice of T.
// If the data interface is not a slice of T, it transforms each element of the data interface into a row of T by calling the given newFunc.
// If the data interface is a slice of T, it appends the slice to the given slice of T.
func (s Slicex[T]) FromList(data interface{}, newFunc func() T) Slicex[T] {
	values, ok := data.([]T)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := newFunc()
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}
