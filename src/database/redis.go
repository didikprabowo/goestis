package database

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/k/config"
)

func ConnectionRedis() *redis.Client {

	co := config.GetConfig()
	c := co.Redis

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", c.Host, c.Port),
		Password: c.Password,
		DB:       0,
	})

	pong, err := client.Ping().Result()

	if err != nil {
		fmt.Println("Gagal konek redis bro", err.Error())
	} else {
		fmt.Println("Berhasil konek Redis...", pong)
	}

	return client

}
