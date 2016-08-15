# factory
为了方便使用，将原包中的常用功能进行再次封装。支持的目标功能有：

1. 支持主从式读写分离，支持多主多从
2. 支持同时使用多种数据库且每一种数据库都支持读写分离
3. 支持缓存数据结果
4. 便捷的分页查询功能
5. 尽量节省代码

以上。

# 用法

```
package main
import (
	"fmt"
	"log"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/mysql"
)

var settings = mysql.ConnectionURL{
	Host:     "localhost",
	Database: "blog",
	User:     "root",
	Password: "root",
}

type Post struct {
    Id      int     `db:"id,omitempty"`
    Title   string  `db:"title"`
    Group   string  `db:"group"`
    Views   int     `db:"views"`
}

type PostCollection struct {
	Post *Post `db:",inline"`
	User *User `db:",inline"`
}

type User struct {
    Id      int     `db:"id,omitempty"`
    Name   string   `db:"name"`
}

func main() {
	database, err := mysql.Open(settings)
	if err != nil {
		log.Fatal(err)
	}
	factory.SetDebug(true)
	factory.AddDB(database).Cluster(0).SetPrefix(`webx_`)
	defer factory.Default().CloseAll()

	var posts []*Post

    err = factory.All(factory.NewParam().SetCollection(`post`).SetPage(1).SetSize(10).SetResult(&posts))
	// 生成SQL：SELECT * FROM `webx_post` LIMIT 10

	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		log.Printf("%q (ID: %d)\n", post.Title, post.Id)
	}
}
```

## 查询多行数据 (使用All方法)

### 方法 1.
```
err = factory.All(factory.NewParam().SetCollection(`post`).SetResult(&posts))
```

也可以附加更多条件（后面介绍的所有方法均支持这种方式）：
```
    err = factory.NewParam().SetCollection(`post`).SetResult(&posts).SetArgs(db.Cond{`title LIKE`:`%test%`}).SetMiddleware(func(r db.Result)db.Result{
        return r.OrderBy(`-id`).Group(`group`)
    }).All()
    // 生成SQL：SELECT * FROM `webx_post` WHERE (`title` LIKE "%test%") GROUP BY `group` ORDER BY `id` DESC
```

### 方法 2.
```
err = factory.NewParam().SetCollection(`post`).SetResult(&posts).All()

```

### 关联查询
```
m := []*PostCollection{}
err = factory.NewParam().SetCollection(`post AS a`).SetCols(db.Raw(`a.*`)).AddJoin(`LEFT`, `user`, `b`, `b.id=a.id`).Select().All(&m)
```

## 查询分页数据 (使用List方法)

### 方法 1.
```
var countFn func()int64
countFn, err = factory.List(factory.NewParam().SetCollection(`post`).SetResult(&posts).SetPage(1).SetSize(10))
```

### 方法 2.
```
countFn, err = factory.NewParam().SetCollection(`post`).SetResult(&posts).SetPage(1).SetSize(10).List()
```

## 查询一行数据 (使用One方法)

### 方法 1.
```
var post Post
err = factory.One(factory.NewParam().SetCollection(`post`).SetResult(&post))
```

### 方法 2.
```
var post Post
err = factory.NewParam().SetCollection(`post`).SetResult(&post).One()
```

## 插入数据 (使用Insert方法)

### 方法 1.
```
var post Post
post=Post{
    Title:`test title`,
}
err = factory.Insert(factory.NewParam().SetCollection(`post`).SetSave(&post))
```

### 方法 2.
```
var post Post
post=Post{
    Title:`test title`,
}
err = factory.NewParam().SetCollection(`post`).SetSave(&post).Insert()
```

## 更新数据 (使用Update方法)

### 方法 1.
```
var post Post
post=Post{
    Title:`test title`,
}
err = factory.Update(factory.NewParam().SetCollection(`post`).SetSave(&post).SetArgs("id",1))
```

### 方法 2.
```
var post Post
post=Post{
    Title:`test title`,
}
err = factory.NewParam().SetCollection(`post`).SetSave(&post).SetArgs("id",1).Update()
```

## 删除数据 (使用Delete方法)

### 方法 1.
```
err = factory.Delete(factory.NewParam().SetCollection(`post`).SetArgs("id",1))
```

### 方法 2.
```
err = factory.NewParam().SetCollection(`post`).SetArgs("id",1).Update()
```

## 使用事务

### 方法 1.
```
	param = factory.NewParam().SetCollection(`post`).SetTxMW(func(t *factory.Transaction) (err error) {
		param := factory.NewParam().SetCollection(`post`).SetSave(map[string]int{
			"views": 1,
		}).SetArgs("id", 1)
		err = t.Update(param)
		// err=fmt.Errorf(`failured`)
		// 当返回 nil 时，自动执行Commit，否则自动执行Rollback
		return
	})
	factory.Tx(param)
```

### 方法 2.
```
    param = factory.NewParam().Begin().SetCollection(`post`)
    err:=param.SetSave(map[string]int{"views": 1}).SetArgs("id", 1).Update()
    if err!=nil {
        param.End(err)
        return
    }
    err=factory.NewParam().TransFrom(param).SetCollection(`post`).SetSave(map[string]int{"views": 2}).SetArgs("id", 1).Update()
    if err!=nil {
        param.End(err)
        return
    }
    err=factory.NewParam().TransFrom(param).SetCollection(`post`).SetSave(map[string]int{"views": 3}).SetArgs("id", 1).Update()
    param.End(err)
```
