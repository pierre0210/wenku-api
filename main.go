package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pierre0210/wenku-api/api/novel"
)

func main() {
	port := flag.Int("p", 5000, "Port")
	flag.Parse()
	router := gin.Default()

	router.GET("/novel/volume/:aid/:vid", novel.HandleGetVolume)
	//router.GET("/novel/chapter")

	addr := fmt.Sprintf("localhost:%d", *port)
	router.Run(addr)
}
