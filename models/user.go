package models

import (
	"torch/torch-server/db"
	o "torch/torch-server/optional"
)

type User struct {
	UserID      uint64           `json:"userID"`
	ClerkID     string           `json:"-"`
	Username    string           `json:"username"`
	Email       string           `json:"email"`
	Birthday    o.NullString     `json:"birthday"`
	Gender      o.NullString     `json:"gender"`
	Country     o.NullString     `json:"country"`
	City        o.NullString     `json:"city"`
	Description o.NullString     `json:"description"`
	UpdatedAt   string           `json:"-"`
	CreatedAt   string           `json:"createdAt"`
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

func GetUserByUserID(userID string) (user User, err error) {
	db.GetDB().Raw(`
		SELECT u.user_id, u.clerk_id, u.username, u.email, u.birthday, u.gender, c.country, u.city, u.description, u.created_at 
		FROM users u
		LEFT JOIN countries c ON u.country_id = c.country_id
		WHERE u.user_id = ? LIMIT 1
	`, userID).Scan(&user)
	return user, err
}
