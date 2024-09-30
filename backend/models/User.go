package models

import (
	"api/api/security"
	"errors"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"gorm.io/gorm"
)

type User struct {
	ID        uint32    `gorm:"primarykey;autoIncrement" json:"id"`
	Username  string    `gorm:"size:255;not null;unique" json:"username"`
	FirstName string    `gorm:"size:255;not null;unique" json:"first_name"`
	LastName  string    `gorm:"size:255;not null;unique" json:"last_name"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Avatar    string    `gorm:"size:255;null" json:"avatar"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (u *User) BeforeSave() error {
	hashedPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.FirstName = html.EscapeString(strings.TrimSpace(u.FirstName))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error
	switch strings.ToLower(action) {
	case "create":
		if u.FirstName == "" {
			err = errors.New("required firstname")
			errorMessages["Required_firstname"] = err.Error()
		}
		if u.LastName == "" {
			err = errors.New("required lastname")
			errorMessages["Required_lastname"] = err.Error()
		}
		if u.Username == "" {
			err = errors.New("required username")
			errorMessages["Required_username"] = err.Error()
		}
		if u.Password == "" {
			err = errors.New("required password")
			errorMessages["Required_password"] = err.Error()
		}
		if u.Password != "" && len(u.Password) < 6 {
			err = errors.New("password should be atleast 6 characters")
			errorMessages["Invalid_password"] = err.Error()
		}
		if u.Email == "" {
			err = errors.New("required email")
			errorMessages["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("invalid email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	case "login":
		if u.Password == "" {
			err = errors.New("password required")
			errorMessages["Required_password"] = err.Error()
		}
		if u.Email == "" {
			err = errors.New("required email")
			errorMessages["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("invalid email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	case "update":
		// if u.Email == "" {
		// 	err = errors.New("required email")
		// 	errorMessages["Required_email"] = err.Error()
		// }
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("invalid email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	case "forgetpassword":
		if u.Email == "" {
			err = errors.New("required email")
			errorMessages["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("invalid email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	default:
		break
	}
	return errorMessages
}

// CREATE USER AND SAVE IN THE DATABASE
func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	u.Prepare()
	tx := db.Begin()
	if err := tx.Create(&u).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error saving user: %w", err)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error comminting transactions: %w", err)
	}
	return u, nil
}

// FIND ALL USERS
func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	users := []User{}
	tx := db.Begin()
	if err := tx.Model(&User{}).Limit(100).Find(&users).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error finding users: %w", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error commiting transaction: %w", err)
	}
	return &users, nil
}

// FIND SINGLE USER
func (u *User) FindUserById(db *gorm.DB, uid uint32) (*User, error) {
	tx := db.Begin()
	if err := tx.Model(&User{}).Where("id = ?", uid).Take(&u).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error commiting transaction: %w", err)
	}
	return u, nil
}

func (u *User) UpdateUser(db *gorm.DB, uid uint32) (*User, error) {
	updates := map[string]interface{}{
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"email":      u.Email,
		"avatar":     u.Avatar,
	}
	tx := db.Begin()
	if err := tx.Model(&User{}).Where("id = ?", uid).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error updating the user: %w", err)
	}
	if err := tx.Model(&User{}).Where("id = ?", uid).Take(&u).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error retrieving updated user: %w", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transactions: %w", err)
	}
	return u, nil
}
