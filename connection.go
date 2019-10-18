package torm

import (
	"database/sql"
	"errors"
	"log"
	"reflect"

	"github.com/thinkoner/torm/grammar"
)

type Connection struct {
	DB          *sql.DB
	tablePrefix string
}

// Table Begin a fluent query against a database table.
func (c *Connection) Table(table string) *Builder {
	return c.Query().From(table)
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

	if kind == reflect.Slice {
		resultValue = reflect.New(resultType).Elem()
	}

	var fields []*Field

	for i := 0; i < resultType.NumField(); i++ {
		fields = append(fields, NewField(resultValue.Field(i), resultType.Field(i)))
	}

	for rows.Next() {
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
	return c.affectingStatement(query, args...)
}

// Update Run an update statement against the database.
func (c *Connection) Update(query string, args ...interface{}) (int64, error) {
	_, affected, err := c.affectingStatement(query, args...)
	return affected, err
}

// Delete Run a delete statement against the database.
func (c *Connection) Delete(query string, args ...interface{}) (int64, error) {
	_, affected, err := c.affectingStatement(query, args...)
	return affected, err
}

// Statement Execute an SQL statement and return the boolean result.
func (c *Connection) Statement(query string, args ...interface{}) error {
	var err error

	stmt, err := c.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)

	return err
}

func (c *Connection) GetQueryGrammar() grammar.Grammar {
	return &grammar.MySqlGrammar{}
}

// AffectingStatement Run an SQL statement and get the number of rows affected.
func (c *Connection) affectingStatement(query string, args ...interface{}) (int64, int64, error) {
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

func (c *Connection) scan(rows *sql.Rows, columns []string, fields []*Field) error {
	count := len(columns)
	resets := make(map[int]*Field)
	values := make([]interface{}, count)
	args := make([]interface{}, count)

	for i, column := range columns {
		args[i] = &values[i]
		for _, field := range fields {
			if field.Ignored {
				continue
			}
			if field.Name != column {
				continue
			}

			if field.Value.Kind() == reflect.Ptr {
				args[i] = field.Value.Addr().Interface()
			} else {
				reflectValue := reflect.New(reflect.PtrTo(field.StructField.Type))
				reflectValue.Elem().Set(field.Value.Addr())
				args[i] = reflectValue.Interface()
				resets[i] = field
			}
		}
	}

	err := rows.Scan(args...)

	if err != nil {
		return err
	}

	for index, field := range resets {
		if v := reflect.ValueOf(args[index]).Elem().Elem(); v.IsValid() {
			field.Value.Set(v)
		}
	}

	return nil
}
