package items

import (
	"errors"
	"net/http"
	"strconv"
	a "torch/torch-server/auth"
	"torch/torch-server/db"

	"github.com/gin-gonic/gin"
)

func RemoveItem(c *gin.Context) {
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

	err = remove(userID, itemID)
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

func remove(userID, itemID uint64) error {
	err := db.GetDB().Exec(`
		CALL DeleteItem(?, ?)
	`, userID, itemID).Error

	return err
}
