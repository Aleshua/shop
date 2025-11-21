package domain

import "time"

type User struct {
	Id              int32
	Email           string
	PasswordHash    string
	EmailVerifiedAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewUser(email, passwordHash string) User {
	now := time.Now()
	return User{
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
