package redis

import (
	"fmt"
	"strconv"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/anticensority/address-to-category-via-dns/pkg/types"
)

func Push(addressToIntCat *types.AddressToIntCat, redisAddress string, redisPassword string) {

	fmt.Println("Connecting to Redis...")
	redisCon, err := redigo.Dial("tcp", redisAddress, redigo.DialPassword(redisPassword))
	if err != nil {
		panic(err)
	}
	defer redisCon.Close()
	fmt.Println("Connected to the redis server.")

	msetArgs := make([]interface{}, 0, 2*len(addressToIntCat.Hostnames))
	for hostname, cat := range addressToIntCat.Hostnames {
		msetArgs = append(msetArgs, hostname, strconv.Itoa(cat))
	}
	s, err := redigo.String(redisCon.Do("MSET", msetArgs...))
	fmt.Println(s)

	//msetArgs = make([]interface{}, 0, 2*len(addressToIntCat.Ipv4Subnets))
	//for ipv4Subnet, cat := range addressToIntCat.Ipv4Subnets {
	//	msetArgs = append(msetArgs, ipv4Subnet, strconv.Itoa(cat))
	//}
	//s, err := redigo.String(redisCon.Do("MSET", msetArgs...))
	//fmt.Println(s)
	
}
