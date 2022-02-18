package model

import (
	"APIGenerator/util"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"os/exec"
	"time"
)

func GetPool(addr string) *redis.Pool {
	return &redis.Pool{
		Dial: func() (Conn redis.Conn, error error) {
			return redis.Dial("tcp", addr)
		},
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 86400 * time.Second,
	}
}

func InitRedis() *redis.Pool {
	config := util.LoadRedisConfig()
	if config.RedisRequire {
		if config.RedisStart {
			fmt.Println("启动 Redis ...")
			err := launchRedis()
			if err != nil {
				fmt.Println("启动 redis 失败", err)
				return nil
			}
			fmt.Println("Redis 已启动")
		} else {
			fmt.Println("Redis 无需启动")
		}
		pool := GetPool(config.RedisHost + ":" + config.RedisPort)
		pool.Get()
		return pool
	} else {
		fmt.Println("无需初始化redis")
		return nil
	}
}

func launchRedis() error {
	command := "nohup redis-server > redis.out 2>&1 &"
	cmd := exec.Command("/bin/sh", "-c", command)
	bytes, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return nil
}
