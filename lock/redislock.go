package lock

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"log"
)

const (
	KeyFormat = "lock_%v"
)

func RedisLock(conn redis.Conn, key string, expireMs uint64) string {
	finalKey := fmt.Sprintf(KeyFormat, key)
	keyValUuid, err := uuid.NewUUID()
	if err != nil {
		log.Printf("NewUUID got error %v", err)
		return ""
	}
	keyValue := keyValUuid.String()
	result, err := redis.String(conn.Do("SET", finalKey, keyValue, "PX", expireMs, "NX"))
	if err != nil {
		return ""
	}
	if result != "OK" {
		return ""
	}
	return keyValue
}

func RedisUnlock(conn redis.Conn, key, value string) {
	finalKey := fmt.Sprintf(KeyFormat, key)
	luaScript := "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
	lua := redis.NewScript(1, luaScript)
	_, err := lua.Do(conn, finalKey, value)
	if err != nil {
		log.Printf("Executing lua script redis unlock got error %v", err)
	}
}
