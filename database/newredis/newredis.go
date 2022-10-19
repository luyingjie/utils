package newredis


import (
"errors"
"utils/util"
"fmt"
goredis "github.com/go-redis/redis/v7"
"sync"
"time"
)

var redisMap map[string]*Redis
var mut sync.Mutex

func init() {
	redisMap = make(map[string]*Redis)
}

type Redis struct {
	Addr     string
	Password string
	Index    int
	PoolSize int
	redisIns *goredis.Client
}

func (ins *Redis) init() error {
	ins.redisIns = goredis.NewClient(&goredis.Options{
		Addr:         ins.Addr,
		Password:     ins.Password,
		DB:           ins.Index,
		MaxRetries:   5,
		PoolSize:     ins.PoolSize,
		MinIdleConns: 5,
	})
	if _, err := ins.redisIns.Ping().Result(); err != nil {
		return errors.New("redis connect error")
	}
	return nil
}

// string
func (ins *Redis) Get(key string) (string, error) {
	return ins.redisIns.Get(key).Result()
}

func (ins *Redis) Set(key string, value interface{}, ttl time.Duration) error {
	return ins.redisIns.Set(key, value, ttl).Err()
}

func (ins *Redis) Decr(key string, decrement int64) (int64, error) {
	return ins.redisIns.DecrBy(key, decrement).Result()
}

func (ins *Redis) Incr(key string, increment int64) (int64, error) {
	return ins.redisIns.IncrBy(key, increment).Result()
}

func (ins *Redis) Setnx(key string, value int64, ttl time.Duration) (bool, error) {
	return ins.redisIns.SetNX(key, value, ttl).Result()
}

func (ins *Redis) IfExist(key string) (bool, error) {
	val, err := ins.redisIns.Exists(key).Result()
	if err != nil {
		return false, err
	} else {
		return val == 1, nil
	}
}
func (ins *Redis) Expire(key string, ttl time.Duration) (bool, error) {
	return ins.redisIns.Expire(key, ttl).Result()
}

func (ins *Redis) TTL(key string) (time.Duration, error) {
	return ins.redisIns.TTL(key).Result()
}

func (ins *Redis) MGet(key ...string) ([]interface{}, error) {
	return ins.redisIns.MGet(key...).Result()
}

func (ins *Redis) MSet(value ...interface{}) (string, error) {
	return ins.redisIns.MSet(value...).Result()
}

func (ins *Redis) StrLen(key string) (int64, error) {
	return ins.redisIns.StrLen(key).Result()
}

func (ins *Redis) Del(key ...string) error {
	_, err := ins.redisIns.Del(key...).Result()
	return err
}

// TODO hash and others

func (ins *Redis) HGet(key, field string) (string, error) {
	return ins.redisIns.HGet(key, field).Result()
}

func (ins *Redis) HSet(key string, values ...interface{}) (int64, error) {
	return ins.redisIns.HSet(key, values...).Result()
}

func (ins *Redis) HDel(key, field string) (int64, error) {
	return ins.redisIns.HDel(key, field).Result()
}

func (ins *Redis) HExists(key, field string) (bool, error) {
	return ins.redisIns.HExists(key, field).Result()
}

func (ins *Redis) HLen(key string) (int64, error) {
	return ins.redisIns.HLen(key).Result()
}

func (ins *Redis) HMSet(key string, values ...interface{}) (bool, error) {
	return ins.redisIns.HMSet(key, values...).Result()
}

func (ins *Redis) HMGet(key string, fields ...string) ([]interface{}, error) {
	return ins.redisIns.HMGet(key, fields...).Result()
}

func (ins *Redis) HKeys(key string) ([]string, error) {
	return ins.redisIns.HKeys(key).Result()
}

func (ins *Redis) HVals(key string) ([]string, error) {
	return ins.redisIns.HVals(key).Result()
}

func (ins *Redis) HGetAll(key string) (map[string]string, error) {
	return ins.redisIns.HGetAll(key).Result()
}

func (ins *Redis) HSetNX(key, field string, values interface{}) (bool, error) {
	return ins.redisIns.HSetNX(key, field, values).Result()
}

func (ins *Redis) HincrBy(key, field string, incr int64) (int64, error) {
	return ins.redisIns.HIncrBy(key, field, incr).Result()
}

func (ins *Redis) HincrByFloat(key, field string, incr float64) (float64, error) {
	return ins.redisIns.HIncrByFloat(key, field, incr).Result()
}

func (ins *Redis) HScan(key string, cur uint64, match string, count int64) (keys []string, cursor uint64, err error) {
	return ins.redisIns.HScan(key, cur, match, count).Result()
}

// list
func (ins *Redis) LPush(key string, value ...interface{}) (int64, error) {
	return ins.redisIns.LPush(key, value...).Result()
}

func GetRedisClient(addr, password string, index, poolsize int) *Redis {
	md5 := util.GetMd5(fmt.Sprintf("%s:%s:%d", addr, password, index))
	redisIns, ok := redisMap[md5]
	if !ok {
		mut.Lock()
		defer mut.Unlock()
		// check again
		redisIns, ok = redisMap[md5]
		if !ok {
			redisIns = new(Redis)
			redisIns.PoolSize = poolsize
			redisIns.Addr = addr
			redisIns.Password = password
			redisIns.Index = index
			redisIns.PoolSize = poolsize
			if err := redisIns.init(); err != nil {
				fmt.Printf("error exist:%v", err)
				panic(err)
			}
			redisMap[md5] = redisIns
		}

	}
	return redisIns
}
