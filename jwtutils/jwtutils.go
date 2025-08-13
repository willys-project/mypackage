package jwtutils

import (
	"encoding/json"
	"net/http"
	"strings"
)

// handleJWTErrorJSON memeriksa error dari JWT dan memberikan response JSON yang sesuai
func HandleJWTErrorJSON(w http.ResponseWriter, err error) bool {
	if err != nil {
		msg := err.Error()
		var statusCode int
		var message string

		switch {
		case strings.Contains(msg, "signature is invalid"):
			statusCode = http.StatusForbidden // 403
			message = "Signature is invalid"
		case strings.Contains(msg, "token is expired"):
			statusCode = http.StatusUnauthorized // 401
			message = "Token expired"
		case strings.Contains(msg, "token not valid yet"):
			statusCode = http.StatusUnauthorized // 401
			message = "Token not valid yet"
		case strings.Contains(msg, "malformed"):
			statusCode = http.StatusBadRequest // 400
			message = "Malformed token"
		case strings.Contains(msg, "unexpected signing method"):
			statusCode = http.StatusUnauthorized // 401
			message = "Invalid signing method"
		default:
			statusCode = http.StatusUnauthorized
			message = msg
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		resp := map[string]interface{}{
			"error":   true,
			"message": message,
			"data":    nil,
		}

		json.NewEncoder(w).Encode(resp)
		return true
	}
	return false
}
