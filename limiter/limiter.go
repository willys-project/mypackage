package limiter

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis" // v6
	"github.com/willys-project/mypackage/auth"
	"github.com/willys-project/mypackage/daterange"
	"github.com/willys-project/mypackage/handler"
)

var mu sync.Mutex

// LimitAccess limits the number of requests per month for a given API key and path.
func LimitAccess(maxRequests int, RedisClient *redis.Client, Debug bool, Secret string, next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if Debug {
			fmt.Println("proses limitAccess")
		}

		var Bearer string
		apiKey := req.Header.Get("X-Auth-Key")

		// Mendapatkan header Authorization
		authorizationHeader := req.Header.Get("Authorization")
		if Debug {
			fmt.Println("proses Authorization Header", authorizationHeader)
		}
		if strings.HasPrefix(authorizationHeader, "Bearer ") {
			Bearer = strings.TrimPrefix(authorizationHeader, "Bearer ")
		}

		if Debug {
			fmt.Println("proses Bearer", Bearer)
		}

		// Verifikasi JWT Token
		token, err := auth.VerifyJWTToken(Bearer, Secret)
		if handler.HandleJWTErrorJSON(res, err) {
			return // hentikan eksekusi kalau error
		}
		if err != nil {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Mengambil klaim "iat" (issued at) dari JWT
		scientificNotation, err := auth.GetJWTClaim(token, "iat")
		if err != nil {
			http.Error(res, "Claim not found", http.StatusUnauthorized)
			return
		}

		// Mengonversi klaim "iat" ke float64
		start, ok := scientificNotation.(float64)
		if !ok {
			http.Error(res, "Invalid claim type", http.StatusUnauthorized)
			return
		}

		// Konversi epoch timestamp (notasi ilmiah) ke waktu
		startTime := time.Unix(int64(start), 0)

		// Generate the monthly interval for the start time
		intervals := daterange.GenerateMonthlyInterval(startTime)

		// Generate key for monthly access count based on path, X-Auth-Key, and current year-month
		key := fmt.Sprintf("access_count:%s:%s:%s", req.URL.Path, apiKey, intervals["end"])

		// Get the current time
		now := time.Now()

		// Parse the 'end' date from intervals
		endDate, err := time.Parse("2006-01-02", intervals["end"])
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// If current time is after the end date, flush Redis count
		if now.After(endDate) {
			// Perform the flush operation only once
			go func() {
				// Check if the key exists in Redis
				exists, err := RedisClient.Exists(key).Result()
				if err != nil {
					log.Printf("Error checking if key exists in Redis: %v", err)
					return
				}

				// If the key exists, flush it
				if exists > 0 {
					err := RedisClient.Del(key).Err()
					if err != nil {
						log.Printf("Error flushing Redis key: %v", err)
					}
				}
			}()
		}

		// Mutex untuk memastikan operasi inkrement dan get dari Redis bersifat atomic
		mu.Lock()
		defer mu.Unlock()

		// Mengambil jumlah akses bulanan dari Redis
		count, err := RedisClient.Get(key).Int()
		if err != nil && err != redis.Nil {
			if Debug {
				log.Printf("Error getting access count from Redis: %v", err)
			}
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Jika jumlah akses belum ada, set ke 0
		if err == redis.Nil {
			count = 0
		}

		// Jika jumlah akses melebihi batas, kembalikan 429 Too Many Requests
		if count >= maxRequests {
			http.Error(res, fmt.Sprintf("Kuota Sudah Habis. Terpakai: %d, Sisa: %d", count, maxRequests-count), http.StatusTooManyRequests)
			return
		}

		// Inkrement jumlah akses bulanan
		err = RedisClient.Incr(key).Err()
		if err != nil {
			log.Printf("Error incrementing access count: %v", err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Tambahkan header dengan jumlah kuota yang digunakan dan sisa kuota
		res.Header().Set("X-Quota-Used", strconv.Itoa(count+1))                  // +1 karena kita sudah increment
		res.Header().Set("X-Quota-Remaining", strconv.Itoa(maxRequests-count-1)) // Sisa kuota setelah increment
		if Debug {
			// Logging untuk memeriksa nilai header
			log.Printf("Setting headers: X-Quota-Used = %d, X-Quota-Remaining = %d", count+1, maxRequests-count-1)
		}

		// Lanjutkan ke handler berikutnya
		next(res, req.WithContext(req.Context()))
		if Debug {
			// Logging response time
			log.Printf("Request processed in %v", time.Since(startTime))
		}
	}
}
