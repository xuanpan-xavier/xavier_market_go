package common

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

var rds *redis.Conn

func RedisPollInit() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5, //最大空闲数
		MaxActive:   0, //最大连接数，0不设上
		Wait:        true,
		IdleTimeout: time.Duration(1) * time.Second, //空闲等待时间
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379") //redis IP地址
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			redis.DialDatabase(0)
			rds = &c
			return c, err
		},
	}
}

func GetRedis() redis.Conn {
	return *rds
}
