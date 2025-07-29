package handlers

import (
	"net/http"
	"strconv"
	"truora-backend/internal/pkg/service"

	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	stockService service.StockService
}

// NewStockHandler creates a new stock handler
func NewStockHandler(stockService service.StockService) *StockHandler {
	return &StockHandler{
		stockService: stockService,
	}
}

// GetStocks handles GET /api/stocks
func (h *StockHandler) GetStocks(c *gin.Context) {
	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	query := c.Query("q")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var stocks interface{}
	var serviceErr error

	// If query parameter is provided, search stocks
	if query != "" {
		stocks, serviceErr = h.stockService.SearchStocks(query, limit, offset)
	} else {
		stocks, serviceErr = h.stockService.GetAllStocks(limit, offset)
	}

	if serviceErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve stocks",
			"details": serviceErr.Error(),
		})
		return
	}

	// Get total count for pagination
	totalCount, err := h.stockService.GetStockCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get stock count",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stocks,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"total":  totalCount,
		},
	})
}

// GetStockByTicker handles GET /api/stocks/:ticker
func (h *StockHandler) GetStockByTicker(c *gin.Context) {
	ticker := c.Param("ticker")
	if ticker == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ticker parameter is required",
		})
		return
	}

	stock, err := h.stockService.GetByTicker(ticker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve stock",
			"details": err.Error(),
		})
		return
	}

	if stock == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Stock not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stock,
	})
}

// FetchStocks handles POST /api/stocks/fetch
func (h *StockHandler) FetchStocks(c *gin.Context) {
	err := h.stockService.FetchAndStoreStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch and store stocks",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stocks fetched and stored successfully",
	})
}

// GenerateRecommendations handles POST /api/recommendations/generate
func (h *StockHandler) GenerateRecommendations(c *gin.Context) {
	err := h.stockService.GenerateRecommendations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate recommendations",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Recommendations generated successfully",
	})
}

// GetRecommendations handles GET /api/recommendations
func (h *StockHandler) GetRecommendations(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	recommendations, err := h.stockService.GetTopRecommendations(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve recommendations",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  recommendations,
		"count": len(recommendations),
	})
}

// HealthCheck handles GET /api/health
func (h *StockHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "truora-stock-api",
		"timestamp": gin.H{
			"unix": gin.H{},
		},
	})
}
