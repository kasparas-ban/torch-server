package models

import (
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

// === GET ===

func GetAllItemsByUser(userID uint64) ([]Item, error) {
	items := []Item{}
	err := db.GetDB().Raw(`
		SELECT public_item_id, title, type, target_date, priority, duration, rec_times, rec_period, rec_progress, rec_updated_at, parent_id, time_spent, created_at
		FROM items WHERE user_id = ?
	`, userID).Scan(&items).Error
	return items, err
}

func GetItemByUser(userID uint64, publicItemID string) (Item, error) {
	var item Item
	err := db.GetDB().Raw(`
		SELECT public_item_id, title, type, target_date, priority, duration, rec_times, rec_period, rec_progress, rec_updated_at, parent_id, time_spent, created_at
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

// === DELETE ===

func RemoveItem(userID uint64, publicItemID string) error {
	err := db.GetDB().Exec(`
		CALL DeleteItem(?, ?)
	`, userID, publicItemID).Error

	return err
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
