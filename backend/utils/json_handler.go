package utils

import (
	"encoding/json"
	"fmt"
)

// converts incoming data into JSON encoded format
func JsonHandler(data interface{}) ([]byte, error) {

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("unable to convert data into JSON format: %v", err)
	}

	return jsonData, nil
}
