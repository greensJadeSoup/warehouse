// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package redis for cache provider
//
// depend on github.com/gomodule/redigo/redis
//
// go install github.com/gomodule/redigo/redis
//
// Usage:
// import(
//   _ "github.com/astaxie/beego/cache/redis"
//   "github.com/astaxie/beego/cache"
// )
//
//  bm, err := cache.NewCache("redis", `{"conn":"127.0.0.1:11211"}`)
//
//  more docs http://beego.me/docs/module/cache.md
package cp_cache

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
	"strings"
	"github.com/garyburd/redigo/redis"

	r "github.com/garyburd/redigo/redis"
)

var (
	// DefaultKey the collection name of redis for cache adapter.
	DefaultKey = "backend"
	NilErr = redis.ErrNil
)

// Cache is Redis cache adapter.
type Cache struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	dbNum    int
	key      string
	password string
	maxIdle  int
	maxActive int
}


func init() {
	Register("redis", NewRedisCache)
}

// NewRedisCache create new redis cache with default collection name.
func NewRedisCache() ICache {
	return &Cache{key: DefaultKey}
}

// actually do the redis cmds, args[0] must be the key name.
func (rc *Cache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	//这里会把
	args[0] = rc.associate(args[0])
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// associate with config key.
func (rc *Cache) associate(originKey interface{}) string {
	return originKey.(string)
}

// Get cache from redis.
func (rc *Cache) Get(key string) (string, error) {
	v1, err := rc.do("GET", key)
	if err != nil {
		return "", err
	}

	v2, err := r.String(v1, nil)
	if err != nil {
		return "", err
	}

	return v2, nil
}

// GetMulti get cache from redis.
func (rc *Cache) GetMulti(keys []string) []interface{} {
	c := rc.p.Get()
	defer c.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, rc.associate(key))
	}
	values, err := redis.Values(c.Do("MGET", args...))
	if err != nil {
		return nil
	}
	return values
}

// Put put cache to redis.
func (rc *Cache) Put(key string, val interface{}, timeout time.Duration) error {
	var err error
	if timeout > 0 {
		_, err = rc.do("SET", key, val, "EX", int64(timeout/time.Second))
	} else {
		_, err = rc.do("SET", key, val)
	}

	return err
}

// Delete delete cache in redis.
func (rc *Cache) Delete(key string) error {
	_, err := rc.do("DEL", key)
	return err
}

// IsExist check cache's existence in redis.
func (rc *Cache) IsExist(key string) bool {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

// Incr increase counter in redis.
func (rc *Cache) Incr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, 1))
	return err
}

// Decr decrease counter in redis.
func (rc *Cache) Decr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, -1))
	return err
}

//添加一个元素，如果存在，则更新
func (rc *Cache) ZAdd(key string, member string, score int64) error {
	_, err := rc.do("ZADD", key, score, member)
	return err
}

//添加一个元素，如果存在，则更新
func (rc *Cache) ZScore(key string, member string) (int64, error) {
	score, err := redis.Int64(rc.do("ZScore", key, member))
	return score, err
}

//返回范围内的元素列表，-inf +inf 为上下无限 (min (max 表示小于，默认是小于等于
func (rc *Cache) ZRangeByScore(key string, min, max int64) ([]interface{}, error) {
	return redis.Values(rc.do("ZRANGEBYSCORE", key, min, max))
}

//删除元素
func (rc *Cache) ZRem(key string, member string) error {
	_, err := rc.do("ZREM", key, member)
	return err
}

//返回区间内的数目
func (rc *Cache) ZCount(key string, min, max int64) (int, error) {
	c, err := redis.Int(rc.do("ZCOUNT", key, min, max))
	return c, err
}

//删除score在指定范围内的元素
func (rc *Cache) ZRemRangeByScore(key string, min, max int64) error {
	_, err := rc.do("ZREMRANGEBYSCORE", key, min, max)
	return err
}

//为key设置过期时间
func (rc *Cache) Expire(key string, second int64) error {
	_, err := rc.do("EXPIRE", key, second)
	return err
}

//push到队列头
func (rc *Cache) LPUSH(key string, member string) error {
	_, err := rc.do("LPUSH", key, member)
	return err
}

//push到队列尾
func (rc *Cache) RPUSH(key string, member string) error {
	_, err := rc.do("RPUSH", key, member)
	return err
}

//队列长度
func (rc *Cache) LLEN(key string) (int, error) {
	c, err := redis.Int(rc.do("LLEN", key))
	return c, err
}

//从队列头取出
func (rc *Cache) LPOP(key string) (string, error) {
	v1, err := rc.do("LPOP", key)
	if err != nil {
		return "", err
	}

	v2, err := r.String(v1, nil)
	if err != nil {
		return "", err
	}

	return v2, nil
}

//从队列尾取出
func (rc *Cache) RPOP(key string) (string, error) {
	v1, err := rc.do("RPOP", key)
	if err != nil {
		return "", err
	}

	v2, err := r.String(v1, nil)
	if err != nil {
		return "", err
	}

	return v2, nil
}

//取出指定范围的元素
func (rc *Cache) LRANGE(key string, min, max int64) ([]string, error) {
	return redis.Strings(rc.do("LRANGE", key, min, max))
}

//添加元素
func (rc *Cache) SADD(key string, member string) error {
	_, err := rc.do("SADD", key, member)
	return err
}

//所有元素
func (rc *Cache) SMEMBERS(key string) ([]string, error) {
	return redis.Strings(rc.do("SMEMBERS", key))
}

//移除元素
func (rc *Cache) SREM(key string, member string) error {
	_, err := rc.do("SREM", key, member)
	return err
}

//元素数目
func (rc *Cache) SCARD(key string) (int, error) {
	c, err := redis.Int(rc.do("SCARD", key))
	return c, err
}

//取出并移除元素
func (rc *Cache) SPOP(key string) (string, error) {
	c, err := redis.String(rc.do("SPOP", key))
	return c, err
}


// ClearAll clean all cache in redis. delete this redis collection.
func (rc *Cache) ClearAll() error {
	c := rc.p.Get()
	defer c.Close()
	cachedKeys, err := redis.Strings(c.Do("KEYS", rc.key+":*"))
	if err != nil {
		return err
	}
	for _, str := range cachedKeys {
		if _, err = c.Do("DEL", str); err != nil {
			return err
		}
	}
	return err
}

// StartAndGC start redis cache adapter.
// config is like {"key":"collection key","conn":"connection info","dbNum":"0"}
// the cache item in redis are stored forever,
// so no gc operation.
func (rc *Cache) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)

	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}

	// Format redis://<password>@<host>:<port>
	cf["conn"] = strings.Replace(cf["conn"], "redis://", "", 1)
	if i := strings.Index(cf["conn"], "@"); i > -1 {
		cf["password"] = cf["conn"][0:i]
		cf["conn"] = cf["conn"][i+1:]
	}

	if _, ok := cf["dbNum"]; !ok {
		cf["dbNum"] = "0"
	}
	if _, ok := cf["password"]; !ok {
		cf["password"] = ""
	}
	rc.key = cf["key"]
	rc.conninfo = cf["conn"]
	rc.dbNum, _ = strconv.Atoi(cf["dbNum"])
	rc.password = cf["password"]
	rc.maxIdle = 5
	rc.maxActive = 0

	rc.connectInit()

	c := rc.p.Get()
	defer c.Close()

	return c.Err()
}

// connect to redis.
func (rc *Cache) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.conninfo)
		if err != nil {
			return nil, err
		}

		if rc.password != "" {
			if _, err := c.Do("AUTH", rc.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", rc.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxActive: rc.maxActive,
		MaxIdle:     rc.maxIdle,
		IdleTimeout: 600 * time.Second,
		Dial:        dialFunc,
	}
}