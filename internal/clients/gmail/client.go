// gmail.go
package gmail

import (
	"fmt"
	"net/smtp"
)

type Config struct {
	Username string
	Password string
}

type Client struct {
	auth     smtp.Auth
	addr     string
	sendMail func(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

const (
	defaultHost = "smtp.gmail.com"
	defaultPort = "587"
)

func NewClient(cfg Config) *Client {
	return &Client{
		auth:     smtp.PlainAuth("", cfg.Username, cfg.Password, defaultHost),
		addr:     fmt.Sprintf("%s:%s", defaultHost, defaultPort),
		sendMail: smtp.SendMail,
	}
}

type Message struct {
	From    string
	To      string
	Subject string
	Text    string
}

func (c *Client) Send(msg Message) error {
	rawMsg := []byte(
		"To: " + msg.To + "\r\n" +
			"From: " + msg.From + "\r\n" +
			"Subject: " + msg.Subject + "\r\n" +
			"\r\n" + msg.Text,
	)

	if err := c.sendMail(c.addr, c.auth, msg.From, []string{msg.To}, rawMsg); err != nil {
		return fmt.Errorf("send email error: %w", err)
	}

	return nil
}
