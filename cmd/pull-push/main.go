package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/anticensority/address-to-category-via-dns/inputs/zapretinfo"
	"github.com/anticensority/address-to-category-via-dns/outputs/redis"
)

func main() {

	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	if REDIS_PASSWORD == "" {
		panic("Provide REDIS_PASSWORD environment variable! The secure way of providing it is by using .env file (e.g. see https://github.com/joho/godotenv).")
	}

	addressToIntCat := zapretInfo.Pull()
	redis.Push(addressToIntCat, "127.0.0.1:6379", REDIS_PASSWORD)
	fmt.Println("Done.")
}
