package factory

var (
	databases = map[string]*DBI{
		DefaultDBKey: DefaultDBI,
	}
)

func DBIRegister(dbi *DBI, keys ...string) {
	key := DefaultDBKey
	if len(keys) > 0 {
		key = keys[0]
	}
	if _, y := databases[key]; y {
		panic(`DBI key already exists, please do not duplicate registrations`)
	}
	databases[key] = dbi
}

func DBIGet(keys ...string) *DBI {
	if len(keys) > 0 {
		return databases[keys[0]]
	}
	return databases[DefaultDBKey]
}

func DBIExists(key string) bool {
	_, ok := databases[key]
	return ok
}

func NewDBI() *DBI {
	return &DBI{
		StructToTable: map[string]string{},
		Fields:        map[string]map[string]*FieldInfo{},
		Models:        ModelInstancers{},
		TableNamers:   map[string]func(obj interface{}) string{},
	}
}

// DBI 数据库信息
type DBI struct {
	// 结构体名与表名对照
	StructToTable map[string]string
	// Fields {table:{field:FieldInfo}}
	Fields FieldValidator
	// Models {StructName:ModelInstancer}
	Models ModelInstancers
	// TableNamers {table:NewName}
	TableNamers TableNamers
}

func (d *DBI) TableName(structName string) string {
	tableName, _ := d.StructToTable[structName]
	return tableName
}
