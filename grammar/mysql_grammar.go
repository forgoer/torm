package grammar

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/thinkoner/torm/query"
)

var selectComponents = [] string{
	"aggregate",
	"columns",
	"from",
	"joins",
	"wheres",
	"groups",
	"havings",
	"orders",
	"limit",
	"offset",
	// "unions",
	// "lock",
}

type MySqlGrammar struct {
	BaseGrammar
	tablePrefix string
}

// CompileSelect Compile a select query into SQL.
func (g *MySqlGrammar) CompileSelect(query *query.Query) string {
	// If the query does not have any columns set, we'll set the columns to the
	// * character to just get all of the columns from the database. Then we
	// can build the query and concatenate all the pieces together as one.
	original := query.Columns

	if query.Columns == nil {
		query.Columns = [] string{"*"}
	}

	// To compile the query, we'll spin through each component of the query and
	// see if that component exists. If it does we'll just call the compiler
	// function for the component which is responsible for making the SQL.
	sql := strings.TrimSpace(
		g.concatenate(
			g.compileComponents(query),
		),
	)
	query.Columns = original
	return sql
}

func (g *MySqlGrammar) WrapTable(table string) string {
	return g.tablePrefix + table
}

func (g *MySqlGrammar) Wrap(value string, prefixAlias bool) string {
	return g.wrapSegments(strings.Split(value, "."))
}

func (g *MySqlGrammar) wrapSegments(segments []string) string {
	for key, segment := range segments {
		if key == 0 && len(segments) > 1 {
			segments[key] = g.WrapTable(segment)
		} else {
			segments[key] = g.wapValue(segment)
		}
	}
	return strings.Join(segments, ".")
}

func (g *MySqlGrammar) wapValue(value string) string {
	if value != "*" {
		return "`" + strings.Replace(value, "`", "``", -1) + "`"
	}
	return value
}

func (g *MySqlGrammar) compileComponents(query *query.Query) []string {
	var sql []string
	for _, component := range selectComponents {
		switch component {
		case "aggregate":
			if query.Aggregate != nil {
				sql = append(sql, g.compileAggregate(query, query.Aggregate))
			}
		case "columns":
			if len(query.Columns) > 0 {
				sql = append(sql, g.compileColumns(query, query.Columns))
			}
		case "from":
			if len(query.From) > 0 {
				sql = append(sql, g.compileFrom(query, query.From))
			}
		case "joins":
			if len(query.Joins) > 0 {
				sql = append(sql, g.compileJoins(query, query.Joins))
			}
		case "wheres":
			if len(query.Wheres) > 0 {
				sql = append(sql, g.compileWheres(query))
			}
		case "groups":
			if len(query.Groups) > 0 {
				sql = append(sql, g.compileGroups(query, query.Groups))
			}
		case "havings":
			if len(query.Havings) > 0 {
				sql = append(sql, g.compileHavings(query, query.Havings))
			}
		case "orders":
			if len(query.Orders) > 0 {
				sql = append(sql, g.compileOrders(query, query.Orders))
			}

		case "limit":
			if query.Limit > 0 {
				sql = append(sql, g.compileLimit(query, query.Limit))
			}
		case "offset":
			if query.Limit > 0 {
				sql = append(sql, g.compileOffset(query, query.Offset))
			}
		}
	}
	return sql
}

func (g *MySqlGrammar) concatenate(segments []string) string {
	s := ""
	for _, segment := range segments {
		if len(segment) == 0 {
			continue
		}
		if len(s) > 0 {
			s = s + " "
		}
		s = s + segment
	}
	return s
}

func (g *MySqlGrammar) compileAggregate(query *query.Query, aggregate *query.Aggregate) string {
	column := strings.Join(aggregate.Columns, ", ")
	if query.Distinct && column != "*" {
		column = "DISTINCT " + column
	}

	return "SELECT " + aggregate.Function + "(" + column + ") AS aggregate"
}

func (g *MySqlGrammar) compileColumns(query *query.Query, columns []string) string {
	// If the query is actually performing an aggregating select, we will let that
	// compiler handle the building of the select clauses, as it will need some
	// more syntax that is best handled by that function to keep things neat.
	if query.Aggregate != nil {
		return ""
	}

	sel := "SELECT "
	if query.Distinct {
		sel = "SELECT DISTINCT "
	}

	return sel + strings.Join(columns, ", ")
}

func (g *MySqlGrammar) compileFrom(query *query.Query, table string) string {
	return "FROM " + g.WrapTable(table)
}

func (g *MySqlGrammar) compileJoins(query *query.Query, joins []*query.Join) string {
	var segments []string
	for _, join := range joins {
		table := g.WrapTable(join.Table)
		segments = append(segments, strings.TrimSpace(join.Type+" JOIN "+table+" "+g.compileWheres(join.Query)))
	}
	return strings.Join(segments, " ")
}

func (g *MySqlGrammar) compileWheres(query *query.Query) string {

	sql := g.compileWheresToArray(query)
	if len(sql) > 0 {
		return g.concatenateWhereClauses(query, sql)
	}

	return ""
}

func (g *MySqlGrammar) compileWheresToArray(query *query.Query) []string {
	var sql []string
	for _, where := range query.Wheres {
		w := ""
		switch where.Type {
		case "Basic":
			w = where.Boolean + " " + g.whereBasic(query, where)
		case "Column":
			w = where.Boolean + " " + g.whereColumn(query, where)
		}
		sql = append(sql, w)
	}
	return sql
}

func (g *MySqlGrammar) concatenateWhereClauses(query *query.Query, sql []string) string {
	conjunction := "WHERE"
	if query.JoinClause {
		conjunction = "ON"
	}
	return conjunction + " " + removeLeadingBoolean(strings.Join(sql, " "))
}

func (g *MySqlGrammar) whereBasic(query *query.Query, where *query.Where) string {
	// value = where.Value
	return g.Wrap(where.Column, false) + " " + where.Operator + " " + "?"
}

func (g *MySqlGrammar) whereColumn(query *query.Query, where *query.Where) string {
	return g.Wrap(where.First, false) + " " + where.Operator + " " + g.Wrap(where.Second, false)
}

func (g *MySqlGrammar) compileGroups(query *query.Query, groups []string) string {
	return "GROUP BY " + strings.Join(groups, ", ")
}

func (g *MySqlGrammar) compileHavings(query *query.Query, havings []*query.Having) string {
	sqls := make([]string, 0)
	for _, having := range havings {
		sqls = append(sqls, g.compileHaving(having))
	}
	sql := strings.Join(sqls, " ")
	return "HAVING " + removeLeadingBoolean(sql)
}

func (g *MySqlGrammar) compileHaving(having *query.Having) string {
	if having.Type == "Raw" {
		return having.Boolean + " " + having.Sql
	}
	return g.compileBasicHaving(having)
}

func (g *MySqlGrammar) compileBasicHaving(having *query.Having) string {
	return having.Boolean + " " + having.Column + " " + having.Operator + " " + "?"
}

func (g *MySqlGrammar) compileOrders(query *query.Query, orders []*query.Order) string {
	var sql []string

	for _, order := range orders {
		s := ""
		if len(order.Sql) > 0 {
			s = order.Sql
		} else {
			s = order.Column + " " + order.Direction
		}
		sql = append(sql, s)
	}

	if len(sql) == 0 {
		return ""
	}

	return "ORDER BY " + strings.Join(sql, ", ")
}

func (g *MySqlGrammar) compileLimit(query *query.Query, limit uint64) string {
	return fmt.Sprintf("LIMIT %v", limit)
}

func (g *MySqlGrammar) compileOffset(query *query.Query, offset uint64) string {
	return fmt.Sprintf("OFFSET %v", offset)
}

// CompileInsert Compile an insert statement into SQL.
func (g *MySqlGrammar) CompileInsert(query *query.Query, values []map[string]interface{}) (string, []interface{}) {
	table := g.WrapTable(query.From)
	var columns []string
	var parameters [] string
	var bindings []interface{}

	if len(values) > 0 {
		first := values[0]
		for k := range first {
			columns = append(columns, k)
		}
		for _, val := range values {
			var vals []string
			for _, column := range columns {
				if col, ok := val[column]; ok {
					bindings = append(bindings, col)
					vals = append(vals, "?")
				}
			}
			if len(vals) > 0 {
				parameters = append(parameters, fmt.Sprintf("(%s)", strings.Join(vals, ", ")))
			}
		}
	}

	return fmt.Sprintf("INSERT INTO %s (%s) values %s", table, strings.Join(columns, ", "), strings.Join(parameters, ",")), bindings
}

// CompileUpdate Compile an update statement into SQL.
func (g *MySqlGrammar) CompileUpdate(query *query.Query, values map[string]interface{}) string {
	table := g.WrapTable(query.From)
	columns := g.compileUpdateColumns(values)

	joins := ""
	if len(query.Joins) > 0 {
		joins = " " + g.compileJoins(query, query.Joins)
	}

	where := g.compileWheres(query)

	sql := strings.TrimSpace(fmt.Sprintf("UPDATE %s%s SET %s %s", table, joins, columns, where))

	if len(query.Orders) > 0 {
		sql = sql + " " + g.compileOrders(query, query.Orders)
	}

	if query.Limit > 0 {
		sql = sql + " " + g.compileLimit(query, query.Limit)
	}

	return sql
}

// compileUpdateColumns Compile all of the columns for an update statement.
func (g *MySqlGrammar) compileUpdateColumns(values map[string]interface{}) string {
	var columns []string

	var keys []string
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		columns = append(columns, fmt.Sprintf("%s = %s", g.Wrap(key, false), g.Parameter(values[key])))
	}
	return strings.Join(columns, ", ")
}

// PrepareBindingsForUpdate Prepare the bindings for an update statement.
func (g *MySqlGrammar) PrepareBindingsForUpdate(bindings map[string][]interface{}, values map[string]interface{}) []interface{} {
	return g.BaseGrammar.PrepareBindingsForUpdate(bindings, values)
}

// CompileDelete Compile a delete statement into SQL.
func (g *MySqlGrammar) CompileDelete(query *query.Query) string {
	table := g.WrapTable(query.From)
	where := ""

	if len(query.Wheres) > 0 {
		where = g.compileWheres(query)
	}

	if len(query.Joins) > 0 {
		return g.compileDeleteWithJoins(query, table, where)
	} else {
		return g.compileDeleteWithoutJoins(query, table, where)
	}
}

// compileDeleteWithoutJoins Compile a delete query that does not use joins.
func (g *MySqlGrammar) compileDeleteWithoutJoins(query *query.Query, table string, where string) string {
	sql := strings.TrimSpace(fmt.Sprintf("DELETE FROM %s %s", table, where))

	if len(query.Orders) > 0 {
		sql = sql + " " + g.compileOrders(query, query.Orders)
	}

	if query.Limit > 0 {
		sql = sql + " " + g.compileLimit(query, query.Limit)
	}

	return sql
}

// compileDeleteWithJoins Compile a delete query that uses joins.
func (g *MySqlGrammar) compileDeleteWithJoins(query *query.Query, table string, where string) string {
	joins := " " + g.compileJoins(query, query.Joins)

	alias := table

	if strings.Contains(strings.ToLower(table), " as ") {
		alias = strings.Split(table, " as ")[1]
	}

	return strings.TrimSpace(fmt.Sprintf("DELETE %s FROM %s%s %s", alias, table, joins, where))
}

func removeLeadingBoolean(value string) string {
	reg := regexp.MustCompile(`(?i:and |or )`)
	n := 0
	b := reg.ReplaceAllFunc([]byte(value), func(bytes []byte) []byte {
		n = n + 1
		if n > 1 {
			return bytes
		}
		return []byte("")
	})
	return string(b)
	// return reg.ReplaceAllString(value, "")
}
