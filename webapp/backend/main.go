package main

import (
	// Simple HTTP Framework
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// routes
	// r.POST("/generate", handlers.GeneratePodcast)

	// start server on port 8080
	r.Run(":8080")
}
