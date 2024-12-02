package utils

import (
	"testing"
)

// TestJsonHandler tests the JsonHandler function
func TestJsonHandler(t *testing.T) {
	// Test case 1: Valid data
	data := map[string]interface{}{"key": "value"}
	expected := `{"key":"value"}`

	result, err := JsonHandler(data)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if string(result) != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}

	// Test case 2: Invalid data (circular reference)
	circularData := map[string]interface{}{}
	circularData["self"] = circularData

	_, err = JsonHandler(circularData)
	if err == nil {
		t.Error("expected an error, got none")
	}
} 