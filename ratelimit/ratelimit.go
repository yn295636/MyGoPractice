package ratelimit

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/yn295636/MyGoPractice/common"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	LuaScriptName = "redis_rate_limit.lua"
	KeyFormat     = "ratelimit_%v"
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

	file, err := os.Open(fmt.Sprintf("../../%v", LuaScriptName))
	if err != nil {
		log.Printf("Opening %v got error %v", LuaScriptName, err)
		return false, err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Reading %v got error %v", LuaScriptName, err)
		return false, err
	}
	luaContent := common.UnsafeBytesToString(bytes)
	lua := redis.NewScript(2, luaContent)
	result, err := redis.Uint64(lua.Do(redisConn, fmt.Sprintf(KeyFormat, key), "window", time.Now().Unix(), windowMs))
	if err != nil {
		log.Printf("Executing lua script %v got error %v", LuaScriptName, err)
		return false, err
	}
	if result == 0 {
		return false, errors.New(fmt.Sprintf("Executing lua script %v got wrong result 0", LuaScriptName))
	}
	return true, nil
}
