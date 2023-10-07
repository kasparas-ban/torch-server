package users

import (
	"errors"
	"net/http"
	"torch/torch-server/db"
	o "torch/torch-server/optional"

	"github.com/gin-gonic/gin"
)

type NewUser struct {
	User
	CountryID o.NullUint8 `json:"countryID"`
}

func AddNewUser(c *gin.Context) {
	clerkID := c.GetString("clerkID")
	if clerkID != "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Could not find clerkID")},
		)
		c.Abort()
	}

	var userReq NewUser
	if err := c.BindJSON(&userReq); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Invalid user object")},
		)
		c.Abort()
	}

	user, err := addUser(clerkID, userReq)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}

func addUser(clerkID string, u NewUser) (ExistingUser, error) {
	var newUser ExistingUser
	err := db.GetDB().Raw(`
		CALL AddUser(?, ?, ?, ?, ?, ?, ?, ?)
	`, clerkID, u.Username, u.Email, u.Birthday, u.Gender, u.CountryID, u.City, u.Description).Scan(&newUser).Error

	return newUser, err
}
