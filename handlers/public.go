package handlers

import (
	"net/http"
	"public-service/models"

	"github.com/gin-gonic/gin"
)

// GetPublicPastes renders the HTML view for public pastes.
func GetPublicPastes(c *gin.Context) {
	pastes, err := models.GetPublicPastes()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "public.html", gin.H{
			"error": "Failed to fetch public pastes",
		})
		return
	}
	c.HTML(http.StatusOK, "public.html", gin.H{
		"pastes": pastes,
	})
}
