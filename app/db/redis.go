package db

import (
	"github.com/up-finder/silk.web/app"
	"github.com/up-finder/silk.web/app/setup"
	"gopkg.in/redis.v3"
)

var Redis = NewRedis(app.Env)

type RedisDB struct {
	redis.Client
}

func NewRedis(env *setup.Env) *RedisDB {
	return NewCustomRedis(env.DB.Redis.Host, env.DB.Redis.Password, env.DB.Redis.DB)
}

func NewCustomRedis(addr, pass string, db int64) *RedisDB {
	redisOptions := redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	}
	return &RedisDB{*redis.NewClient(&redisOptions)}
}
