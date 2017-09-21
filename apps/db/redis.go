package db

import (
	"fmt"
	"github.com/go-redis/redis"
	appConfig "SSO/apps/config"
)

func NewRedisClient() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     appConfig.Config.REDIS.Url,
		//Password: "cjaeks^&*()", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return client
}