<h1 align="center">
  ThinkORM
</h1>

<p align="center">
	<strong>ThinkORM is a simple and Powerful ORM for Go.</strong>
</p>

<p align="center">
    <a href="https://travis-ci.org/thinkoner/torm">
		<img src="https://travis-ci.org/thinkoner/torm.svg?branch=master" alt="Build Status">
  	</a>
  	<a href="https://coveralls.io/github/thinkoner/torm?branch=master">
    	<img src="https://coveralls.io/repos/github/thinkoner/torm/badge.svg?branch=master" alt="Go Report Card">
    </a>
	<a href="https://goreportcard.com/report/github.com/thinkoner/torm">
		<img src="https://goreportcard.com/badge/github.com/thinkoner/torm" alt="Go Report Card">
  	</a>
	<a href="https://godoc.org/github.com/thinkoner/torm">
		<img src="https://godoc.org/github.com/thinkoner/torm?status.svg" alt="GoDoc">
  	</a>
	<a href="https://gitter.im/think-go/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge">
		<img src="https://badges.gitter.im/think-go/community.svg" alt="Join the chat">
  	</a>
	<a href="https://github.com/thinkoner/torm/releases">
		<img src="https://img.shields.io/github/release/thinkoner/torm.svg" alt="Latest Stable Version">
	</a>
	<a href="LICENSE">
		<img src="https://img.shields.io/github/license/thinkoner/torm.svg" alt="License">
	</a>
</p>


## Installation

```
go get github.com/thinkoner/torm
```

## Quick start

##### Get the database connection:

```go
config := torm.Config{
		Driver: "mysql",
		Dsn:    "root:abc-123@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true",
	}
db, _ := torm.Open(config)
```

##### Query Builder:

```go
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

var users []User
conn.Table("users").
	Where("gender", "M").
	Where("name", "like", "%a%").
	WhereColumn("created_at", "!=", "updated_at").
	WhereIn("addr", []interface{}{"Columbia", "Alaska"}).
	WhereBetween("birth_date", []interface{}{"1990-01-01", "1999-12-31"}).
	Get(&users)
```

## License

This project is licensed under the [Apache 2.0 license](LICENSE).

## Contact

If you have any issues or feature requests, please contact us. PR is welcomed.
- https://github.com/thinkoner/torm/issues
- duanpier@gmail.com
