package postgresql

import (
	"encoding/json"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory/mysql"
)

func FindInSet(fieldName string, value string, useFulltextIndex ...bool) db.RawValue {
	if len(useFulltextIndex) > 0 && useFulltextIndex[0] {
		v := mysql.CleanFulltextOperator(value)
		v = strings.ReplaceAll(v, `,`, ``)
		return match(v, fieldName)
	}
	fieldName = com.AddCSlashes(fieldName, '"')
	return db.Raw(`? = ANY (string_to_array("`+fieldName+`",','))`, value)
}

func FindInJSON(fieldName string, value interface{}, jsonFields ...string) db.RawValue {
	fieldName = com.AddCSlashes(fieldName, '"')
	arr := []interface{}{value}
	b, _ := json.Marshal(arr)
	return db.Raw(`"` + fieldName + `" ?& array` + string(b))
}

func match(safelyMatchValue string, keys ...string) db.RawValue {
	rawValues := make([]db.RawValue, 0, len(keys))
	for _, key := range keys {
		if len(key) == 0 {
			continue
		}
		key = com.AddCSlashes(key, '"')
		rawValues = append(rawValues, db.Raw(`"`+key+`" @@ '`+safelyMatchValue+`'`))
	}
	return db.OrRawValues(rawValues...)
}

func Match(value string, keys ...string) db.RawValue {
	value = mysql.CleanFulltextOperator(value)
	return match(value, keys...)
}
