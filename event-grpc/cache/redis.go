package cache

import (
	"github.com/go-redis/redis/v7"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"sync"
)

var (
	coreRedis *redis.ClusterClient
	redisInit sync.Once
)

type RedisCluster string

const (
	RedisClusterCore RedisCluster = "CORE"
)

func Redis(config config.Config, cluster RedisCluster) *redis.ClusterClient {
	redisInit.Do(func() {
		coreRedis = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    config.RedisClusterAddresses,
			PoolSize: 100,
			Password: config.REDIS_CLUSTER_PASSWORD,
		})

	})
	return coreRedis
}
