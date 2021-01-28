package mongo

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type M bson.M
type RegEx bson.RegEx
type D bson.D

type DB struct {
	session *mgo.Session
	db      *mgo.Database
	isClose bool
}

func NewId() bson.ObjectId {
	// Id  bson.ObjectId `bson:"_id"`
	// 转字符 Hex()
	// 转 objectid
	// 	id := "5204af979955496907000001"
	// objectId := bson.ObjectIdHex(id)
	return bson.NewObjectId()
}

func ToId(id string) bson.ObjectId {
	return bson.ObjectIdHex(id)
}

func NewStrId() string {
	return bson.NewObjectId().Hex()
}

// NewDB 创建新连接。 每次都要创建连接使用完后关闭。pool应该交给包，不应该在外面单独做。
func NewDB(_url, dbName string) (*DB, error) {
	_session, err := mgo.Dial(_url)
	if err != nil {
		return nil, err
	}

	// session 的读操作开始是向其他服务器发起（且通过一个唯一的连接），只要出现了一次写操作，session 的连接就会切换至主服务器。由此可见此模式下，能够分散一些读操作到其他服务器，但是读操作不一定能够获得最新的数据。
	// _session.SetMode(mgo.Monotonic, true)

	c := &DB{
		session: _session,
		db:      _session.DB(dbName),
		isClose: true,
	}

	return c, nil
}

func (db *DB) Close() {
	db.session.Close()
}

func (db *DB) GetSession() *mgo.Session {
	return db.session
}

func (db *DB) GetDB() *mgo.Database {
	return db.db
}

func (db *DB) SetClose(i bool) {
	db.isClose = i
}

func (db *DB) SetDB(dbName string) {
	db.db = db.session.DB(dbName)
}

// Insert : 添加数据
func (db *DB) Insert(C string, Data interface{}) error {
	if C == "" || Data == nil {
		return errors.New("确少必要的参数")
	}
	c := db.db.C(C)
	defer func() {
		if db.isClose {
			db.session.Close()
		}
	}()

	err := c.Insert(Data)
	return err
}

// Remove : 删除数据
func (db *DB) Remove(C string, Query *M) error {
	if C == "" || Query == nil {
		return errors.New("确少必要的参数")
	}

	c := db.db.C(C)
	defer func() {
		if db.isClose {
			db.session.Close()
		}
	}()

	err := c.Remove(Query)
	return err
}

func (db *DB) RemoveAll(C string, Query *M) error {
	if C == "" || Query == nil {
		return errors.New("确少必要的参数")
	}

	c := db.db.C(C)
	defer func() {
		if db.isClose {
			db.session.Close()
		}
	}()

	_, err := c.RemoveAll(Query)
	return err
}

func (db *DB) RemoveId(C string, id interface{}) error {
	if C == "" || id == nil {
		return errors.New("确少必要的参数")
	}

	c := db.db.C(C)
	defer func() {
		if db.isClose {
			db.session.Close()
		}
	}()

	err := c.RemoveId(id)
	return err
}

// Update : 修改数据
func (db *DB) Update(C string, Query *M, Data interface{}) error {
	if C == "" || Query == nil || Data == nil {
		return errors.New("确少必要的参数")
	}

	c := db.db.C(C)
	defer func() {
		if db.isClose {
			db.session.Close()
		}
	}()

	err := c.Update(Query, Data) // bson.M{"$set": Data}
	return err
}

func (db *DB) UpdateAll(C string, Query *M, Data interface{}) error {
	if C == "" || Query == nil || Data == nil {
		return errors.New("确少必要的参数")
	}

	c := db.db.C(C)
	defer func() {
		if db.isClose {
			db.session.Close()
		}
	}()

	_, err := c.UpdateAll(Query, Data)
	return err
}

func (db *DB) UpdateId(C string, id interface{}, Data interface{}) error {
	if C == "" || id == nil || Data == nil {
		return errors.New("确少必要的参数")
	}

	c := db.db.C(C)
	defer func() {
		if db.isClose {
			db.session.Close()
		}
	}()

	err := c.UpdateId(id, Data)
	return err
}

// Select : 查询数据
func (db *DB) Select(C string, Query *M, ResultModel interface{}, Limit, Skip int, Sort string) error {
	if C == "" || Query == nil || ResultModel == nil {
		return errors.New("确少必要的参数")
	}

	c := db.db.C(C)
	defer func() {
		if db.isClose {
			db.session.Close()
		}
	}()

	query := c.Find(Query)
	if Limit != 0 && Skip != 0 {
		query = query.Limit(Limit).Skip(Skip)
	}
	if Sort != "" {
		query = query.Sort(Sort)
	}

	err := query.All(ResultModel)
	return err
}

// 这个方法调用有问题，找不到对应的id string。
func (db *DB) FindId(C string, id interface{}, ResultModel interface{}) error {
	if C == "" || id == nil || ResultModel == nil {
		return errors.New("确少必要的参数")
	}

	c := db.db.C(C)
	defer func() {
		if db.isClose {
			db.session.Close()
		}
	}()

	query := c.FindId(id)

	err := query.One(ResultModel)
	return err
}

func (db *DB) Find(C string, Query *M, Sort string, ResultModel interface{}) error {
	if C == "" || Query == nil || ResultModel == nil {
		return errors.New("确少必要的参数")
	}

	c := db.db.C(C)
	defer func() {
		if db.isClose {
			db.session.Close()
		}
	}()

	query := c.Find(Query)
	if Sort != "" {
		query = query.Sort(Sort)
	}

	err := query.One(ResultModel)
	return err
}

// Count : 查询总数
func (db *DB) Count(C string, Query *M) (int, error) {
	if C == "" || Query == nil {
		return 0, errors.New("确少必要的参数")
	}

	c := db.db.C(C)
	defer func() {
		if db.isClose {
			db.session.Close()
		}
	}()

	return c.Count()
}
