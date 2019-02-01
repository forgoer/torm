package torm

type JoinClause struct {
	Builder
	Type          string
	Table         string
	ParentBuilder *Builder
}

func NewJoinClause(parentBuilder *Builder, typeStr string, table string) *JoinClause {
	j := &JoinClause{
		Type:          typeStr,
		Table:         table,
		ParentBuilder: parentBuilder,
	}
	j.Connection = parentBuilder.Connection
	j.grammar = parentBuilder.GetGrammar()
	return j
}
