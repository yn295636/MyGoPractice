package redis

import (
	rredis "github.com/gomodule/redigo/redis"
	"time"
)

var (
	redisPool *rredis.Pool
)

func InitRedisPool(redisAddr string) {
	redisPool = &rredis.Pool{
		MaxIdle:     10,
		MaxActive:   50,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (rredis.Conn, error) {
			c, err := rredis.Dial("tcp", redisAddr)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c rredis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func GetConn() rredis.Conn {
	return redisPool.Get()
}
