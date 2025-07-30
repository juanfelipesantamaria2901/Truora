package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"truora-backend/internal/pkg/repository"
	"truora-backend/internal/pkg/service"
	"truora-backend/internal/platform/cockroachdb"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database connection
	db, err := cockroachdb.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository and service
	stockRepo := repository.NewStockRepository(db.DB)
	apiURL := getEnv("STOCK_API_URL", "https://api")
	apiKey := getEnv("STOCK_API_KEY", "Bearer ")
	stockService := service.NewStockService(stockRepo, apiURL, apiKey)

	log.Println("Starting Truora Stock Worker...")

	// Create a channel to listen for interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Create a ticker for periodic tasks
	dataFetchInterval := getEnvDuration("DATA_FETCH_INTERVAL", 6*time.Hour)           // Default: every 6 hours
	recommendationInterval := getEnvDuration("RECOMMENDATION_INTERVAL", 24*time.Hour) // Default: daily

	dataFetchTicker := time.NewTicker(dataFetchInterval)
	recommendationTicker := time.NewTicker(recommendationInterval)

	defer dataFetchTicker.Stop()
	defer recommendationTicker.Stop()

	// Run initial tasks
	log.Println("Running initial data fetch...")
	if err := stockService.FetchAndStoreStocks(); err != nil {
		log.Printf("Initial data fetch failed: %v", err)
	} else {
		log.Println("Initial data fetch completed successfully")
	}

	log.Println("Generating initial recommendations...")
	if err := stockService.GenerateRecommendations(); err != nil {
		log.Printf("Initial recommendation generation failed: %v", err)
	} else {
		log.Println("Initial recommendations generated successfully")
	}

	log.Printf("Worker started. Data fetch interval: %v, Recommendation interval: %v", dataFetchInterval, recommendationInterval)

	// Main worker loop
	for {
		select {
		case <-dataFetchTicker.C:
			log.Println("Starting scheduled data fetch...")
			if err := stockService.FetchAndStoreStocks(); err != nil {
				log.Printf("Scheduled data fetch failed: %v", err)
			} else {
				log.Println("Scheduled data fetch completed successfully")
			}

		case <-recommendationTicker.C:
			log.Println("Starting scheduled recommendation generation...")
			if err := stockService.GenerateRecommendations(); err != nil {
				log.Printf("Scheduled recommendation generation failed: %v", err)
			} else {
				log.Println("Scheduled recommendations generated successfully")
			}

		case <-c:
			log.Println("Received interrupt signal, shutting down worker...")
			return
		}
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvDuration gets environment variable as duration with fallback
func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		log.Printf("Invalid duration format for %s: %s, using fallback", key, value)
	}
	return fallback
}
