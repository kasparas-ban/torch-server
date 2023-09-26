package models

import (
	"torch/torch-server/db"

	"gorm.io/gorm"
)

type TimerHistory struct {
	UserID    uint64 `json:"-"`
	StartTime string `json:"startDate"`
	EndTime   string `json:"endDate"`
	ItemID    uint64 `json:"itemID"`
}

type UpsertReq struct {
	StartTime string `json:"startDate"`
	EndTime   string `json:"endDate"`
	ItemID    uint64 `json:"itemID"`
}

func GetTimerHistory(userID uint64) (data []TimerHistory, err error) {
	err = db.GetDB().Raw(`
		SELECT user_id, start_time, end_time, item_id
		FROM timer_history WHERE user_id = ?
	`, userID).Scan(&data).Error
	return data, err
}

func UpsertTimerHistory(record UpsertReq, userID uint64) (err error) {
	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Add record into the Timer History table
		err = tx.Exec(`
			INSERT INTO timer_history (user_id, start_time, end_time, item_id) VALUES (?, ?, ?, ?)
		`, userID, record.StartTime, record.EndTime, record.ItemID).Error
		if err != nil {
			return err
		}

		// Remove the oldest record
		err = tx.Exec(`
			DELETE FROM timer_history
    	WHERE user_id = ?
    	ORDER BY end_time
    	LIMIT 1 OFFSET 5;
			`, userID).Error
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
