package items

import (
	"errors"
	"net/http"
	a "torch/torch-server/auth"
	"torch/torch-server/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExistingItem interface {
	Update(userID uint64) (updatedItem Item, err error)
}

type ExistingTask struct {
	Task
	ItemID uint64 `json:"itemID"`
}

type ExistingGoal struct {
	Goal
	ItemID uint64 `json:"itemID"`
}

type ExistingDream struct {
	Dream
	ItemID uint64 `json:"itemID"`
}

func UpdateItem(c *gin.Context) {
	userID := a.GetUserID(c)
	itemType := c.Param("type")

	switch itemType {
	case "task":
		UpdateExistingItem[ExistingTask](c, userID)
	case "goal":
		UpdateExistingItem[ExistingGoal](c, userID)
	case "dream":
		UpdateExistingItem[ExistingDream](c, userID)
	default:
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid item type")},
		)
		c.Abort()
		return
	}
}

func UpdateExistingItem[T ExistingItem](c *gin.Context, userID uint64) {

	var item T
	if err := c.BindJSON(&item); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid item object")},
		)
		c.Abort()
	}

	updatedItem, err := item.Update(userID)
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

func (item ExistingTask) Update(userID uint64) (Item, error) {
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

	return updatedItem, err
}

func (item ExistingGoal) Update(userID uint64) (Item, error) {
	var updatedItem Item

	err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Update item in the items table
		err := tx.Exec(`
			UPDATE items
			SET
				title = ?, 
				target_date = ?,
				priority = ?,
				parent_id = ?
			WHERE
				user_id = ? AND item_id = ?
		`, item.Title, item.TargetDate, item.Priority, item.ParentID, userID, item.ItemID).Error
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
			SELECT item_id, title, type, target_date, priority, parent_id, time_spent, created_at
			FROM items WHERE item_id = LAST_INSERT_ID()
		`).Scan(&updatedItem).Error
		if err != nil {
			return err
		}

		return nil
	})

	return updatedItem, err
}

func (item ExistingDream) Update(userID uint64) (Item, error) {
	var updatedItem Item

	err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Update item in the items table
		err := tx.Exec(`
			UPDATE items
			SET
				title = ?, 
				target_date = ?,
				priority = ?
			WHERE
				user_id = ? AND item_id = ?
		`, item.Title, item.TargetDate, item.Priority, userID, item.ItemID).Error
		if err != nil {
			return err
		}

		// Select the updated item
		err = tx.Raw(`
			SELECT item_id, title, type, target_date, priority, time_spent, created_at
			FROM items WHERE item_id = LAST_INSERT_ID()
		`).Scan(&updatedItem).Error
		if err != nil {
			return err
		}

		return nil
	})

	return updatedItem, err
}
