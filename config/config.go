package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis"
	"github.com/willys-project/mypackage/secrets"
)

type Env struct {
	ProjectID      string
	EnvFlag        string
	Path           string
	Debug          bool
	BaseURL        string
	BigQueryClient *bigquery.Client
	RedisClient    *redis.Client
	SentryDSN      string
}

// String yang sebelumnya hardcoded dijadikan parameter
type InitOptions struct {
	SentrySecretName string // ex: "sentry-dsn-ticmidata-directapi"
	BaseURLProd      string // ex: "https://api2.ticmidata.co.id/direct/v1/"
	BaseURLDefault   string // ex: "https://api.ticmidata.co.id/direct/v1/"
	RedisSecretName  string // ex: "cp-redis-host-development"
}

func Init(projectID, envFlag, path string, debug bool, opt InitOptions) (*Env, error) {
	if debug {
		fmt.Println("proses Init() config")
	}

	ctx := context.Background()

	// ===== Secret Manager client (reusable) =====
	secClient, err := secrets.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("secrets.New: %w", err)
	}
	defer secClient.Close()

	// ===== BigQuery =====
	bqClient, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Println(err, "failed to create BigQuery client")
		return nil, err
	}

	// ===== Sentry DSN dari Secret Manager =====
	dsn, err := secClient.Get(ctx, projectID, opt.SentrySecretName)
	if debug && err != nil {
		log.Println("getSecret DSN error:", err)
	}
	if debug && dsn != "" {
		fmt.Println("Sentry DSN loaded")
	}

	e := &Env{
		ProjectID:      projectID,
		EnvFlag:        envFlag,
		Path:           path,
		Debug:          debug,
		BigQueryClient: bqClient,
		SentryDSN:      dsn,
	}

	// ===== Konfigurasi BaseURL & Redis =====
	switch envFlag {
	case "local":
		e.BaseURL = opt.BaseURLDefault + path + "/?secCode="
		e.RedisClient = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   1,
		})

	case "production":
		e.BaseURL = opt.BaseURLProd + path + "/?secCode="
		redisHost, err := secClient.Get(ctx, projectID, opt.RedisSecretName)
		if err != nil {
			return nil, fmt.Errorf("secret redis host: %w", err)
		}
		e.RedisClient = redis.NewClient(&redis.Options{
			Addr: redisHost + ":6379",
			DB:   1,
		})
		if _, err = e.RedisClient.Ping().Result(); err != nil {
			return nil, fmt.Errorf("redis ping: %w", err)
		}

	default:
		e.BaseURL = opt.BaseURLDefault + "saham/" + path + "/?secCode="
		redisHost, err := secClient.Get(ctx, projectID, opt.RedisSecretName)
		if err != nil {
			return nil, fmt.Errorf("secret redis host: %w", err)
		}
		e.RedisClient = redis.NewClient(&redis.Options{
			Addr: redisHost + ":6379",
			DB:   1,
		})
		if _, err = e.RedisClient.Ping().Result(); err != nil {
			return nil, fmt.Errorf("redis ping: %w", err)
		}
	}

	// ===== Init Sentry =====
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              e.SentryDSN,
		Environment:      e.EnvFlag,
		TracesSampleRate: 1.0,
		EnableTracing:    true,
		Debug:            e.Debug,
	}); err != nil {
		if envFlag == "local" {
			go func() {
				sentry.ConfigureScope(func(scope *sentry.Scope) {
					scope.SetTag("functionTag", "InitConfig")
					scope.SetLevel(sentry.LevelError)
				})
				sentry.CaptureMessage(err.Error())
			}()
		} else {
			return nil, fmt.Errorf("sentry.Init: %w", err)
		}
	}

	return e, nil
}

func (e *Env) Close() {
	sentry.Flush(2 * time.Second)
	if e.RedisClient != nil {
		_ = e.RedisClient.Close()
	}
}
