package main

import (
	"fmt"
	"log"
	"time"

	"database/sql"

	"github.com/webx-top/cache/ttlmap"
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

type PostCollection2 struct {
	P *dbschema.Post `db:",inline"`
	U *dbschema.User `db:",inline"`
}

func main() {
	database, err := mysql.Open(settings)
	if err != nil {
		log.Fatal(err)
	}
	factory.SetDebug(true)
	cacher := ttlmap.New(1000000)
	factory.SetCacher(cacher)
	factory.AddDB(database).Cluster(0).SetPrefix(`webx_`)
	defer factory.Default().CloseAll()

	var posts []*dbschema.Post

	rows, err := factory.NewParam(nil).DB().Query(`SELECT * FROM webx_post ORDER BY id DESC`)
	if err != nil {
		log.Fatal(err)
	}
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		recv := make([]interface{}, len(columns))
		for k := range columns {
			recv[k] = &sql.NullString{}
		}
		err = rows.Scan(recv...)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("=======================================")
		for k, v := range recv {
			colName := columns[k]
			log.Printf("%v: %#v\n", colName, v.(*sql.NullString).String)
		}
	}

	_, err = factory.QueryTo(factory.NewParam(nil).SetCollection(`SELECT * FROM webx_post ORDER BY id DESC`).SetRecv(&posts).SetPage(1).SetSize(10))
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}

	//err = db.Find("webx_post").All(&posts)
	log.Println(`查询方式1：使用Factory查询`)
	err = factory.All(factory.NewParam(nil).SetCollection(`post`).SetRecv(&posts).SetPage(2).SetSize(10))
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}

	fmt.Println(``)
	fmt.Println(``)
	log.Println(`查询方式2：使用Param查询`)
	err = factory.NewParam().SetCollection(`post`).SetRecv(&posts).All()
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}

	fmt.Println(``)
	fmt.Println(``)
	log.Println(`查询方式3：使用dbschema的List方法查询`)

	post := &dbschema.Post{} //<------------------ define post

	_, err = post.List(nil, nil, 1, 100000)
	if err != nil {
		log.Fatal(err)
	}
	objects := post.Objects()
	for _, post := range objects {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}

	fmt.Println(``)
	fmt.Println(``)
	log.Println(`查询方式3.1：使用dbschema的List方法查询`)
	_, err = post.List(nil, nil, 1, 100000, db.Cond{`id >`: 1})
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range post.Objects() {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}

	log.Println(`-----old----------------`)
	for _, post := range objects {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}

	fmt.Println(``)
	fmt.Println(``)
	log.Println(`查询方式4：使用LeftJoin关联查询`)
	m := []*PostCollection{}
	err = factory.NewParam().SetCollection(`post`, `a`).SetCols(db.Raw(`a.*`)).AddJoin(`LEFT`, `user`, `b`, `b.id=a.id`).Select().All(&m)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range m {
		log.Printf("%q (ID: %d)\n", post.Post.Title, post.Post.Id)
	}

	fmt.Println(``)
	fmt.Println(``)
	log.Println(`查询方式4.1：使用LeftJoin关联查询，并测试字段转换`)
	m2 := []*PostCollection2{}
	structField := `P.Title`
	tableName := ``
	structField2 := `U.Id`
	tableName2 := ``
	paramJoin := factory.NewParam().SetCollection(`post`, `a`).SetCols(db.Raw(`a.*`)).AddJoin(`LEFT`, `user`, `b`, `b.id=a.id`)
	paramJoin.TableField(&PostCollection2{}, &structField, &tableName)
	paramJoin.TableField(&PostCollection2{}, &structField2, &tableName2)
	err = paramJoin.Select().All(&m2)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range m2 {
		log.Printf("%q (ID: %d)\n", post.P.Title, post.P.Id)
	}
	log.Println(`=== 字段转换结果：===================`)
	log.Println(structField, `====>`, tableName)
	log.Println(structField2, `====>`, tableName2)

	if tableName != `a.title` {
		panic(`字段名称转换错误:` + structField + `应该转换为a.title，实际却被转为了` + tableName)
	}
	if tableName2 != `b.id` {
		panic(`字段名称转换错误:` + structField2 + `应该转换为b.id，实际却被转为了` + tableName2)
	}

	fmt.Println(``)
	fmt.Println(``)
	log.Println(`查询方式5：使用Next查询大结果集`)
	res := factory.NewParam().SetCollection(`post`).Result()
	defer res.Close() //操作结束后别忘了执行关闭操作
	for res.Next(post) {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}

	fmt.Println(``)
	fmt.Println(``)
	log.Println(`查询方式6：使用Next查询大结果集`)
	res = post.Param().Result()
	for res.Next(post) {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}
	res.Close() //操作结束后别忘了执行关闭操作

	fmt.Println(``)
	fmt.Println(``)
	log.Println(`查询方式7：使用使用事务`)
	param := post.Param().Begin()
	res = param.Result()
	for res.Next(post) {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}

	res.Close() //操作结束后别忘了执行关闭操作
	err = param.SetSend(map[string]int{
		"views": 818,
	}).SetArgs("id", 1).Update()
	//err = fmt.Errorf(`failured`)
	param.End(err)

	fmt.Println(``)
	fmt.Println(``)
	log.Println(`查询方式8：使用Param的Setter查询`)
	err = factory.NewParam().Setter().Collection(`post`).Recv(&posts).All()
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}

	fmt.Println(``)
	fmt.Println(``)
	log.Println(`查询方式9：使用Param的Model查询`)
	recv := post.NewObjects()
	_, err = post.Param().Model().List(recv, nil, 1, 999, db.And(
		db.Cond{`id >`: 1},
		db.Cond{`id <`: 10},
	))
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range *recv {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}

	fmt.Println(``)
	fmt.Println(``)
	log.Println(`测试Upsert：`)
	err = post.Get(nil, db.Cond{"id": 12})
	if err != nil {
		log.Fatal(err)
	}
	post.Content += ` by Upsert!`
	post.Id = 13
	_, err = post.Upsert(nil, db.Cond{"id": post.Id})
	//err = post.Param().Setter().Args(db.Cond{"id": post.Id}).Send(post).Upsert()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(``)
	fmt.Println(``)
	log.Println(`测试缓存：`)
	recv = post.NewObjects()
	for i := 0; i < 5; i++ {
		_, err = post.Param().SetCache(10*time.Minute, `testCaching`).SetRecv(recv).SetArgs(db.And(
			db.Cond{`id >`: 1},
			db.Cond{`id <`: 10},
		)).SetPage(1).SetSize(10).List()
		if err != nil {
			log.Fatal(err)
		}
		for _, post := range *recv {
			log.Printf("%d => %q (ID: %d)\n", i+1, post.Title, post.Id)
		}
	}

	fmt.Println(``)
	fmt.Println(``)
	log.Println(`测试缓存2：`)
	for i := 0; i < 5; i++ {
		_, err = post.NewParam().SetCache(10*time.Minute, `testCaching`).Model().List(recv, nil, 1, 10, db.And(
			db.Cond{`id >`: 1},
			db.Cond{`id <`: 10},
		))
		if err != nil {
			log.Fatal(err)
		}
		for _, post := range *recv {
			log.Printf("%d => %q (ID: %d)\n", i+1, post.Title, post.Id)
		}
	}
	return

	param = post.Param()
	factory.Tx(param.SetTxMW(func(t *factory.Transaction) (err error) {
		param = param.SetSend(map[string]int{
			"views": 1,
		}).SetArgs("id", 1)
		err = t.Update(param)
		//err=fmt.Errorf(`failured`)
		return
	}))

}
