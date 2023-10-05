package items

import (
	"net/http"

	"torch/torch-server/db"
	o "torch/torch-server/optional"
	r "torch/torch-server/recurring"

	"github.com/gin-gonic/gin"
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

func GetAllItems(c *gin.Context) {
	userID := c.GetUint64("userID")
	if userID == 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Could not find userID"},
		)
		c.Abort()
		return
	}

	items, err := getAllItemsByUser(userID)
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

func getAllItemsByUser(userID uint64) (items []Item, err error) {
	err = db.GetDB().Raw(`
		SELECT item_id, title, type, target_date, priority, duration, rec_times, rec_period, rec_progress, rec_updated_at, parent_id, time_spent, created_at
		FROM items WHERE user_id = ?
	`, userID).Scan(&items).Error
	return items, err
}
