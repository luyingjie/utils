package mongo

import (
	myError "utils/error"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// 和其他共用
// type M bson.M

// Connect : 连接数据库。连接池的session释放还需要研究。
func Connect(Url string) (*mgo.Session, error) {
	session, err := mgo.Dial(Url)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// Insert : 添加数据
func Insert(Url, DB, C string, Data interface{}) error {
	if Url == "" || DB == "" || C == "" || Data == nil {
		return myError.New("确少必要的参数")
	}
	session, err := Connect(Url)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DB).C(C)
	defer session.Close()

	err = c.Insert(Data)
	if err != nil {
		return err
	}
	return nil
}

// Remove : 删除数据
func Remove(Url, DB, C string, Query *M) error {
	if Url == "" || DB == "" || C == "" || Query == nil {
		return myError.New("确少必要的参数")
	}
	session, err := Connect(Url)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DB).C(C)
	defer session.Close()

	_, err = c.RemoveAll(Query)
	if err != nil {
		return err
	}
	return nil
}

// Update : 修改数据
func Update(Url, DB, C string, Query *M, Data interface{}) error {
	if Url == "" || DB == "" || C == "" || Query == nil || Data == nil {
		return myError.New("确少必要的参数")
	}
	session, err := Connect(Url)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DB).C(C)
	defer session.Close()

	err = c.Update(Query, bson.M{"$set": Data})
	if err != nil {
		return err
	}
	return nil
}

// Select : 查询数据
func Select(Url, DB, C string, Query *M, ResultModel interface{}, Limit, Skip int, Sort string) error {
	if Url == "" || DB == "" || C == "" || Query == nil || ResultModel == nil {
		return myError.New("确少必要的参数")
	}
	session, err := Connect(Url)
	if err != nil {
		return err
	}
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

	err = query.All(ResultModel)
	if err != nil {
		return err
	}
	return nil
}
