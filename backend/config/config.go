package config

import (
	"os"
	"strconv"
)

type productInfo struct {
	ProductName    string
	CompanyName    string
	CompanyAddress string
	SupportEmail   string
}

type postgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type objStorageConfig struct {
	AccessKey string
	SecretKey string
	Endpoint  string
	Region    string
	UseSSL    bool
}

type pgAdminConfig struct {
	Email    string
	Password string
	Port     int
}

type emailConfig struct {
	FromAddress string
	SMTPHost    string
	SMTPPort    int
	SMTPUser    string
	SMTPPass    string
	SMTPTLS     bool
}

type axiomConfig struct {
	Dataset string
	Token   string
}

type rateLimitConfig struct {
	PublicRequestsPerInterval          int
	PublicMaxTokens                    int
	PublicRefillIntervalSeconds        int
	AuthenticatedRequestsPerInterval   int
	AuthenticatedMaxTokens             int
	AuthenticatedRefillIntervalSeconds int
}

type config struct {
	Env         string
	BaseURL     string
	Port        int
	AdminUserId string
	Postgres    postgresConfig
	ObjStorage  objStorageConfig
	PgAdmin     pgAdminConfig
	Email       emailConfig
	ProductInfo productInfo
	Axiom       axiomConfig
	RateLimit   rateLimitConfig
}

func New() *config {
	cfg := &config{
		Env:         getEnv("ENV"),
		BaseURL:     getEnv("BASE_URL"),
		Port:        getEnvAsInt("PORT"),
		AdminUserId: getEnv("ADMIN_USER_ID"),
		Postgres: postgresConfig{
			Host:     getEnv("POSTGRES_HOST"),
			Port:     getEnvAsInt("POSTGRES_PORT"),
			User:     getEnv("POSTGRES_USER"),
			Password: getEnv("POSTGRES_PASSWORD"),
			Name:     getEnv("POSTGRES_NAME"),
		},
		ObjStorage: objStorageConfig{
			AccessKey: getEnv("MINIO_ACCESS_KEY"),
			SecretKey: getEnv("MINIO_SECRET_KEY"),
			Endpoint:  getEnv("MINIO_ENDPOINT"),
			Region:    getEnv("MINIO_REGION"),
			UseSSL:    getEnvAsBool("MINIO_USE_SSL"),
		},
		PgAdmin: pgAdminConfig{
			Email:    getEnv("PGADMIN_DEFAULT_EMAIL"),
			Password: getEnv("PGADMIN_DEFAULT_PASSWORD"),
			Port:     getEnvAsInt("PGADMIN_PORT"),
		},
		Email: emailConfig{
			FromAddress: getEnvWithDefault("EMAIL_FROM_ADDRESS", ""),
			SMTPHost:    getEnvWithDefault("SMTP_HOST", ""),
			SMTPPort:    getEnvAsIntWithDefault("SMTP_PORT", 0),
			SMTPUser:    getEnvWithDefault("SMTP_USER", ""),
			SMTPPass:    getEnvWithDefault("SMTP_PASS", ""),
			SMTPTLS:     getEnvAsBoolWithDefault("SMTP_TLS", false),
		},
		ProductInfo: productInfo{
			ProductName:    "Announcable",
			CompanyName:    getEnvWithDefault("COMPANY_NAME", ""),
			CompanyAddress: getEnvWithDefault("COMPANY_ADDRESS", ""),
			SupportEmail:   getEnvWithDefault("SUPPORT_EMAIL", ""),
		},
		Axiom: axiomConfig{
			Dataset: getEnvWithDefault("AXIOM_DATASET", ""),
			Token:   getEnvWithDefault("AXIOM_TOKEN", ""),
		},
		RateLimit: rateLimitConfig{
			PublicRequestsPerInterval:          getEnvAsIntWithDefault("RATE_LIMIT_PUBLIC_REQUESTS", 100),
			PublicMaxTokens:                    getEnvAsIntWithDefault("RATE_LIMIT_PUBLIC_MAX_TOKENS", 100),
			PublicRefillIntervalSeconds:        getEnvAsIntWithDefault("RATE_LIMIT_PUBLIC_REFILL_SECONDS", 60),
			AuthenticatedRequestsPerInterval:   getEnvAsIntWithDefault("RATE_LIMIT_AUTHENTICATED_REQUESTS", 1000),
			AuthenticatedMaxTokens:             getEnvAsIntWithDefault("RATE_LIMIT_AUTHENTICATED_MAX_TOKENS", 1000),
			AuthenticatedRefillIntervalSeconds: getEnvAsIntWithDefault("RATE_LIMIT_AUTHENTICATED_REFILL_SECONDS", 60),
		},
	}

	// Validate rate limit configuration
	if cfg.RateLimit.PublicMaxTokens <= 0 {
		panic("RATE_LIMIT_PUBLIC_MAX_TOKENS must be positive")
	}
	if cfg.RateLimit.PublicRequestsPerInterval <= 0 {
		panic("RATE_LIMIT_PUBLIC_REQUESTS must be positive")
	}
	if cfg.RateLimit.PublicRefillIntervalSeconds <= 0 {
		panic("RATE_LIMIT_PUBLIC_REFILL_SECONDS must be positive")
	}
	if cfg.RateLimit.AuthenticatedMaxTokens <= 0 {
		panic("RATE_LIMIT_AUTHENTICATED_MAX_TOKENS must be positive")
	}
	if cfg.RateLimit.AuthenticatedRequestsPerInterval <= 0 {
		panic("RATE_LIMIT_AUTHENTICATED_REQUESTS must be positive")
	}
	if cfg.RateLimit.AuthenticatedRefillIntervalSeconds <= 0 {
		panic("RATE_LIMIT_AUTHENTICATED_REFILL_SECONDS must be positive")
	}

	return cfg
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic("Environment variable " + key + " not set")
}

func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string) int {
	valueStr := getEnv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	panic("Environment variable " + name + " is not an integer")
}

func getEnvAsIntWithDefault(name string, defaultValue int) int {
	if value, exists := os.LookupEnv(name); exists && value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		panic("Environment variable " + name + " is not an integer")
	}
	return defaultValue
}

func getEnvAsBool(name string) bool {
	valStr := getEnv(name)
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	panic("Environment variable " + name + " is not a boolean")
}

func getEnvAsBoolWithDefault(name string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(name); exists && value != "" {
		if val, err := strconv.ParseBool(value); err == nil {
			return val
		}
		panic("Environment variable " + name + " is not a boolean")
	}
	return defaultValue
}

// IsEmailEnabled returns true if email is configured
func (c *config) IsEmailEnabled() bool {
	return c.Email.SMTPHost != ""
}
