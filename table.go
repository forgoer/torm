package torm

import (
	"github.com/thinkoner/torm/utils"
	"reflect"
)

type TableName interface {
	TableName() string
}

// Model Begin a fluent query against a database Model.
func (c *Connection) Model(model interface{}) *Builder {
	var table string

	kind := reflect.Indirect(reflect.ValueOf(model)).Kind()

	if kind != reflect.Struct {
		panic("unsupported model, should be slice or struct")
	}

	if t, ok := model.(TableName); ok {
		table = t.TableName()
	} else {
		modelType := reflect.ValueOf(model).Type()
		if t, ok := reflect.New(modelType).Interface().(TableName); ok {
			table = t.TableName()
		} else {
			table = utils.SnakeCase(modelType.Name())
		}
	}

	return c.Table(table)
}

// Create Save a new model to the database.
func (c *Connection) Create(model interface{}) error {
	schema, err := NewSchema(model)
	if err != nil {
		return err
	}

	attributes := schema.Attributes()
	insertId, _, err := c.Model(model).Insert(attributes)

	if err != nil {
		return err
	}

	schema.SetId(insertId)

	return nil
}

// Save Save the model to the database.
func (c *Connection) Save(model interface{}) error {
	var err error
	var schema *Schema
	schema, err = NewSchema(model)
	if err != nil {
		return err
	}

	attributes := schema.Attributes()

	if schema.PrimaryField != nil && !schema.PrimaryField.IsBlank {
		_, err = c.Model(model).Where(schema.PrimaryField.Name, schema.PrimaryField.Value.Addr().Interface()).Update(attributes)
	} else {
		var insertId int64
		insertId, _, err = c.Model(model).Insert(attributes)
		if err == nil {
			schema.SetId(insertId)
		}
	}

	return err
}

// Destroy Destroy the model.
func (c *Connection) Destroy(model interface{}) error {
	var err error
	var schema *Schema
	schema, err = NewSchema(model)
	if err != nil {
		return err
	}

	if schema.PrimaryField != nil && !schema.PrimaryField.IsBlank {
		_, err = c.Model(model).Where(schema.PrimaryField.Name, schema.PrimaryField.Value.Addr().Interface()).Delete()
	} else {
		_, err = c.Model(model).Delete()
	}

	return err
}
