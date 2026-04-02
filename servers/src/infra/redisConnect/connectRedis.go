package redisConnect

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	// 创建 Redis 客户端
	RD := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // 刚才我们启动的地址
		Password: "",               // 默认没设密码
		DB:       0,                // 默认数据库

		PoolSize:     10, // 最大连接数，10个完全够你的卡牌游戏了
		MinIdleConns: 2,  // 最小空闲连接
		MaxRetries:   3,  // 命令执行失败后的重试次数
		DialTimeout:  5 * time.Second,
	})
	ctx := context.Background()
	_, err := RD.Ping(ctx).Result()
	if err != nil {
		log.Fatalln(err.Error())
		return nil
	}
	fmt.Println("Successfully connected to redis")
	return RD
}
