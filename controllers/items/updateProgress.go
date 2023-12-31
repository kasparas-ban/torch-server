package items

import (
	"net/http"
	a "torch/torch-server/auth"
	m "torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

type UpdateItemProgressReq struct {
	PublicItemID string `json:"itemID"`
	TimeSpent    uint   `json:"timeSpent"`
}

type UpdateUserProgressReq struct {
	TimeSpent uint `json:"timeSpent"`
}

func HandleUpdateItemProgress(c *gin.Context) {
	userID, err := a.GetUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	var reqBody UpdateItemProgressReq
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid request payload"},
		)
		c.Abort()
		return
	}

	err = m.UpdateProgress(userID, reqBody.PublicItemID, reqBody.TimeSpent)
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

func HandleUpdateUserProgress(c *gin.Context) {
	userID, err := a.GetUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	var reqBody UpdateUserProgressReq
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid request payload"},
		)
		c.Abort()
		return
	}

	err = m.UpdateUserProgress(userID, reqBody.TimeSpent)
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
