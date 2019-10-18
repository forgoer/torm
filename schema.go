package torm

import (
	"errors"
	"reflect"
)

// Schema model definition
type Schema struct {
	Fields       []*Field
	PrimaryField *Field
	attributes   map[string]interface{}
}

func (s *Schema) SetId(id int64) error {
	for _, f := range s.Fields {
		if f.Primary {
			f.SetValue(id)
			break
		}
	}

	return nil
}

func (s *Schema) Attributes() map[string]interface{} {
	if s.attributes == nil {
		s.attributes = make(map[string]interface{})
		for _, field := range s.Fields {
			if !field.Ignored {
				s.attributes[field.Name] = field.Value.Addr().Interface()
			}
		}
	}

	return s.attributes
}

func NewSchema(model interface{}) (*Schema, error) {
	results := reflect.Indirect(reflect.ValueOf(model))

	kind := results.Kind()
	if kind != reflect.Struct {
		return nil, errors.New("unsupported value, should be struct")
	}

	var resultType reflect.Type
	var resultValue reflect.Value

	resultType = results.Type()
	resultValue = results

	var schema Schema

	for i := 0; i < resultType.NumField(); i++ {
		field := NewField(resultValue.Field(i), resultType.Field(i))
		if schema.PrimaryField == nil && field.Primary {
			schema.PrimaryField = field
		}
		schema.Fields = append(schema.Fields, field)
	}

	return &schema, nil
}
