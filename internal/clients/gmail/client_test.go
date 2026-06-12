package gmail

import (
	"errors"
	"net/smtp"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient(Config{})

	if client == nil {
		t.Fatal("want client to be initialized, got nil")
	}
}

func TestClient_Send(t *testing.T) {
	newTestClient := func(fn func(string, smtp.Auth, string, []string, []byte) error) *Client {
		return &Client{
			auth:     smtp.PlainAuth("", "user", "pass", "smtp.gmail.com"),
			addr:     "smtp.gmail.com:587",
			sendMail: fn,
		}
	}

	t.Run("Success", func(t *testing.T) {
		var capturedMsg []byte

		client := newTestClient(func(_ string, _ smtp.Auth, _ string, _ []string, msg []byte) error {
			capturedMsg = msg
			return nil
		})

		err := client.Send(Message{
			From:    "sender@example.com",
			To:      "receiver@example.com",
			Subject: "Hello",
			Text:    "World",
		})

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		for _, want := range []string{"To: receiver@example.com", "Subject: Hello", "World"} {
			if !strings.Contains(string(capturedMsg), want) {
				t.Errorf("message missing %q\ngot: %s", want, capturedMsg)
			}
		}
	})

	t.Run("SMTP Error", func(t *testing.T) {
		client := newTestClient(func(_ string, _ smtp.Auth, _ string, _ []string, _ []byte) error {
			return errors.New("connection refused")
		})

		err := client.Send(Message{From: "a@b.com", To: "c@d.com"})

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "send email error") {
			t.Errorf("error not wrapped correctly: %v", err)
		}
	})
}
