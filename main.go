package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pierre0210/wenku-api/api/auth"
	"github.com/pierre0210/wenku-api/api/novel"
	"github.com/pierre0210/wenku-api/internal/database"
)

func main() {
	port := flag.Int("p", 5000, "Port")
	flag.Parse()

	err := godotenv.Load()
	router := gin.Default()
	database.InitDatabase()

	if err != nil {
		log.Fatalln("Failed to load .env")
	}

	novelRouter := router.Group("/novel")
	novelRouter.GET("/volume/:aid/:vid", novel.HandleGetVolume)
	novelRouter.GET("/chapter/:aid/:vid/:cid", novel.HandleGetChapter)
	novelRouter.GET("/index/:aid", novel.HandleGetIndex)

	authRouter := router.Group("/auth")
	//authRouter.POST("/signup")
	authRouter.POST("/signin", auth.HandleLogin)

	addr := fmt.Sprintf("localhost:%d", *port)
	router.Run(addr)
}
