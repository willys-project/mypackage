package handler

import (
	"fmt"
	"log"
	"runtime"
)

func LogErrorWithLine(err error) string {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		return fmt.Sprintf("error: %s (%s:%d)", err, file, line)
	}
	return fmt.Sprintf("error: %s", err)
}

func LogError(err error) {
	log.Printf("Error: %v", err)
}
