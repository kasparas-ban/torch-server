package items

import (
	"errors"
	"net/http"
	a "torch/torch-server/auth"
	m "torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

type UpdateItemProgressReq struct {
	PublicItemID string `json:"itemID"`
	TimeSpent    uint   `json:"timeSpent"`
}

func HandleUpdateItemProgress(c *gin.Context) {
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

	err := m.UpdateProgress(userID, reqBody.PublicItemID, reqBody.TimeSpent)
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
