package repository

import (
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/database"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *database.DB
}

func New(db *database.DB) *Repository {
	return &Repository{db: db}
}

type UserRepository struct {
	*Repository
}

func NewUserRepository(repo *Repository) *UserRepository {
	return &UserRepository{Repository: repo}
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type MediaRepository struct {
	*Repository
}

func NewMediaRepository(repo *Repository) *MediaRepository {
	return &MediaRepository{Repository: repo}
}

func (r *MediaRepository) Create(media *models.Media) error {
	return r.db.Create(media).Error
}

func (r *MediaRepository) FindByID(id string) (*models.Media, error) {
	var media models.Media
	err := r.db.Preload("User").Where("id = ?", id).First(&media).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *MediaRepository) FindPublic(page, limit int, genre, search string) ([]models.Media, int64, error) {
	var media []models.Media
	var total int64

	query := r.db.Model(&models.Media{}).Preload("User").Where("is_public = ?", true)

	if genre != "" {
		query = query.Where("genre = ?", genre)
	}

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ? OR tags ILIKE ?", like, like, like)
	}

	offset := (page - 1) * limit
	query.Count(&total)
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&media).Error

	return media, total, err
}

func (r *MediaRepository) DeleteByIDAndUserID(id, userID string) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Media{}).Error
}

func (r *MediaRepository) IncrementViewCount(id string) error {
	return r.db.Model(&models.Media{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}
