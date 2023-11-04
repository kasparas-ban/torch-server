package models

import (
	"torch/torch-server/db"
	o "torch/torch-server/optional"
	"torch/torch-server/util"

	"gorm.io/gorm"
)

type User struct {
	UserID       uint64       `json:"-"`
	PublicUserID string       `json:"userID"`
	ClerkID      string       `json:"-"`
	Username     string       `json:"username"`
	Email        string       `json:"email"`
	Birthday     o.NullString `json:"birthday"`
	Gender       o.NullString `json:"gender"`
	City         o.NullString `json:"city"`
	Description  o.NullString `json:"description"`
	UpdatedAt    string       `json:"-"`
	CreatedAt    string       `json:"createdAt"`
}

type ExistingUser struct {
	User
	Country o.NullString `json:"country"`
}

type NewUser struct {
	PublicUserID string       `json:"userID" validate:"required"`
	Username     string       `json:"username" validate:"required"`
	Email        string       `json:"email" validate:"required,email"`
	Birthday     o.NullString `json:"birthday"`
	Gender       o.NullString `json:"gender"`
	City         o.NullString `json:"city"`
	Description  o.NullString `json:"description"`

	CountryCode o.NullString `json:"countryCode"`
}

func GetUserInfo(userID uint64) (ExistingUser, error) {
	var user ExistingUser
	err := db.GetDB().Raw(`
		SELECT u.public_user_id, u.clerk_id, u.username, u.email, u.birthday, u.gender, c.country, u.city, u.description, u.created_at 
		FROM users u
		LEFT JOIN countries c ON u.country_id = c.country_id
		WHERE u.user_id = ? LIMIT 1
	`, userID).Scan(&user).Error
	return user, err
}

func GetUserByClerkID(clerkID string) (ExistingUser, error) {
	var user ExistingUser
	err := db.GetDB().Raw(`
		SELECT u.public_user_id, u.clerk_id, u.username, u.email, u.birthday, u.gender, c.country, u.city, u.description, u.created_at 
		FROM users u
		LEFT JOIN countries c ON u.country_id = c.country_id
		WHERE u.clerk_id = ? LIMIT 1
	`, clerkID).Scan(&user).Error
	return user, err
}

func AddUser(clerkID string, u NewUser) (ExistingUser, error) {
	var newUser ExistingUser
	publicUserID, err := util.New()
	if err != nil {
		return newUser, err
	}

	db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Get country ID
		var countryId o.NullUint
		if u.CountryCode.IsValid && u.CountryCode.Val != "" {
			err := tx.Raw(`
				SELECT country_id FROM countries WHERE country_code = ?
			`, u.CountryCode.Val).Scan(countryId).Error
			if err != nil {
				return err
			}
		}

		// Add user
		err = tx.Raw(`
			INSERT INTO users (public_user_id, clerk_id, username, email, birthday, gender, country_id, city, description) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, publicUserID, clerkID, u.Username, u.Email, u.Birthday, u.Gender, countryId, u.City, u.Description, u.UserID).Error
		if err != nil {
			return err
		}

		// Select new user
		err = tx.Raw(`
			SELECT u.public_user_id, u.clerk_id, u.username, u.email, u.birthday, u.gender, c.country, u.city, u.description, u.created_at 
			FROM users u
			LEFT JOIN countries c ON u.country_id = c.country_id
			WHERE u.user_id = LAST_INSERT_ID() LIMIT 1;
		`).Scan(&newUser).Error
		if err != nil {
			return err
		}

		return nil
	})

	return newUser, err
}

func UpdateUser(user NewUser) (ExistingUser, error) {
	var updatedUser ExistingUser
	err := db.GetDB().Raw(`
		CALL UpdateUser(?, ?, ?, ?, ?, ?, ?, ?)
	`, user.UserID, user.Username, user.Email, user.Birthday, user.Gender, user.CountryCode, user.City, user.Description).Scan(&updatedUser).Error

	return updatedUser, err
}

func DeleteUser(userID uint64) error {
	err := db.GetDB().Exec(`
		CALL DeleteUser(?)
	`, userID).Error

	return err
}
