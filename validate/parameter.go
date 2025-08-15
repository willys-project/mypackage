package validate

import (
	"log"
	"net/http"
	"sync"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/willys-project/mypackage/functions"
	"github.com/willys-project/mypackage/handler"
)

const CACHE_PREFIX = "param:"
const (
	DAILY         = "daily"
	WEEKLY        = "weekly"
	MONTHLY       = "monthly"
	REGION        = "asia-southeast2"
	CACHE_VERSION = 3
	projectID     = "ticmidatadev"
)

var (
	BigqueryClient *bigquery.Client
	redisClient    *redis.Client
	mu             sync.Mutex
	envFlag        string
	debug          bool
	err            error
	JwtSecret, _   = functions.GetSecret(projectID, "jwt-secret")
)

func ValidateParameters(req *http.Request) (bool, error) {
	queryParams := req.URL.Query()

	secCode := queryParams.Get("secCode")
	granularity := queryParams.Get("granularity")
	startDateStr := queryParams.Get("startDate")
	endDateStr := queryParams.Get("endDate")

	// Check if any required parameter is missing
	if startDateStr == "" || endDateStr == "" || granularity == "" || secCode == "" {
		err := handler.NewCustomError("periksa lagi parameter")
		return false, err
	}

	// Validate secCode length
	if len(secCode) > 4 {
		err := handler.NewCustomError("secCode must be no longer than 4 characters")
		return false, err
	}

	// Validate granularity value
	if !contains([]string{DAILY, WEEKLY, MONTHLY}, granularity) {
		err := handler.NewCustomError("Invalid granularity value")
		return false, err
	}

	// Parse startDate and endDate into time.Time
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

	// Check if startDate is not more than 1 year from now (only for non-production environments)
	if envFlag != "production" {
		maxStartDate := time.Now().AddDate(-1, 0, 0)
		if startDate.Before(maxStartDate) {
			log.Printf("startDate %s exceeds the limit of 1 year from now", startDateStr)
			return false, handler.NewCustomError("startDate exceeds the limit of 1 year from now")
		}
	}

	// Limit endDate to be within 1 month from startDate
	if endDate.After(startDate.AddDate(0, 1, 0)) {
		err := handler.NewCustomError("endDate %s exceeds the maximum allowed range of 1 month from startDate %s", endDateStr, startDateStr)
		log.Println(err)
		return false, err
	}

	return true, nil
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
