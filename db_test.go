package torm

import (
	"log"
	"testing"

	_ "github.com/thinkoner/torm/driver/mysql"
	_ "github.com/thinkoner/torm/driver/sqlite"
)

var DB *Connection

func init() {
	config := Config{
		Driver: "mysql",
		Dsn:    "root:abc-123@tcp(10.0.4.159:33066)/test?charset=utf8",
		// Driver: "sqlite3",
		// Dsn:    "tests/testDB.db",
	}
	conn, err := Open(config)
	if err != nil {
		log.Println("connect err：", err)
	}
	DB = conn
}

//
// func TestSelect(test *testing.T) {
//
// 	result, err := DB.Select("SELECT * FROM USER WHERE id = ?", []interface{}{"1"})
//
// 	if err != nil {
// 		test.Error(err)
// 	}
//
// 	log.Println(result)
// }
//

type User struct {
	Id   int64  `torm:"column:id"`
	Name string `torm:"column:name"`
	Sex  string `torm:"column:sex"`
	Addr string `torm:""`
}

func (u *User) TableName() string {
	return "users"
}

func TestTableSelect(test *testing.T) {
	var users []User
	err := DB.Table("user").Get(&users)

	if err != nil {
		test.Error(err)
	}
	for _, user := range users {
		log.Println(user)
	}
}
func TestTableSelectOne(test *testing.T) {
	var user User
	err := DB.Table("user").First(&user)

	if err != nil {
		test.Error(err)
	}
	log.Println(user)
	// for _, user := range users{
	// 	log.Println(user)
	// }
}

func TestModelSelect(test *testing.T) {
	var users []User

	err := DB.Model(User{}).Get(&users)

	if err != nil {
		test.Error(err)
	}

	for _, user := range users {
		log.Println(user)
	}
}

func TestModelInsert(test *testing.T) {

	var data []map[string]interface{}

	data = append(data, map[string]interface{}{
		"name": "Alice",
		"sex":  "fmale",
	})

	data = append(data, map[string]interface{}{
		"name": "Martinez",
		"sex":  "male",
	})
	data = append(data, map[string]interface{}{
		"name": "Martin",
		"sex":  "male",
	})

	id, affected, err := DB.Model(User{}).Inserts(data)

	if err != nil {
		panic(err)
		test.Error(err)
	}

	test.Log(id, affected, err)
}

func TestModelUpdate(test *testing.T) {

	data := map[string]interface{}{
		"sex": "male",
		"addr": "Columbia",
	}

	affected, err := DB.Model(User{}).Where("name", "Alice").Update(data)

	if err != nil {
		panic(err)
		test.Error(err)
	}

	test.Log(affected, err)
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
