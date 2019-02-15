package torm

import (
	"database/sql"
	"errors"
	"log"
	"reflect"

	"github.com/thinkoner/torm/grammar"
	"github.com/thinkoner/torm/utils"
)

type tabler interface {
	TableName() string
}

type Connection struct {
	DB          *sql.DB
	tablePrefix string
}

// Table Begin a fluent query against a database table.
func (c *Connection) Table(table string) *Builder {
	return c.Query().From(table)
}

// Model Begin a fluent query against a database Model.
func (c *Connection) Model(model interface{}) *Builder {
	var table string

	kind := reflect.Indirect(reflect.ValueOf(model)).Kind()

	if kind != reflect.Struct {
		panic("unsupported model, should be slice or struct")
	}

	if t, ok := model.(tabler); ok {
		table = t.TableName()
	} else {
		modelType := reflect.ValueOf(model).Type()
		if t, ok := reflect.New(modelType).Interface().(tabler); ok {
			table = t.TableName()
		} else {
			table = utils.SnakeCase(modelType.Name())
		}
	}

	return c.Table(table)
}

// Query Get a new query builder instance.
func (c *Connection) Query() *Builder {
	return NewBuilder(c, c.GetQueryGrammar())
}

// SelectOne Run a select statement and return a single result.
func (c *Connection) SelectOne(query string, bindings []interface{}, dest interface{}) error {
	return c.Select(query, bindings, dest)
}

// Select Run a select statement against the database.
func (c *Connection) Select(query string, bindings []interface{}, dest interface{}) error {

	log.Println(query)

	var err error

	stmt, err := c.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(bindings...)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	results := reflect.Indirect(reflect.ValueOf(dest))

	kind := results.Kind()
	if kind != reflect.Slice && kind != reflect.Struct {
		return errors.New("unsupported destination, should be slice or struct")
	}

	var isPtr bool
	var resultType reflect.Type
	var resultValue reflect.Value

	resultType = results.Type()
	resultValue = results

	if kind == reflect.Slice {
		resultType = resultType.Elem()
	}

	if resultType.Kind() == reflect.Ptr {
		resultType = resultType.Elem()
		isPtr = true
	}
	for rows.Next() {
		var fields []*Field

		if kind == reflect.Slice {
			resultValue = reflect.New(resultType).Elem()
		}

		for i := 0; i < resultType.NumField(); i++ {
			fields = append(fields, &Field{
				Value:       resultValue.Field(i),
				StructField: resultType.Field(i),
			})
		}

		err := c.scan(rows, columns, fields)
		if err != nil {
			return err
		}
		if kind == reflect.Slice {
			if isPtr {
				resultValue = resultValue.Addr()
			}
			results.Set(reflect.Append(results, resultValue))
		} else {
			break
		}
	}
	return nil
}

func (c *Connection) Scan(query string, bindings []interface{}, dest ...interface{}) error {
	log.Println(query)

	var err error

	stmt, err := c.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(bindings...)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()

	return rows.Scan(dest...)
}

// Insert Run an insert statement against the database.
func (c *Connection) Insert(query string, args ...interface{}) (int64, int64, error) {
	return c.AffectingStatement(query, args...)
}

// Update Run an update statement against the database.
func (c *Connection) Update(query string, args ...interface{}) (int64, error) {
	_, affected, err := c.AffectingStatement(query, args...)
	return affected, err
}

// Delete Run a delete statement against the database.
func (c *Connection) Delete(query string, args ...interface{}) (int64, error) {
	_, affected, err := c.AffectingStatement(query, args...)
	return affected, err
}

// AffectingStatement Run an SQL statement and get the number of rows affected.
func (c *Connection) AffectingStatement(query string, args ...interface{}) (int64, int64, error) {
	var err error

	log.Println(query)

	stmt, err := c.DB.Prepare(query)

	if err != nil {
		return 0, 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)

	if err != nil {
		return 0, 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, affected, err
	}

	insertId, err := result.LastInsertId()
	return insertId, affected, err
}

// Statement Execute an SQL statement and return the boolean result.
func (c *Connection) Statement(query string, args ...interface{}) (bool, error) {
	var err error

	stmt, err := c.DB.Prepare(query)

	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Connection) GetQueryGrammar() grammar.Grammar {
	return &grammar.MySqlGrammar{}
}

func (c *Connection) scan(rows *sql.Rows, columns []string, fields []*Field) error {
	// result := make(map[string]interface{})

	count := len(columns)
	values := make([]interface{}, count)
	args := make([]interface{}, count)

	for i, column := range columns {
		args[i] = &values[i]

		for _, field := range fields {
			if field.IsColumn(column) {
				args[i] = field.Value.Addr().Interface()
			}
		}
	}

	err := rows.Scan(args...)
	if err != nil {
		return err
	}
	return nil
}
