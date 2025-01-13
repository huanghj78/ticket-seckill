package cache

import (
	"context"
	"log"
	"ticket-seckill/conf"
	"time"

	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

func Init() *redis.Client {
	c := conf.Conf.Redis
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	client := redis.NewClient(&redis.Options{
		Addr:     c.Host,
		Password: c.Password,
		DB:       0,
	})
	Client = client
	if err := Client.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	return client
}
