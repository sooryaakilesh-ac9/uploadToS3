package db

import (
	"backend/pkg/images"
	"backend/utils"
	"fmt"

	"gorm.io/gorm"
)

func FetchImageFromDB(imageId uint) (*images.Flyer, error) {
	// Connect to the database
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	// Declare a variable to hold the fetched image
	var image images.Flyer

	// Fetch the image from the database by ID
	if err := db.First(&image, imageId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("quote with ID %d not found", imageId)
		}
		return nil, fmt.Errorf("failed to fetch quote: %w", err)
	}

	// a standalone module which conveters the given data into JSON format
	utils.JsonHandler(image)

	// Return the fetched image
	return &image, nil
}

func FetchAllImagesFromDB() ([]images.Flyer, error) {
	var images []images.Flyer

	// Connect to the database
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	// Fetch all quotes
	if err := db.Find(&images).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch images: %w", err)
	}

	fmt.Printf("fetched all quotes from DB!\n")

	return images, nil
}
