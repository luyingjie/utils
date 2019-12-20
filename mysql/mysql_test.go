package mysql

import (
	"testing"

	// "fmt"
	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	Id   int
	Age  int
	Name string
}

type user1 struct {
	Id   int `pk orm:"column(id)"`
	age  int
	name string
}

func TestMain(t *testing.T) {

}
