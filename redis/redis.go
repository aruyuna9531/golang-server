package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
)

type RedisCli struct {
	c *redis.Client
}

var redisCli = &RedisCli{}

func GetRedisCli() *RedisCli {
	return redisCli
}

func (r *RedisCli) Init() {
	r.c = redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: "",
		DB:       0,
	})
	log.Println("Redis client inited")
}

func (r *RedisCli) OnClose() {
	err := r.c.Close()
	if err != nil {
		panic(err)
	}
	log.Println("Redis client closed")
}

func (r *RedisCli) Set(K, V string) {
	cmd := r.c.Set(context.Background(), K, V, 0)
	if cmd.Err() != nil {
		panic(cmd.Err())
	}
}

func (r *RedisCli) Get(K string) string {
	cmd := r.c.Get(context.Background(), K)
	if cmd.Err() != nil {
		panic(cmd.Err())
	}
	res, err := cmd.Result()
	if errors.Is(err, redis.Nil) {
		return ""
	} else if err != nil {
		panic(cmd.Err())
	}
	return res
}

func (r *RedisCli) Test() {
	r.Set("1", "1")
	log.Printf("%s", r.Get("1"))
}
