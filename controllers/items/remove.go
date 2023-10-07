package items

import (
	"errors"
	"net/http"
	"strconv"
	a "torch/torch-server/auth"
	m "torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

func HandleRemoveItem(c *gin.Context) {
	userID := a.GetUserID(c)

	itemID, err := strconv.ParseUint(c.Param("itemID"), 10, 64)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid item ID")},
		)
		c.Abort()
		return
	}

	err = m.RemoveItem(userID, itemID)
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
