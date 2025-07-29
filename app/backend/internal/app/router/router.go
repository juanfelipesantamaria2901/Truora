package router

import (
	"truora-backend/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the Gin router
func SetupRouter(stockHandler *handlers.StockHandler) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// Health check endpoint
	r.GET("/health", stockHandler.HealthCheck)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Stock routes
		stocks := v1.Group("/stocks")
		{
			stocks.GET("", stockHandler.GetStocks)                // GET /api/v1/stocks
			stocks.GET("/:ticker", stockHandler.GetStockByTicker) // GET /api/v1/stocks/:ticker
			stocks.POST("/fetch", stockHandler.FetchStocks)       // POST /api/v1/stocks/fetch
		}

		// Recommendation routes
		recommendations := v1.Group("/recommendations")
		{
			recommendations.GET("", stockHandler.GetRecommendations)             // GET /api/v1/recommendations
			recommendations.POST("/generate", stockHandler.GenerateRecommendations) // POST /api/v1/recommendations/generate
		}
	}

	return r
}

// corsMiddleware handles CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
