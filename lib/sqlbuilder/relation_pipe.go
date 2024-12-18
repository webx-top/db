package sqlbuilder

import (
	"reflect"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"
)

// parse neq(field,value)
func parsePipe(pipeName string) Pipe {
	pos := strings.Index(pipeName, "(")
	if pos > 0 {
		param := pipeName[pos+1:]
		param = strings.TrimSuffix(param, ")")
		funcName := pipeName[0:pos]
		if gen, ok := PipeGeneratorList[funcName]; ok {
			return gen(param)
		}
		return nil
	}
	pipe, ok := PipeList[pipeName]
	if !ok {
		return nil
	}
	return pipe
}

type Pipe func(row reflect.Value, val interface{}) interface{}
type Pipes map[string]Pipe
type PipeGenerators map[string]func(params string) Pipe

func (pipes *Pipes) Add(name string, pipe Pipe) {
	(*pipes)[name] = pipe
}

func (gens *PipeGenerators) Add(name string, generator func(params string) Pipe) {
	(*gens)[name] = generator
}

var (
	PipeGeneratorList = PipeGenerators{
		`neq`: func(params string) Pipe { // name:value
			args := strings.SplitN(params, `:`, 2)
			var (
				fieldName     string
				expectedValue string
			)
			switch len(args) {
			case 2:
				fieldName = strings.TrimSpace(args[0])
				expectedValue = strings.TrimSpace(args[1])
			default:
				return nil
			}
			return func(row reflect.Value, v interface{}) interface{} {
				fieldValue := mapper.FieldByName(row, fieldName).Interface()
				if expectedValue != param.AsString(fieldValue) {
					return v
				}
				return nil
			}
		},
		`eq`: func(params string) Pipe { // name:value
			args := strings.SplitN(params, `:`, 2)
			var (
				fieldName     string
				expectedValue string
			)
			switch len(args) {
			case 2:
				fieldName = strings.TrimSpace(args[0])
				expectedValue = strings.TrimSpace(args[1])
			default:
				return nil
			}
			return func(row reflect.Value, v interface{}) interface{} {
				fieldValue := mapper.FieldByName(row, fieldName).Interface()
				if expectedValue == param.AsString(fieldValue) {
					return v
				}
				return nil
			}
		},
		`isNil`: func(structField string) Pipe {
			if len(structField) == 0 {
				return nil
			}
			structFieldPath := strings.Split(structField, `.`)
			return func(row reflect.Value, v interface{}) interface{} {
				fv := reflect.Indirect(row)
				for _, structField := range structFieldPath {
					fv = fv.FieldByName(structField)
					if !fv.IsValid() {
						return v
					}
					fv = reflect.Indirect(fv)
				}
				if fv.Kind() == reflect.Ptr && fv.IsNil() {
					return v
				}
				return nil
			}
		},
		`isZero`: func(structField string) Pipe {
			if len(structField) == 0 {
				return nil
			}
			if !com.IsUpperLetter(rune(structField[0])) {
				return func(row reflect.Value, v interface{}) interface{} {
					fv := mapper.FieldByName(row, structField)
					if reflect.Indirect(fv).IsZero() {
						return v
					}
					return nil
				}
			}
			structFieldPath := strings.Split(structField, `.`)
			return func(row reflect.Value, v interface{}) interface{} {
				fv := reflect.Indirect(row)
				for _, structField := range structFieldPath {
					fv = fv.FieldByName(structField)
					if !fv.IsValid() {
						return v
					}
					fv = reflect.Indirect(fv)
				}
				if fv.IsZero() {
					return v
				}
				return nil
			}
		},
		`notZero`: func(structField string) Pipe {
			if len(structField) == 0 {
				return nil
			}
			if !com.IsUpperLetter(rune(structField[0])) {
				structField = com.PascalCaseWith(structField, '.')
			}
			structFieldPath := strings.Split(structField, `.`)
			return func(row reflect.Value, v interface{}) interface{} {
				fv := reflect.Indirect(row)
				for _, structField := range structFieldPath {
					fv = fv.FieldByName(structField)
					if !fv.IsValid() {
						return nil
					}
					fv = reflect.Indirect(fv)
				}
				if fv.IsZero() {
					return nil
				}
				return v
			}
		},
	}
	PipeList = Pipes{
		`split`: func(_ reflect.Value, v interface{}) interface{} {
			val := v.(string)
			if len(val) == 0 {
				return nil
			}
			if val[0] == '[' {
				val = strings.Trim(val, `[]`)
				if len(val) == 0 {
					return nil
				}
			}
			items := strings.Split(val, `,`)
			result := make([]interface{}, 0, len(items))
			for _, item := range items {
				item = strings.TrimSpace(item)
				if len(item) == 0 {
					continue
				}
				result = append(result, item)
			}
			return result
		},
		`gtZero`: func(_ reflect.Value, v interface{}) interface{} {
			i := param.AsUint64(v)
			if i > 0 {
				return i
			}
			return nil
		},
		`notEmpty`: func(_ reflect.Value, v interface{}) interface{} {
			s := v.(string)
			if len(s) > 0 {
				return s
			}
			return nil
		},
	}
)
