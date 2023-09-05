package controllers

import (
	"errors"
	"net/http"
	a "torch/torch-server/auth"
	m "torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

func GetUserInfoByClerkID(c *gin.Context) {
	clerkID := a.GetClerkID(c)

	user, err := m.GetUserInfoByClerkID(clerkID)
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
	userID := a.GetUserID(c)

	user, err := m.GetUserInfo(userID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": errors.New("Failed to get user info")},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}
