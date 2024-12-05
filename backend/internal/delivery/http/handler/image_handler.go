package handler

import (
	"backend/internal/domain/entity"
	"backend/internal/usecase"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type ImageHandler struct {
	imageUseCase *usecase.ImageUseCase
}

func NewImageHandler(useCase *usecase.ImageUseCase) *ImageHandler {
	return &ImageHandler{
		imageUseCase: useCase,
	}
}

func (h *ImageHandler) HandleImagesUpload(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	flyer, err := h.imageUseCase.UploadImage(file, handler)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Image uploaded successfully",
		"filename": handler.Filename,
		"id":       flyer.Id,
	})
}

func (h *ImageHandler) HandleImagesImport(w http.ResponseWriter, r *http.Request) {
	importDir := os.Getenv("IMPORT_DIR_IMAGES")
	if importDir == "" {
		http.Error(w, "IMPORT_DIR_IMAGES environment variable not set", http.StatusBadRequest)
		return
	}

	successCount, failureCount, err := h.imageUseCase.ImportImages(importDir)
	if err != nil {
		log.Printf("Error importing images: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success_count": successCount,
		"failure_count": failureCount,
		"total":         successCount + failureCount,
	})
} 