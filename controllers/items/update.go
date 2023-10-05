package items

import (
	"errors"
	"net/http"
	a "torch/torch-server/auth"
	"torch/torch-server/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateItem(c *gin.Context) {
	userID := a.GetUserID(c)

	var newItem Item
	if err := c.BindJSON(&newItem); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid item object")},
		)
		c.Abort()
		return
	}

	err := newItem.update(userID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, newItem)
}

func (item Item) update(userID uint64) error {
	var updatedItem Item

	err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Update item in the items table
		err := tx.Exec(`
			UPDATE items
			SET
				title = ?, 
				target_date = ?,
				priority = ?,
				duration = ?,
				parent_id = ?
			WHERE
				user_id = ? AND item_id = ?
		`, item.Title, item.TargetDate, item.Priority, item.Duration, item.ParentID, userID, item.ItemID).Error
		if err != nil {
			return err
		}

		if item.ParentID.Valid == false {
			// Remove from relations table where child is the item
			err = tx.Exec(`
				DELETE FROM item_relations WHERE user_id = ? AND child_id = ?
			`, userID, item.ItemID).Error
			if err != nil {
				return err
			}
		} else {
			// Update relations table where child is the item
			err = tx.Exec(`
				UPDATE item_relations SET parent_id = ? WHERE user_id = ? AND child_id = ?
			`, item.ParentID, userID, item.ItemID).Error
			if err != nil {
				return err
			}
		}

		// Select the updated item
		err = tx.Raw(`
			SELECT item_id, title, type, target_date, priority, duration, rec_times, rec_period, rec_progress, parent_id, time_spent, created_at
			FROM items WHERE item_id = LAST_INSERT_ID()
		`).Scan(&updatedItem).Error
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
