package models

import (
	"errors"
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
	City         o.NullString `json:"city,omitempty"`
	Description  o.NullString `json:"description,omitempty"`
	UpdatedAt    string       `json:"-"`
	CreatedAt    string       `json:"createdAt"`
}

type ExistingUser struct {
	User
	Country o.NullString `json:"country"`
}

type NewUser struct {
	PublicUserID string       `json:"userID" validate:"required"`
	Username     string       `json:"username" validate:"required,gt=5,lt=21"`
	Email        string       `json:"email" validate:"required,email"`
	Birthday     o.NullString `json:"birthday"`
	Gender       o.NullString `json:"gender"`
	City         o.NullString `json:"city"`
	Description  o.NullString `json:"description"`
	CountryCode  o.NullString `json:"countryCode" validate:"lt=3"`
}

type UpdateUserReq struct {
	PublicUserID string       `json:"userID" validate:"required"`
	Username     string       `json:"username" validate:"required,gt=5,lt=21"`
	Birthday     o.NullString `json:"birthday"`
	Gender       o.NullString `json:"gender"`
	City         o.NullString `json:"city"`
	Description  o.NullString `json:"description"`
	CountryCode  o.NullString `json:"countryCode" validate:"lt=3"`
}

func (u *NewUser) Validate() error {
	if err := Validate.Struct(u); err != nil {
		return err
	}

	if u.Gender.IsValid && (u.Gender.Val != "MALE" && u.Gender.Val != "FEMALE" && u.Gender.Val != "OTHER") {
		return errors.New("incorrect gender value")
	}

	if u.CountryCode.IsValid && (len(u.CountryCode.Val) > 2) {
		return errors.New("incorrect country code value")
	}

	return nil
}

func (u *UpdateUserReq) Validate() error {
	if err := Validate.Struct(u); err != nil {
		return err
	}

	if u.Gender.IsValid && (u.Gender.Val != "MALE" && u.Gender.Val != "FEMALE" && u.Gender.Val != "OTHER") {
		return errors.New("incorrect gender value")
	}

	if u.CountryCode.IsValid && (len(u.CountryCode.Val) > 2) {
		return errors.New("incorrect country code value")
	}

	return nil
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

	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Get country ID
		var countryId o.NullUint
		if u.CountryCode.IsValid && u.CountryCode.Val != "" {
			err := tx.Raw(`
				SELECT country_id FROM countries WHERE country_code = ?
			`, u.CountryCode.Val).Scan(&countryId).Error
			if err != nil {
				return err
			}
		}

		// Add user
		err = tx.Exec(`
			INSERT INTO users (public_user_id, clerk_id, username, email, birthday, gender, country_id, city, description) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, publicUserID, clerkID, u.Username, u.Email, u.Birthday, u.Gender, countryId, u.City, u.Description).Error
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

func UpdateUser(userID uint64, u UpdateUserReq) (ExistingUser, error) {
	var updatedUser ExistingUser

	err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Get country ID
		var countryId o.NullUint
		if u.CountryCode.IsValid && u.CountryCode.Val != "" {
			err := tx.Raw(`
				SELECT country_id FROM countries WHERE country_code = ?
			`, u.CountryCode.Val).Scan(&countryId).Error
			if err != nil {
				return err
			}
		}

		// Update user
		err := tx.Exec(`
			UPDATE users 
			SET username = ?, birthday = ?, gender = ?, country_id = ?, city = ?, description = ? 
			WHERE user_id = ?
		`, u.Username, u.Birthday, u.Gender, countryId, u.City, u.Description, userID).Error
		if err != nil {
			return err
		}

		// Select updated user
		err = tx.Raw(`
			SELECT u.public_user_id, u.clerk_id, u.username, u.email, u.birthday, u.gender, c.country, u.city, u.description, u.created_at 
			FROM users u
			LEFT JOIN countries c ON u.country_id = c.country_id
			WHERE u.user_id = ? LIMIT 1;
		`, userID).Scan(&updatedUser).Error
		if err != nil {
			return err
		}

		return nil
	})

	return updatedUser, err
}

func DeleteUser(userID uint64) error {
	err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Delete user
		err := tx.Exec(`
			DELETE FROM users WHERE user_id = ?
		`, userID).Error
		if err != nil {
			return err
		}

		// Delete all items
		err = tx.Exec(`
			DELETE FROM items WHERE user_id = ?
		`, userID).Error
		if err != nil {
			return err
		}

		// Delete timer history records
		err = tx.Exec(`
			DELETE FROM timer_history WHERE user_id = ?
		`, userID).Error
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
