// internal/api/router.go
package api

import (
	"github.com/linkiog/doodocs/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")

	api.POST("/archive/information", handlers.GetArchiveInformation)
	api.POST("/archive/files", handlers.CreateArchive)

	api.POST("/mail/file", handlers.SendFileByEmail)
}
