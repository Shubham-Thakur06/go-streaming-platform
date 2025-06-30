package database

import (
	"fmt"
	"log"

	"github.com/Shubham-Thakur06/go-streaming-platform/internal/config"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/models"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

func Init(cfg config.DatabaseConfig) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	if err := migrateDB(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connected and migrated successfully")
	return &DB{DB: db}, nil
}

// InitWithHost initializes the database connection and creates initial host user
func InitWithHost(dbCfg config.DatabaseConfig, hostCfg config.HostConfig) (*DB, error) {
	db, err := Init(dbCfg)
	if err != nil {
		return nil, err
	}

	// Create initial host user
	if err := CreateInitialHost(db.DB, hostCfg); err != nil {
		return nil, fmt.Errorf("failed to create initial host user: %w", err)
	}

	return db, nil
}

func migrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Media{},
		&models.UserActivity{},
	)
}

// CreateInitialHost creates the initial host user if it doesn't exist
func CreateInitialHost(db *gorm.DB, hostConfig config.HostConfig) error {
	var count int64
	db.Model(&models.User{}).Count(&count)

	if count > 0 {
		log.Println("Host user already exists, skipping creation")
		return nil
	}

	// Hash the password
	hashedPassword, err := HashPassword(hostConfig.Password)
	if err != nil {
		return err
	}

	// Create the host user
	host := models.User{
		Username: hostConfig.Username,
		Email:    hostConfig.Email,
		Password: hashedPassword,
	}

	if err := db.Create(&host).Error; err != nil {
		return err
	}

	log.Println("Initial host user created successfully")
	log.Printf("Username: %s", hostConfig.Username)
	log.Printf("Password: %s", hostConfig.Password)
	log.Println("Please change the password after first login!")

	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}
