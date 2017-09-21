package libs

import (
	"github.com/go-redis/redis"
	appConfig "SSO/apps/config"
	"fmt"
)

func NewRedisClient() *redis.Client {

	fmt.Println(appConfig.Config.REDIS.Url);

	client := redis.NewClient(&redis.Options{
		Addr:     appConfig.Config.REDIS.Url,
		Password: appConfig.Config.REDIS.Password, // no password set
		DB:       1,  // use default DB
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return client
}
