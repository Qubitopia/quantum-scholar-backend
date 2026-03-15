package handlers

import (
	"net/http"

	"github.com/Qubitopia/quantum-scholar-backend/database"
	"github.com/Qubitopia/quantum-scholar-backend/models"

	"github.com/gin-gonic/gin"
)


func ListAssignedTestToUser(c *gin.Context) {
	userID := c.Param("userID")

	var testAssignedToUser []models.TestAssignedToUser
	if err := database.DB.Where("user_id = ?", userID).Find(&testAssignedToUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve assigned tests"})
		return
	}
	c.JSON(http.StatusOK, testAssignedToUser)
}