package main

import (
	"log"

	"github.com/webx-top/db"
	"github.com/webx-top/db/_tools/generator/dbschema"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/mysql"
)

var settings = mysql.ConnectionURL{
	Host:     "localhost",
	Database: "blog",
	User:     "root",
	Password: "root",
}

type PostCollection struct {
	Post *dbschema.Post `db:",inline"`
	User *dbschema.User `db:",inline"`
}

func main() {
	database, err := mysql.Open(settings)
	if err != nil {
		log.Fatal(err)
	}
	factory.SetDebug(true)
	factory.AddDB(database).Cluster(0).SetPrefix(`webx_`)
	defer factory.Default().CloseAll()

	var posts []*dbschema.Post
	//err = db.Find("webx_post").All(&posts)
	log.Println(`查询方式1：使用Factory查询`)
	err = factory.Default().All(factory.NewParam(nil).SetCollection(`post`).SetResult(&posts))
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}

	log.Println(`查询方式2：使用Param查询`)
	err = factory.NewParam().SetCollection(`post`).SetResult(&posts).All()
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}

	log.Println(`查询方式3：使用dbschema的List方法查询`)
	post := &dbschema.Post{}
	posts, _, err = post.List(nil, 1, 100000)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}

	log.Println(`查询方式4：使用LeftJoin关联查询`)
	m := []*PostCollection{}
	err = factory.NewParam().SetCollection(`post AS a`).SetCols(db.Raw(`a.*`)).Select().LeftJoin(`webx_user b`).On(`b.id=a.uid`).All(&m)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range m {
		log.Printf("%q (ID: %d)\n", post.Post.Title, post.Post.Id)
	}
}
