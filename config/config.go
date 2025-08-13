package config

import (
	"context"
	"sync"

	"cloud.google.com/go/bigquery"
	"github.com/go-redis/redis"
	"github.com/willys-project/mypackage/functions"
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
