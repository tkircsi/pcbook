package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User represents the user's informations
type User struct {
	Username       string
	HashedPassword string
	Role           string
}

// NewUser creates a new user with hashed password and returns it
func NewUser(username string, password string, role string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %v", err)
	}

	return &User{
		Username:       username,
		HashedPassword: string(hashedPassword),
		Role:           role,
	}, nil
}

// IsCorrectPassword compares the given password with the user's password
func (u *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}

// Clone makes a copy of the user
func (u *User) Clone() *User {
	return &User{
		Username:       u.Username,
		HashedPassword: u.HashedPassword,
		Role:           u.Role,
	}
}
