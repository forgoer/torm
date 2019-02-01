package torm

import (
	"reflect"
	"strings"

	"github.com/thinkoner/torm/utils"
)

// Field The model field definition
type Field struct {
	Value       reflect.Value
	StructField reflect.StructField
	params      map[string]string
}

// IsColumn Whether the specified column is configured
func (f *Field) IsColumn(column string) bool {
	c, ok := f.Column()
	if !ok {
		return false
	}
	return c == column
}

// Column Get the column name
func (f *Field) Column() (string, bool) {
	column, ok := f.GetParams("column")
	if !ok {
		return utils.SnakeCase(f.StructField.Name), true
	}

	if len(column) > 0 && column != "-" {
		return column, true
	}

	return "", false
}

// GetParams Get the params of tag
func (f *Field) GetParams(key string) (string, bool) {
	if f.params == nil {
		f.parseParams()
	}
	val, ok := f.params[key]
	return val, ok
}

func (f *Field) parseParams() {
	f.params = make(map[string]string)
	tagStr := f.StructField.Tag.Get("torm")
	for _, kv := range strings.Split(tagStr, ";") {
		v := strings.Split(kv, ":")
		k := strings.TrimSpace(strings.ToLower(v[0]))
		if len(v) >= 2 {
			f.params[k] = strings.Join(v[1:], ":")
		} else {
			f.params[k] = ""
		}
	}
}
