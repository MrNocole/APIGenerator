package util

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"os"
)

type RedisConfig struct {
	RedisRequire bool   `json:"redisrequire"`
	RedisHost    string `json:"redis_host"`
	RedisPort    string `json:"redis_port"`
	RedisStart   bool   `json:"redis_start"`
}

func LoadRedisConfig() *RedisConfig {
	if !CheckFileIsExist("redis.json") {
		f, err := os.Create("redis.json")
		if err != nil {
			fmt.Println("redis.json 初始化失败")
			panic(err)
		}
		initInfo := RedisConfig{
			RedisRequire: true,
			RedisHost:    "localhost",
			RedisPort:    "6379",
			RedisStart:   true,
		}
		data, err := json.Marshal(initInfo)
		if err != nil {
			fmt.Println("redis.json 初始化解析失败")
			panic(err)
		}
		_, err = f.Write(data)
		if err != nil {
			fmt.Println("redis.json 初始化写入失败")
			panic(err)
		}
		fmt.Println("redis.json 初始化完成！")
		return &initInfo
	}
	fmt.Println("redis.json 已存在")
	data, err := ioutil.ReadFile("redis.json")
	if err != nil {
		fmt.Println("redis.json 读取失败")
		panic(err)
	}
	var config *RedisConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("redis.json 解析失败")
		panic(err)
	}
	return config
}

func RedisGetSet(Conn redis.Conn, key string) ([]string, error) {
	reply, err := redis.Strings(Conn.Do("SMEMBERS", key))
	if err != nil {
		fmt.Println("查询Redis失败", err)
		return nil, err
	}
	return reply, err
}

func RedisGetHKeys(Conn redis.Conn, key string) ([]string, error) {
	reply, err := redis.Strings(Conn.Do("HKEYS", key))
	if err != nil {
		fmt.Println("查询Redis失败", err)
		return nil, err
	}
	return reply, err
}

func RedisGetHVals(Conn redis.Conn, key string) ([]string, error) {
	reply, err := redis.Strings(Conn.Do("HVALS", key))
	if err != nil {
		fmt.Println("查询Redis失败", err)
		return nil, err
	}
	return reply, err
}

func RedisInsertSet(Conn redis.Conn, key string, args string) {
	_, err := Conn.Do("SADD", key, args)
	if err != nil {
		fmt.Println("redis插入失败", err)
	}
}

func RedisInsertH(Conn redis.Conn, UUid string, key string, arg string) {
	fmt.Println("redis hmap:", UUid, " insert ", key, ":", arg)
	_, err := Conn.Do("HSET", UUid, key, arg)
	if err != nil {
		fmt.Println("redis插入失败", err)
	}
}
