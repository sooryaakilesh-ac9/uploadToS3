package internal

import (
    "backend/internal/delivery/http/handler"
    "backend/internal/infrastructure/metadata"
    "backend/internal/infrastructure/persistence/postgres"
    "backend/internal/infrastructure/s3"
    "backend/internal/usecase/image"
    "gorm.io/gorm"
)

func InitializeImageHandler(db *gorm.DB) (*handler.ImageHandler, error) {
    // Create infrastructure services
    s3Service, err := s3.NewS3Service()
    if err != nil {
        return nil, err
    }
    
    metadataService := metadata.NewMetadataService(s3Service)
    imageRepo := postgres.NewImageRepository(db)
    
    // Create use case
    imageUseCase := image.NewImageUseCase(imageRepo, s3Service, metadataService)
    
    // Create handler
    imageHandler := handler.NewImageHandler(imageUseCase)
    
    return imageHandler, nil
} 