package services

import (
	"cepm-backend/models"
	"cepm-backend/repositories"

	"gorm.io/gorm"
)

type SystemSettingService struct {
	settingRepo *repositories.SystemSettingRepository
}

func NewSystemSettingService(settingRepo *repositories.SystemSettingRepository) *SystemSettingService {
	return &SystemSettingService{settingRepo: settingRepo}
}

func (s *SystemSettingService) GetSetting(key string) (*models.SystemSetting, error) {
	return s.settingRepo.GetSetting(key)
}

func (s *SystemSettingService) UpdateSetting(setting *models.SystemSetting) error {
	return s.settingRepo.UpdateSetting(setting)
}

func (s *SystemSettingService) CreateOrUpdateSetting(key, value string) error {
	setting, err := s.settingRepo.GetSetting(key)
	if err != nil {
		// If not found, create a new setting
		if err == gorm.ErrRecordNotFound {
			newSetting := models.SystemSetting{Key: key, Value: value}
			return s.settingRepo.CreateSetting(&newSetting)
		} else {
			return err
		}
	}
	// If found, update the existing setting
	setting.Value = value
	return s.settingRepo.UpdateSetting(setting)
}
