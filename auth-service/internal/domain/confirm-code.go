package domain

import "time"

type ConfirmCode struct {
	Id        int32
	UserId    int32
	Code      string
	Attempts  int
	CreatedAt time.Time
}

func NewConfirmCode(userId int32, code string) ConfirmCode {
	now := time.Now()
	return ConfirmCode{
		UserId:    userId,
		Code:      code,
		Attempts:  0,
		CreatedAt: now,
	}
}
