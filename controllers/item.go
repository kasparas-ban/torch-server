package controllers

import (
	"errors"
	"net/http"
	"strconv"
	m "torch/torch-server/models"
	o "torch/torch-server/optional"

	"github.com/gin-gonic/gin"
)

type NewItemReq struct {
	Title      string
	Type       string
	TargetDate o.NullString
	Priority   o.NullString
	Duration   o.NullInt
	ParentID   o.NullUint64
}

type RemoveItemReq struct {
	userID uint64
	itemID uint64
}

type UpdateItemProgressReq struct {
	userID uint64
	itemID uint64
	timeSpent int
}

func GetAllItems(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Could not find userID parameter")},
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
	var newItem m.Item
	if err := c.BindJSON(&newItem); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid item object")},
		)
		c.Abort()
		return
	}

	err := m.AddItem(newItem)
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
	var reqBody RemoveItemReq
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid item object")},
		)
		c.Abort()
		return
	}
	
	err := m.RemoveItem(reqBody.userID, reqBody.itemID)
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
	var newItem m.Item
	if err := c.BindJSON(&newItem); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid item object")},
		)
		c.Abort()
		return
	}

	err := m.UpdateItem(newItem)
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
	var reqBody UpdateItemProgressReq
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid request payload")},
		)
		c.Abort()
		return
	}

	err := m.UpdateItemProgress(reqBody.userID, reqBody.itemID, reqBody.timeSpent)
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