package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"

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

var cfg = config.New()

func SendPasswordReset(c *PasswordResetConfig) error {
	data := map[string]string{
		"action_url":      c.ActionURL,
		"product_url":     cfg.BaseURL,
		"product_name":    cfg.ProductInfo.ProductName,
		"support_email":   cfg.ProductInfo.SupportEmail,
		"company_name":    cfg.ProductInfo.CompanyName,
		"company_address": cfg.ProductInfo.CompanyAddress,
	}
	return sendEmail(c.To, "Reset Your Password", passwordResetTmpl, data)
}

func SendEmailConfirm(c *EmailConfirmConfig) error {
	data := map[string]string{
		"action_url":      c.ActionURL,
		"product_url":     cfg.BaseURL,
		"sender_name":     cfg.ProductInfo.PersonalName,
		"product_name":    cfg.ProductInfo.ProductName,
		"support_email":   cfg.ProductInfo.SupportEmail,
		"company_name":    cfg.ProductInfo.CompanyName,
		"company_address": cfg.ProductInfo.CompanyAddress,
	}
	return sendEmail(c.To, "Welcome to "+cfg.ProductInfo.ProductName, welcomeTmpl, data)
}

func SendUserInvite(c *UserInviteConfig) error {
	data := map[string]string{
		"action_url":        c.ActionURL,
		"organisation_name": c.OrganisationName,
		"product_url":       cfg.BaseURL,
		"product_name":      cfg.ProductInfo.ProductName,
		"support_email":     cfg.ProductInfo.SupportEmail,
		"company_name":      cfg.ProductInfo.CompanyName,
		"company_address":   cfg.ProductInfo.CompanyAddress,
	}
	return sendEmail(c.To, "You're Invited to "+c.OrganisationName, userInviteTmpl, data)
}

func sendEmail(to, subject string, tmpl *template.Template, data map[string]string) error {
	var body bytes.Buffer
	if err := tmpl.ExecuteTemplate(&body, "base", data); err != nil {
		return fmt.Errorf("error rendering template: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.Email.FromAddress)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(
		cfg.Email.SMTPHost,
		cfg.Email.SMTPPort,
		cfg.Email.SMTPUser,
		cfg.Email.SMTPPass,
	)

	d.SSL = false
	if cfg.Email.SMTPTLS {
		// Production: enable STARTTLS
		d.TLSConfig = &tls.Config{ServerName: cfg.Email.SMTPHost}
	} else {
		// Local dev (Mailcatcher): skip TLS entirely
		d.TLSConfig = nil
	}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	return nil
}
