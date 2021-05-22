package main

import (
	"os/exec"
	"strings"

	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/sqlbuilder"
)

func TableToStructName(tableName string, prefix string) string {
	if len(prefix) > 0 {
		tableName = strings.TrimPrefix(tableName, prefix)
	}
	return factory.ToCamleCase(tableName)
}

func ZeroValue(typ string) string {
	switch typ {
	case `uint`, `int`, `uint32`, `int32`, `int64`, `uint64`:
		return `0`
	case `float32`, `float64`:
		return `0.0`
	case `bool`:
		return `false`
	case `string`:
		return "``"
	case `byte[]`:
		return `nil`
	case `time.Time`:
		return `time.Time{}`
	default:
		panic(`undefined zero value for ` + typ)
	}
}

func DataType(fieldInfo *factory.FieldInfo) string {
	switch fieldInfo.DataType {
	case `int`, `tinyint`, `smallint`, `mediumint`:
		if fieldInfo.Unsigned {
			return `uint`
		}
		return `int`
	case `bigint`:
		if fieldInfo.Unsigned {
			return `uint64`
		}
		return `int64`
	case `decimal`, `double`:
		return `float64`
	case `float`:
		return `float32`
	case `bit`, `binary`, `varbinary`, `tinyblob`, `blob`, `mediumblob`, `longblob`: //二进制
		return `byte[]`
	case `geometry`, `point`, `linestring`, `polygon`, `multipoint`, `multilinestring`, `multipolygon`, `geometrycollection`: //几何图形
		return `byte[]`

	//postgreSQL
	case `boolean`:
		return `bool`
	case `oid`:
		if fieldInfo.Unsigned {
			return `uint64`
		}
		return `int64`

	default:
		return `string`
	}
}

func GetTableFields(engine string, d sqlbuilder.Database, tableName string, typeMap map[string][]string) ([]string, map[string]factory.FieldInfo, []string) {
	switch engine {
	case "mymysql", "mysql":
		fallthrough
	default:
		return getMySQLTableFields(d, tableName, typeMap)
	}
}

func GetTableComment(engine string, d sqlbuilder.Database, tableName string) (string, error) {
	switch engine {
	case "mymysql", "mysql":
		fallthrough
	default:
		return getMySQLTableComment(d, tableName)
	}
}

func Format(file string) error {
	cmd := exec.Command(`gofmt`, `-l`, `-s`, `-w`, file)
	return cmd.Run()
}
