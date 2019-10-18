package torm

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"

	"github.com/thinkoner/torm/utils"
)

// Field The model field definition
type Field struct {
	Value       reflect.Value
	StructField reflect.StructField
	Attrs       map[string]string
	Primary     bool
	Name        string
	Ignored     bool
	IsBlank     bool
}

// GetParams Get the attr of tag
func (f *Field) GetAttr(key string) (string, bool) {
	key = strings.ToUpper(key)
	val, ok := f.Attrs[key]

	return val, ok
}

// Set set a value to the field
func (f *Field) SetValue(value interface{}) (err error) {

	//if !f.Value.CanAddr() {
	//
	//}

	reflectValue, ok := value.(reflect.Value)
	if !ok {
		reflectValue = reflect.ValueOf(value)
	}

	fieldValue := f.Value
	if reflectValue.IsValid() {
		if reflectValue.Type().ConvertibleTo(fieldValue.Type()) {
			fieldValue.Set(reflectValue.Convert(fieldValue.Type()))
		} else {
			if fieldValue.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(f.StructField.Type.Elem()))
				}
				fieldValue = fieldValue.Elem()
			}

			if reflectValue.Type().ConvertibleTo(fieldValue.Type()) {
				fieldValue.Set(reflectValue.Convert(fieldValue.Type()))
			} else if scanner, ok := fieldValue.Addr().Interface().(sql.Scanner); ok {
				v := reflectValue.Interface()
				if valuer, ok := v.(driver.Valuer); ok {
					if v, err = valuer.Value(); err == nil {
						err = scanner.Scan(v)
					}
				} else {
					err = scanner.Scan(v)
				}
			} else {
				err = fmt.Errorf("could not convert argument of field  from %s to %s", reflectValue.Type(), fieldValue.Type())
			}
		}
	} else {
		f.Value.Set(reflect.Zero(f.Value.Type()))
	}

	f.IsBlank = utils.IsBlank(f.Value)

	return err
}

func (f *Field) initialize() {
	tag := f.StructField.Tag.Get("torm")

	for _, kv := range strings.Split(tag, ";") {
		v := strings.Split(kv, ":")
		k := strings.TrimSpace(strings.ToUpper(v[0]))
		if len(v) >= 2 {
			f.Attrs[k] = strings.Join(v[1:], ":")
		} else {
			f.Attrs[k] = ""
		}
		switch k {
		case "-":
			f.Ignored = true
		case "PRIMARY_KEY":
			f.Primary = true
		case "COLUMN":
			f.Name = f.Attrs[k]
		}
	}

	if f.Name == "" || f.Name == "-" {
		f.Name = utils.SnakeCase(f.StructField.Name)
	}

	f.IsBlank = utils.IsBlank(f.Value)
}

func NewField(value reflect.Value, structField reflect.StructField) *Field {
	field := Field{
		Value:       value,
		StructField: structField,
		Attrs:       map[string]string{},
	}
	field.initialize()

	return &field
}
