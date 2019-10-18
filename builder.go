package torm

import (
	"reflect"
	"strings"

	"errors"
	"github.com/thinkoner/torm/grammar"
	"github.com/thinkoner/torm/query"
)

type Builder struct {
	Connection *Connection
	grammar    grammar.Grammar
	Query      *query.Query
	Bindings   map[string][]interface{}
}

type Binding []interface{}

func NewBuilder(connection *Connection, grammar grammar.Grammar) *Builder {

	return &Builder{
		Connection: connection,
		grammar:    grammar,
		Query:      &query.Query{},
		Bindings:   getDefaultBindings(),
	}
}

// Select the columns to be selected.
func (b *Builder) Select(columns ...string) *Builder {
	if len(columns) == 0 {
		columns = []string{"*"}
	}
	b.Query.Columns = columns

	return b
}

// SelectRaw Add a new "raw" select expression to the query.
func (b *Builder) SelectRaw(expression string, bindings ...interface{}) *Builder {
	b.AddSelect(expression)

	if len(bindings) != 0 {
		b.AddBinding(bindings, "select")
	}

	return b
}

func (b *Builder) AddSelect(columns ...string) *Builder {
	for _, column := range columns {
		b.Query.Columns = append(b.Query.Columns, column)
	}

	return b
}

// From Set the table which the query is targeting.
func (b *Builder) From(table string) *Builder {
	b.Query.From = table
	return b
}

// Join Add a join clause to the query.
func (b *Builder) Join(table string, first string, args ...interface{}) *Builder {
	var (
		operator string
		second   string
		typtStr  string
		where    bool
	)

	typtStr = "INNER"

	count := len(args)

	if count >= 1 {
		operator = args[0].(string)
	}
	if count >= 2 {
		second = args[1].(string)
	}
	if count >= 3 {
		typtStr, _ = args[2].(string)
	}
	if count >= 4 {
		where, _ = args[2].(bool)
	}
	q := &query.Query{}

	if where {
		where := &query.Where{
			Type:     "Basic",
			Column:   first,
			Operator: operator,
			Value:    second,
		}
		q.Wheres = append(q.Wheres, where)
		b.AddBinding(second, "join")
	} else {
		where := &query.Where{
			Type:     "Column",
			First:    first,
			Operator: operator,
			Second:   second,
			Boolean:  "and",
		}
		q.JoinClause = true
		q.Wheres = append(q.Wheres, where)
	}

	b.Query.Joins = append(b.Query.Joins, &query.Join{
		Type:  typtStr,
		Table: table,
		Query: q,
	})
	return b
}

// func (b *Builder) Value(column string) error {
// 	result, err := b.First([]string{column})
// 	if len(result) > 0 {
// 		return result[column], err
// 	}
// 	return nil, err
// }

func (b *Builder) First(dest interface{}, columns ...interface{}) error {
	return b.Take(1).Get(dest, columns...)
}

func (b *Builder) Get(dest interface{}, columns ...interface{}) error {
	var cols []string

	if len(columns) == 0 {
		cols = []string{"*"}
	} else {
		if _, ok := columns[0].([]string); ok {
			cols = columns[0].([]string)
		} else {
			return errors.New("Invalid parameters")
		}
	}
	if len(cols) == 0 {
		cols = []string{"*"}
	}

	original := b.Query.Columns

	if original == nil {
		b.Query.Columns = cols
	}

	err := b.runSelect(dest)

	b.Query.Columns = original

	return err
}

func (b *Builder) Scan(dest ...interface{}) error {
	return b.runScan(dest...)
}

// Where Add a basic where clause to the query.
func (b *Builder) Where(column string, args ...interface{}) *Builder {
	var (
		operator string
		value    interface{}
		boolean  string
	)

	boolean = "and"

	count := len(args)

	if count == 1 {
		operator = "="
		value = args[0]
	}
	if count >= 2 {
		operator, _ = args[0].(string)
		value = args[1]
	}
	if count >= 3 {
		boolean, _ = args[2].(string)
	}

	where := &query.Where{
		Type:     "Basic",
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  boolean,
	}
	b.Query.Wheres = append(b.Query.Wheres, where)

	b.AddBinding(value, "where")

	return b
}

// OrWhere Add an "or where" clause to the query.
func (b *Builder) OrWhere(column string, args ...interface{}) *Builder {
	var (
		operator string
		value    interface{}
	)

	count := len(args)

	if count == 1 {
		operator = "="
		value = args[0]
	}
	if count >= 2 {
		operator, _ = args[0].(string)
		value = args[1]
	}

	return b.Where(column, operator, value, "OR")
}

// WhereColumn Add a "where" clause comparing two columns to the query.
func (b *Builder) WhereColumn(first string, args ...interface{}) *Builder {
	var (
		operator string
		value    string
		boolean  string
	)
	boolean = "AND"
	count := len(args)

	if count == 1 {
		operator = "="
		value = args[0].(string)
	}
	if count >= 2 {
		operator, _ = args[0].(string)
		value = args[1].(string)
	}
	if count >= 3 {
		boolean, _ = args[2].(string)
	}

	where := &query.Where{
		Type:     "Column",
		First:    first,
		Operator: operator,
		Second:   value,
		Boolean:  boolean,
	}
	b.Query.Wheres = append(b.Query.Wheres, where)

	return b
}

// OrWhereColumn Add an "or where" clause comparing two columns to the query.
func (b *Builder) OrWhereColumn(first string, args ...interface{}) *Builder {
	var (
		operator string
		value    string
	)

	count := len(args)

	if count == 1 {
		operator = "="
		value = args[0].(string)
	}
	if count >= 2 {
		operator, _ = args[0].(string)
		value = args[1].(string)
	}
	return b.WhereColumn(first, operator, value, "OR")
}

// WhereRaw Add a raw where clause to the query.
func (b *Builder) WhereRaw(sql string, bindings ...interface{}) *Builder {

	b.Query.Wheres = append(
		b.Query.Wheres,
		&query.Where{
			Type:    "Raw",
			Sql:     sql,
			Boolean: "AND",
		},
	)

	if len(bindings) != 0 {
		b.AddBinding(bindings, "where")
	}

	return b
}

// OrWhereRaw Add a raw or where clause to the query.
func (b *Builder) OrWhereRaw(sql string, bindings ...interface{}) *Builder {

	b.Query.Wheres = append(
		b.Query.Wheres,
		&query.Where{
			Type:    "Raw",
			Sql:     sql,
			Boolean: "OR",
		},
	)

	if len(bindings) != 0 {
		b.AddBinding(bindings, "where")
	}

	return b
}

// WhereIn Add a "where in" clause to the query.
func (b *Builder) WhereIn(column string, values []interface{}, args ...interface{}) *Builder {
	var (
		boolean string
		not     bool
	)
	count := len(args)

	if count == 0 {
		boolean = "and"
		not = false
	}
	if count == 1 {
		boolean = args[0].(string)
		not = false
	}
	if count >= 2 {
		boolean = args[0].(string)
		not = args[1].(bool)
	}
	t := "In"
	if not {
		t = "NotIn"
	}
	b.Query.Wheres = append(
		b.Query.Wheres,
		&query.Where{
			Type:    t,
			Column:  column,
			Values:  values,
			Boolean: boolean,
		},
	)
	for _, value := range values {
		b.AddBinding(value, "where")
	}

	return b
}

// OrWhereIn Add an "or where in" clause to the query.
func (b *Builder) OrWhereIn(column string, values []interface{}) *Builder {
	return b.WhereIn(column, values, "OR")
}

// WhereNotIn Add a "where not in" clause to the query.
func (b *Builder) WhereNotIn(column string, values []interface{}, args ...interface{}) *Builder {
	var boolean string
	count := len(args)

	if count == 0 {
		boolean = "AND"
	}
	if count == 1 {
		boolean = args[0].(string)
	}
	return b.WhereIn(column, values, boolean, true)
}

// OrWhereNotIn Add an "or where not in" clause to the query.
func (b *Builder) OrWhereNotIn(column string, values []interface{}) *Builder {
	return b.WhereNotIn(column, values, "OR")
}

// WhereNull Add a "where null" clause to the query.
func (b *Builder) WhereNull(column string, args ...interface{}) *Builder {
	var (
		boolean string
		not     bool
	)
	count := len(args)

	if count == 0 {
		boolean = "and"
		not = false
	}
	if count == 1 {
		boolean = args[0].(string)
		not = false
	}
	if count >= 2 {
		boolean = args[0].(string)
		not = args[1].(bool)
	}
	t := "Null"
	if not {
		t = "NotNull"
	}
	b.Query.Wheres = append(
		b.Query.Wheres,
		&query.Where{
			Type:    t,
			Column:  column,
			Boolean: boolean,
		},
	)
	return b
}

// OrWhereNull Add an "or where null" clause to the query.
func (b *Builder) OrWhereNull(column string) *Builder {
	return b.WhereNull(column, "OR")
}

// WhereNotNull Add a "where not null" clause to the query.
func (b *Builder) WhereNotNull(column string, args ...interface{}) *Builder {
	var (
		boolean string
	)
	count := len(args)

	if count == 0 {
		boolean = "and"
	}
	if count == 1 {
		boolean = args[0].(string)
	}
	return b.WhereNull(column, boolean, true)
}

// OrWhereNotNull Add an "or where not null" clause to the query.
func (b *Builder) OrWhereNotNull(column string) *Builder {
	return b.WhereNotNull(column, "OR")
}

// WhereBetween Add a where between statement to the query.
func (b *Builder) WhereBetween(column string, values []interface{}, args ...interface{}) *Builder {
	var (
		boolean string
		not     bool
	)
	count := len(args)

	if count == 0 {
		boolean = "and"
		not = false
	}
	if count == 1 {
		boolean = args[0].(string)
		not = false
	}
	if count >= 2 {
		boolean = args[0].(string)
		not = args[1].(bool)
	}
	b.Query.Wheres = append(
		b.Query.Wheres,
		&query.Where{
			Type:    "Between",
			Column:  column,
			Boolean: boolean,
			Not:     not,
		},
	)
	b.AddBinding(values, "where")

	return b
}

// OrWhereBetween Add an or where between statement to the query.
func (b *Builder) OrWhereBetween(column string, values []interface{}) *Builder {
	return b.WhereBetween(column, values, "OR")
}

// WhereNotBetween Add a where not between statement to the query.
func (b *Builder) WhereNotBetween(column string, values []interface{}, args ...interface{}) *Builder {
	var (
		boolean string
	)
	count := len(args)

	if count == 0 {
		boolean = "and"
	}
	if count == 1 {
		boolean = args[0].(string)
	}
	return b.WhereBetween(column, values, boolean, true)
}

// OrWhereNotBetween Add an or where not between statement to the query.
func (b *Builder) OrWhereNotBetween(column string, values []interface{}) *Builder {
	return b.WhereNotBetween(column, values, "OR")
}

// GroupBy Add a "group by" clause to the query.
func (b *Builder) GroupBy(groups ...string) *Builder {
	if len(groups) != 0 {
		for _, group := range groups {
			b.Query.Groups = append(b.Query.Groups, group)
		}
	}

	return b
}

// Having Add a "having" clause to the query.
func (b *Builder) Having(column string, args ...interface{}) *Builder {
	var (
		operator string
		value    interface{}
		boolean  string
	)

	boolean = "and"

	count := len(args)

	if count == 1 {
		operator = "="
		value = args[0]
	}
	if count >= 2 {
		operator, _ = args[0].(string)
		value = args[1]
	}
	if count >= 3 {
		boolean, _ = args[2].(string)
	}

	having := &query.Having{
		Type:     "Basic",
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  boolean,
	}
	b.Query.Havings = append(b.Query.Havings, having)

	b.AddBinding(value, "having")

	return b
}

// OrderBy Add an "order by" clause to the query.
func (b *Builder) OrderBy(column string, args ...string) *Builder {
	direction := "asc"

	if len(args) > 0 {
		direction = args[0]
	}

	if len(b.Query.Unions) == 0 {
		b.Query.Orders = append(
			b.Query.Orders,
			&query.Order{
				Column:    column,
				Direction: strings.ToUpper(direction),
			},
		)
	} else {
		b.Query.UnionOrders = append(
			b.Query.UnionOrders,
			&query.UnionOrder{
				Column:    column,
				Direction: strings.ToUpper(direction),
			},
		)
	}
	return b
}

// OrderByDesc Add a descending "order by" clause to the query.
func (b *Builder) OrderByDesc(column string) *Builder {
	return b.OrderBy(column, "desc")
}

// OrderByRaw Add a raw "order by" clause to the query.
func (b *Builder) OrderByRaw(sql string, bindings ...interface{}) *Builder {
	if len(b.Query.Unions) == 0 {
		b.Query.Orders = append(
			b.Query.Orders,
			&query.Order{
				Type: "Raw",
				Sql:  sql,
			},
		)
	} else {
		b.Query.UnionOrders = append(
			b.Query.UnionOrders,
			&query.UnionOrder{
				Type: "Raw",
				Sql:  sql,
			},
		)
	}

	if len(bindings) != 0 {
		b.AddBinding(bindings, "order")
	}

	return b
}

// Take Alias to set the "limit" value of the query.
func (b *Builder) Take(value uint64) *Builder {
	b.Limit(value)
	return b
}

// Limit Set the "limit" value of the query.
func (b *Builder) Limit(value uint64) *Builder {

	if len(b.Query.Unions) > 0 {
		b.Query.UnionLimit = value
	} else {
		b.Query.Limit = value
	}

	return b
}

// Skip Alias to set the "offset" value of the query.
func (b *Builder) Skip(value uint64) *Builder {
	return b.Offset(value)
}

// Offset Set the "offset" value of the query.
func (b *Builder) Offset(value uint64) *Builder {
	if len(b.Query.Unions) > 0 {
		b.Query.UnionOffset = value
	} else {
		b.Query.Offset = value
	}

	return b
}

// Count Retrieve the "count" result of the query.
func (b *Builder) Count(dest ...interface{}) error {
	return b.Aggregate("COUNT", []string{"*"}, dest...)
}

// Min Retrieve the minimum value of a given column.
func (b *Builder) Min(columns string, dest ...interface{}) error {
	return b.Aggregate("MIN", []string{columns}, dest...)
}

// Max Retrieve the maximum value of a given column.
func (b *Builder) Max(columns string, dest ...interface{}) error {
	return b.Aggregate("MAX", []string{columns}, dest...)
}

// Sum Retrieve the sum of the values of a given column.
func (b *Builder) Sum(columns string, dest ...interface{}) error {
	return b.Aggregate("SUM", []string{columns}, dest...)
}

// Avg Retrieve the average of the values of a given column.
func (b *Builder) Avg(columns string, dest ...interface{}) error {
	return b.Aggregate("AVG", []string{columns}, dest...)
}

// Aggregate Execute an aggregate function on the database.
func (b *Builder) Aggregate(function string, columns []string, dest ...interface{}) error {
	cols := []string{"*"}

	if len(columns) > 0 {
		cols = columns
	}

	nb := b.CloneWithout("columns").CloneWithoutBindings("select").setAggregate(function, cols)

	original := nb.Query.Columns

	if original == nil {
		nb.Query.Columns = cols
	}

	return nb.Scan(dest...)
}

func (b *Builder) ToSql() string {
	return b.grammar.CompileSelect(b.Query)
}

// AddBinding Add a binding to the query.
func (b *Builder) AddBinding(value interface{}, segment string) {
	if _, ok := b.Bindings[segment]; !ok {
		return
	}
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(value)
		for i := 0; i < s.Len(); i++ {
			b.Bindings[segment] = append(b.Bindings[segment], s.Index(i).Interface())
		}
	default:
		b.Bindings[segment] = append(b.Bindings[segment], value)
	}
}

func (b *Builder) GetBindings() []interface{} {
	var bindings []interface{}
	for _, val := range b.Bindings {
		for _, v := range val {
			bindings = append(bindings, v)
		}
	}
	return bindings
}

// CloneWithout Clone the query without the given properties.
func (b *Builder) CloneWithout(properties ...string) *Builder {
	builder := b
	for _, property := range properties {
		switch property {
		case "columns":
			builder.Query.Columns = nil
		case "orders":
			// builder.Query.Columns = nil
		case "limit":
			builder.Query.Limit = 0
		case "offset":
			// builder.Query.Columns = nil
		}
	}
	return builder
}

// CloneWithoutBindings Clone the query without the given bindings.
func (b *Builder) CloneWithoutBindings(except ...string) *Builder {
	builder := b
	for _, val := range except {
		builder.Bindings[val] = nil
	}
	return builder
}

// GetConnection Get the database connection instance.
func (b *Builder) GetConnection() *Connection {
	return b.Connection
}

// GetGrammar Get the query grammar instance.
func (b *Builder) GetGrammar() grammar.Grammar {
	return b.grammar
}

func (b *Builder) runSelect(dest interface{}) error {
	return b.Connection.Select(
		b.ToSql(),
		b.GetBindings(),
		dest,
	)
}

func (b *Builder) runScan(dest ...interface{}) error {
	return b.Connection.Scan(
		b.ToSql(),
		b.GetBindings(),
		dest...,
	)
}

func (b *Builder) setAggregate(function string, columns []string) *Builder {
	b.Query.Aggregate = &query.Aggregate{
		Function: function,
		Columns:  columns,
	}
	return b
}

// Insert Insert a new record into the database.
func (b *Builder) Insert(value map[string]interface{}) (int64, int64, error) {
	var values []map[string]interface{}
	values = append(values, value)
	return b.Inserts(values)
}

// Inserts Bulk insert records into the database.
func (b *Builder) Inserts(values []map[string]interface{}) (int64, int64, error) {
	sql, bindings := b.GetGrammar().CompileInsert(b.Query, values)
	return b.Connection.Insert(sql, bindings...)
}

// Update a record in the database.
func (b *Builder) Update(value map[string]interface{}) (int64, error) {
	sql := b.GetGrammar().CompileUpdate(b.Query, value)
	cleanBindings := cleanBindings(b.GetGrammar().PrepareBindingsForUpdate(b.Bindings, value))
	return b.Connection.Update(
		sql,
		cleanBindings...,
	)
}

// Delete a record from the database.
func (b *Builder) Delete(args ...interface{}) (int64, error) {
	if len(args) > 0 {
		b.Where(b.Query.From+".id", args[0])
	}

	return b.Connection.Delete(
		b.GetGrammar().CompileDelete(b.Query),
		b.GetBindings()...,
	)
}

// cleanBindings Remove all of the expressions from a list of bindings.
func cleanBindings(bindings []interface{}) []interface{} {
	var result []interface{}
	for _, b := range bindings {
		if _, ok := b.(grammar.Expression); !ok {
			result = append(result, b)
		}
	}
	return result
}

func getDefaultBindings() map[string][]interface{} {
	return map[string][]interface{}{
		"select": make([]interface{}, 0),
		"join":   make([]interface{}, 0),
		"where":  make([]interface{}, 0),
		"having": make([]interface{}, 0),
		"order":  make([]interface{}, 0),
		"union":  make([]interface{}, 0),
	}
}
