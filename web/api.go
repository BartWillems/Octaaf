package web

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gobuffalo/pop"
)

type Connections struct {
	Octaaf      *tgbotapi.BotAPI
	Postgres    *pop.Connection
	Redis       *redis.Client
	KaliID      int64
	Environment string
}

var conn Connections

func New(c Connections) *gin.Engine {
	conn = c

	if conn.Environment == "production" {
		gin.SetMode("release")
	}

	router := gin.Default()
	api := router.Group("/api/v1")
	{
		api.GET("/health", health)
	}

	kali := api.Group("/kali")
	{
		kali.GET("/quote", quote)
	}

	return router
}
