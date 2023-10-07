package models

import (
	"torch/torch-server/db"
	o "torch/torch-server/optional"
	r "torch/torch-server/recurring"

	"gorm.io/gorm"
)

type Item struct {
	ItemID     uint64       `json:"itemID"`
	UserID     uint64       `json:"-"`
	Title      string       `json:"title"`
	Type       string       `json:"type"`
	TargetDate o.NullString `json:"targetDate"`
	Priority   o.NullString `json:"priority"`
	Duration   o.NullUint   `json:"duration"`
	Recurring  r.Recurring  `gorm:"embedded" json:"recurring,omitempty"`
	ParentID   o.NullUint64 `json:"parentID"`
	TimeSpent  uint         `json:"timeSpent"`
	CreatedAt  string       `json:"createdAt"`
}

type GeneralItem interface {
	AddItem(userID uint64) (addedItem Item, err error)
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

type ExistingItem interface {
	UpdateItem(userID uint64) (updatedItem Item, err error)
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

type UpdateItemProgressReq struct {
	ItemID    uint64 `json:"itemID"`
	TimeSpent uint   `json:"timeSpent"`
}

// === GET ===

func GetAllItemsByUser(userID uint64) ([]Item, error) {
	items := []Item{}
	err := db.GetDB().Raw(`
		SELECT item_id, title, type, target_date, priority, duration, rec_times, rec_period, rec_progress, rec_updated_at, parent_id, time_spent, created_at
		FROM items WHERE user_id = ?
	`, userID).Scan(&items).Error
	return items, err
}

func GetItemByUser(userID, itemID uint64) (item Item, err error) {
	err = db.GetDB().Raw(`
		SELECT item_id, title, type, target_date, priority, duration, rec_times, rec_period, rec_progress, rec_updated_at, parent_id, time_spent, created_at
		FROM items WHERE user_id = ? AND item_id = ?
	`, userID, itemID).Scan(&item).Error
	return item, err
}

// === CREATE ===

func (item Task) AddItem(userID uint64) (addedItem Item, err error) {
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

func (item Goal) AddItem(userID uint64) (addedItem Item, err error) {
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

func (item Dream) AddItem(userID uint64) (addedItem Item, err error) {
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

// === UPDATE ===

func (item ExistingTask) UpdateItem(userID uint64) (Item, error) {
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

func (item ExistingGoal) UpdateItem(userID uint64) (Item, error) {
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

func (item ExistingDream) UpdateItem(userID uint64) (Item, error) {
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

// === DELETE ===

func RemoveItem(userID, itemID uint64) error {
	err := db.GetDB().Exec(`
		CALL DeleteItem(?, ?)
	`, userID, itemID).Error

	return err
}

// === UPDATE PROGRESS ===

func UpdateProgress(userID, itemID uint64, timeSpent uint) error {
	err := db.GetDB().Exec(`
		UPDATE items
		SET time_spent = time_spent + ?
		WHERE user_id = ? AND item_id = ?
	`, timeSpent, userID, itemID).Error
	return err
}
