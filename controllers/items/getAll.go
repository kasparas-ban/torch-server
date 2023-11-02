package items

import (
	"errors"
	"fmt"
	"net/http"

	m "torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

func HandleGetAllItems(c *gin.Context) {
	userID := c.GetUint64("userID")
	if userID == 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Could not find userID"},
		)
		c.Abort()
		return
	}

	items, err := m.GetAllItemsByUser(userID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, items)
}

func HandleGetItem(c *gin.Context) {
	publicItemID := c.Param("itemID")
	if publicItemID == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid item ID"},
		)
		c.Abort()
		return
	}

	userID := c.GetUint64("userID")
	if userID == 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Could not find userID"},
		)
		c.Abort()
		return
	}

	item, err := m.GetItemByUser(userID, publicItemID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		c.Abort()
		return
	}
	if item.PublicItemID == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("Could not find item with ID %s", publicItemID)},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, item)
}
