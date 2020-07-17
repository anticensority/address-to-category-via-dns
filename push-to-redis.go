package main

import (
	"fmt"
	"os"
	"github.com/gomodule/redigo/redis"
	_ "github.com/joho/godotenv/autoload"
)

func pushToRedis(redisAddress string, redisPassword string) {
	fmt.Println("Connecting to Redis...")
	_, err := redis.Dial("tcp", redisAddress, redis.DialPassword(redisPassword))
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to the redis server.")
}

func main() {
	pushToRedis("127.0.0.1:6379", os.Getenv("REDIS_PASSWORD"))
}
