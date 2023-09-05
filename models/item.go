package models

import (
	"torch/torch-server/db"
	o "torch/torch-server/optional"

	"gorm.io/gorm"
)

type Item struct {
	ItemID     uint64        `json:"itemID"`
	UserID     uint64        `json:"-"`
	Title      string       `json:"title"`
	Progress   float32      `json:"progress"`
	Type       string       `json:"type"`
	TargetDate o.NullString `json:"targetDate"`
	Priority   o.NullString `json:"priority"`
	Duration   o.NullInt    `json:"duration"`
	ParentID   o.NullUint64 `json:"parentID"`
	TimeSpent  int          `json:"timeSpent"`
	TimeLeft   int          `json:"timeLeft"`
	CreatedAt  string       `json:"createdAt"`
}

func GetAllItemsByUser(userID int64) (items []Item, err error) {
	err = db.GetDB().Raw(`
		SELECT item_id, title, progress, type, target_date, priority, duration, time_spent, time_left, created_at
		FROM items WHERE user_id = ?
	`, userID).Scan(&items).Error
	return items, err
}

func AddItem(item Item) (err error) {
	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		err = tx.Exec(`
			INSERT INTO items (user_id, title, type, target_date, priority, duration, parent, time_left) VALUES (?, ?, ?, NULL, ?, ?, ?, ?)
		`, item.UserID, item.Title, item.Type, item.Priority, item.Duration, item.ParentID, item.TimeLeft).Error
		if err != nil {
			return err
		}

		if item.ParentID.IsValid == true {
			err = tx.Exec(`
				INSERT INTO item_relations (user_id, parent_id, child_id) VALUES (?, ?, ?)
			`, item.UserID, item.ParentID, item.ItemID).Error
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
		err = tx.Exec(`
			DELETE FROM items WHERE user_id = ? AND item_id = ?
		`, userID, itemID).Error
		if err != nil {
			return err
		}

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

func UpdateItem(item Item) (err error) {
	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		err = tx.Exec(`
			UPDATE items
				title = ?, 
				progress = ?, 
				target_date = ?,
				priority = ?,
				duration = ?,
				parent_id = ?,
				time_spent = ?,
				time_left = ? 
			WHERE
				user_id = ? AND item_id = ?
		`, item.Title, item.Progress, item.TargetDate, item.Priority, item.Duration, item.ParentID, item.TimeSpent, item.TimeLeft, item.UserID, item.ItemID).Error
		if err != nil {
			return err
		}

		if item.ParentID.IsValid == false {
			err = tx.Exec(`
				DELETE FROM item_relations WHERE user_id = ? AND (parent_id = ? OR child_id = ?)
			`, item.UserID, item.ItemID, item.ItemID).Error
			if err != nil {
				return err
			}
		} else {
			err = tx.Exec(`
				UPDATE item_relations SET parent_id = ? WHERE user_id = ? AND child_id = ?
			`, item.ParentID, item.UserID, item.ItemID).Error
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
			progress = (time_spent + ?) / (time_left - ?),
			time_spent = time_spent + ?,
			time_left = time_left - ?
		WHERE
			user_id = ? AND item_id = ?
	`, timeSpent, timeSpent, timeSpent, timeSpent, userID, itemID).Error
	return err
}