package main

import (
	"log"
	"os"
	"truora-backend/internal/app/handlers"
	"truora-backend/internal/app/router"
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

	// Run migrations
	if err := cockroachdb.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repository
	stockRepo := repository.NewStockRepository(db.DB)

	// Initialize service
	apiURL := getEnv("STOCK_API_URL", "https://api.karenai.click/swechallenge/list")
	apiKey := getEnv("STOCK_API_KEY", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdHRlbXB0cyI6MSwiZW1haWwiOiJqZnNnMjkwMUBnbWFpbC5jb20iLCJleHAiOjE3NTMxNTA3MTYsImlkIjoiIiwicGFzc3dvcmQiOiJ1c2VyLyoqL0ZST00vKiovdXNlcnMvKiovV0hFUkUvKiovJzEnPScxJzstLSJ9.2fiNZtZsYux1n8hQf5oeTlrSQ0p0v2zYc6gkJmWli7k")
	stockService := service.NewStockService(stockRepo, apiURL, apiKey)

	// Initialize handlers
	stockHandler := handlers.NewStockHandler(stockService)

	// Setup router
	r := router.SetupRouter(stockHandler)

	// Get port from environment
	port := getEnv("PORT", "8000")

	log.Printf("Starting Truora Stock API server on port %s", port)
	log.Printf("Health check available at: http://localhost:%s/health", port)
	log.Printf("API documentation: http://localhost:%s/api/v1/stocks", port)

	// Start server
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
