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
	PersonalName   string
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
	FromAddress       string
	McServer          string
	McPort            int
	PostmarkServerUrl string
	PostmarkToken     string
}

type legal struct {
	ToSVersion string
	PPVersion  string
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
	Env         string
	BaseURL     string
	Port        int
	AdminUserId string
	Postgres    postgresConfig
	ObjStorage  objStorageConfig
	PgAdmin     pgAdminConfig
	Email       emailConfig
	ProductInfo productInfo
	Legal       legal
	Payment     payment
	Axiom       axiomConfig
}

func New() *config {
	return &config{
		Env:         getEnv("ENV"),
		BaseURL:     getEnv("BASE_URL"),
		Port:        getEnvAsInt("PORT"),
		AdminUserId: getEnv("ADMIN_USER_ID"),
		Legal: legal{
			ToSVersion: getEnv("TOS_VERSION"),
			PPVersion:  getEnv("PP_VERSION"),
		},
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
			FromAddress:       getEnv("EMAIL_FROM_ADDRESS"),
			McServer:          getEnv("MAILCATCHER_SERVER"),
			McPort:            getEnvAsInt("MAILCATCHER_PORT"),
			PostmarkServerUrl: getEnv("POSTMARK_SERVER_URL"),
			PostmarkToken:     getEnv("POSTMARK_TOKEN"),
		},
		ProductInfo: productInfo{
			ProductName:    "Release Notes",
			CompanyName:    "Release Notes Inc.",
			CompanyAddress: "1234 Elm St, Springfield, IL 62701",
			SupportEmail:   "support@test.de",
			PersonalName:   "Daniel",
		},
		Axiom: axiomConfig{
			Dataset: getEnv("AXIOM_DATASET"),
			Token:   getEnv("AXIOM_TOKEN"),
		},
	}
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic("Environment variable " + key + " not set")
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
