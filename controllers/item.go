package controllers

import (
	"errors"
	"net/http"
	"strconv"
	m "torch/torch-server/models"
	o "torch/torch-server/optional"

	"github.com/gin-gonic/gin"
)

type NewItemRequest struct {
	Title      string
	Type       string
	TargetDate o.NullString
	Priority   o.NullString
	Duration   o.NullInt
	ParentID   o.NullUint64
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

	var newItem m.Item
	if err := c.BindJSON(&newItem); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid item object")},
		)
		c.Abort()
		return
	}

	newItem.UserID = userID

	err = m.AddItem(newItem)
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

func UpdateItem(c *gin.Context) {
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

	var newItem m.Item
	if err := c.BindJSON(&newItem); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid item object")},
		)
		c.Abort()
		return
	}

	newItem.UserID = userID

	err = m.UpdateItem(newItem)
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
