package cache

import (
	"github.com/go-redis/redis/v7"
	"init.bulbul.tv/bulbul-backend/saas-gateway/config"
	"sync"
)

var (
	singletonRedis *redis.ClusterClient
	redisInit      sync.Once
)

func Redis(config config.Config) *redis.ClusterClient {
	redisInit.Do(func() {
		singletonRedis = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    config.RedisClusterAddresses,
			PoolSize: 100,
			Password: config.REDIS_CLUSTER_PASSWORD,
		})
	})
	return singletonRedis
}
