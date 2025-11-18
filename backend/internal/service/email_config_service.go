package service

import (
    "errors"
    "strings"

    "github.com/example/golang-email/internal/model"
    "github.com/example/golang-email/internal/pkg/crypto"
    "github.com/example/golang-email/internal/repository"
)

type EmailConfigService interface {
    List() ([]model.EmailConfig, error)
    Get(id uint) (*model.EmailConfig, error)
    Create(input *model.EmailConfig) (*model.EmailConfig, error)
    Update(id uint, update *model.EmailConfig) (*model.EmailConfig, error)
    Delete(id uint) error
}

type emailConfigService struct {
    repo repository.EmailConfigRepository
}

func NewEmailConfigService(repo repository.EmailConfigRepository) EmailConfigService {
    return &emailConfigService{repo: repo}
}

func (s *emailConfigService) List() ([]model.EmailConfig, error) {
    return s.repo.List()
}

func (s *emailConfigService) Get(id uint) (*model.EmailConfig, error) {
    return s.repo.GetByID(id)
}

func (s *emailConfigService) Create(input *model.EmailConfig) (*model.EmailConfig, error) {
    if strings.TrimSpace(input.Provider) == "" || strings.TrimSpace(input.Host) == "" || input.Port <= 0 || strings.TrimSpace(input.Username) == "" || strings.TrimSpace(input.Password) == "" {
        return nil, errors.New("参数不合法")
    }
    enc, err := crypto.EncryptString(input.Password)
    if err != nil {
        return nil, err
    }
    input.Password = enc
    if err := s.repo.Create(input); err != nil {
        return nil, err
    }
    input.Password = ""
    return input, nil
}

func (s *emailConfigService) Update(id uint, update *model.EmailConfig) (*model.EmailConfig, error) {
    existing, err := s.repo.GetByID(id)
    if err != nil {
        return nil, err
    }
    if existing == nil {
        return nil, errors.New("记录不存在")
    }
    if update.Provider != "" {
        existing.Provider = update.Provider
    }
    if update.Host != "" {
        existing.Host = update.Host
    }
    if update.Port != 0 {
        existing.Port = update.Port
    }
    if update.Username != "" {
        existing.Username = update.Username
    }
    if update.Password != "" {
        enc, err := crypto.EncryptString(update.Password)
        if err != nil {
            return nil, err
        }
        existing.Password = enc
    }
    existing.IsActive = update.IsActive
    if update.DailyLimit != 0 {
        existing.DailyLimit = update.DailyLimit
    }
    if err := s.repo.Update(existing); err != nil {
        return nil, err
    }
    existing.Password = ""
    return existing, nil
}

func (s *emailConfigService) Delete(id uint) error {
    return s.repo.Delete(id)
}