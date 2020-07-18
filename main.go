package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/kenchan0130/twitter-like-feed/routers"
)

func main() {
	port := 8080
	if len(os.Args) > 1 {
		p, _ := strconv.Atoi(os.Args[1])
		port = p
	}

	router := routers.InitRouter()
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}
}
