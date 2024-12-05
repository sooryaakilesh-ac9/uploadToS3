package service

import "backend/internal/domain/entity"

type MetadataService interface {
    UpdateImageMetadata(images []entity.Flyer) error
    UpdateQuoteMetadata(quotes []entity.Quote) error
} 