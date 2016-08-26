//Generated by webx-top/db
package dbschema

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	
)

type Album struct {
	trans	*factory.Transaction
	
	Id              	int     	`db:"id,omitempty,pk" comment:"ID"`
	Title           	string  	`db:"title" comment:"标题"`
	Description     	string  	`db:"description" comment:"简介"`
	Content         	string  	`db:"content" comment:"正文"`
	Created         	int     	`db:"created" comment:"创建时间"`
	Updated         	int     	`db:"updated" comment:"编辑时间"`
	Views           	int     	`db:"views" comment:"浏览次数"`
	Comments        	int     	`db:"comments" comment:"评论次数"`
	Likes           	int     	`db:"likes" comment:"喜欢次数"`
	Display         	string  	`db:"display" comment:"显示"`
	Deleted         	int     	`db:"deleted" comment:"删除时间"`
	AllowComment    	string  	`db:"allow_comment" comment:"是否允许评论"`
	Tags            	string  	`db:"tags" comment:"标签"`
	Catid           	int     	`db:"catid" comment:"分类ID"`
}

func (this *Album) SetTrans(trans *factory.Transaction) *Album {
	this.trans = trans
	return this
}

func (this *Album) Param() *factory.Param {
	return factory.NewParam(Factory).SetTrans(this.trans).SetCollection("album")
}

func (this *Album) Get(mw func(db.Result) db.Result) error {
	return this.Param().SetRecv(this).SetMiddleware(mw).One()
}

func (this *Album) List(mw func(db.Result) db.Result, page, size int) ([]*Album, func() int64, error) {
	r := []*Album{}
	counter, err := this.Param().SetPage(page).SetSize(size).SetRecv(&r).SetMiddleware(mw).List()
	return r, counter, err
}

func (this *Album) ListByOffset(mw func(db.Result) db.Result, offset, size int) ([]*Album, func() int64, error) {
	r := []*Album{}
	counter, err := this.Param().SetOffset(offset).SetSize(size).SetRecv(&r).SetMiddleware(mw).List()
	return r, counter, err
}

func (this *Album) Add(args ...*Album) (interface{}, error) {
	var data = this
	if len(args)>0 {
		data = args[0]
	}
	return this.Param().SetSend(data).Insert()
}

func (this *Album) Edit(mw func(db.Result) db.Result, args ...*Album) error {
	var data = this
	if len(args)>0 {
		data = args[0]
	}
	return this.Param().SetSend(data).SetMiddleware(mw).Update()
}

func (this *Album) Delete(mw func(db.Result) db.Result) error {
	return this.Param().SetMiddleware(mw).Delete()
}

