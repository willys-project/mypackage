package functions

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"time"

	"cloud.google.com/go/storage"
	"github.com/willys-project/mypackage/goresponse"
	"github.com/willys-project/mypackage/handler"
)

var ctx = context.Background() // Pastikan kamu sudah mengimpor context package

// Pastikan kamu sudah menginstal package ini

// IsEmpty checks if the provided data is empty.
func IsEmpty(data interface{}) bool {
	switch v := data.(type) {
	case string:
		return v == ""
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	default:
		return false
	}
}

func GetAppName() string {
	_, file, _, _ := runtime.Caller(0)                             // Ambil nama file sumber kode
	return filepath.Base(file[:len(file)-len(filepath.Ext(file))]) // Ambil nama file tanpa ekstensi
}

// ValidateParameters checks if the required "secCode" parameter is present in the request and writes an error response if missing.
func ValidateParameters(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Query().Get("secCode") == "" {
		goresponse.ApiResUnprocEntity(w, "secCode is required")
		return false
	}
	return true
}

func GetFile(secCode string, bucketName string, objectPrefix string) (map[string]interface{}, error) {

	// Establish a connection to Google Cloud Storage
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Query to get objects from the bucket with a specific prefix
	query := client.Bucket(bucketName).Objects(ctx, &storage.Query{Prefix: objectPrefix + secCode})

	attrs, err := query.Next()
	if err == storage.ErrObjectNotExist {
		// Object not found
		handler.HandleError("getFile", err)
		return nil, errors.New("Object not found")
	}
	if err != nil {
		return nil, err
	}

	// Read the contents of the object
	rc, err := client.Bucket(bucketName).Object(attrs.Name).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	// Read JSON data from the object
	dataBytes, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a map
	var dataMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		return nil, err
	}

	return dataMap, nil
}

// LimitToLastMonth checks if the given date is within the last month and writes an error to the response if not.
func LimitToLastMonth(ipoUpdatedSince string, w http.ResponseWriter) bool {
	// Parse the date in the expected format (yyyy-MM-dd)
	ipoUpdatedTime, err := time.Parse("2006-01-02", ipoUpdatedSince)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return false
	}

	// Get the current date and time
	currentTime := time.Now()

	// Calculate one month ago
	oneMonthAgo := currentTime.AddDate(0, -1, 0) // Menambahkan -1 bulan

	// If the ipoUpdatedSince date is earlier than one month ago, return forbidden
	if ipoUpdatedTime.Before(oneMonthAgo) {
		http.Error(w, "Forbidden: updated more than 1 month ago", http.StatusForbidden)
		return false
	}

	return true
}
