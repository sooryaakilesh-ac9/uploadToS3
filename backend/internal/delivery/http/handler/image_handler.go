package handler

import (
	"backend/internal/delivery/http/response"
	"backend/internal/usecase/image"
	"log"
	"net/http"
	"os"
)

type ImageHandler struct {
	imageUseCase *image.ImageUseCase
}

func NewImageHandler(useCase *image.ImageUseCase) *ImageHandler {
	return &ImageHandler{
		imageUseCase: useCase,
	}
}

func (h *ImageHandler) HandleImagesUpload(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("image")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to retrieve file")
		return
	}
	defer file.Close()

	flyer, err := h.imageUseCase.UploadImage(file, handler)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, map[string]interface{}{
		"message":  "Image uploaded successfully",
		"filename": handler.Filename,
		"id":       flyer.Id,
	})
}

func (h *ImageHandler) HandleImagesImport(w http.ResponseWriter, r *http.Request) {
	importDir := os.Getenv("IMPORT_DIR_IMAGES")
	if importDir == "" {
		response.Error(w, http.StatusBadRequest, "IMPORT_DIR_IMAGES environment variable not set")
		return
	}

	successCount, failureCount, err := h.imageUseCase.ImportImages(importDir)
	if err != nil {
		log.Printf("Error importing images: %v", err)
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, map[string]interface{}{
		"success_count": successCount,
		"failure_count": failureCount,
		"total":         successCount + failureCount,
	})
} 