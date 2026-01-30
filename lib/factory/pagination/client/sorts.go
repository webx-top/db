/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package client

import (
	"slices"
	"strings"

	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
)

// CanSort 是否允许排序
// ctx: echo.Context
// table: factory.Short | string | ICanSortFields
// field: 字段
func CanSort(ctx echo.Context, table interface{}, field string) bool {
	if csf, ok := table.(ICanSortFields); ok {
		return slices.Contains(csf.CanSortFields(ctx), field)
	}

	noPrefixTableName := factory.NoPrefixTableName(table)
	return factory.ExistField(noPrefixTableName, field)
}

// Sorts 获取数据查询时的排序方式
// ctx: echo.Context
// table: factory.Short | string | ICanSortFields
// defaultSorts: 默认排序
func Sorts(ctx echo.Context, table interface{}, defaultSorts ...string) []interface{} {
	sorts := make([]interface{}, 0, len(defaultSorts)+1)
	sort := ctx.Form(`sort`)
	field := strings.TrimPrefix(sort, `-`)
	if len(field) > 0 && CanSort(ctx, table, field) {
		sorts = append(sorts, sort)
		for _, defaultSort := range defaultSorts {
			if field != strings.TrimPrefix(defaultSort, `-`) {
				sorts = append(sorts, defaultSort)
			}
		}
	} else {
		for _, defaultSort := range defaultSorts {
			sorts = append(sorts, defaultSort)
		}
	}
	return sorts
}

type CanSortFieldsFunc func(echo.Context) []string

// CanSortFields returns the fields that can be sorted.
// The returned fields will be checked in the Sorts function to determine if a field can be sorted.
// If the returned fields do not contain the field, it will not be sorted.
func (f CanSortFieldsFunc) CanSortFields(ctx echo.Context) []string {
	return f(ctx)
}

// CanSortFields returns a function that implements the ICanSortFields interface.
// The returned function returns the fields that can be sorted.
// The returned fields will be checked in the Sorts function to determine if a field can be sorted.
// If the returned fields do not contain the field, it will not be sorted.
func CanSortFields(fields ...string) CanSortFieldsFunc {
	return func(_ echo.Context) []string {
		return fields
	}
}

type ICanSortFields interface {
	CanSortFields(echo.Context) []string
}
