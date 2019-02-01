package grammar

import "sort"

type BaseGrammar struct {
}

// Parameter Get the appropriate query parameter place-holder for a value.
func (g *BaseGrammar) Parameter(value interface{}) interface{} {
	if g.IsExpression(value) {
		return value.(Expression).GetValue()
	}
	return "?"
}

// IsExpression Determine if the given value is a raw expression.
func (g *BaseGrammar) IsExpression(value interface{}) bool {
	_, ok := value.(Expression)
	return ok
}

// PrepareBindingsForUpdate Prepare the bindings for an update statement.
func (g *BaseGrammar) PrepareBindingsForUpdate(bindings map[string][]interface{}, values map[string]interface{}) []interface{} {

	var results []interface{}

	for _, val := range bindings["join"] {
		results = append(results, val)
	}

	var keys []string
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		results = append(results, values[key])
	}

	for kb, vb := range bindings {
		if kb == "join" {
			continue
		}
		if kb == "select" {
			continue
		}
		for _, val := range vb {
			results = append(results, val)
		}
	}
	return results
}
