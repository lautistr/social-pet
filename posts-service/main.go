package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/lautistr/social-pet/database"
	"github.com/lautistr/social-pet/enums"
	"github.com/lautistr/social-pet/events"
	"github.com/lautistr/social-pet/models"
)

func setupRouter() (router *gin.Engine) {
	router = gin.Default()
	router.POST(enums.ENDPOINTS_POSTS, createPostHandler)
	return router
}

func main() {
	var cfg models.Config
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatal(err)
	}

	database.SetupPostgresRepository(cfg)
	events.SetupNatsMessageBroker(cfg.NatsAddress)

	router := setupRouter()
	router.Run(enums.POSTS_PORT)

}
