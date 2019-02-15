package torm

import (
	"log"
	"strings"
	"testing"
	"time"

	_ "github.com/thinkoner/torm/driver/mysql"
	_ "github.com/thinkoner/torm/driver/sqlite"
)

var DB *Connection

func init() {
	DB = initDb()
}

func initDb() *Connection {
	dsn := "root:@tcp(127.0.0.1:3306)/"
	driver := "mysql"
	db, err := Open(Config{
		Driver: driver,
		Dsn:    dsn + "?charset=utf8&parseTime=true",
	})
	if err != nil {
		panic(err)
	}
	_, _, err = db.AffectingStatement(`CREATE DATABASE IF NOT EXISTS torm_test`)

	if err != nil {
		panic(err)
	}

	db, err = Open(Config{
		Driver: driver,
		Dsn:    dsn + "torm_test?charset=utf8&parseTime=true",
	})
	if err != nil {
		panic(err)
	}

	db.AffectingStatement(`DROP TABLE IF EXISTS users;`)
	db.AffectingStatement(`
CREATE TABLE users (
  name varchar(255) DEFAULT NULL,
  gender varchar(255) DEFAULT NULL,
  addr varchar(255) DEFAULT NULL,
  birth_date date DEFAULT NULL,
  balance decimal(15,4) DEFAULT '0.0000',
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;
`)
	db.AffectingStatement(`
INSERT INTO users VALUES ('Jenny', 'M', 'California', '2010-02-11', '8.5000', '2018-11-26 14:42:59', '2019-02-11 15:00:41'),
 ('Andrew', 'M', 'Boston', '2010-02-11', '6.3400', '2018-11-26 14:42:59', '2019-02-11 15:00:43'),
 ('Alex', 'F', 'Alaska', '2010-03-01', '1.2465', '2018-11-26 14:42:59', '2019-02-11 15:01:01'),
 ('Adrian', 'M', 'Chicago', '2012-05-18', '0.0000', '2018-11-26 14:42:59', '2019-02-11 15:00:44'),
 ('Simon', 'F', 'Columbia', '2008-06-30', '82.2360', '2018-11-26 14:42:59', '2019-02-11 15:01:00'),
 ('Neil', 'F', 'California', '2001-10-22', '10.0000', '2018-11-26 14:42:59', '2019-02-11 15:00:58'),
 ('Richard', 'F', 'Frankfort', '2010-12-11', '1256.0200', '2018-11-26 14:42:59', '2019-02-11 15:00:59'),
 ('Ann', 'F', 'Columbia', '2010-11-11', '236.0000', '2018-11-26 14:42:59', '2019-02-11 15:00:56'),
 ('Christine', 'F', 'Alaska', '2010-02-11', '36.2100', '2018-11-26 14:42:59', '2019-02-11 15:00:56'),
 ('Mike', 'M', 'Columbia', '2010-02-11', '365.0221', '2018-11-26 14:42:59', '2019-02-11 15:00:46'),
 ('Dave', 'M', 'Boston', '2010-02-11', '36.5546', '2018-11-26 14:42:59', '2019-02-11 15:00:47'),
 ('Richard', 'M', 'Alaska', '1999-02-11', '5.4500', '2018-11-26 14:42:59', '2019-02-11 15:00:47'),
 ('Laura', 'F', 'Columbia', '1995-08-22', '96.4560', '2018-11-26 14:42:59', '2019-02-11 15:00:54'),
 ('Bill', 'M', 'Columbia', '2010-07-16', '6.5465', '2018-11-26 14:42:59', '2019-02-11 15:00:49'),
 ('David', 'M', 'Boston', '2003-02-14', '7561.5540', '2018-11-26 14:42:59', '2019-02-11 15:00:50'),
 ('Alice', 'M', 'Columbia', '2014-05-26', '0.0000', '2019-02-11 15:01:31', '2019-02-11 15:00:50'),
 ('Patricio', 'M', 'Columbia', '2005-09-19', '0.0000', '2019-02-11 15:01:31', '2019-02-11 15:00:51'),
 ('Martinez', 'M', '', '2000-02-11', '0.0000', '2019-02-11 15:01:31', '2019-02-11 15:00:51'),
 ('Martin', 'M', '', '1996-10-01', '0.0000', '2019-02-11 15:01:31', '2019-02-11 15:01:31');
`)

	return db
}

type User struct {
	Name      string `torm:"column:name"`
	Gender    string `torm:"column:gender"`
	Addr      string
	BirthDate string
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) TableName() string {
	return "users"
}

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
		Addr  string
		Count int64 `torm:"column:ct"`
	}
	var userAddr []UserAddr
	err := DB.Table("users").
		Select("addr", "count(name) as ct").
		GroupBy("gender").
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

//
// func TestTableSelectRaw(test *testing.T) {
//
// 	result, err := DB.Table("user").SelectRaw("account * ? as price_with_min, account * ? as price_with_max", 10, 100).Get()
//
// 	if err != nil {
// 		test.Error(err)
// 	}
//
// 	log.Println(result)
// }
//
// func TestTableGroupBy(test *testing.T) {
//
// 	result, err := DB.Table("user").SelectRaw("sex, count(id) as count").GroupBy("sex").Having("count", ">", "6").Get()
//
// 	if err != nil {
// 		test.Error(err)
// 	}
//
// 	log.Println(result)
// }
// func TestTableOrderBy(test *testing.T) {
//
// 	result, err := DB.Table("user").OrderByDesc("name").Get()
//
// 	if err != nil {
// 		test.Error(err)
// 	}
//
// 	log.Println(result)
// }
//
// func TestInsert(test *testing.T) {
//
// 	result, err := DB.Insert("INSERT INTO user VALUES (0, 'Bob', 'male', 'California', '16.62', '2018-11-26 14:42:59', '2018-11-26 14:43:01')")
//
// 	if err != nil {
// 		test.Error(err)
// 	}
//
// 	log.Println(result)
// }
//
// func TestUpdate(test *testing.T) {
//
// 	result1, err := DB.Update("UPDATE user SET addr = 'Los Angeles' WHERE name = 'Bob'")
//
// 	if err != nil {
// 		test.Error(err)
// 	}
//
// 	log.Println(result1)
// }
//
// func TestDelete(test *testing.T) {
// 	result1, err := DB.Delete("DELETE FROM user WHERE name = 'Bob'")
//
// 	if err != nil {
// 		test.Error(err)
// 	}
//
// 	log.Println(result1)
// }
//
// func TestTableGet(test *testing.T) {
// 	result, err := DB.Table("user").Where("sex", "male").Get()
//
// 	if err != nil {
// 		test.Error(err)
// 	}
//
// 	log.Println(result)
// }
//
// func TestTableFirst(test *testing.T) {
// 	result, err := DB.Table("user").Where("sex", "male").First()
//
// 	if err != nil {
// 		test.Error(err)
// 	}
//
// 	log.Println(result)
// }
//
// func TestTableValue(test *testing.T) {
//
// 	result, err := DB.Table("user").Where("sex", "male").Value("name")
//
// 	if err != nil {
// 		test.Error(err)
// 	}
//
// 	log.Println(result)
// }
//
// func TestTableAggregate(test *testing.T) {
//
// 	_, err := DB.Table("user").Where("sex", "male").Count()
//
// 	result, err := DB.Table("user").Where("sex", "male").Avg("account")
//
// 	if err != nil {
// 		test.Error(err)
// 	}
//
// 	log.Println(result)
// }
//
// func TestTableLimit(test *testing.T) {
// 	result, err := DB.Table("user").Skip(2).Take(2).Get()
// 	if err != nil {
// 		test.Error(err)
// 	}
// 	log.Println(result)
// }
//
// func TestTableJoin(test *testing.T) {
// 	result, err := DB.Table("user").Join("posts", "user.id", "=", "posts.user_id").Select("user.*", "posts.title").Get()
// 	if err != nil {
// 		test.Error(err)
// 	}
// 	log.Println(result)
// }
