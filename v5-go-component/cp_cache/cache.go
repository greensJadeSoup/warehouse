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

// Package cache provide a Cache interface and some implement engine
// Usage:
//
// import(
//   "github.com/astaxie/beego/cache"
// )
//
// bm, err := cache.NewCache("memory", `{"interval":60}`)
//
// Use it like this:
//
//	bm.Put("astaxie", 1, 10 * time.Second)
//	bm.Get("astaxie")
//	bm.IsExist("astaxie")
//	bm.Delete("astaxie")
//
//  more docs http://beego.me/docs/module/cache.md
package cp_cache

import (
	"fmt"
	"time"
	"warehouse/v5-go-component/cp_dc"
)

// Cache interface contains all behaviors for cache adapter.
// usage:
//	cache.Register("file",cache.NewFileCache) // this operation is run in init method of file.go.
//	c,err := cache.NewCache("file","{....}")
//	c.Put("key",value, 3600 * time.Second)
//	v := c.Get("key")
//
//	c.Incr("counter")  // now is 1
//	c.Incr("counter")  // now is 2
//	count := c.Get("counter").(int)

var cp_cache ICache

type ICache interface {
	// get cached value by key.
	Get(key string) (string, error)
	// GetMulti is a batch version of Get.
	GetMulti(keys []string) []interface{}
	// set cached value with key and expire time.
	Put(key string, val interface{}, timeout time.Duration) error
	// delete cached value by key.
	Delete(key string) error
	// increase cached int value by key, as a counter.
	Incr(key string) error
	// decrease cached int value by key, as a counter.
	Decr(key string) error
	//添加一个元素，如果存在，则更新
	ZAdd(key string, member string, score int64) error
	//查看一个成员的分数值
	ZScore(key string, member string) (int64, error)
	//返回范围内的元素列表，-inf +inf 为上下无限 (min (max 表示小于，默认是小于等于
	ZRangeByScore(key string, min, max int64) ([]interface{}, error)
	//删除元素
	ZRem(key string, member string) error
	//返回区间内的数目
	ZCount(key string, min, max int64) (int, error)
	//删除score在指定范围内的元素
	ZRemRangeByScore(key string, min, max int64) error
	//为key设置过期时间
	Expire(key string, second int64) error
	// check if cached value exists or not.
	IsExist(key string) bool
	// clear all cache.
	ClearAll() error
	// start gc routine based on config string settings.
	StartAndGC(config string) error

	// list
	LPUSH(key string, member string) error
	RPUSH(key string, member string) error
	LPOP(key string) (string, error)
	RPOP(key string) (string, error)
	LLEN(key string) (int, error)
	LRANGE(key string, min, max int64) ([]string, error)

	// set
	SADD(key string, member string) error
	SREM(key string, member string) error
	SMEMBERS(key string) ([]string, error)
	SCARD(key string) (int, error)
	SPOP(key string) (string, error)
}

// Instance is a function create a new Cache Instance
type Instance func() ICache

var adapters = make(map[string]Instance)

// Register makes a cache adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewCache Create a new cache driver by adapter name and config string.
// config need to be correct JSON as string: {"interval":360}.
// it will start gc automatically.
func NewCache(adapterName, config string) (adapter ICache, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}

func InitCache(cacheConf *cp_dc.DcCacheConfig) error {
	var err error

	if cacheConf.Type == "single" {
		cp_cache, err = NewRedisSingel(cacheConf.Server, cacheConf.Port, cacheConf.Password)
	} else {
		cp_cache, err = NewRedisCluster(cacheConf.Server, cacheConf.Port, cacheConf.Password)
	}

	return err
}

func GetCache() ICache {
	return cp_cache
}