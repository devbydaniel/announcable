package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/devbydaniel/release-notes-go/config"
	"gopkg.in/gomail.v2"
)

type PasswordResetConfig struct {
	To        string
	ActionURL string
}

type EmailConfirmConfig struct {
	To        string
	ActionURL string
}

type UserInviteConfig struct {
	To               string
	OrganisationName string
	ActionURL        string
}

type templateEmail struct {
	From          string      `json:"From"`
	To            string      `json:"To"`
	TemplateAlias string      `json:"TemplateAlias"`
	TemplateModel interface{} `json:"TemplateModel"`
	MessageStream string      `json:"MessageStream"`
}

var cfg = config.New()

func SendPasswordReset(c *PasswordResetConfig) error {
	model := map[string]string{
		"action_url":      c.ActionURL,
		"product_url":     cfg.BaseURL,
		"product_name":    cfg.ProductInfo.ProductName,
		"support_email":   cfg.ProductInfo.SupportEmail,
		"company_name":    cfg.ProductInfo.CompanyName,
		"company_address": cfg.ProductInfo.CompanyAddress,
	}

	if cfg.Env == "production" {
		return sendTemplate("password-reset", c.To, model)
	}
	return sendToMailcatcher(c.To, "Password Reset", "Click here to reset your password: "+c.ActionURL)
}

func SendEmailConfirm(c *EmailConfirmConfig) error {
	model := map[string]string{
		"action_url":      c.ActionURL,
		"product_url":     cfg.BaseURL,
		"sender_name":     cfg.ProductInfo.PersonalName,
		"product_name":    cfg.ProductInfo.ProductName,
		"support_email":   cfg.ProductInfo.SupportEmail,
		"company_name":    cfg.ProductInfo.CompanyName,
		"company_address": cfg.ProductInfo.CompanyAddress,
	}

	if cfg.Env == "production" {
		return sendTemplate("welcome", c.To, model)
	}
	return sendToMailcatcher(c.To, "Welcome to "+cfg.ProductInfo.ProductName, "Please confirm your email address by clicking here: "+c.ActionURL)
}

func SendUserInvite(c *UserInviteConfig) error {
	model := map[string]string{
		"action_url":        c.ActionURL,
		"organisation_name": c.OrganisationName,
		"product_url":       cfg.BaseURL,
		"product_name":      cfg.ProductInfo.ProductName,
		"support_email":     cfg.ProductInfo.SupportEmail,
		"company_name":      cfg.ProductInfo.CompanyName,
		"company_address":   cfg.ProductInfo.CompanyAddress,
	}
	if cfg.Env == "production" {
		return sendTemplate("user-invitation", c.To, model)
	}
	return sendToMailcatcher(c.To, "Invitation to join", fmt.Sprintf("You have been invited to join %s. Click here to accept: %s", c.OrganisationName, c.ActionURL))
}

func sendTemplate(templateAlias string, to string, templateModel interface{}) error {
	email := templateEmail{
		From:          cfg.Email.FromAddress,
		To:            to,
		TemplateAlias: templateAlias,
		TemplateModel: templateModel,
		MessageStream: "outbound",
	}

	jsonData, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", cfg.Email.PostmarkServerUrl+"/withTemplate", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", cfg.Email.PostmarkToken)

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

func sendToMailcatcher(to, subject, text string) error {
	fromAddr := cfg.Email.FromAddress
	mcServer := cfg.Email.McServer
	mcPort := cfg.Email.McPort

	m := gomail.NewMessage()
	m.SetHeader("From", fromAddr)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", text)

	d := gomail.NewDialer(mcServer, mcPort, "", "")
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	return nil
}
