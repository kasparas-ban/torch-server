package controllers

import (
	"errors"
	"net/http"
	a "torch/torch-server/auth"
	m "torch/torch-server/models"
	o "torch/torch-server/optional"
	r "torch/torch-server/recurring"

	"github.com/gin-gonic/gin"
)

type NewItemReq struct {
	Title      string
	Type       string
	TargetDate o.NullString
	Priority   o.NullString
	Duration   o.NullUint
	Recurring  r.Recurring
	ParentID   o.NullUint64
}

type RemoveItemReq struct {
	itemID uint64
}

type UpdateItemProgressReq struct {
	itemID    uint64
	timeSpent int
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

	items, err := m.GetAllItemsByUser(userID)
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

func AddItem(c *gin.Context) {
	userID := a.GetUserID(c)

	var newItem m.AddItemReq
	if err := c.BindJSON(&newItem); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid item object")},
		)
		c.Abort()
		return
	}

	err := m.AddItem(newItem, userID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, newItem)
}

func RemoveItem(c *gin.Context) {
	userID := a.GetUserID(c)

	var reqBody RemoveItemReq
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid item object")},
		)
		c.Abort()
		return
	}

	err := m.RemoveItem(userID, reqBody.itemID)
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

func UpdateItem(c *gin.Context) {
	userID := a.GetUserID(c)

	var newItem m.UpdateItemReq
	if err := c.BindJSON(&newItem); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid item object")},
		)
		c.Abort()
		return
	}

	err := m.UpdateItem(newItem, userID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, newItem)
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

	err := m.UpdateItemProgress(userID, reqBody.itemID, reqBody.timeSpent)
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
