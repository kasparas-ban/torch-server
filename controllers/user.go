package controllers

import (
	"errors"
	"net/http"
	m "torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

type UserPost struct {
	ClerkID string
}

func GetUserInfoByClerkID(c *gin.Context) {
	clerkID := c.Param("clerkID")
	if clerkID == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("No clerkID found")},
		)
		c.Abort()
		return
	}

	user, err := m.GetUserByClerkID(clerkID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUserInfo(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("No userID found")},
		)
		c.Abort()
		return
	}

	user, err := m.GetUserByUserID(userID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}
