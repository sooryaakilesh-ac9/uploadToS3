package repository

import "backend/internal/domain/entity"

type ImageRepository interface {
    Store(image *entity.Flyer) (uint, error)
    Update(image *entity.Flyer) error
    FindByID(id uint) (*entity.Flyer, error)
    FindAll() ([]entity.Flyer, error)
} 