package torm

import (
	"log"
	"strings"
	"testing"
	"time"
)

func TestTableSelectBasic(t *testing.T) {
	var users []User
	err := DB.Table("users").Where("gender", "M").Get(&users)

	if err != nil {
		t.Error(err)
	}
	for _, user := range users {
		if user.Gender != "M" {
			t.Error("Expect: user's gender should be ", "M")
		}
	}
}

func TestTableSelectOne(test *testing.T) {
	var user User
	err := DB.Table("users").First(&user)

	if err != nil {
		test.Error(err)
	}
	log.Println(user)
}

func TestTableSelectAdvance(t *testing.T) {
	var users []User
	err := DB.Table("users").
		Where("gender", "M").
		Where("name", "like", "%a%").
		WhereColumn("created_at", "!=", "updated_at").
		WhereIn("addr", []interface{}{"Columbia", "Alaska"}).
		WhereBetween("birth_date", []interface{}{"1990-01-01", "1999-12-31"}).
		Get(&users)

	if err != nil {
		t.Error(err)
	}
	for _, user := range users {
		t.Log(user)
		if user.Gender != "M" {
			t.Error("Expect: user's gender should be ", "M")
		}
		if !strings.Contains(user.Name, "a") {
			t.Error("Expect: user's name should contains ", "a")
		}
		if user.CreatedAt == user.UpdatedAt {
			t.Error("Expect: user's CreatedAt should != UpdatedAt ")
		}
		if user.Addr != "Columbia" && user.Addr != "Alaska" {
			t.Error("Expect: user's addr should be Columbia or Alaska ")
		}
		if !strings.HasPrefix(user.BirthDate, "199") {
			t.Error("Expect: user's BirthDate should start with '199' ")
		}
	}
}

func TestTableGroup(t *testing.T) {
	type UserAddr struct {
		Name  string
		Count int64 `torm:"column:ct"`
	}
	var userAddr []UserAddr
	err := DB.Table("users").
		Select("name", "count(name) as ct").
		GroupBy("name").
		Having("ct", ">", 1).
		Get(&userAddr)

	if err != nil {
		t.Error(err)
	}
	t.Log(userAddr)
}

func TestTableOrderBy(t *testing.T) {
	var users []User
	err := DB.Table("users").
		OrderByDesc("balance").
		Get(&users)

	if err != nil {
		t.Error(err)
	}
	last := users[len(users)-1]
	if last.Balance != 0.00 {
		t.Error("Expect: users should be order by `balance` desc")
	}
}

func TestTableAggregates(t *testing.T) {
	var count int64
	err := DB.Table("users").Count(&count)
	if err != nil {
		t.Error(err)
	}
	t.Log(count)
}

func TestModelSelect(t *testing.T) {
	var users []User

	err := DB.Model(User{}).Get(&users)

	if err != nil {
		t.Error(err)
	}
}

func TestModelInsert(t *testing.T) {

	var data []map[string]interface{}

	data = append(data, map[string]interface{}{
		"name":   "Alice",
		"gender": "F",
	})

	data = append(data, map[string]interface{}{
		"name":   "Martine",
		"gender": "M",
	})

	id, affected, err := DB.Model(User{}).Inserts(data)

	if err != nil {
		panic(err)
		t.Error(err)
	}

	t.Log(id, affected, err)
}

func TestModelUpdate(t *testing.T) {

	data := map[string]interface{}{
		"gender": "F",
		"addr":   "Columbia",
	}

	_, err := DB.Model(User{}).Where("name", "Martine").Update(data)

	if err != nil {
		t.Error(err)
	}
}

func TestModelCreate(t *testing.T) {
	u := &User{
		Name:      "testing",
		Gender:    "F",
		Addr:      "Columbia",
		BirthDate: "1990-01-01",
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := DB.Create(u)
	if err != nil {
		t.Error(err)
	}

	t.Log(u.Id)
}
