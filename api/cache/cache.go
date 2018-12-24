package cache

import (
	"github.com/ReneKroon/ttlcache"
	"github.com/hoisie/redis"
	"time"
)

type Cache struct {
	redisStore bool
	rds        *redis.Client
	mem        *ttlcache.Cache
	expire     int64
}

func New(redisStore bool, expire int64) *Cache {
	var c Cache
	c.redisStore = redisStore
	c.mem = ttlcache.NewCache()
	c.mem.SetTTL(time.Second * time.Duration(expire))
	c.expire = expire
	return &c
}

func (c *Cache) InitRedisCache(dial, pswd string, db, connections int) {
	var client redis.Client
	client.Addr = dial
	client.Password = pswd
	client.Db = db
	client.MaxPoolSize = connections
	c.rds = &client
}

func (c *Cache) SetExNx(key string, value []byte) (bool, error) {
	if c.redisStore {
		b, err := c.rds.Setnx(key, value)
		if err != nil {
			return b, err
		}
		if b {
			b, err = c.rds.Expire(key, c.expire)
			return b, err
		}
		return false, nil
	}
	_, finded := c.mem.Get(key)
	if !finded {
		c.mem.Set(key, value)
		return true, nil
	}
	return false, nil
}
