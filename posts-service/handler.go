package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lautistr/social-pet/enums"
	"github.com/lautistr/social-pet/events"
	"github.com/lautistr/social-pet/models"
	"github.com/lautistr/social-pet/repository"
	"github.com/segmentio/ksuid"
)

type createPostRequest struct {
	Body string `json:"body"`
}

func createPostHandler(c *gin.Context) {
	var req createPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdAt := time.Now().UTC()
	id, err := ksuid.NewRandom()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": enums.ERROR_CREATE_POST})
		return
	}
	Post := models.Post{
		ID:        id.String(),
		Body:      req.Body,
		CreatedAt: createdAt,
	}

	if err := repository.InsertPost(c.Request.Context(), &Post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := events.PublishCreatedPost(c.Request.Context(), &Post); err != nil {
		log.Printf("%s: %v", enums.ERROR_PUBLISH_POST_EVENT, err)
	}

	c.JSON(http.StatusCreated, Post)
}
