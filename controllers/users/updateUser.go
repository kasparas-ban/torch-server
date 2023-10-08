package users

import (
	"net/http"
	m "torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

func HandleUpdateUser(c *gin.Context) {
	var userReq m.NewUser
	if err := c.BindJSON(&userReq); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid user object"},
		)
		c.Abort()
	}
	userReq.UserID = c.GetUint64("userID")

	updatedUser, err := m.UpdateUser(userReq)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to update the user"},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}
