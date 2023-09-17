package models

import (
	"torch/torch-server/db"
	o "torch/torch-server/optional"

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
	Recurring  Recurring    `gorm:"embedded" json:"recurring"`
	ParentID   o.NullUint64 `json:"parentID"`
	TimeSpent  uint         `json:"timeSpent"`
	CreatedAt  string       `json:"createdAt"`
}

type Recurring struct {
	Times    uint   `gorm:"column:rec_times" json:"times"`
	Period   string `gorm:"column:rec_period" json:"period"`
	Progress uint   `gorm:"column:rec_progress" json:"progress"`
}

type AddItemReq struct {
	Title      string       `json:"title"`
	Type       string       `json:"type"`
	Recurring  Recurring    `gorm:"embedded" json:"recurring,omitempty"`
	TargetDate o.NullString `json:"targetDate"`
	Priority   o.NullString `json:"priority"`
	Duration   o.NullUint   `json:"duration"`
	ParentID   o.NullUint64 `json:"parentID"`
}

type UpdateItemReq struct {
	ItemID     uint64       `json:"itemID"`
	Title      string       `json:"title"`
	Type       string       `json:"type"`
	TargetDate o.NullString `json:"targetDate"`
	Priority   o.NullString `json:"priority"`
	Duration   o.NullUint   `json:"duration"`
	ParentID   o.NullUint64 `json:"parentID"`
	TimeSpent  uint         `json:"timeSpent"`
}

func GetAllItemsByUser(userID uint64) (items []Item, err error) {
	err = db.GetDB().Raw(`
		SELECT item_id, title, type, target_date, priority, duration, rec_times, rec_period, rec_progress, parent_id, time_spent, created_at
		FROM items WHERE user_id = ?
	`, userID).Scan(&items).Error
	return items, err
}

func AddItem(item AddItemReq, userID uint64) (err error) {
	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Add item into the items table
		var rec_period *string
		if item.Recurring.Period == "" {
			rec_period = nil
		} else {
			// fmt.Println("attempt", item.Recurring.Period)
			rec_period = new(string)
			*rec_period = item.Recurring.Period
		}
		err = tx.Exec(`
			INSERT INTO items (user_id, title, type, target_date, priority, duration, rec_times, rec_period, parent_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, userID, item.Title, item.Type, item.TargetDate, item.Priority, item.Duration, item.Recurring.Times, rec_period, item.ParentID).Error
		if err != nil {
			return err
		}

		if item.ParentID.IsValid == true {
			// Add item to the relations table
			err = tx.Exec(`
				INSERT INTO item_relations (user_id, parent_id, child_id) VALUES (?, ?, LAST_INSERT_ID())
			`, userID, item.ParentID).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func RemoveItem(userID, itemID uint64) (err error) {
	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Remove item from the items table
		err = tx.Exec(`
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

func UpdateItem(item UpdateItemReq, userID uint64) (err error) {
	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Update item in the items table
		err = tx.Exec(`
			UPDATE items
				title = ?, 
				target_date = ?,
				priority = ?,
				duration = ?,
				parent_id = ?,
				time_spent = ?
			WHERE
				user_id = ? AND item_id = ?
		`, item.Title, item.TargetDate, item.Priority, item.Duration, item.ParentID, item.TimeSpent, userID, item.ItemID).Error
		if err != nil {
			return err
		}

		if item.ParentID.IsValid == false {
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

		return nil
	})

	return err
}

func UpdateItemProgress(userID, itemID uint64, timeSpent int) (err error) {
	err = db.GetDB().Raw(`
		UPDATE items
		SET
			time_spent = time_spent + ?
		WHERE
			user_id = ? AND item_id = ?
	`, timeSpent, timeSpent, userID, itemID).Error
	return err
}
