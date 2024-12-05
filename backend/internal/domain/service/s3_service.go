package service

type S3Service interface {
    UploadImage(filePath string, fileName string) error
    UploadMetadata(filePath string, fileName string) error
} 