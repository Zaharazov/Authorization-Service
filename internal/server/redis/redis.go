package redis

import (
	"Authorization-Service/internal/server/configs"

	"github.com/go-redis/redis"
)

var Client *redis.Client = redis.NewClient(&redis.Options{
	Addr: configs.Redis_Addr,
	DB:   configs.DB_id,
})
