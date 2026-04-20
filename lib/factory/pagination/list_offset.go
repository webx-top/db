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

package pagination

import (
	"github.com/webx-top/db"
	clientPagination "github.com/webx-top/db/lib/factory/pagination/client"
	"github.com/webx-top/echo"
	"github.com/webx-top/pagination"
)

// OffsetLister 偏移值分页列表查询接口
type OffsetLister interface {
	ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error)
}

var _ OffsetChunkLister = (*OffsetList)(nil)
var _ OffsetLister = (*OffsetList)(nil)

type OffsetChunkLister interface {
	SetProg(prog func(int64, int64))
	ChunkList(eachPageCallback func() error, size int, offset int) error
	ChunkListNoOffset(eachPageCallback func() (nextCond db.Compound, err error), size int) error
}

// NewOffsetLister 创建偏移值分页列表查询
func NewOffsetLister(list OffsetLister, recv interface{}, mw func(db.Result) db.Result, args ...interface{}) *OffsetList {
	return &OffsetList{
		ListParam: NewListParam(recv, mw, args...),
		ls:        list,
	}
}

// OffsetList 偏移值分页列表查询封装
type OffsetList struct {
	*ListParam
	ls OffsetLister
}

// ListByOffset 分页查询
func (f *OffsetList) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = f.recv
	} else {
		f.recv = recv
	}
	if mw == nil {
		mw = f.mw
	}
	if len(args) < 1 {
		args = f.args
	}
	return f.ls.ListByOffset(recv, mw, offset, size, args...)
}

// ListSize 列表尺寸
func (f *OffsetList) ListSize() int {
	if sz, ok := f.ls.(ListSizer); ok {
		return sz.ListSize()
	}
	if f.recv == nil {
		return ObjectsSize(f.ls)
	}
	return ObjectsSize(f.recv)
}

// ChunkList 分批查询列表
func (f *OffsetList) ChunkList(eachPageCallback func() error, size int, offset int) error {
	cnt, err := f.ListByOffset(f.recv, f.mw, offset, size)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil
		}
		return err
	}
	initOffset := offset
	for total := cnt(); int64(offset) < total; offset += size {
		if f.prog != nil {
			f.prog(int64(offset), total)
		}
		if offset > initOffset {
			_, err = f.ListByOffset(f.recv, f.mw, offset, size)
			if err != nil {
				if err == db.ErrNoMoreRows {
					return nil
				}
				return err
			}
		}
		err = eachPageCallback()
		if err != nil {
			return err
		}
	}
	return err
}

// ChunkList 分批查询列表
func (f *OffsetList) ChunkListNoOffset(eachPageCallback func() (nextCond db.Compound, err error), size int) error {
	cnt, err := f.ListByOffset(f.recv, f.mw, 0, size)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil
		}
		return err
	}
	args := f.args
	var cond db.Compound
	var offset int64
	initOffset := offset
	step := int64(size)
	for total := cnt(); offset < total; offset += step {
		if f.prog != nil {
			f.prog(0, total)
		}
		if offset > initOffset {
			_, err = f.ListByOffset(f.recv, f.mw, 0, size, args...)
			if err != nil {
				if err == db.ErrNoMoreRows {
					return nil
				}
				return err
			}
		}
		cond, err = eachPageCallback()
		if err != nil {
			return err
		}
		if cond != nil {
			args = []interface{}{cond}
		}
	}
	return err
}

// Paging 分页信息
func (f *OffsetList) Paging(ctx echo.Context, varSuffix ...string) (*pagination.Pagination, error) {
	return PagingWithOffsetLister(ctx, f, varSuffix...)
}

// DataTable 分页信息
func (f *OffsetList) DataTable(ctx echo.Context, args ...string) (map[string]interface{}, error) {
	return clientPagination.DataTable(ctx, f, args...)
}

// DataGrid 分页信息
func (f *OffsetList) DataGrid(ctx echo.Context, searchField string, args ...string) (map[string]interface{}, error) {
	return clientPagination.DataGrid(ctx, f, searchField, args...)
}

// JqGrid 分页信息
func (f *OffsetList) JqGrid(ctx echo.Context, searchField string, args ...string) (map[string]interface{}, error) {
	return clientPagination.JqGrid(ctx, f, searchField, args...)
}
