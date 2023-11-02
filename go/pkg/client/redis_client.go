package client

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var (
	// Pool - redis pool
	Pool *redis.Client
)

var ctx = context.Background()

/*func redisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", common.RedisIP),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

*/

func GetKeyKVS(key string) []byte {
	val, err := Pool.Get(ctx, key).Bytes()
	if err != nil {
		panic(err)
	}
	return val
}

func SetKeyKVS(key string, content []byte) {
	err := Pool.Set(ctx, key, content, 0).Err()
	if err != nil {
		panic(err)
	}
}
