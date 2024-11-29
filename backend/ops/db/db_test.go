package db

import (
	"testing"
)

func TestGetDB(t *testing.T) {
	result, err := GetDB()

	// Check if the result is nil
	if result != nil {
		t.Errorf("Expected result to be nil, but got %v", result)
	}

	// Check if the error matches the expected message
	expectedErr := "database connection not initialized, call initDB() first"
	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected error %q, but got %v", expectedErr, err)
	}
}