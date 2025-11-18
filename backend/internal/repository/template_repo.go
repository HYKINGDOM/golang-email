package repository

import (
    "github.com/example/golang-email/internal/model"
    "gorm.io/gorm"
)

type EmailTemplateRepository interface {
    List() ([]model.EmailTemplate, error)
    GetByID(id uint) (*model.EmailTemplate, error)
    Create(t *model.EmailTemplate) error
    Update(t *model.EmailTemplate) error
    Delete(id uint) error
}

type GormEmailTemplateRepository struct{ db *gorm.DB }

func NewEmailTemplateRepository(db *gorm.DB) EmailTemplateRepository { return &GormEmailTemplateRepository{db: db} }

func (r *GormEmailTemplateRepository) List() ([]model.EmailTemplate, error) {
    var items []model.EmailTemplate
    if err := r.db.Order("id DESC").Find(&items).Error; err != nil { return nil, err }
    return items, nil
}

func (r *GormEmailTemplateRepository) GetByID(id uint) (*model.EmailTemplate, error) {
    var t model.EmailTemplate
    if err := r.db.First(&t, id).Error; err != nil { return nil, err }
    return &t, nil
}

func (r *GormEmailTemplateRepository) Create(t *model.EmailTemplate) error { return r.db.Create(t).Error }
func (r *GormEmailTemplateRepository) Update(t *model.EmailTemplate) error { return r.db.Save(t).Error }
func (r *GormEmailTemplateRepository) Delete(id uint) error { return r.db.Delete(&model.EmailTemplate{}, id).Error }