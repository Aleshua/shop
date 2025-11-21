package services

import "context"

type IEmailSenderService interface {
	Send(ctx context.Context, to string, body string) error
}
