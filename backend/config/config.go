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

type payment struct {
	StripeKey     string
	WebhookSecret string
}

type axiomConfig struct {
	Dataset string
	Token   string
}

type config struct {
	Env            string
	AppEnvironment string
	BaseURL        string
	Port           int
	AdminUserId    string
	Postgres       postgresConfig
	ObjStorage     objStorageConfig
	PgAdmin        pgAdminConfig
	Email       emailConfig
	ProductInfo productInfo
	Payment     payment
	Axiom       axiomConfig
}

func New() *config {
	cfg := &config{
		Env:            getEnv("ENV"),
		AppEnvironment: getEnvWithDefault("APP_ENVIRONMENT", "self-hosted"),
		BaseURL:        getEnv("BASE_URL"),
		Port:           getEnvAsInt("PORT"),
		AdminUserId: getEnv("ADMIN_USER_ID"),
		Payment: payment{
			StripeKey:     getEnv("STRIPE_KEY"),
			WebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET"),
		},
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
			FromAddress: getEnv("EMAIL_FROM_ADDRESS"),
			SMTPHost:    getEnv("SMTP_HOST"),
			SMTPPort:    getEnvAsInt("SMTP_PORT"),
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
	}

	// Validate APP_ENVIRONMENT
	if cfg.AppEnvironment != "cloud" && cfg.AppEnvironment != "self-hosted" {
		panic("APP_ENVIRONMENT must be either 'cloud' or 'self-hosted', got: " + cfg.AppEnvironment)
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

func getEnvAsBool(name string) bool {
	valStr := getEnv(name)
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	panic("Environment variable " + name + " is not a boolean")
}

func getEnvAsBoolWithDefault(name string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(name); exists {
		if val, err := strconv.ParseBool(value); err == nil {
			return val
		}
		panic("Environment variable " + name + " is not a boolean")
	}
	return defaultValue
}

// IsCloud returns true if the application is running in cloud mode
func (c *config) IsCloud() bool {
	return c.AppEnvironment == "cloud"
}

// IsSelfHosted returns true if the application is running in self-hosted mode
func (c *config) IsSelfHosted() bool {
	return c.AppEnvironment == "self-hosted"
}
