package query

type Query struct {
	Distinct    bool
	Columns     []string
	From        string
	Joins       []*Join
	Wheres      []*Where
	Groups      []string
	Havings     []*Having
	Orders      []*Order
	Unions      []*Union
	UnionOrders []*UnionOrder
	Limit       uint64
	Offset      uint64
	UnionLimit  uint64
	UnionOffset uint64
	Aggregate   *Aggregate
	JoinClause  bool
}

type Aggregate struct {
	Function string
	Columns  []string
}

type Join struct {
	Type  string
	Table string
	Query *Query
}

type Where struct {
	Type     string
	Sql      string
	Column   string
	First    string
	Second   string
	Operator string
	Value    interface{}
	Values   []interface{}
	Boolean  string
	Not      bool
}

type Having struct {
	Type     string
	Sql      string
	Column   string
	Operator string
	Value    interface{}
	Boolean  string
}

type Order struct {
	Type      string
	Sql       string
	Column    string
	Direction string
}

type Union struct {
}

type UnionOrder struct {
	Type      string
	Sql       string
	Column    string
	Direction string
}
