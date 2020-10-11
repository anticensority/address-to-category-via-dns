package main

import (
	"os"
	_ "github.com/joho/godotenv/autoload"
	"github.com/anticensority/address-to-category-via-dns/types"
	"github.com/anticensority/address-to-category-via-dns/outputs/redis"
)

func main() {

	addressToIntCat := &types.AddressToIntCat{
		Hostnames: map[string]int{
			"kasparov.ru": 1,
		},
	}

	redis.Push(addressToIntCat, "127.0.0.1:6379", os.Getenv("REDIS_PASSWORD"))
}
