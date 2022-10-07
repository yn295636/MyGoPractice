package ratelimit

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	KeyFormat  = "ratelimit_%v"
	luaContent = `
if (redis.call('exists', KEYS[1]) == 0) then
    redis.call('rpush', KEYS[1], ARGV[1]);
    return redis.call('pexpire', KEYS[1], ARGV[2]);
else
    return redis.call('rpushx', KEYS[1], ARGV[1]);
end;
	`
)

func Ratelimit(redisConn redis.Conn, key string, windowMs, threshold uint) (bool, error) {
	redisKey := fmt.Sprintf(KeyFormat, key)
	curLen, err := redis.Uint64(redisConn.Do("LLEN", redisKey))
	if err != nil {
		log.Printf("Getting %v length from redis got error %v", redisKey, err)
		return false, err
	}
	if curLen >= uint64(threshold) {
		return false, nil
	}

	lua := redis.NewScript(2, luaContent)
	result, err := redis.Uint64(lua.Do(redisConn, fmt.Sprintf(KeyFormat, key), "window", time.Now().Unix(), windowMs))
	if err != nil {
		log.Printf("Executing lua script got error %v", err)
		return false, err
	}
	if result == 0 {
		return false, errors.New("executing lua script got wrong result 0")
	}
	return true, nil
}
