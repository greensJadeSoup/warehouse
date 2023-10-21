package cp_cache

import (
	"fmt"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
)

type redisSingel struct {
	ICache
}

func NewRedisSingel(ip, port string, password string) (ICache, error) {
	c, err := NewCache("redis", fmt.Sprintf(`{"conn":"%s","password":"%s"}`, ip + ":" + port, password))
	if err != nil {
		return nil, cp_error.NewSysError("初始化 Redis 错误：" + err.Error())
	}

	if c == nil {
		return nil, cp_error.NewSysError("初始化 Redis 错误")
	}

	cp_log.Info("redis single connect success ...")

	return redisSingel{ICache: c}, nil
}

//redisSingel直接继承ICACHE, 方法只需要在redis.go中写即可
//以下如果有需要，可以重载redis.go中的方法

//
//func (this redisSingel) Add(key string, value interface{}, expire time.Duration) error {
//	err := this.ICache.Put(key, value, expire)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (this redisSingel) Del(key string) error {
//	err := this.ICache.Delete(key)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (this redisSingel) IsExist(key string) bool {
//	return this.ICache.IsExist(key)
//}
//
//func (this redisSingel) Zadd(key, member string, score int64) error {
//	return this.ICache.ZAdd(key, member, score)
//}
//
//func (this redisSingel) ZScore(key, member string) (int64, error) {
//	return this.ICache.ZScore(key, member)
//}
//
//func (this redisSingel) ZRangeByScore(key string, min, max int64) ([]interface{}, error) {
//	return this.ICache.ZRangeByScore(key, min, max)
//}
//
//func (this redisSingel) ZRem(key, member string) error {
//	return this.ICache.ZRem(key, member)
//}
//
//func (this redisSingel) ZCount(key string, min, max int64) (int, error) {
//	return this.ICache.ZCount(key, min, max)
//}
//
//func (this redisSingel) ZRemRangeByScore(key string, min, max int64) error {
//	return this.ICache.ZRemRangeByScore(key, min, max)
//}
//
//func (this redisSingel) Expire(key string, second int64) error {
//	return this.ICache.Expire(key, second)
//}
//
//func (this redisSingel) LPUSH(key string, member string) error {
//	return this.ICache.LPUSH(key, member)
//}
//
//func (this redisSingel) LLEN(key string) (int, error) {
//	return this.ICache.LLEN(key)
//}
