package postgres

import (
    "backend/internal/domain/entity"
    "gorm.io/gorm"
)

type ImageRepository struct {
    db *gorm.DB
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
    return &ImageRepository{db: db}
}

func (r *ImageRepository) Store(image *entity.Flyer) (uint, error) {
    result := r.db.Create(image)
    if result.Error != nil {
        return 0, result.Error
    }
    return image.Id, nil
}

func (r *ImageRepository) Update(image *entity.Flyer) error {
    return r.db.Save(image).Error
}

func (r *ImageRepository) FindByID(id uint) (*entity.Flyer, error) {
    var image entity.Flyer
    if err := r.db.First(&image, id).Error; err != nil {
        return nil, err
    }
    return &image, nil
}

func (r *ImageRepository) FindAll() ([]entity.Flyer, error) {
    var images []entity.Flyer
    if err := r.db.Find(&images).Error; err != nil {
        return nil, err
    }
    return images, nil
} 