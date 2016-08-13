package main

import (
	"log"

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
	db := factory.SetDebug(true)
	db.AddDB(database).Cluster(0).SetPrefix(`webx_`)
	defer db.CloseAll()

	var posts []*dbschema.Post
	//err = db.Find("webx_post").All(&posts)
	err = db.All(factory.NewParam(nil).SetCollection(`post`).SetResult(&posts))
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}

	err = factory.NewParam(db).SetCollection(`post`).SetResult(&posts).All()
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}

	post := &dbschema.Post{}
	posts, _, err = post.List(nil, 1, 100000)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}
}
