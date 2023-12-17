package users

import (
	"net/http"
	a "torch/torch-server/auth"
	m "torch/torch-server/models"

	"github.com/clerkinc/clerk-sdk-go/clerk"
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

	// Delete user from Clerk
	clerkID, clerkIDExists := c.Get("clerkID")
	clerkClient, exists := c.Get("clerkData")
	client, ok := clerkClient.(clerk.Client)
	if !ok || !exists || !clerkIDExists {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Unexpected error occured"},
		)
		c.Abort()
	}
	_, err = client.Users().Delete(clerkID.(string))
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Unexpected error occured"},
		)
		c.Abort()
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
