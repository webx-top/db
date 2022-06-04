# webx-top/db

Forked from upper/db

## 本仓库新增功能

1. 支持表前缀
2. 支持ForceIndex查询
3. 支持关联查询（感谢gosql提供灵感，[查看用法](https://github.com/webx-top/db/blob/master/_tools/test/relation/main.go)）
4. 新增[lib/factory](https://github.com/webx-top/db/tree/master/lib/factory)包
5. 增加MySQL表结构体生成工具（安装命令：`go install github.com/webx-top/db/cmd/dbgenerator`，使用命令`dbgenerator -h`查看用法）
6. 其它改进

[原始文档](http://www.admpub.com:8080/upper-db-manual/en/)

## Cases
- [Nging](https://github.com/admpub/nging)

<p align="center">
  <img src="https://upper.io/img/gopher.svg" width="256">
</p>

<p align="center">
  <a href="https://github.com/upper/db/actions?query=workflow%3Aunit-tests"><img alt="upper/db unit tests status" src="https://github.com/upper/db/workflows/unit-tests/badge.svg"></a>
</p>

# upper/db

`upper/db` is a productive data access layer (DAL) for [Go](https://golang.org)
that provides agnostic tools to work with different data sources, such as:

* [PostgreSQL](https://upper.io/v4/adapter/postgresql)
* [MySQL](https://upper.io/v4/adapter/mysql)
* [MSSQL](https://upper.io/v4/adapter/mssql)
* [CockroachDB](https://upper.io/v4/adapter/cockroachdb)
* [MongoDB](https://upper.io/v4/adapter/mongo)
* [QL](https://upper.io/v4/adapter/ql)
* [SQLite](https://upper.io/v4/adapter/sqlite)

See [upper.io/v4](//upper.io/v4) for documentation and code samples.

## The tour

![tour](https://user-images.githubusercontent.com/385670/91495824-c6fabb00-e880-11ea-925b-a30b94474610.png)

Take the [tour](https://tour.upper.io) to see real live examples in your
browser.

## License

Licensed under [MIT License](./LICENSE)

## Contributors

See the [list of contributors](https://github.com/upper/db/graphs/contributors).
