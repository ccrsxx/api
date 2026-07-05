package contacts

import (
	"context"

	"github.com/ccrsxx/api/internal/clients/pushover"
)

type querier interface {
}

type cloudflareClient interface {
	VerifyTurnstile(ctx context.Context, token string, remoteIP string) error
}

type pushoverClient interface {
	SendMessage(ctx context.Context, messageRequest pushover.MessageRequest) error
}

type emailClient interface {
	SendEmail(ctx context.Context, to string, subject string, body string) error
}

type Service struct {
	emailClient      emailClient
	pushoverClient   pushoverClient
	cloudflareClient cloudflareClient
}

type ServiceConfig struct {
	EmailClient      emailClient
	PushoverClient   pushoverClient
	CloudflareClient cloudflareClient
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		emailClient:      cfg.EmailClient,
		pushoverClient:   cfg.PushoverClient,
		cloudflareClient: cfg.CloudflareClient,
	}
}

type CreateContactInput struct {
	Name    string `json:"name" validate:"required"`
	Token   string `json:"token" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Message string `json:"message" validate:"required"`
}

func (s *Service) CreateContact(ctx context.Context, input CreateContactInput, ipAddress string) error {
	return nil
}
