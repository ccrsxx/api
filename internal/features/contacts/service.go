package contacts

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ccrsxx/api/internal/clients/gmail"
	"github.com/ccrsxx/api/internal/clients/pushover"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type querier interface {
	CreateContact(ctx context.Context, arg sqlc.CreateContactParams) (sqlc.Contact, error)
	UpdateContactDeliveredAtByID(ctx context.Context, id pgtype.UUID) (sqlc.Contact, error)
}

type cloudflareClient interface {
	VerifyTurnstile(ctx context.Context, token string, remoteIP string) error
}

type pushoverClient interface {
	SendMessage(ctx context.Context, messageRequest pushover.MessageRequest) error
}

type emailClient interface {
	Send(msg gmail.Message) error
}

type Service struct {
	db               querier
	emailClient      emailClient
	emailTarget      string
	emailAddress     string
	pushoverClient   pushoverClient
	cloudflareClient cloudflareClient
}

type ServiceConfig struct {
	Database         querier
	EmailClient      emailClient
	EmailTarget      string
	EmailAddress     string
	PushoverClient   pushoverClient
	CloudflareClient cloudflareClient
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		db:               cfg.Database,
		emailClient:      cfg.EmailClient,
		emailTarget:      cfg.EmailTarget,
		emailAddress:     cfg.EmailAddress,
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
	contact, err := s.db.CreateContact(ctx, sqlc.CreateContactParams{
		Name:    input.Name,
		Email:   input.Email,
		Message: input.Message,
	})

	if err != nil {
		return fmt.Errorf("create contact create record error: %w", err)
	}

	title := fmt.Sprintf("New contact from %s (%s)", contact.Name, contact.Email)
	description := input.Message

	err = s.pushoverClient.SendMessage(ctx, pushover.MessageRequest{
		Title:   title,
		Message: description,
	})

	if err != nil {
		return fmt.Errorf("create contact pushover notification error: %w", err)
	}

	go s.sendNewContactEmail(title, description)

	_, err = s.db.UpdateContactDeliveredAtByID(ctx, contact.ID)

	if err != nil {
		return fmt.Errorf("create contact update record error: %w", err)
	}

	return nil
}

func (s *Service) sendNewContactEmail(title string, description string) {
	err := s.emailClient.Send(gmail.Message{
		From:    s.emailAddress,
		To:      s.emailTarget,
		Subject: title,
		Text:    description,
	})

	// Ignore email sending error, just log it
	if err != nil {
		slog.Error("create contact email notification error", "error", err)
	}
}
