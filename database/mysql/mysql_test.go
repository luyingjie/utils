package mysql

import (
	"fmt"
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
	db := DBConn("cmpadmin:CMP_Zhu88jie@tcp(139.198.190.114:3306)/testing_v1.8.5_20191211?charset=utf8")
	rows, _ := db.Query("select id from e_platform_node where cloud_resource_id = '/service/sites/43FC07EB/hosts/165' and is_deleted = 0")
	ttt := ParseRows(rows)
	// fmt.Println(string(ttt[0]["id"].([]uint8)))
	// fmt.Println(ttt[0])
	fmt.Println(string(ttt[0]["id"].([]uint8)))
}
