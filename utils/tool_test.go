package utils

import (
	"reflect"
	"testing"
)

func TestIsBlank(t *testing.T) {
	value := reflect.Indirect(reflect.ValueOf(0))
	res := IsBlank(value)
	if !res {
		t.Errorf("%d should by blank", 0)
	}

	value = reflect.Indirect(reflect.ValueOf(0.00))
	res = IsBlank(value)
	if !res {
		t.Errorf("%f should by blank", 0.00)
	}

	value = reflect.Indirect(reflect.ValueOf(""))
	res = IsBlank(value)
	if !res {
		t.Errorf("%s should by blank", `""`)
	}
}
