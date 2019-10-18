package torm

import "github.com/thinkoner/torm/field"

type Model struct {
	field.IDAttr
	field.CreatedAtAttr
	field.UpdatedAtAttr
}

type SoftDeletes struct {
	field.DeletedAtAttr
}
