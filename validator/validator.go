package validator

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"

	// Sesuaikan path berikut dengan repo Anda

	"github.com/willys-project/mypackage/goresponse"
	"github.com/willys-project/mypackage/handler"
)

// ===== Constants =====

const (
	DAILY   = "daily"
	WEEKLY  = "weekly"
	MONTHLY = "monthly"
)

// ===== Globals (sesuai pola yang Anda kirim) =====

var (
	// BigqueryClient *bigquery.Client
	redisClient *redis.Client
	mu          sync.Mutex
	envFlag     string
	debug       bool
	err         error
	// JwtSecret, _   = functions.GetSecret(projectID, "jwt-secret")
)

// ValidateSecCode memeriksa apakah query param "secCode" ada.
// Jika kosong, langsung mengembalikan response 422 Unprocessable Entity.
func ValidateSecCode(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Query().Get("secCode") == "" {
		goresponse.ApiResUnprocEntity(w, "secCode is required")
		return false
	}
	return true
}

func ValidsecCodegranularitystartDateendDate(req *http.Request) (bool, error) {
	queryParams := req.URL.Query()

	secCode := queryParams.Get("secCode")
	granularity := queryParams.Get("granularity")
	startDateStr := queryParams.Get("startDate")
	endDateStr := queryParams.Get("endDate")

	// Cek parameter wajib
	if startDateStr == "" || endDateStr == "" || granularity == "" || secCode == "" {
		err := handler.NewCustomError("periksa lagi parameter")
		return false, err
	}

	// Validasi panjang secCode
	if len(secCode) > 4 {
		err := handler.NewCustomError("secCode must be no longer than 4 characters")
		return false, err
	}

	// Validasi nilai granularity
	if !contains([]string{DAILY, WEEKLY, MONTHLY}, granularity) {
		err := handler.NewCustomError("Invalid granularity value")
		return false, err
	}

	// Parse tanggal
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		handler.LogErrorWithLine(errors.Wrap(err, "failed to parse startDate"))
		return false, handler.NewCustomError("Invalid startDate format")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		handler.LogErrorWithLine(errors.Wrap(err, "failed to parse endDate"))
		return false, handler.NewCustomError("Invalid endDate format")
	}

	// startDate tidak boleh lebih dari 1 tahun yang lalu (untuk non-production)
	if envFlag != "production" {
		maxStartDate := time.Now().AddDate(-1, 0, 0)
		if startDate.Before(maxStartDate) {
			log.Printf("startDate %s exceeds the limit of 1 year from now", startDateStr)
			return false, handler.NewCustomError("startDate exceeds the limit of 1 year from now")
		}
	}

	// endDate maksimal 1 bulan dari startDate
	if endDate.After(startDate.AddDate(0, 1, 0)) {
		err := handler.NewCustomError("endDate %s exceeds the maximum allowed range of 1 month from startDate %s", endDateStr, startDateStr)
		log.Println(err)
		return false, err
	}

	return true, nil
}

// ===== Helper =====

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
