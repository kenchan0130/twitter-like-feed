package main

import (
	"fmt"
	"github.com/kenchan0130/twitter-like-feed/routers"
	"log"
	"os"
	"strconv"
)

func main() {
	port := 8080
	if len(os.Args) > 1 {
		p, _ := strconv.Atoi(os.Args[1])
		port = p
	}

	router := routers.Init()
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}
}
