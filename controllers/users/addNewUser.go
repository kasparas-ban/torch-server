package users

import (
	"net/http"
	m "torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

func HandleAddNewUser(c *gin.Context) {
	clerkID := c.GetString("clerkID")
	if clerkID == "" {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": "Could not find clerkID"},
		)
		return
	}

	var userReq m.NewUser
	if err := c.BindJSON(&userReq); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid user object"},
		)
		return
	}

	if err := userReq.Validate(); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid user object"},
		)
		return
	}

	user, err := m.AddUser(clerkID, userReq)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(http.StatusOK, user)
}
