package mongo

import (
	"fmt"
	"log"
	"testing"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Student struct {
	// Id_   bson.ObjectId `bson:"_id"`
	Id    int       `bson:"_id"`
	Name  string    `bson:"name"`
	Phone string    `bson:"phone"`
	Email string    `bson:"email"`
	Sex   string    `bson:"sex"`
	Time  time.Time `bson:"time"`
}

func ConnecToDB() *mgo.Collection {
	session, err := mgo.Dial("192.168.182.11:27017")
	if err != nil {
		panic(err)
	}
	//defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("user")
	return c
}

func InsertToMogoTest() {
	c := ConnecToDB()
	stu1 := map[string]interface{}{
		"_id":   5,
		"name":  "liss",
		"phone": "13980989767",
		"email": "12832984@qq.com",
		"sex":   "M",
		"time":  time.Now(),
	}
	var stu2 interface{} = Student{
		Id:    6,
		Name:  "liss",
		Phone: "13980989767",
		Email: "12832984@qq.com",
		Sex:   "M",
		Time:  time.Now(),
	}
	// stu1 := Student{
	// 	Id:    3,
	// 	Name:  "zhangsan",
	// 	Phone: "13480989765",
	// 	Email: "329832984@qq.com",
	// 	Sex:   "F",
	// 	Time:  time.Now(),
	// }
	// stu2 := Student{
	// 	Id:    4,
	// 	Name:  "liss",
	// 	Phone: "13980989767",
	// 	Email: "12832984@qq.com",
	// 	Sex:   "M",
	// 	Time:  time.Now(),
	// }
	err := c.Insert(&stu1, &stu2)
	if err != nil {
		log.Fatal(err)
	}
}

// 普通查询
func GetDataViaSexTest() {
	c := ConnecToDB()
	result := Student{}
	err := c.Find(bson.M{"sex": "M"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("student", result)
	students := make([]Student, 20)
	err = c.Find(nil).All(&students)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(students)

}

func UpdateDBViaIdTest() {
	//id := bson.ObjectIdHex("5a66a96306d2a40a8b884049")
	c := ConnecToDB()
	err := c.Update(bson.M{"email": "12832984@qq.com"}, bson.M{"$set": bson.M{"name": "haha", "phone": "37848"}})
	if err != nil {
		log.Fatal(err)
	}
}

func RemoveFromMgoTest() {
	c := ConnecToDB()
	_, err := c.RemoveAll(bson.M{"phone": "13480989765"})
	if err != nil {
		log.Fatal(err)
	}
}

// 查询
func SelectTest() {
	c := ConnecToDB()
	// 按条件获取一条记录
	// result := Student{}
	// err := c.Find(bson.M{"sex": "M"}).One(&result)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("student", result)

	// 获取20条
	// students := make([]Student, 20)
	// err = c.Find(nil).All(&students)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(students)

	// 注意： 时间类型存数据库的时候是转了UTC时间，-8小时存。取的时候回自动转换回来 +8小时。查询的时候也不需要特别处理就回转成UTC时间取查数据库。
	//        但是使用数据库工具的时候需要将当前时间-8小时去查询，直接用UTC查询。new Date("2019-08-30T07:25:00Z")
	// 按时间查询 Find(bson.M{}).Select()
	// c.Find(bson.M{"time": bson.M{"$gt": startTime, "$lt": endTime}}).All(&items)
	// result := Student{}
	// 按指定时间查询
	// t := time.Date(2019, 8, 30, 15, 25, 00, 00, time.Local) // 正常当前时间
	// t := time.Date(2019, 8, 30, 07, 25, 00, 00, time.Local).UTC() // UTC时间, 时间不对，查询不正确。
	// err := c.Find(bson.M{"time": bson.M{"$lt": t}}).Sort("-time").One(&result)
	// 按当前时间查询
	// err := c.Find(bson.M{"time": bson.M{"$lt": time.Now().UTC()}}).Sort("-time").One(&result) // time.Now().UTC() 是不对的，应该直接使用当前时间time.Now()
	// 单条查询
	// err := c.Find(bson.M{"_id": bson.M{"$gt": 0}}).Sort("-_id").One(&result)
	// 多条查询
	// err := c.Find(bson.M{"_id": bson.M{"$gt": 0}}).Sort("-_id").All(&result)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("student", result)

	// 模糊查询 ,  正则实现。
	// 	regex操作符的介绍

	// MongoDB使用$regex操作符来设置匹配字符串的正则表达式，使用PCRE(Pert Compatible Regular Expression)作为正则表达式语言。

	// regex操作符
	// {<field>:{$regex:/pattern/，$options:’<options>’}}
	// {<field>:{$regex:’pattern’，$options:’<options>’}}
	// {<field>:{$regex:/pattern/<options>}}
	// 正则表达式对象
	// {<field>: /pattern/<options>}
	// $regex与正则表达式对象的区别:

	// 在$in操作符中只能使用正则表达式对象，例如:{name:{$in:[/^joe/i,/^jack/}}
	// 在使用隐式的$and操作符中，只能使用$regex，例如:{name:{$regex:/^jo/i, $nin:['john']}}
	// 当option选项中包含X或S选项时，只能使用$regex，例如:{name:{$regex:/m.*line/,$options:"si"}}
	// $regex操作符的使用

	// $regex操作符中的option选项可以改变正则匹配的默认行为，它包括i, m, x以及S四个选项，其含义如下

	// i 忽略大小写，{<field>{$regex/pattern/i}}，设置i选项后，模式中的字母会进行大小写不敏感匹配。
	// m 多行匹配模式，{<field>{$regex/pattern/,$options:'m'}，m选项会更改^和$元字符的默认行为，分别使用与行的开头和结尾匹配，而不是与输入字符串的开头和结尾匹配。
	// x 忽略非转义的空白字符，{<field>:{$regex:/pattern/,$options:'m'}，设置x选项后，正则表达式中的非转义的空白字符将被忽略，同时井号(#)被解释为注释的开头注，只能显式位于option选项中。
	// s 单行匹配模式{<field>:{$regex:/pattern/,$options:'s'}，设置s选项后，会改变模式中的点号(.)元字符的默认行为，它会匹配所有字符，包括换行符(\n)，只能显式位于option选项中。
	// 使用$regex操作符时，需要注意下面几个问题:

	// i，m，x，s可以组合使用，例如:{name:{$regex:/j*k/,$options:"si"}}
	// 在设置索弓}的字段上进行正则匹配可以提高查询速度，而且当正则表达式使用的是前缀表达式时，查询速度会进一步提高，例如:{name:{$regex: /^joe/}
	result := Student{}
	// 简单匹配
	// err := c.Find(bson.M{"name": bson.M{"$regex": "san"}}).Sort("-time").One(&result)
	// 不区分大小写
	// err := c.Find(bson.M{"name": bson.M{"$regex": "SAN", "$options": "i"}}).Sort("-time").One(&result)
	// 数组查找  find({tags:{$regex:"run"}})
	// 使用mgo的正则
	err := c.Find(M{"name": M{"$regex": bson.RegEx{Pattern: "san", Options: "i"}}}).Sort("-time").One(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("student", result)

	// 分页查询
	// con.Find(nil).Sort("-_id").Limit(5).Skip(0).All(&user).Count()
}

func TestMain(t *testing.T) {
	// InsertToMogoTest()
	// GetDataViaSexTest()
	// SelectTest()

	// 查询的封装方法测试
	// result := []Student{}
	// query := M{"name": M{"$regex": bson.RegEx{Pattern: "san", Options: "i"}}}
	// Select("192.168.182.11:27017", "test", "user", &query, &result, 0, 0, "-time")
	// fmt.Println("student", result)

	// 插入数据的封装方法测试
	// var stu interface{} = Student{
	// 	Id:    8,
	// 	Name:  "liss",
	// 	Phone: "13980989767",
	// 	Email: "12832984@qq.com",
	// 	Sex:   "M",
	// 	Time:  time.Now(),
	// }
	// Insert("192.168.182.11:27017", "test", "user", &stu)
	// 测试下日志

}
