package main

import (
	"log"

	"github.com/webx-top/db/_tools/generator/model"
	"github.com/webx-top/db/mysql"
)

var settings = mysql.ConnectionURL{
	Host:     "localhost",
	Database: "blog",
	User:     "root",
	Password: "root",
}

type Book struct {
	ID        int    `db:"id"`
	Title     string `db:"title"`
	AuthorID  int    `db:"author_id"`
	SubjectID int    `db:"subject_id"`
}

func main() {
	sess, err := mysql.Open(settings)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	var posts []model.Post
	err = sess.Collection("webx_post").Find().All(&posts)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}
}
