// Package main is the entry point for the spice API application.
// It sets up a Gin-based HTTP server and defines routes for the API.
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello mom!",
		})
	})
	r.Run()
}
