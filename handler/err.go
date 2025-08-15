package handler

import (
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
)

// CustomError defines a custom error type with a message
type CustomError struct {
	message string
}

// Error implements the error interface for CustomError
func (e *CustomError) Error() string {
	return e.message
}

// HandleError logs the error and sends it to Sentry
func HandleError(functionName string, err error) {
	log.Println(functionName, ":", err)
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("functionTag", functionName)
		scope.SetLevel(sentry.LevelError)
	})
	// Replace goresponse.LogErrorWithLine with a local implementation
	sentry.CaptureMessage(LogErrorWithLine(err))
}

// NewCustomError creates a new CustomError instance
func NewCustomError(format string, args ...interface{}) *CustomError {
	return &CustomError{
		message: fmt.Sprintf(format, args...),
	}
}
