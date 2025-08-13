package handler

import (
	//	"goresponse"
	"log"

	"github.com/getsentry/sentry-go"
)

func HandleError(functionName string, err error) {
	log.Println(functionName, ":", err)
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("functionTag", functionName)
		scope.SetLevel(sentry.LevelError)
	})
	// Replace goresponse.LogErrorWithLine with a local implementation
	sentry.CaptureMessage(LogErrorWithLine(err))
}
