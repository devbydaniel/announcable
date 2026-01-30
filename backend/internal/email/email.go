package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"time"

	"github.com/devbydaniel/announcable/config"
	mail "github.com/wneessen/go-mail"
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

	m := mail.NewMsg()
	if err := m.From(cfg.Email.FromAddress); err != nil {
		return fmt.Errorf("error setting from address: %w", err)
	}
	if err := m.To(to); err != nil {
		return fmt.Errorf("error setting to address: %w", err)
	}
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextHTML, body.String())

	opts := []mail.Option{
		mail.WithPort(cfg.Email.SMTPPort),
		mail.WithTimeout(10 * time.Second),
	}

	if cfg.Email.SMTPUser != "" {
		opts = append(opts,
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(cfg.Email.SMTPUser),
			mail.WithPassword(cfg.Email.SMTPPass),
		)
	}

	if cfg.Email.SMTPTLS {
		opts = append(opts,
			mail.WithTLSPolicy(mail.TLSMandatory),
			mail.WithTLSConfig(&tls.Config{ServerName: cfg.Email.SMTPHost}),
		)
	} else {
		opts = append(opts,
			mail.WithTLSPolicy(mail.TLSOpportunistic),
			mail.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		)
	}

	c, err := mail.NewClient(cfg.Email.SMTPHost, opts...)
	if err != nil {
		return fmt.Errorf("error creating mail client: %w", err)
	}

	if err := c.DialAndSend(m); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	return nil
}
