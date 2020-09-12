package mongo

import (
	myError "utils/error"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type M bson.M

// Connect : 连接数据库。连接池的session释放还需要研究。
func Connect(Url string) *mgo.Session {
	session, err := mgo.Dial(Url)
	if err != nil {
		myError.Try(2000, 3, err)
	}
	return session
}

// Insert : 添加数据
func Insert(Url, DB, C string, Data interface{}) {
	if Url == "" || DB == "" || C == "" || Data == nil {
		myError.Trys(1000, 2, "确少必要的参数")
		return
	}
	session := Connect(Url)
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DB).C(C)
	defer session.Close()

	err := c.Insert(Data)
	if err != nil {
		myError.Try(2000, 3, err)
	}
}

// Remove : 删除数据
func Remove(Url, DB, C string, Query *M) {
	if Url == "" || DB == "" || C == "" || Query == nil {
		myError.Trys(1000, 2, "确少必要的参数")
		return
	}
	session := Connect(Url)
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DB).C(C)
	defer session.Close()

	_, err := c.RemoveAll(Query)
	if err != nil {
		myError.Try(2000, 3, err)
	}
}

// Update : 修改数据
func Update(Url, DB, C string, Query *M, Data interface{}) {
	if Url == "" || DB == "" || C == "" || Query == nil || Data == nil {
		myError.Trys(1000, 2, "确少必要的参数")
		return
	}
	session := Connect(Url)
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DB).C(C)
	defer session.Close()

	err := c.Update(Query, bson.M{"$set": Data})
	if err != nil {
		myError.Try(2000, 3, err)
	}
}

// Select : 查询数据
func Select(Url, DB, C string, Query *M, ResultModel interface{}, Limit, Skip int, Sort string) {
	if Url == "" || DB == "" || C == "" || Query == nil || ResultModel == nil {
		myError.Trys(1000, 2, "确少必要的参数")
		return
	}
	session := Connect(Url)
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DB).C(C)
	defer session.Close()

	query := c.Find(Query)
	if Limit != 0 && Skip != 0 {
		query = query.Limit(Limit).Skip(Skip)
	}
	if Sort != "" {
		query = query.Sort(Sort)
	}

	err := query.All(ResultModel)
	if err != nil {
		myError.Try(2000, 3, err)
	}
}
