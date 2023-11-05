package users

import (
	"net/http"
	a "torch/torch-server/auth"
	m "torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

func HandleUpdateUser(c *gin.Context) {
	var userData m.UpdateUserReq
	if err := c.BindJSON(&userData); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid user object"},
		)
		return
	}

	userID, err := a.GetUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	if err := userData.Validate(); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid user object"},
		)
		return
	}

	updatedUser, err := m.UpdateUser(userID, userData)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to update the user"},
		)
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}
