package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/devbydaniel/release-notes-go/config"
	"gopkg.in/gomail.v2"
)

type Config struct {
	To      string
	Subject string
	Body    string
}

var cfg = config.New()

func Send(e *Config) error {
	if cfg.Env == "production" {
		return sendViaPostmark(e)
	}
	return sendToMailCatcher(e)
}

func sendToMailCatcher(e *Config) error {
	addr := cfg.Email.FromAddress
	server := cfg.Email.McServer
	port := cfg.Email.McPort

	m := gomail.NewMessage()
	m.SetHeader("From", addr)
	m.SetHeader("To", e.To)
	m.SetHeader("Subject", e.Subject)
	m.SetBody("text/html", e.Body)
	d := gomail.NewDialer(server, port, "", "")
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	return nil
}

func sendViaPostmark(e *Config) error {
	server := cfg.Email.PostmarkServerUrl
	token := cfg.Email.PostmarkToken
	fromAddr := cfg.Email.FromAddress

	payload := struct {
		From          string `json:"From"`
		To            string `json:"To"`
		Subject       string `json:"Subject"`
		HtmlBody      string `json:"HtmlBody"`
		MessageStream string `json:"MessageStream"`
	}{
		From:          fromAddr,
		To:            e.To,
		Subject:       e.Subject,
		HtmlBody:      e.Body,
		MessageStream: "outbound",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", server, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
