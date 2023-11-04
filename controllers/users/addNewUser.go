package users

import (
	"net/http"
	m "torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

func HandleAddNewUser(c *gin.Context) {
	// clerkID := c.GetString("clerkID")
	clerkID := "123"
	if clerkID == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Could not find clerkID"},
		)
		c.Abort()
		return
	}

	var userReq m.NewUser
	if err := c.BindJSON(&userReq); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid user object"},
		)
		c.Abort()
		return
	}

	if err := userReq.Validate(); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid user object"},
		)
		c.Abort()
		return
	}

	user, err := m.AddUser(clerkID, userReq)
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
