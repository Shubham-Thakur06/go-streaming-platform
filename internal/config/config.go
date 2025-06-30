package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Storage  StorageConfig
	JWT      JWTConfig
	Host     HostConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type StorageConfig struct {
	Provider string // "aws", "azure", "gcp"
	AWS      AWSConfig
	Azure    AzureConfig
	GCP      GCPConfig
}

type AWSConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	BucketName      string
}

type AzureConfig struct {
	AccountName   string
	AccountKey    string
	ContainerName string
}

type GCPConfig struct {
	ProjectID       string
	BucketName      string
	CredentialsFile string
}

type JWTConfig struct {
	SecretKey string
	Expiry    int // in hours
}

type HostConfig struct {
	Username string
	Password string
	Email    string
}

func LoadEnv() error {
	return godotenv.Load()
}

func New() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "streaming_platform"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Storage: StorageConfig{
			Provider: getEnv("STORAGE_PROVIDER", "aws"),
			AWS: AWSConfig{
				AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
				SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
				Region:          getEnv("AWS_REGION", "us-east-1"),
				BucketName:      getEnv("AWS_BUCKET_NAME", ""),
			},
			Azure: AzureConfig{
				AccountName:   getEnv("AZURE_ACCOUNT_NAME", ""),
				AccountKey:    getEnv("AZURE_ACCOUNT_KEY", ""),
				ContainerName: getEnv("AZURE_CONTAINER_NAME", ""),
			},
			GCP: GCPConfig{
				ProjectID:       getEnv("GCP_PROJECT_ID", ""),
				BucketName:      getEnv("GCP_BUCKET_NAME", ""),
				CredentialsFile: getEnv("GCP_CREDENTIALS_FILE", ""),
			},
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET_KEY", "your-secret-key"),
			Expiry:    getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		},
		Host: HostConfig{
			Username: getEnv("HOST_USERNAME", "host"),
			Password: getEnv("HOST_PASSWORD", "host123"),
			Email:    getEnv("HOST_EMAIL", "host@streaming-platform.com"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
