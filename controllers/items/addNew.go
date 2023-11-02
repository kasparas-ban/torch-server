package items

import (
	"net/http"
	a "torch/torch-server/auth"
	m "torch/torch-server/models"
	"torch/torch-server/util"

	"github.com/gin-gonic/gin"
)

func HandleAddItem(c *gin.Context) {
	userID := a.GetUserID(c)
	itemType := c.Param("type")

	switch itemType {
	case "task":
		SaveItem[m.Task](c, userID)
	case "goal":
		SaveItem[m.Goal](c, userID)
	case "dream":
		SaveItem[m.Dream](c, userID)
	default:
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid item type"},
		)
		c.Abort()
		return
	}
}

func SaveItem[T m.GeneralItem](c *gin.Context, userID uint64) {
	var newItem T
	if err := c.BindJSON(&newItem); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid item object"},
		)
		c.Abort()
		return
	}

	publicItemID, err := util.New()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to save the item"},
		)
		c.Abort()
		return
	}

	item, err := newItem.AddItem(userID, publicItemID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, item)
}
