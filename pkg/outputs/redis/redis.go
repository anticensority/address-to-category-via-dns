package redis

import (
	"fmt"
	"os"
	"strconv"
	redigo "github.com/gomodule/redigo/redis"
	_ "github.com/joho/godotenv/autoload"
	"github.com/anticensority/address-to-category-via-dns/pkg/types"
)

func PushToRedis(addressToIntCat *types.AddressToIntCat, redisAddress string, redisPassword string) {
	fmt.Println("Connecting to Redis...")
	redisCon, err := redigo.Dial("tcp", redisAddress, redigo.DialPassword(redisPassword))
	if err != nil {
		panic(err)
	}
	defer redisCon.Close()
	msetArgs := make([]interface{}, 0, 2*len(addressToIntCat.Hostnames))
	for hostname, cat := range addressToIntCat.Hostnames {
		msetArgs = append(msetArgs, hostname, strconv.Itoa(cat))
	}
	s, err := redigo.String(redisCon.Do("MSET", msetArgs...))
	fmt.Println(s)
	
	fmt.Println("Connected to the redis server.")
}
