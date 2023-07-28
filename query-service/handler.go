package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lautistr/social-pet/enums"
	"github.com/lautistr/social-pet/events"
	"github.com/lautistr/social-pet/models"
	"github.com/lautistr/social-pet/repository"
	"github.com/lautistr/social-pet/search"
)

func onCreatedPost(m events.CreatedPostMessage) {
	Post := models.Post{
		ID:        m.Post.ID,
		Body:      m.Post.Body,
		CreatedAt: m.Post.CreatedAt,
	}
	if err := search.IndexPost(context.Background(), Post); err != nil {
		log.Println(err)
	}
}

func listPostsHandler(c *gin.Context) {
	posts, err := repository.ListPosts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func searchHandler(c *gin.Context) {
	query := c.Query(enums.PARAMS_QUERY)
	if len(query) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": enums.ERROR_QUERY_REQUIRED})
		return
	}

	posts, err := search.SearchPost(c, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}
