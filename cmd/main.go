package main

import (
	"log"
	"os"

	"github.com/linkiog/doodocs/internal/api"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()
	api.SetupRoutes(router)

	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}
