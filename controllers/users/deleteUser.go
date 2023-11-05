package users

import (
	"net/http"
	a "torch/torch-server/auth"
	m "torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

func HandleDeleteUser(c *gin.Context) {
	userID, err := a.GetUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	err = m.DeleteUser(userID)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to delete user"},
		)
		return
	}

	c.JSON(http.StatusOK, nil)
}
