package grammar

import "github.com/thinkoner/torm/query"

type Grammar interface {
	// CompileSelect Compile a select query into SQL.
	CompileSelect(query *query.Query) string

	// CompileInsert Compile an insert statement into SQL.
	CompileInsert(query *query.Query, values []map[string]interface{}) (string, []interface{})

	// CompileUpdate Compile an update statement into SQL.
	CompileUpdate(query *query.Query, values map[string]interface{}) string

	// CompileDelete Compile a delete statement into SQL.
	CompileDelete(query *query.Query) string

	// PrepareBindingsForUpdate Prepare the bindings for an update statement.
	PrepareBindingsForUpdate(bindings map[string][]interface{}, values map[string]interface{}) []interface{}
}
