package items

import (
	"errors"
	"net/http"
	a "torch/torch-server/auth"
	"torch/torch-server/db"
	o "torch/torch-server/optional"
	r "torch/torch-server/recurring"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GeneralItem interface {
	Save(userID uint64) (addedItem Item, err error)
}

type CommonItem struct {
	Title      string       `json:"title"`
	Type       string       `json:"type"`
	TargetDate o.NullString `json:"targetDate"`
	Priority   o.NullString `json:"priority"`
}

type Task struct {
	CommonItem
	Duration  o.NullUint   `json:"duration"`
	Recurring r.Recurring  `gorm:"embedded" json:"recurring,omitempty"`
	ParentID  o.NullUint64 `json:"parentID"`
}

type Goal struct {
	CommonItem
	ParentID o.NullUint64 `json:"parentID"`
}

type Dream struct {
	CommonItem
}

func AddItem(c *gin.Context) {
	userID := a.GetUserID(c)
	itemType := c.Param("type")

	switch itemType {
	case "task":
		SaveItem[Task](c, userID)
	case "goal":
		SaveItem[Goal](c, userID)
	case "dream":
		SaveItem[Dream](c, userID)
	default:
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid item type")},
		)
		c.Abort()
		return
	}
}

func SaveItem[T GeneralItem](c *gin.Context, userID uint64) {
	var newItem T
	if err := c.BindJSON(&newItem); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid item object")},
		)
		c.Abort()
	}

	item, err := newItem.Save(userID)
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

func (item Task) Save(userID uint64) (addedItem Item, err error) {
	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Add item into the items table
		err = tx.Exec(`
			INSERT INTO items (user_id, title, type, target_date, priority, duration, rec_times, rec_period, parent_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, userID, item.Title, item.Type, item.TargetDate, item.Priority, item.Duration, item.Recurring.Times, item.Recurring.Period, item.ParentID).Error
		if err != nil {
			return err
		}

		if item.ParentID.Valid == true {
			// Add item to the relations table
			err = tx.Exec(`
				INSERT INTO item_relations (user_id, parent_id, child_id) VALUES (?, ?, LAST_INSERT_ID())
			`, userID, item.ParentID).Error
			if err != nil {
				return err
			}
		}

		// Select the newly added item
		err = tx.Raw(`
			SELECT item_id, title, type, target_date, priority, duration, rec_times, rec_period, rec_progress, parent_id, time_spent, created_at
			FROM items WHERE item_id = LAST_INSERT_ID()
		`).Scan(&addedItem).Error
		if err != nil {
			return err
		}

		return nil
	})

	return addedItem, err
}

func (item Goal) Save(userID uint64) (addedItem Item, err error) {
	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Add item into the items table
		err = tx.Exec(`
			INSERT INTO items (user_id, title, type, target_date, priority, parent_id) VALUES (?, ?, ?, ?, ?, ?)
		`, userID, item.Title, item.Type, item.TargetDate, item.Priority, item.ParentID).Error
		if err != nil {
			return err
		}

		if item.ParentID.Valid == true {
			// Add item to the relations table
			err = tx.Exec(`
				INSERT INTO item_relations (user_id, parent_id, child_id) VALUES (?, ?, LAST_INSERT_ID())
			`, userID, item.ParentID).Error
			if err != nil {
				return err
			}
		}

		// Select the newly added item
		err = tx.Raw(`
			SELECT item_id, title, type, target_date, priority, duration, rec_times, rec_period, rec_progress, parent_id, time_spent, created_at
			FROM items WHERE item_id = LAST_INSERT_ID()
		`).Scan(&addedItem).Error
		if err != nil {
			return err
		}

		return nil
	})

	return addedItem, err
}

func (item Dream) Save(userID uint64) (addedItem Item, err error) {
	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Add item into the items table
		err = tx.Exec(`
			INSERT INTO items (user_id, title, type, target_date, priority) VALUES (?, ?, ?, ?, ?)
		`, userID, item.Title, item.Type, item.TargetDate, item.Priority).Error
		if err != nil {
			return err
		}

		// Select the newly added item
		err = tx.Raw(`
			SELECT item_id, title, type, target_date, priority, duration, rec_times, rec_period, rec_progress, parent_id, time_spent, created_at
			FROM items WHERE item_id = LAST_INSERT_ID()
		`).Scan(&addedItem).Error
		if err != nil {
			return err
		}

		return nil
	})

	return addedItem, err
}
