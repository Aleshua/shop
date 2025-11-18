package adapters

import (
	"fmt"
	"net/smtp"

	"auth/internal/config"
)

type ConfirmEmailCodeSenderService struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewConfirmEmailCodeSenderService(params config.Email) *ConfirmEmailCodeSenderService {
	return &ConfirmEmailCodeSenderService{
		Host:     params.Host,
		Port:     params.Port,
		Username: params.Username,
		Password: params.Password,
		From:     params.From,
	}
}

func (s *ConfirmEmailCodeSenderService) Send(to string, body string) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	msg := []byte(
		fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Код подтверждения\r\n\r\nВаш код подтверждения: %s",
			s.From,
			to,
			body,
		),
	)

	return smtp.SendMail(addr, auth, s.From, []string{to}, msg)
}
