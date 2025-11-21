package domain

import "time"

type Session struct {
	Id        int32
	UserId    int32
	TokenHash string
	ExpiresAt int64
	Revoked   bool
	CreatedAt time.Time
}

func NewSession(userId int32, tokenHash string, ttl int64) Session {
	now := time.Now()
	return Session{
		UserId:    userId,
		TokenHash: tokenHash,
		ExpiresAt: now.Unix() + ttl,
		Revoked:   false,
		CreatedAt: now,
	}
}
