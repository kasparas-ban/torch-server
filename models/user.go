package models

import (
	"torch/torch-server/db"
	o "torch/torch-server/optional"
)

type User struct {
	UserID      uint64       `json:"userID"`
	ClerkID     string       `json:"-"`
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	Birthday    o.NullString `json:"birthday"`
	Gender      o.NullString `json:"gender"`
	City        o.NullString `json:"city"`
	Description o.NullString `json:"description"`
	UpdatedAt   string       `json:"-"`
	CreatedAt   string       `json:"createdAt"`
}

type ExistingUser struct {
	User
	Country o.NullString `json:"country"`
}

type NewUser struct {
	User
	CountryID o.NullUint8 `json:"countryID"`
}

func GetUserInfo(userID uint64) (ExistingUser, error) {
	var user ExistingUser
	err := db.GetDB().Raw(`
		SELECT u.user_id, u.clerk_id, u.username, u.email, u.birthday, u.gender, c.country, u.city, u.description, u.created_at 
		FROM users u
		LEFT JOIN countries c ON u.country_id = c.country_id
		WHERE u.user_id = ? LIMIT 1
	`, userID).Scan(&user).Error
	return user, err
}

func GetUserByClerkID(clerkID string) (ExistingUser, error) {
	var user ExistingUser
	err := db.GetDB().Raw(`
		SELECT u.user_id, u.clerk_id, u.username, u.email, u.birthday, u.gender, c.country, u.city, u.description, u.created_at 
		FROM users u
		LEFT JOIN countries c ON u.country_id = c.country_id
		WHERE u.clerk_id = ? LIMIT 1
	`, clerkID).Scan(&user).Error
	return user, err
}

func AddUser(clerkID string, u NewUser) (ExistingUser, error) {
	var newUser ExistingUser
	err := db.GetDB().Raw(`
		CALL AddUser(?, ?, ?, ?, ?, ?, ?, ?)
	`, clerkID, u.Username, u.Email, u.Birthday, u.Gender, u.CountryID, u.City, u.Description).Scan(&newUser).Error

	return newUser, err
}

func UpdateUser(user NewUser) (ExistingUser, error) {
	var updatedUser ExistingUser
	err := db.GetDB().Raw(`
		CALL UpdateUser(?, ?, ?, ?, ?, ?, ?, ?)
	`, user.UserID, user.Username, user.Email, user.Birthday, user.Gender, user.CountryID, user.City, user.Description).Scan(&updatedUser).Error

	return updatedUser, err
}

func DeleteUser(userID uint64) error {
	err := db.GetDB().Exec(`
		CALL DeleteUser(?)
	`, userID).Error

	return err
}
