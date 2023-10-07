package items

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	itemID, err := strconv.ParseUint(c.Param("itemID"), 10, 64)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid item ID")},
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

	item, err := m.GetItemByUser(userID, itemID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		c.Abort()
		return
	}
	if item.ItemID == 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("Could not find item with ID %d", itemID)},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, item)
}
