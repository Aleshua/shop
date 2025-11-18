package mocks

type MockEmailSenderService struct {
	Sent []string
}

func NewMockEmailSender() *MockEmailSenderService {
	return &MockEmailSenderService{}
}

func (m *MockEmailSenderService) Send(to string, body string) error {
	m.Sent = append(m.Sent, to)
	return nil
}
