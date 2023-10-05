package items

import (
	"errors"
	"net/http"
	"strconv"
	a "torch/torch-server/auth"
	"torch/torch-server/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Remove item from the items table
		err := tx.Exec(`
			DELETE FROM items WHERE user_id = ? AND item_id = ?
		`, userID, itemID).Error
		if err != nil {
			return err
		}

		// Remove all relationships with the item
		err = tx.Exec(`
			DELETE FROM item_relations WHERE user_id = ? AND (parent_id = ? OR child_id = ?)
		`, userID, itemID, itemID).Error
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
