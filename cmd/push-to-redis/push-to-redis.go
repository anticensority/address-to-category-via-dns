package main

import (
	"fmt"
	"os"
	"strconv"
	"github.com/gomodule/redigo/redis"
	_ "github.com/joho/godotenv/autoload"
	"github.com/anticensority/address-to-category-via-dns/pkg/types"
)

func pushToRedis(addressToIntCat *types.AddressToIntCat, redisAddress string, redisPassword string) {
	fmt.Println("Connecting to Redis...")
	redisCon, err := redis.Dial("tcp", redisAddress, redis.DialPassword(redisPassword))
	if err != nil {
		panic(err)
	}
	defer redisCon.Close()
	msetArgs := make([]interface{}, 0, 2*len(addressToIntCat.Hostnames))
	for hostname, cat := range addressToIntCat.Hostnames {
		msetArgs = append(msetArgs, hostname, strconv.Itoa(cat))
	}
	s, err := redis.String(redisCon.Do("MSET", msetArgs...))
	fmt.Println(s)
	
	fmt.Println("Connected to the redis server.")
}

func main() {

	addressToIntCat := &types.AddressToIntCat{
		Hostnames: map[string]int{
			"kasparov.ru": 1,
		},
	}

	pushToRedis(addressToIntCat, "127.0.0.1:6379", os.Getenv("REDIS_PASSWORD"))
}
