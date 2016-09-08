package factory

type FieldInfo struct {
	DataType      string   `json:"dataType" xml:"dataType" bson:"dataType"`
	Unsigned      bool     `json:"unsigned" xml:"unsigned" bson:"unsigned"`
	PrimaryKey    bool     `json:"primaryKey" xml:"primaryKey" bson:"primaryKey"`
	AutoIncrement bool     `json:"autoIncrement" xml:"autoIncrement" bson:"autoIncrement"`
	Min           int      `json:"min" xml:"min" bson:"min"`
	Max           int      `json:"max" xml:"max" bson:"max"`
	MaxSize       int      `json:"maxSize" xml:"maxSize" bson:"maxSize"`
	Options       []string `json:"options" xml:"options" bson:"options"`
	DefaultValue  string   `json:"defaultValue" xml:"defaultValue" bson:"defaultValue"`
	Comment       string   `json:"comment" xml:"comment" bson:"comment"`
}

type FieldValidator map[string]map[string]*FieldInfo

func (f FieldValidator) ValidField(table string, field string) bool {
	if tb, ok := f[table]; ok {
		_, ok = tb[field]
		return ok
	}
	return false
}

func (f FieldValidator) ValidTable(table string) bool {
	_, ok := f[table]
	return ok
}

var Fields FieldValidator = map[string]map[string]*FieldInfo{}

func ValidField(table string, field string) bool {
	return Fields.ValidField(table, field)
}

func ValidTable(table string) bool {
	return Fields.ValidTable(table)
}
