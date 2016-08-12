package main

import (
	"log"

	"github.com/webx-top/db/_tools/generator/model"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/mysql"
)

var settings = mysql.ConnectionURL{
	Host:     "localhost",
	Database: "blog",
	User:     "root",
	Password: "root",
}

func main() {
	database, err := mysql.Open(settings)
	if err != nil {
		log.Fatal(err)
	}
	db := factory.New().SetDebug(true)
	db.AddDB(database).Cluster(0).SetPrefix(`webx_`)
	defer db.CloseAll()

	var posts []model.Post
	//err = db.Find("webx_post").All(&posts)
	err = db.All("post", nil, &posts)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}
}
