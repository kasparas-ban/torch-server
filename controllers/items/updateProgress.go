package items

import (
	"errors"
	"net/http"
	a "torch/torch-server/auth"
	"torch/torch-server/db"

	"github.com/gin-gonic/gin"
)

type UpdateItemProgressReq struct {
	ItemID    uint64 `json:"itemID"`
	TimeSpent uint   `json:"timeSpent"`
}

func UpdateItemProgress(c *gin.Context) {
	userID := a.GetUserID(c)

	var reqBody UpdateItemProgressReq
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid request payload")},
		)
		c.Abort()
		return
	}

	err := updateProgress(userID, reqBody.ItemID, reqBody.TimeSpent)
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

func updateProgress(userID, itemID uint64, timeSpent uint) error {
	err := db.GetDB().Exec(`
		UPDATE items
		SET time_spent = time_spent + ?
		WHERE user_id = ? AND item_id = ?
	`, timeSpent, userID, itemID).Error
	return err
}
