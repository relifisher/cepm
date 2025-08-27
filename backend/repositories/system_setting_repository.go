package repositories

import (
	"cepm-backend/models"
	"gorm.io/gorm"
)

type SystemSettingRepository struct {
	db *gorm.DB
}

func NewSystemSettingRepository(db *gorm.DB) *SystemSettingRepository {
	return &SystemSettingRepository{db: db}
}

func (r *SystemSettingRepository) GetSetting(key string) (*models.SystemSetting, error) {
	var setting models.SystemSetting
	err := r.db.Where("key = ?", key).First(&setting).Error
	return &setting, err
}

func (r *SystemSettingRepository) UpdateSetting(setting *models.SystemSetting) error {
	return r.db.Save(setting).Error
}

func (r *SystemSettingRepository) CreateSetting(setting *models.SystemSetting) error {
	return r.db.Create(setting).Error
}
