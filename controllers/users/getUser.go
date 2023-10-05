package users

import (
	"errors"
	"net/http"
	"torch/torch-server/db"
	o "torch/torch-server/optional"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserID      uint64       `json:"userID"`
	ClerkID     string       `json:"-"`
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	Birthday    o.NullString `json:"birthday"`
	Gender      o.NullString `json:"gender"`
	Country     o.NullString `json:"country"`
	City        o.NullString `json:"city"`
	Description o.NullString `json:"description"`
	UpdatedAt   string       `json:"-"`
	CreatedAt   string       `json:"createdAt"`
}

func GetUserInfoByClerkID(c *gin.Context) {
	clerkID := c.GetString("clerkID")
	if clerkID != "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Could not find clerkID")},
		)
		c.Abort()
	}

	user, err := GetUserByClerkID(clerkID)
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

func GetUserInfo(c *gin.Context) {
	userID := c.GetUint64("userID")
	if userID == 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Could not find userID"},
		)
		c.Abort()
	}

	user, err := getUserInfo(userID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": errors.New("Failed to get user info")},
		)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}

func getUserInfo(userID uint64) (user User, err error) {
	db.GetDB().Raw(`
		SELECT u.user_id, u.clerk_id, u.username, u.email, u.birthday, u.gender, c.country, u.city, u.description, u.created_at 
		FROM users u
		LEFT JOIN countries c ON u.country_id = c.country_id
		WHERE u.user_id = ? LIMIT 1
	`, userID).Scan(&user)
	return user, err
}

func GetUserByClerkID(clerkID string) (user User, err error) {
	db.GetDB().Raw(`
		SELECT u.user_id, u.clerk_id, u.username, u.email, u.birthday, u.gender, c.country, u.city, u.description, u.created_at 
		FROM users u
		LEFT JOIN countries c ON u.country_id = c.country_id
		WHERE u.clerk_id = ? LIMIT 1
	`, clerkID).Scan(&user)
	return user, err
}
