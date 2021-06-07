package entity

import (
	"html"
	"strings"
	"time"

	"github.com/9sarkan/notes/infrastructure/security"
)

type User struct {
	ID        uint64    `gorm:"primary_key;auto_increment;" json:"id"`
	FirstName string    `gorm:"size:100;not null;" json:"first_name"`
	LastName  string    `gorm:"size:100;not null;" json:"last_name"`
	Username  string    `gorm:"size:100;not null;unique;" json:"username"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type PublicUser struct {
	ID        uint64 `gorm:"primary_key;auto_increment;" json:"id"`
	FirstName string `gorm:"size:100;not null;" json:"first_name"`
	LastName  string `gorm:"size:100;not null;" json:"last_name"`
	Username  string `gorm:"size:100;not null;unique;" json:"username"`
}

func (u *User) BeforeSave() error {
	hash_password, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hash_password)
	return nil
}

type Users []User

func (user *User) PublicUser() interface{} {
	return &PublicUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}
}

func (users Users) PublicUsers() []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.PublicUser()
	}
	return result
}

func (user *User) Prepare() {
	user.FirstName = html.EscapeString(strings.TrimSpace(user.FirstName))
	user.LastName = html.EscapeString(strings.TrimSpace(user.LastName))
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}
func (user *User) Validate(action string) map[string]string {
	errorMessage := make(map[string]string)

	if user.Username == "" {
		errorMessage["username_required"] = "username required"
	}
	switch strings.ToLower(action) {
	case "login":
		if user.Password == "" {
			errorMessage["password_required"] = "password_required"
		}
	case "update":
	case "forgotPassword":
	default:

	}

	return errorMessage
}
