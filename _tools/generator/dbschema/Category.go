//Generated by webx-top/db
package dbschema

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	
	"time"
)

type Category struct {
	trans	*factory.Transaction
	
	Id         	uint    	`db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	Pid        	uint    	`db:"pid" bson:"pid" comment:"上级分类" json:"pid" xml:"pid"`
	Name       	string  	`db:"name" bson:"name" comment:"分类名称" json:"name" xml:"name"`
	Description	string  	`db:"description" bson:"description" comment:"说明" json:"description" xml:"description"`
	Haschild   	string  	`db:"haschild" bson:"haschild" comment:"是否有子分类" json:"haschild" xml:"haschild"`
	Updated    	uint    	`db:"updated" bson:"updated" comment:"更新时间" json:"updated" xml:"updated"`
	RcType     	string  	`db:"rc_type" bson:"rc_type" comment:"关联类型" json:"rc_type" xml:"rc_type"`
	Sort       	int     	`db:"sort" bson:"sort" comment:"排序" json:"sort" xml:"sort"`
	Tmpl       	string  	`db:"tmpl" bson:"tmpl" comment:"模板" json:"tmpl" xml:"tmpl"`
}

func (this *Category) Trans() *factory.Transaction {
	return this.trans
}

func (this *Category) Use(trans *factory.Transaction) *Category {
	this.trans = trans
	return this
}

func (this *Category) Param() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetTrans(this.trans).SetCollection("category")
}

func (this *Category) Get(mw func(db.Result) db.Result) error {
	return this.Param().SetRecv(this).SetMiddleware(mw).One()
}

func (this *Category) List(mw func(db.Result) db.Result, page, size int) ([]*Category, func() int64, error) {
	r := []*Category{}
	counter, err := this.Param().SetPage(page).SetSize(size).SetRecv(&r).SetMiddleware(mw).List()
	return r, counter, err
}

func (this *Category) ListByOffset(mw func(db.Result) db.Result, offset, size int) ([]*Category, func() int64, error) {
	r := []*Category{}
	counter, err := this.Param().SetOffset(offset).SetSize(size).SetRecv(&r).SetMiddleware(mw).List()
	return r, counter, err
}

func (this *Category) Add() (interface{}, error) {
	
	return this.Param().SetSend(this).Insert()
}

func (this *Category) Edit(mw func(db.Result) db.Result) error {
	this.Updated = uint(time.Now().Unix())
	return this.Param().SetSend(this).SetMiddleware(mw).Update()
}

func (this *Category) Delete(mw func(db.Result) db.Result) error {
	
	return this.Param().SetMiddleware(mw).Delete()
}

