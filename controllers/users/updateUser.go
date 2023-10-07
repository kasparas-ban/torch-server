package users

import (
	"errors"
	"net/http"
	m "torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

func HandleUpdateUser(c *gin.Context) {
	var userReq m.NewUser
	if err := c.BindJSON(&userReq); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid user object")},
		)
		c.Abort()
	}
	if userReq.UserID == 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("User ID is missing")},
		)
		c.Abort()
	}

	updatedUser, err := m.UpdateUser(userReq)
	if err != nil || updatedUser.UserID == 0 {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": errors.New("Failed to update the user")},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}
