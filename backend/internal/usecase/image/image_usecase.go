package image

import (
	"backend/internal/domain/entity"
	"backend/internal/domain/repository"
	"backend/internal/domain/service"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type ImageUseCase struct {
	imageRepo       repository.ImageRepository
	s3Service      service.S3Service
	metadataService service.MetadataService
}

func NewImageUseCase(repo repository.ImageRepository, s3 service.S3Service, meta service.MetadataService) *ImageUseCase {
	return &ImageUseCase{
		imageRepo:       repo,
		s3Service:      s3,
		metadataService: meta,
	}
}

func (uc *ImageUseCase) UploadImage(file multipart.File, header *multipart.FileHeader) (*entity.Flyer, error) {
	// Create temporary file
	tempFile, err := os.CreateTemp("", "upload-*.tmp")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy uploaded file to temp file
	if _, err := io.Copy(tempFile, file); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	flyer, err := uc.createFlyer(tempFile.Name(), header.Filename)
	if err != nil {
		return nil, err
	}

	id, err := uc.imageRepo.Store(flyer)
	if err != nil {
		return nil, err
	}
	flyer.Id = id

	// Update the URL with the ID
	flyer.Url = fmt.Sprintf("%s/%d_%s", os.Getenv("S3_BUCKET_NAME"), id, header.Filename)
	if err := uc.imageRepo.Update(flyer); err != nil {
		return nil, err
	}

	if err := uc.s3Service.UploadImage(tempFile.Name(), fmt.Sprintf("%d_%s", id, header.Filename)); err != nil {
		return nil, err
	}

	if err := uc.updateMetadata(); err != nil {
		return nil, err
	}

	return flyer, nil
}

func (uc *ImageUseCase) createFlyer(filePath, filename string) (*entity.Flyer, error) {
	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer imgFile.Close()

	img, format, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	fileFormat := strings.ToUpper(format)
	if fileFormat == "JPG" {
		fileFormat = "JPEG"
	}

	orientation := "portrait"
	if width > height {
		orientation = "landscape"
	}

	return &entity.Flyer{
		Design: entity.Design{
			Resolution: entity.Resolution{
				Width:  width,
				Height: height,
				Unit:   1,
			},
			Type:        "image",
			Tags:        uc.extractImageTags(filename),
			FileFormat:  fileFormat,
			Orientation: orientation,
			FileName:    filename,
		},
		Lang: "en-US",
		Url:  fmt.Sprintf("%s/%s", os.Getenv("S3_BUCKET_NAME"), filename),
	}, nil
}

func (uc *ImageUseCase) extractImageTags(filename string) []string {
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	return strings.Split(name, "_")
}

func (uc *ImageUseCase) updateMetadata() error {
	images, err := uc.imageRepo.FindAll()
	if err != nil {
		return err
	}
	return uc.metadataService.UpdateImageMetadata(images)
}

func (uc *ImageUseCase) ImportImages(importDir string) (successCount, failureCount int, err error) {
	// Validate import directory exists
	if _, err := os.Stat(importDir); os.IsNotExist(err) {
		return 0, 0, fmt.Errorf("import directory %s does not exist", importDir)
	}

	// Process each file in the directory
	entries, err := os.ReadDir(importDir)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !isValidImageFile(filename) {
			failureCount++
			continue
		}

		if err := uc.processImageFile(importDir, filename); err != nil {
			log.Printf("Failed to process %s: %v", filename, err)
			failureCount++
			continue
		}
		successCount++
	}

	// Update metadata after importing all images
	if err := uc.updateMetadata(); err != nil {
		return successCount, failureCount, fmt.Errorf("failed to update metadata: %w", err)
	}

	return successCount, failureCount, nil
}

func (uc *ImageUseCase) processImageFile(importDir, filename string) error {
	sourcePath := filepath.Join(importDir, filename)
	flyer, err := uc.createFlyer(sourcePath, filename)
	if err != nil {
		return fmt.Errorf("failed to create flyer: %w", err)
	}

	id, err := uc.imageRepo.Store(flyer)
	if err != nil {
		return fmt.Errorf("failed to store in database: %w", err)
	}

	flyer.Id = id
	flyer.Url = fmt.Sprintf("%s/%d_%s", os.Getenv("S3_BUCKET_NAME"), id, filename)
	
	if err := uc.imageRepo.Update(flyer); err != nil {
		return fmt.Errorf("failed to update flyer URL: %w", err)
	}

	if err := uc.s3Service.UploadImage(sourcePath, fmt.Sprintf("%d_%s", id, filename)); err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	return nil
}

func isValidImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

// Additional methods... 