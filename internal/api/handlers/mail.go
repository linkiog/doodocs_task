package handlers

import (
	"net/http"

	"github.com/linkiog/doodocs/internal/services"

	"github.com/gin-gonic/gin"
)

func SendFileByEmail(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read file"})
		return
	}
	defer file.Close()

	emails := c.PostForm("emails")
	if emails == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Emails are required"})
		return
	}

	err = services.SendFileByEmail(file, header, emails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
