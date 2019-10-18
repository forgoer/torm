package torm

import (
	"testing"
)

func TestNewSchema(t *testing.T) {
	schema, err := NewSchema(&User{})

	if err != nil {
		t.Error(err)
	}

	t.Log(schema.Fields)
}
