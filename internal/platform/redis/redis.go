package redis

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"fmt"
	"github.com/go-redis/redis/extra/redisotel/v9"
	"github.com/go-redis/redis/v9"
)

func New(cnf config.Config) *redis.Client {
	addr := fmt.Sprintf("%v:%v", cnf.RedisHost, cnf.RedisPort)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cnf.RedisPassword,
	})

	client.AddHook(redisotel.NewTracingHook())

	return client
}
