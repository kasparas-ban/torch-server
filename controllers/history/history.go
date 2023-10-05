package history

import (
	"errors"
	"net/http"
	a "torch/torch-server/auth"
	"torch/torch-server/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TimerHistory struct {
	UserID    uint64 `json:"-"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	ItemID    uint64 `json:"itemID"`
}

type UpsertReq struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	ItemID    uint64 `json:"itemID"`
}

func GetTimerHistory(c *gin.Context) {
	userID := c.GetUint64("userID")
	if userID == 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Could not find userID"},
		)
		c.Abort()
		return
	}

	items, err := getTimerHistory(userID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, items)
}

func UpsertTimerHistory(c *gin.Context) {
	userID := a.GetUserID(c)

	var timerData UpsertReq
	if err := c.BindJSON(&timerData); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid timer object")},
		)
		c.Abort()
		return
	}

	err := upsertTimerHistory(timerData, userID)
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

func getTimerHistory(userID uint64) (data []TimerHistory, err error) {
	err = db.GetDB().Raw(`
		SELECT user_id, start_time, end_time, item_id
		FROM timer_history WHERE user_id = ?
	`, userID).Scan(&data).Error
	return data, err
}

func upsertTimerHistory(record UpsertReq, userID uint64) (err error) {
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
