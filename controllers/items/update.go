package items

import (
	"net/http"
	a "torch/torch-server/auth"
	m "torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

func HandleUpdateItem(c *gin.Context) {
	userID, err := a.GetUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	itemType := c.Param("type")

	switch itemType {
	case "task":
		UpdateExistingItem[m.ExistingTask](c, userID)
	case "goal":
		UpdateExistingItem[m.ExistingGoal](c, userID)
	case "dream":
		UpdateExistingItem[m.ExistingDream](c, userID)
	default:
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid item type"},
		)
		c.Abort()
		return
	}
}

func UpdateExistingItem[T m.ExistingItem](c *gin.Context, userID uint64) {

	var item T
	if err := c.BindJSON(&item); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid item object"},
		)
		c.Abort()
		return
	}

	updatedItem, err := item.UpdateItem(userID)
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
