package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/anticensority/address-to-category-via-dns/pkg/types"
	"github.com/anticensority/address-to-category-via-dns/pkg/inputs/zapretinfo"
)

func main() {

	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	if REDIS_PASSWORD == "" {
		panic("Provide REDIS_PASSWORD environment variable! The secure way of providing it is by using .env file (e.g. see https://github.com/joho/godotenv).")
	}

	result := zapretInfo.Pull()
	fmt.Println(result)

	fmt.Println("Done.")
}
