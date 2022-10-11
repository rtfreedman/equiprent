package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password,omitempty" gorm:"-"`
	PasswordHash string `json:"passwordHash,omitempty"`
	PermLevel    int    `json:"permLevel"`
	Active       bool   `json:"active"`
}

func (u *User) Register() (err error) {
	return
}

func (u *User) Deactivate() (err error) {
	return
}