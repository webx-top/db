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

package db

import (
	"fmt"
	"strings"
)

// RawValue interface represents values that can bypass SQL filters. This is an
// exported interface but it's rarely used directly, you may want to use
// the `db.Raw()` function instead.
type RawValue interface {
	fmt.Stringer
	Compound

	// Raw returns the string representation of the value that the user wants to
	// pass without any escaping.
	Raw() string

	// Arguments returns the arguments to be replaced on the query.
	Arguments() []interface{}
}

type rawValue struct {
	v string
	a *[]interface{} // This may look ugly but allows us to use db.Raw() as keys for db.Cond{}.
}

func (r rawValue) Arguments() []interface{} {
	if r.a != nil {
		return *r.a
	}
	return nil
}

func (r rawValue) Raw() string {
	return r.v
}

func (r rawValue) String() string {
	return r.Raw()
}

// Sentences return each one of the map records as a compound.
func (r rawValue) Sentences() []Compound {
	return []Compound{r}
}

// Operator returns the default compound operator.
func (r rawValue) Operator() CompoundOperator {
	return OperatorNone
}

// Empty return false if this struct holds no value.
func (r rawValue) Empty() bool {
	return r.v == ""
}

// Raw marks chunks of data as protected, so they pass directly to the query
// without any filtering. Use with care.
//
// Example:
//
//	// SOUNDEX('Hello')
//	Raw("SOUNDEX('Hello')")
//
// Raw returns a value that satifies the db.RawValue interface.
func Raw(value string, args ...interface{}) RawValue {
	r := rawValue{v: value, a: nil}
	if len(args) > 0 {
		r.a = &args
	}
	return r
}

var _ = RawValue(&rawValue{})
var _ = RawValue(&rawValues{})

func AndRawValues(values ...RawValue) *rawValues {
	return &rawValues{values: values, operator: OperatorAnd}
}

func OrRawValues(values ...RawValue) *rawValues {
	return &rawValues{values: values, operator: OperatorOr}
}

type rawValues struct {
	values   []RawValue
	operator CompoundOperator
}

func (r *rawValues) Add(v RawValue) *rawValues {
	r.values = append(r.values, v)
	return r
}

func (r rawValues) Arguments() []interface{} {
	var n int
	for i := range r.values {
		n += len(r.values[i].Arguments())
	}
	args := make([]interface{}, 0, n)
	for i := range r.values {
		args = append(args, r.values[i].Arguments()...)
	}
	return args
}

func (r rawValues) Raw() string {
	sqls := make([]string, len(r.values))
	for i := range r.values {
		sqls[i] = r.values[i].Raw()
	}
	oprt := `AND`
	if r.operator == OperatorOr {
		oprt = `OR`
	}
	if len(sqls) > 1 {
		return `(` + strings.Join(sqls, `) `+oprt+` (`) + `)`
	}
	return strings.Join(sqls, ` `+oprt+` `)
}

func (r rawValues) String() string {
	return r.Raw()
}

// Sentences return each one of the map records as a compound.
func (r rawValues) Sentences() []Compound {
	compounds := make([]Compound, len(r.values))
	for i := range r.values {
		compounds[i] = r.values[i]
	}
	return compounds
}

// Operator returns the default compound operator.
func (r rawValues) Operator() CompoundOperator {
	return r.operator
}

// Empty return false if this struct holds no value.
func (r rawValues) Empty() bool {
	return len(r.values) == 0
}
