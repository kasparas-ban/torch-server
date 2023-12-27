package items

import (
	"net/http"
	a "torch/torch-server/auth"
	m "torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

func HandleUpdateItemStatus(c *gin.Context) {
	userID, err := a.GetUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	var body m.UpdateItemStatusReq
	if err := c.BindJSON(&body); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid item object"},
		)
		c.Abort()
		return
	}

	updatedItem, err := m.UpdateItemStatus(userID, body)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, updatedItem)
}
