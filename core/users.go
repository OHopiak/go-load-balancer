package core

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type (
	User struct {
		ID           uint       `json:"id" gorm:"primary_key"`
		CreatedAt    time.Time  `json:"created_at"`
		UpdatedAt    time.Time  `json:"updated_at"`
		DeletedAt    *time.Time `json:"deleted_at,omitempty" sql:"index"`
		Username     string     `json:"username"`
		PasswordHash string     `json:"-"`
		ConcurrentTasksLimit uint `json:"concurrent_tasks_limit"`
	}
)

func GetPasswordHash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	return string(hash), err
}
