package config

import (
	"context"
	"mypackage/functions"
	"sync"

	"cloud.google.com/go/bigquery"
	"github.com/go-redis/redis"
)

// Global variables
var (
	ctx            = context.Background()
	BigqueryClient *bigquery.Client
	projectID      = "ticmidatadev" // Replace with your actual project ID
	redisClient    *redis.Client
	mu             sync.Mutex
	envFlag        string
	debug          bool
	JwtSecret, _   = functions.GetSecret(projectID, "jwt-secret")
	baseURL        string
)
