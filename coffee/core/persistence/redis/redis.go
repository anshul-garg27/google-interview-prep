package redis

import (
	"strings"
	"sync"

	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
)

var singletonRedis *redis.ClusterClient
var redisInit sync.Once

func Redis() *redis.ClusterClient {
	redisInit.Do(func() {
		redisClusterAddressesStr := viper.GetString("REDIS_CLUSTER_ADDRESSES")
		redisClusterAddresses := strings.Split(redisClusterAddressesStr, ",")
		password := viper.GetString("REDIS_CLUSTER_PASSWORD")
		singletonRedis = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    redisClusterAddresses,
			PoolSize: 100,
			Password: password,
		})
	})
	return singletonRedis
}
func GetValueByKey(key string) (string, error) {
	redisClient := Redis()
	value, err := redisClient.Get(key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}
func DeleteValueByKey(key string) error {
	redisClient := Redis()
	return redisClient.Del(key).Err()

}
