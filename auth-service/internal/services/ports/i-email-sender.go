package ports

type IEmailSenderService interface {
	Send(to string, body string) error
}
