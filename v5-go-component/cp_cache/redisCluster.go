package cp_cache

import (
	"context"
	r "github.com/go-redis/redis/v8"
	"golang.org/x/exp/errors/fmt"
	"strings"
	"time"
	"warehouse/v5-go-component/cp_log"
)

type redisCluster struct {
	Ctx context.Context
	*r.ClusterClient
}

var ctx = context.Background()

func NewRedisCluster(ip, port, password string) (ICache, error) {
	var ipList, portList, svrList []string

	ipList = strings.Split(ip, ",")
	portList = strings.Split(port, ",")

	for i, v := range ipList {
		svrList = append(svrList, v + ":" + portList[i])
	}

	cluster := r.NewClusterClient(&r.ClusterOptions{
		Addrs:  svrList,
		Password: password,
		DialTimeout: 3 * time.Second,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	_, err := cluster.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("[redis-cluster] [%v] ping err:%v", svrList, err.Error()))
	}

	cp_log.Info("[redis-cluster] ping success.")

	return &redisCluster{ClusterClient: cluster, Ctx: ctx}, nil
}

func (this redisCluster) Put(key string, value interface{}, expire time.Duration) error {
	_, err := this.ClusterClient.Set(ctx, key, value, expire).Result()
	if err != nil {
		return err
	}

	return nil
}

func (this redisCluster) Delete(key string) error {
	_, err := this.ClusterClient.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (this redisCluster) Get(key string) (string, error) {
	return this.ClusterClient.Get(ctx, key).Result()
}

func (this redisCluster) IsExist(key string) bool {
	result := this.ClusterClient.Exists(ctx, key).String()

	if result == "true"{
		return true
	}

	return false
}

func (this redisCluster) ZAdd(key, member string, score int64) error {
	return nil
}

func (this redisCluster) ZScore(key, member string) (int64, error) {
	return 0, nil
}

func (this redisCluster) ZRangeByScore(key string, min, max int64) ([]interface{}, error) {
	return nil, nil
}

func (this redisCluster) ZRem(key, member string) error {
	return nil
}

func (this redisCluster) ZCount(key string, min, max int64) (int, error) {
	return 0, nil
}

func (this redisCluster) ZRemRangeByScore(key string, min, max int64) error {
	return nil
}

func (this redisCluster) Expire(key string, second int64) error {
	return nil
}

func (this redisCluster) ClearAll() error {
	return nil
}

func (this redisCluster) Decr(key string) error {
	return nil
}

func (this redisCluster) GetMulti(keys []string) []interface{} {
	return nil
}

func (this redisCluster) Incr(key string) error {
	return nil
}

func (this redisCluster) StartAndGC(config string) error {
	return nil
}

func (rc redisCluster) LPUSH(key string, member string) error {
	return nil
}

func (rc redisCluster) RPUSH(key string, member string) error {
	return nil
}

func (rc *redisCluster) LLEN(key string) (int, error) {
	return 0, nil
}

func (rc *redisCluster) LPOP(key string) (string, error) {
	return "", nil
}

func (rc *redisCluster) RPOP(key string) (string, error) {
	return "", nil
}

func (rc *redisCluster) LRANGE(key string, min, max int64) ([]string, error) {
	return nil, nil
}

func (rc redisCluster) SADD(key string, member string) error {
	return nil
}

func (rc redisCluster) SREM(key string, member string) error {
	return nil
}

func (this redisCluster) SMEMBERS(key string) ([]string, error) {
	return nil, nil
}

func (rc *redisCluster) SCARD(key string) (int, error) {
	return 0, nil
}

func (rc *redisCluster) SPOP(key string) (string, error) {
	return "", nil
}

