package session

import (
	"fmt"
	"time"
	vmap "utils/container/map"
	"utils/database/redis"
	verror "utils/os/error"
	"utils/util/json"

	"utils/os/timer"
)

// StorageRedis implements the Session Storage interface with redis.
type StorageRedis struct {
	redis         *redis.Redis    // Redis client for session storage.
	prefix        string          // Redis key prefix for session id.
	updatingIdMap *vmap.StrIntMap // Updating TTL set for session id.
}

var (
	// DefaultStorageRedisLoopInterval is the interval updating TTL for session ids
	// in last duration.
	DefaultStorageRedisLoopInterval = time.Minute
)

// NewStorageRedis creates and returns a redis storage object for session.
func NewStorageRedis(redis *redis.Redis, prefix ...string) *StorageRedis {
	if redis == nil {
		panic("redis instance for storage cannot be empty")
		return nil
	}
	s := &StorageRedis{
		redis:         redis,
		updatingIdMap: vmap.NewStrIntMap(true),
	}
	if len(prefix) > 0 && prefix[0] != "" {
		s.prefix = prefix[0]
	}
	// Batch updates the TTL for session ids timely.
	timer.AddSingleton(DefaultStorageRedisLoopInterval, func() {
		fmt.Print("StorageRedis.timer start")
		var id string
		var err error
		var ttlSeconds int
		for {
			if id, ttlSeconds = s.updatingIdMap.Pop(); id == "" {
				break
			} else {
				if err = s.doUpdateTTL(id, ttlSeconds); err != nil {
					verror.Try(5000, 3, err)
				}
			}
		}
		fmt.Print("StorageRedis.timer end")
	})
	return s
}

// New creates a session id.
// This function can be used for custom session creation.
func (s *StorageRedis) New(ttl time.Duration) (id string) {
	return ""
}

// Get retrieves session value with given key.
// It returns nil if the key does not exist in the session.
func (s *StorageRedis) Get(id string, key string) interface{} {
	return nil
}

// GetMap retrieves all key-value pairs as map from storage.
func (s *StorageRedis) GetMap(id string) map[string]interface{} {
	return nil
}

// GetSize retrieves the size of key-value pairs from storage.
func (s *StorageRedis) GetSize(id string) int {
	return -1
}

// Set sets key-value session pair to the storage.
// The parameter <ttl> specifies the TTL for the session id (not for the key-value pair).
func (s *StorageRedis) Set(id string, key string, value interface{}, ttl time.Duration) error {
	return ErrorDisabled
}

// SetMap batch sets key-value session pairs with map to the storage.
// The parameter <ttl> specifies the TTL for the session id(not for the key-value pair).
func (s *StorageRedis) SetMap(id string, data map[string]interface{}, ttl time.Duration) error {
	return ErrorDisabled
}

// Remove deletes key with its value from storage.
func (s *StorageRedis) Remove(id string, key string) error {
	return ErrorDisabled
}

// RemoveAll deletes all key-value pairs from storage.
func (s *StorageRedis) RemoveAll(id string) error {
	return ErrorDisabled
}

// GetSession returns the session data as *vmap.StrAnyMap for given session id from storage.
//
// The parameter <ttl> specifies the TTL for this session, and it returns nil if the TTL is exceeded.
// The parameter <data> is the current old session data stored in memory,
// and for some storage it might be nil if memory storage is disabled.
//
// This function is called ever when session starts.
func (s *StorageRedis) GetSession(id string, ttl time.Duration, data *vmap.StrAnyMap) (*vmap.StrAnyMap, error) {
	fmt.Printf("StorageRedis.GetSession: %s, %v", id, ttl)
	r, err := s.redis.DoVar("GET", s.key(id))
	if err != nil {
		return nil, err
	}
	content := r.Bytes()
	if len(content) == 0 {
		return nil, nil
	}
	var m map[string]interface{}
	if err = json.Unmarshal(content, &m); err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	if data == nil {
		return vmap.NewStrAnyMapFrom(m, true), nil
	} else {
		data.Replace(m)
	}
	return data, nil
}

// SetSession updates the data map for specified session id.
// This function is called ever after session, which is changed dirty, is closed.
// This copy all session data map from memory to storage.
func (s *StorageRedis) SetSession(id string, data *vmap.StrAnyMap, ttl time.Duration) error {
	fmt.Printf("StorageRedis.SetSession: %s, %v, %v", id, data, ttl)
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = s.redis.DoVar("SETEX", s.key(id), int64(ttl.Seconds()), content)
	return err
}

// UpdateTTL updates the TTL for specified session id.
// This function is called ever after session, which is not dirty, is closed.
// It just adds the session id to the async handling queue.
func (s *StorageRedis) UpdateTTL(id string, ttl time.Duration) error {
	fmt.Printf("StorageRedis.UpdateTTL: %s, %v", id, ttl)
	if ttl >= DefaultStorageRedisLoopInterval {
		s.updatingIdMap.Set(id, int(ttl.Seconds()))
	}
	return nil
}

// doUpdateTTL updates the TTL for session id.
func (s *StorageRedis) doUpdateTTL(id string, ttlSeconds int) error {
	fmt.Printf("StorageRedis.doUpdateTTL: %s, %d", id, ttlSeconds)
	_, err := s.redis.DoVar("EXPIRE", s.key(id), ttlSeconds)
	return err
}

func (s *StorageRedis) key(id string) string {
	return s.prefix + id
}
