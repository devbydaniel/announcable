package user

import (
	"fmt"

	"github.com/devbydaniel/announcable/config"
	"github.com/devbydaniel/announcable/internal/email"
	"github.com/devbydaniel/announcable/internal/password"
	"github.com/google/uuid"
)

type service struct {
	repo repository
}

func NewService(r repository) *service {
	log.Trace().Msg("NewService")
	return &service{repo: r}
}

func (s *service) Create(email, pw string, emailVerified bool) (*User, error) {
	log.Trace().Str("email", email).Msg("Create")
	hashedPassword, err := password.HashPassword(pw)
	if err != nil {
		return nil, err
	}
	user := User{Email: email, Password: hashedPassword, EmailVerified: emailVerified}
	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *service) UpdatePassword(id uuid.UUID, pw string) error {
	log.Trace().Str("id", id.String()).Msg("UpdatePassword")
	hashedPassword, err := password.HashPassword(pw)
	if err != nil {
		return err
	}
	return s.repo.Update(id, &User{Password: hashedPassword})
}

func (s *service) GetByEmail(email string) (*User, error) {
	log.Trace().Str("email", email).Msg("GetByEmail")
	return s.repo.FindByEmail(email)
}

func (s *service) GetById(id uuid.UUID) (*User, error) {
	log.Trace().Str("id", id.String()).Msg("GetById")
	return s.repo.FindById(id)
}

func (s *service) Delete(id uuid.UUID) error {
	log.Trace().Str("id", id.String()).Msg("Delete")
	return s.repo.Delete(id)
}

func (s *service) SendVerifcationEmail(u *User, token string) error {
	log.Trace().Str("email", u.Email).Msg("SendVerifcationEmail")
	baseUrl := config.New().BaseURL
	protocol := "https"
	if config.New().Env == "development" {
		protocol = "http"
	}
	verifyUrl := fmt.Sprintf("%s://%s/verify-email?token=%s", protocol, baseUrl, token)
	config := email.EmailConfirmConfig{
		To:        u.Email,
		ActionURL: verifyUrl,
	}

	if err := email.SendEmailConfirm(&config); err != nil {
		log.Error().Err(err).Msg("Failed to send email")
		return err
	}
	return nil
}

func (s *service) VerifyEmail(id uuid.UUID) error {
	log.Trace().Msg("VerifyEmail")
	return s.repo.Update(id, &User{EmailVerified: true})
}

func (s *service) SendPwResetEmail(u *User, token string) error {
	log.Trace().Str("email", u.Email).Msg("SendPwResetEmail")
	baseUrl := config.New().BaseURL
	protocol := "https"
	if config.New().Env == "development" {
		protocol = "http"
	}
	url := fmt.Sprintf("%s://%s/reset-pw/%s", protocol, baseUrl, token)
	config := email.PasswordResetConfig{
		To:        u.Email,
		ActionURL: url,
	}

	if err := email.SendPasswordReset(&config); err != nil {
		log.Error().Err(err).Msg("Failed to send email")
		return err
	}
	return nil
}
