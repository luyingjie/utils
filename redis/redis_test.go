package redis

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestMain(t *testing.T) {
	// t.Fatal(err)
	c, errc := redis.Dial("tcp", "192.168.182.11:6379")
	if errc != nil {
		// t.Fatal(errc)
		// return
	}
	defer c.Close()

	// 若value的类型为int，则用redis.Int转换

	// 若value的类型为string，则用redis.String转换

	// 若value的类型为json，则用redis.Byte转换

	// n,_ := c.Do("lpush","name","Luke","Luke1")
	// fmt.Println(n)
	// result,_ := redis.Values(c.Do("lpop","abc"))
	// values, _ := redis.Values(c.Do("lrange", "name",0,1))
	// values, _ := c.Do("llen", "name")
	// values, _ := c.Do("lpop", "name")
	// sss, s := c.Do("lpop","loglist_test")
	// fmt.Println(sss)
	// fmt.Println(s)
	// values,_ := redis.String(sss,s)
	// values==""

	// c.Do("SET", "go_key", "redigo")
	// values, _ := redis.String(c.Do("GET", "go_key"))
	// fmt.Println("取到的结果")
	// fmt.Println(values)
	// for _, v := range values {
	// 	fmt.Println(v.([]byte))
	// }

	//  订阅者模式
	psc := redis.PubSubConn{Conn: c}
	psc.Subscribe()
}

// 2.6 发布和订阅(Pub/Sub)
// 使用Send，Flush和Receive方法来是实现Pub/Sb订阅者。

// c.Send("SUBSCRIBE", "example")
// c.Flush()
// for {
// reply, err := c.Receive()
// if err != nil {
// return err
// }
// // process pushed message
// }
// PubSubConn类型封装了Conn提供了便捷的方法来实现订阅者模式。Subscribe,PSubscribe,Unsubscribe和PUnsubscribe方法发送和清空订阅管理命令。Receive将一个推送消息
// 转化为一个在type switch更为方便使用的类型。

// psc := redis.PubSubConn{c}
// psc.Subscribe("example")
// for {
// switch v := psc.Receive().(type) {
// case redis.Message:
// fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
// case redis.Subscription:
// fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
// case error:
// return v
// }
// }
