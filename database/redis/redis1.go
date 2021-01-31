package redis

import (
	"github.com/garyburd/redigo/redis"
)

//连接 redis
func connect(url string) (redis.Conn, error) {
	c, err := redis.Dial("tcp", url)
	if err != nil {
		return nil, err
	}
	return c, nil
	// defer c.Close()
}

func PushList(url, name, value string) error {
	c, err := connect(url)
	if err != nil {
		return err
	}
	defer c.Close()
	c.Do("rpush", name, value)
	// 报错会无意义的开销，并且会产生大量的垃圾日志，如果抛出还会中断监听。
	// _, err := c.Do("rpush",name, value)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func GetList(url, name string) (string, error) {
	c, err := connect(url)
	if err != nil {
		return "", err
	}
	defer c.Close()
	value, err1 := redis.String(c.Do("lpop", name))
	if err1 != nil {
		return "", err1
	}
	return value, nil
}

func Set(url, key string, val interface{}) (interface{}, error) {
	c, err := connect(url)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	value, err1 := c.Do("SET", key, val)
	if err1 != nil {
		return nil, err1
	}
	return value, nil
}

func GetToStr(url, key string) (string, error) {
	c, err := connect(url)
	if err != nil {
		return "", err
	}
	defer c.Close()
	value, err1 := redis.String(c.Do("GET", key))
	if err1 != nil {
		return "", err
	}
	return value, nil
}

func GetToInt(url, key string) (int, error) {
	c, err := connect(url)
	if err != nil {
		return 0, err
	}
	defer c.Close()
	value, err1 := redis.Int(c.Do("GET", key))
	if err1 != nil {
		return 0, err1
	}
	return value, nil
}
