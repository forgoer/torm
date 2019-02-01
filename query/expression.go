package query

// Expression query expression
type Expression struct {
	value interface{}
}

// Expr Create a new raw query expression.
func Expr(value interface{}) *Expression {
	return &Expression{
		value: value,
	}
}

// GetValue Get the value of the expression.
func (e *Expression) GetValue() interface{} {
	return e.value
}
