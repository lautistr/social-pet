package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/lautistr/social-pet/database"
	"github.com/lautistr/social-pet/enums"
	"github.com/lautistr/social-pet/events"
	"github.com/lautistr/social-pet/repository"
)

type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
}

func setupRouter() (router *gin.Engine) {
	router = gin.Default()
	router.POST("/posts", createPostHandler)
	return router
}

func setupRepository(config Config) {
	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", config.PostgresUser, config.PostgresPassword, config.PostgresDB)
	repo, err := database.NewPostgresRepository(addr)
	if err != nil {
		log.Fatal(err)
	}
	repository.SetRepository(repo)
}

func setupMessageBroker(address string) {
	n, err := events.NewNats(fmt.Sprintf("nats://%s", address))
	if err != nil {
		log.Fatal(err)
	}
	events.SetEventStore(n)

	defer events.Close()
}

func main() {
	var cfg Config
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatal(err)
	}

	setupMessageBroker(cfg.NatsAddress)
	setupRepository(cfg)

	router := setupRouter()
	router.Run(enums.POSTS_PORT)

}
