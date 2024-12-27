package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, code int, message string) {
	JSON(w, code, Response{
		Success: false,
		Error:   message,
	})
} 