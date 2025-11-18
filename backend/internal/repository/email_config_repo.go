package repository

import (
    "errors"

    "github.com/example/golang-email/internal/model"
    "gorm.io/gorm"
)

type EmailConfigRepository interface {
    List() ([]model.EmailConfig, error)
    GetByID(id uint) (*model.EmailConfig, error)
    Create(cfg *model.EmailConfig) error
    Update(cfg *model.EmailConfig) error
    Delete(id uint) error
}

type GormEmailConfigRepository struct {
    db *gorm.DB
}

func NewEmailConfigRepository(db *gorm.DB) EmailConfigRepository {
    return &GormEmailConfigRepository{db: db}
}

func (r *GormEmailConfigRepository) List() ([]model.EmailConfig, error) {
    var items []model.EmailConfig
    if err := r.db.Order("id DESC").Find(&items).Error; err != nil {
        return nil, err
    }
    return items, nil
}

func (r *GormEmailConfigRepository) GetByID(id uint) (*model.EmailConfig, error) {
    var item model.EmailConfig
    if err := r.db.First(&item, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &item, nil
}

func (r *GormEmailConfigRepository) Create(cfg *model.EmailConfig) error {
    return r.db.Create(cfg).Error
}

func (r *GormEmailConfigRepository) Update(cfg *model.EmailConfig) error {
    return r.db.Save(cfg).Error
}

func (r *GormEmailConfigRepository) Delete(id uint) error {
    return r.db.Delete(&model.EmailConfig{}, id).Error
}