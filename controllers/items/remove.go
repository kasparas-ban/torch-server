package items

import (
	"net/http"
	a "torch/torch-server/auth"
	m "torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

func HandleRemoveItem(c *gin.Context) {
	userID, err := a.GetUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	publicItemID := c.Param("itemID")
	if publicItemID == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid item ID"},
		)
		c.Abort()
		return
	}

	err = m.RemoveItem(userID, publicItemID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}
