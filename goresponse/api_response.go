package goresponse

import (
	"encoding/json"
	"net/http"

	"github.com/willys-project/mypackage/model"
)

// ApiResNotFound mengirimkan response 404 Not Found
func ApiResNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found"))
}

// ApiResMethodNotAllowed mengirimkan response 405 Method Not Allowed
func ApiResMethodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Method Not Allowed"))
}

// ApiResUnprocEntity mengirimkan response 422 Unprocessable Entity dengan message
func ApiResUnprocEntity(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	response := model.APIResponse{
		Message: message,
		Status:  http.StatusUnprocessableEntity,
	}
	json.NewEncoder(w).Encode(response)
}

// ApiResOK mengirimkan response 200 OK dengan data, atau 204 No Content jika data kosong
func ApiResOK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	// Check if data is nil or empty
	if data == nil || isEmpty(data) {
		// If data is empty, return a 204 No Content status
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// If data is not empty, return the data with a 200 OK status
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// APIResUnauthorized mengirimkan response 401 Unauthorized dengan message
func APIResUnauthorized(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	response := model.APIResponse{
		Message: message,
		Status:  http.StatusUnauthorized,
	}
	json.NewEncoder(w).Encode(response)
}

// isEmpty memeriksa apakah data kosong (untuk penggunaan di ApiResOK)
func isEmpty(data interface{}) bool {
	// Add custom logic for determining if data is considered "empty"
	// For example, for slices and maps:
	switch v := data.(type) {
	case nil:
		return true
	case string:
		return v == ""
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	}
	return false
}

// APIResBadRequest sends a 400 Bad Request response with a message.
func APIResBadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	response := model.APIResponse{
		Message: message,
		Status:  http.StatusBadRequest,
	}
	json.NewEncoder(w).Encode(response)
}
