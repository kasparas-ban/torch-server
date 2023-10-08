package users

import (
	"net/http"
	"torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

func HandleGetUserInfo(c *gin.Context) {
	userID := c.GetUint64("userID")
	if userID == 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Could not find userID"},
		)
		c.Abort()
	}

	user, err := models.GetUserInfo(userID)
	if err != nil || user.PublicUserID == "" {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get user info"},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUserInfoByClerkID(c *gin.Context) {
	clerkID := c.GetString("clerkID")
	if clerkID != "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Could not find clerkID"},
		)
		c.Abort()
	}

	user, err := models.GetUserByClerkID(clerkID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		c.Abort()
		return
	}
	if user.UserID == 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{"error": "User not found"},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}
