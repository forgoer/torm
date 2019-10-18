package torm

import (
	"reflect"
	"testing"
)

type StructForField struct {
	ID      uint   `torm:"primary_key" json:"id"`
	Name    string `torm:"column:name"`
	Gender  string `torm:"column:sex"`
	Addr    string
	Balance float64     `torm:"NOT NULL;DEFAULT:0"`
	Ignored interface{} `torm:"-"`
}

func TestNewField(t *testing.T) {
	s := &StructForField{
		ID:      1,
		Name:    "testing",
		Gender:  "F",
		Balance: 366.36,
	}

	results := reflect.Indirect(reflect.ValueOf(s))

	resultType := results.Type()
	resultValue := results

	fields := make(map[string]*Field)

	for i := 0; i < resultType.NumField(); i++ {
		fields[resultType.Field(i).Name] = NewField(resultValue.Field(i), resultType.Field(i))
	}

	f := fields["ID"]
	if f.Name != "id" || f.Ignored || !f.Primary {
		t.Errorf(`TestNewField: fields["ID"]'s name should be %s and cannot be ignored and is the primary key.`, f.Name)
	}
	err := f.SetValue(8)
	if err != nil {
		t.Error(err)
	}

	f = fields["Gender"]
	if f.Name != "sex" {
		t.Errorf(`TestNewField: fields["Gender"]'s name should be %s and cannot be ignored and is the primary key.`, f.Name)
	}

	f = fields["Ignored"]
	if !f.Ignored {
		t.Errorf(`TestNewField: fields["Ignored"]'s name should be ignored.`)
	}

	f = fields["Balance"]
	_, ok := f.GetAttr("not null")
	if !ok {
		t.Errorf(`TestNewField: fields["Balance"]'s 'not null' attribution should find.`)
	}
	str, _ := f.GetAttr("default")
	if str != "0" {
		t.Errorf(`TestNewField: fields["Balance"]'s 'default' attribution should be '0'.`)
	}

}
