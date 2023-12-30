package models

import (
	"errors"
	"torch/torch-server/db"
	o "torch/torch-server/optional"
	r "torch/torch-server/recurring"
)

type Item struct {
	ItemID       uint64       `json:"-"`
	PublicItemID string       `json:"itemID"`
	UserID       uint64       `json:"-"`
	Title        string       `json:"title"`
	Type         string       `json:"type"`
	TargetDate   o.NullString `json:"targetDate"`
	Priority     o.NullString `json:"priority"`
	Duration     o.NullUint   `json:"duration"`
	Recurring    r.Recurring  `gorm:"embedded" json:"recurring,omitempty"`
	ParentID     o.NullString `json:"parentID"`
	TimeSpent    uint         `json:"timeSpent"`
	Status       string       `json:"status"`
	CreatedAt    string       `json:"createdAt"`
}

type GeneralItem interface {
	AddItem(userID uint64, publicItemID string) (Item, error)
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
	ParentID  o.NullString `json:"parentID"`
}

type Goal struct {
	CommonItem
	ParentID o.NullString `json:"parentID"`
}

type Dream struct {
	CommonItem
}

type ExistingItem interface {
	UpdateItem(userID uint64) (updatedItem Item, err error)
}

type ExistingTask struct {
	Task
	ItemID       uint64 `json:"-"`
	PublicItemID string `json:"itemID"`
}

type ExistingGoal struct {
	Goal
	ItemID       uint64 `json:"-"`
	PublicItemID string `json:"itemID"`
}

type ExistingDream struct {
	Dream
	ItemID       uint64 `json:"-"`
	PublicItemID string `json:"itemID"`
}

type UpdateItemProgressReq struct {
	ItemID    uint64 `json:"itemID"`
	TimeSpent uint   `json:"timeSpent"`
}

type UpdateItemStatusReq struct {
	PublicItemID     string `json:"itemID"`
	Status           string `json:"status"`
	UpdateAssociated bool   `json:"updateAssociated"`
	ItemType         string `json:"itemType"`
}

type DeleteItemStatusReq struct {
	PublicItemID     string `json:"itemID"`
	ItemType         string `json:"itemType"`
	DeleteAssociated bool   `json:"deleteAssociated"`
}

// === GET ===

func GetAllItemsByUser(userID uint64) ([]Item, error) {
	items := []Item{}
	err := db.GetDB().Raw(`
		SELECT public_item_id, title, type, target_date, priority, duration, rec_times, rec_period, rec_progress, rec_updated_at, parent_id, time_spent, status, created_at
		FROM items WHERE user_id = ?
	`, userID).Scan(&items).Error
	return items, err
}

func GetItemByUser(userID uint64, publicItemID string) (Item, error) {
	var item Item
	err := db.GetDB().Raw(`
		SELECT public_item_id, title, type, target_date, priority, duration, rec_times, rec_period, rec_progress, rec_updated_at, parent_id, time_spent, status, created_at
		FROM items WHERE user_id = ? AND public_item_id = ?
	`, userID, publicItemID).Scan(&item).Error
	return item, err
}

// === CREATE ===

func (item Task) AddItem(userID uint64, publicItemID string) (Item, error) {
	var addedItem Item
	err := db.GetDB().Raw(`
		CALL AddItem(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, userID, publicItemID, item.Title, "TASK", item.TargetDate, item.Priority, item.Duration, item.Recurring.Times, item.Recurring.Progress, item.Recurring.Period, item.ParentID).Scan(&addedItem).Error

	return addedItem, err
}

func (item Goal) AddItem(userID uint64, publicItemID string) (Item, error) {
	var addedItem Item
	err := db.GetDB().Raw(`
		CALL AddItem(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, userID, publicItemID, item.Title, "GOAL", item.TargetDate, item.Priority, nil, nil, nil, nil, item.ParentID).Scan(&addedItem).Error
	return addedItem, err
}

func (item Dream) AddItem(userID uint64, publicItemID string) (Item, error) {
	var addedItem Item
	err := db.GetDB().Raw(`
		CALL AddItem(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, userID, publicItemID, item.Title, "DREAM", item.TargetDate, item.Priority, nil, nil, nil, nil, nil).Scan(&addedItem).Error

	return addedItem, err
}

// === UPDATE ===

func (item ExistingTask) UpdateItem(userID uint64) (Item, error) {
	var updatedItem Item
	err := db.GetDB().Raw(`
		CALL UpdateItem(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, userID, item.PublicItemID, item.Title, item.TargetDate, item.Priority, item.Duration, item.Recurring.Times, item.Recurring.Progress, item.Recurring.Period, item.ParentID).Scan(&updatedItem).Error

	return updatedItem, err
}

func (item ExistingGoal) UpdateItem(userID uint64) (Item, error) {
	var updatedItem Item
	err := db.GetDB().Raw(`
		CALL UpdateItem(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, userID, item.PublicItemID, item.Title, item.TargetDate, item.Priority, nil, nil, nil, nil, item.ParentID).Scan(&updatedItem).Error

	return updatedItem, err
}

func (item ExistingDream) UpdateItem(userID uint64) (Item, error) {
	var updatedItem Item
	err := db.GetDB().Raw(`
		CALL UpdateItem(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, userID, item.PublicItemID, item.Title, item.TargetDate, item.Priority, nil, nil, nil, nil, nil).Scan(&updatedItem).Error

	return updatedItem, err
}

// === UPDATE STATUS ===

func UpdateItemStatus(userID uint64, body UpdateItemStatusReq) (Item, error) {
	var updatedItem Item

	if body.ItemType == "TASK" {
		err := db.GetDB().Raw(`
			UPDATE items
			SET status = ?
			WHERE user_id = ? AND public_item_id = ?
		`, body.Status, userID, body.PublicItemID).Scan(&updatedItem).Error
		return updatedItem, err
	}

	if body.ItemType == "GOAL" {
		if body.UpdateAssociated {
			err := db.GetDB().Raw(`
				UPDATE items
				SET status = ?
				WHERE user_id = ? AND (
					public_item_id = ? OR parent_id = ?
				)
			`, body.Status, userID, body.PublicItemID, body.PublicItemID).Scan(&updatedItem).Error
			return updatedItem, err
		} else {
			err := db.GetDB().Raw(`
			UPDATE items
			SET status = ?
			WHERE user_id = ? AND public_item_id = ?
		`, body.Status, userID, body.PublicItemID).Scan(&updatedItem).Error
			return updatedItem, err
		}
	}

	if body.ItemType == "DREAM" {
		if body.UpdateAssociated {
			err := db.GetDB().Raw(`
				UPDATE items
				SET status = ?
				WHERE user_id = ? AND (
					public_item_id = ? OR 
					parent_id = ? OR
					parent_id IN (SELECT derived.public_item_id FROM (SELECT public_item_id FROM items WHERE parent_id = ?) AS derived)
				)
			`, body.Status, userID, body.PublicItemID, body.PublicItemID, body.PublicItemID).Scan(&updatedItem).Error
			return updatedItem, err
		} else {
			err := db.GetDB().Raw(`
			UPDATE items
			SET status = ?
			WHERE user_id = ? AND public_item_id = ?
		`, body.Status, userID, body.PublicItemID).Scan(&updatedItem).Error
			return updatedItem, err
		}
	}

	return updatedItem, errors.New("failed to update item status")
}

// === DELETE ===

func RemoveItem(userID uint64, body DeleteItemStatusReq) error {
	if !body.DeleteAssociated || body.ItemType == "TASK" {
		err := db.GetDB().Exec(`
			CALL DeleteOneItem(?, ?)
		`, userID, body.PublicItemID).Error
		return err
	}

	if body.ItemType == "GOAL" {
		err := db.GetDB().Exec(`
			CALL DeleteGoalAll(?, ?)
		`, userID, body.PublicItemID).Error
		return err
	}

	if body.ItemType == "DREAM" {
		err := db.GetDB().Exec(`
			CALL DeleteDreamAll(?, ?)
		`, userID, body.PublicItemID).Error
		return err
	}

	return errors.New("failed to delete item")
}

// === UPDATE PROGRESS ===

func UpdateProgress(userID uint64, publicItemID string, timeSpent uint) error {
	err := db.GetDB().Exec(`
		UPDATE items
		SET time_spent = time_spent + ?
		WHERE user_id = ? AND public_item_id = ?
	`, timeSpent, userID, publicItemID).Error
	return err
}

func UpdateUserProgress(userID uint64, timeSpent uint) error {
	err := db.GetDB().Exec(`
		UPDATE users
		SET focus_time = focus_time + ?
		WHERE user_id = ?
	`, timeSpent, userID).Error
	return err
}
