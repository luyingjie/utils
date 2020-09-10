package redis

import (
	"utils/error"

	"github.com/garyburd/redigo/redis"
)

//连接 redis
func connect(url string) redis.Conn {
	c, err := redis.Dial("tcp", url)
	if err != nil {
		error.Try(2000, 3, "utils/database/redis/redis/connect/Dial", err)
	}
	return c
	// defer c.Close()
}

func PushList(url, name, value string) {
	c := connect(url)
	defer c.Close()
	c.Do("rpush", name, value)
	// 报错会无意义的开销，并且会产生大量的垃圾日志，如果抛出还会中断监听。
	// _, err := c.Do("rpush",name, value)
	// if err != nil {
	// 	// 没有预估使用 go的作用
	// 	error.TryWarning(err)
	// }
}

func GetList(url, name string) string {
	c := connect(url)
	defer c.Close()
	value, _ := redis.String(c.Do("lpop", name))
	// if err != nil {
	// 	error.TryWarning(err)
	// }
	return value
}

func Set(url, key string, val interface{}) interface{} {
	c := connect(url)
	defer c.Close()
	value, _ := c.Do("SET", key, val)
	// if err != nil {
	// 	error.TryWarning(err)
	// }
	return value
}

func GetToStr(url, key string) string {
	c := connect(url)
	defer c.Close()
	value, _ := redis.String(c.Do("GET", key))
	// if err != nil {
	// 	error.TryWarning(err)
	// }
	return value
}

func GetToInt(url, key string) int {
	c := connect(url)
	defer c.Close()
	value, _ := redis.Int(c.Do("GET", key))
	// if err != nil {
	// 	error.TryWarning(err)
	// }
	return value
}
