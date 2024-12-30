package sms

import "context"

type MockSMSService struct {
}

func (m *MockSMSService) Send(ctx context.Context, templateId string, params []string, phones []string) error {
	return nil
}

func NewMockSMSService() *MockSMSService {
	return &MockSMSService{}
}
