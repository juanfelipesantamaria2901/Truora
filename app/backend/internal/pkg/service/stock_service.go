package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"truora-backend/internal/pkg/models"
	"truora-backend/internal/pkg/repository"
)

type StockService interface {
	FetchAndStoreStocks() error
	GetAllStocks(limit, offset int) ([]models.Stock, error)
	GetByTicker(ticker string) (*models.Stock, error)
	SearchStocks(query string, limit, offset int) ([]models.Stock, error)
	GenerateRecommendations() error
	GetTopRecommendations(limit int) ([]models.StockRecommendation, error)
	GetStockCount() (int64, error)
}

type stockService struct {
	repo   repository.StockRepository
	apiURL string
	apiKey string
}

// ExternalStockData represents the structure of data from external API
type ExternalStockData struct {
	Ticker     string `json:"ticker"`
	TargetFrom string `json:"target_from"`
	TargetTo   string `json:"target_to"`
	Company    string `json:"company"`
	Action     string `json:"action"`
	Brokerage  string `json:"brokerage"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	Time       string `json:"time"`
}

// ExternalAPIResponse represents the API response structure
type ExternalAPIResponse struct {
	Items    []ExternalStockData `json:"items"`
	NextPage string              `json:"next_page"`
}

// NewStockService creates a new stock service
func NewStockService(repo repository.StockRepository, apiURL, apiKey string) StockService {
	return &stockService{
		repo:   repo,
		apiURL: apiURL,
		apiKey: apiKey,
	}
}

// FetchAndStoreStocks fetches stocks from external API and stores them
func (s *stockService) FetchAndStoreStocks() error {
	// Create HTTP client with API key
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", s.apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key to headers if available
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	// Make HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch stocks from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Parse response
	var apiResponse ExternalAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return fmt.Errorf("failed to decode API response: %w", err)
	}

	// Convert and store stocks
	var stocks []models.Stock
	for _, stockData := range apiResponse.Items {
		// Parse time
		parsedTime, err := time.Parse(time.RFC3339, stockData.Time)
		if err != nil {
			parsedTime = time.Now()
		}

		stock := models.Stock{
			Ticker:      stockData.Ticker,
			Company:     stockData.Company,
			TargetFrom:  stockData.TargetFrom,
			TargetTo:    stockData.TargetTo,
			Action:      stockData.Action,
			Brokerage:   stockData.Brokerage,
			RatingFrom:  stockData.RatingFrom,
			RatingTo:    stockData.RatingTo,
			Time:        parsedTime,
			LastUpdated: time.Now(),
		}
		stocks = append(stocks, stock)
	}

	// Bulk create stocks
	return s.repo.BulkCreate(stocks)
}

// fetchStockPage fetches a single page of stock data
func (s *stockService) fetchStockPage(nextPage string) ([]models.Stock, *string, error) {
	url := s.apiURL
	if nextPage != "" {
		url += "?next_page=" + nextPage
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResponse ExternalAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	stocks := make([]models.Stock, len(apiResponse.Items))
	for i, stockData := range apiResponse.Items {
		stocks[i] = models.Stock{
			Ticker:      stockData.Ticker,
			Company:     stockData.Company,
			TargetFrom:  stockData.TargetFrom,
			TargetTo:    stockData.TargetTo,
			Action:      stockData.Action,
			Brokerage:   stockData.Brokerage,
			RatingFrom:  stockData.RatingFrom,
			RatingTo:    stockData.RatingTo,
			Time:        time.Now(),
			LastUpdated: time.Now(),
		}
	}

	nextPageValue := apiResponse.NextPage
	var nextPagePtr *string
	if nextPageValue != "" {
		nextPagePtr = &nextPageValue
	}
	return stocks, nextPagePtr, nil
}

// GetAllStocks retrieves all stocks with pagination
func (s *stockService) GetAllStocks(limit, offset int) ([]models.Stock, error) {
	return s.repo.GetAll(limit, offset)
}

// GetByTicker retrieves a stock by its ticker
func (s *stockService) GetByTicker(ticker string) (*models.Stock, error) {
	return s.repo.GetByTicker(ticker)
}

// SearchStocks searches stocks by query
func (s *stockService) SearchStocks(query string, limit, offset int) ([]models.Stock, error) {
	return s.repo.SearchStocks(query, limit, offset)
}

// GetStockCount returns the total number of stocks
func (s *stockService) GetStockCount() (int64, error) {
	return s.repo.GetStockCount()
}

// GenerateRecommendations generates stock recommendations based on analyst ratings and actions
func (s *stockService) GenerateRecommendations() error {
	log.Println("Generating stock recommendations...")

	// Get all stocks
	stocks, err := s.repo.GetAll(10000, 0) // Get a large number to analyze all
	if err != nil {
		return fmt.Errorf("failed to get stocks for analysis: %w", err)
	}

	if len(stocks) == 0 {
		return fmt.Errorf("no stocks available for analysis")
	}

	// Group stocks by ticker to analyze multiple analyst opinions
	stockGroups := make(map[string][]models.Stock)
	for _, stock := range stocks {
		stockGroups[stock.Ticker] = append(stockGroups[stock.Ticker], stock)
	}

	// Generate recommendations for each ticker
	for ticker, tickerStocks := range stockGroups {
		if len(tickerStocks) == 0 {
			continue
		}

		// Use the most recent stock data
		latestStock := tickerStocks[0]
		for _, stock := range tickerStocks {
			if stock.Time.After(latestStock.Time) {
				latestStock = stock
			}
		}

		score := s.calculateRecommendationScore(tickerStocks)
		riskLevel := s.calculateRiskLevel(tickerStocks)
		expectedReturn := s.calculateExpectedReturn(tickerStocks)
		reason := s.generateReason(tickerStocks, score)
		sentiment := s.calculateAnalystSentiment(tickerStocks)
		upgradeCount, downgradeCount := s.countUpgradesDowngrades(tickerStocks)

		recommendation := &models.StockRecommendation{
			StockID:             latestStock.ID,
			RecommendationScore: score,
			RiskLevel:           riskLevel,
			ExpectedReturn:      expectedReturn,
			TimeHorizon:         "medium",
			Reason:              reason,
			AnalystSentiment:    sentiment,
			UpgradeCount:        upgradeCount,
			DowngradeCount:      downgradeCount,
		}

		if err := s.repo.CreateRecommendation(recommendation); err != nil {
			log.Printf("Failed to store recommendation for %s: %v", ticker, err)
		}
	}

	log.Printf("Generated recommendations for %d tickers", len(stockGroups))
	return nil
}

// calculateRecommendationScore calculates a recommendation score based on analyst ratings and actions
func (s *stockService) calculateRecommendationScore(stocks []models.Stock) float64 {
	score := 50.0 // Base score
	upgradeCount := 0
	downgradeCount := 0
	buyRatings := 0
	sellRatings := 0
	holdRatings := 0

	// Analyze all analyst actions for this ticker
	for _, stock := range stocks {
		// Count upgrades and downgrades
		if strings.Contains(strings.ToLower(stock.Action), "upgrade") {
			upgradeCount++
		} else if strings.Contains(strings.ToLower(stock.Action), "downgrade") {
			downgradeCount++
		}

		// Analyze target ratings
		ratingTo := strings.ToLower(stock.RatingTo)
		if strings.Contains(ratingTo, "buy") || strings.Contains(ratingTo, "outperform") || strings.Contains(ratingTo, "strong buy") {
			buyRatings++
		} else if strings.Contains(ratingTo, "sell") || strings.Contains(ratingTo, "underperform") || strings.Contains(ratingTo, "strong sell") {
			sellRatings++
		} else if strings.Contains(ratingTo, "hold") || strings.Contains(ratingTo, "neutral") {
			holdRatings++
		}
	}

	// Calculate score based on analyst sentiment (60% weight)
	if upgradeCount > downgradeCount {
		score += float64(upgradeCount-downgradeCount) * 10
	} else if downgradeCount > upgradeCount {
		score -= float64(downgradeCount-upgradeCount) * 10
	}

	// Rating distribution (40% weight)
	totalRatings := buyRatings + sellRatings + holdRatings
	if totalRatings > 0 {
		buyPercentage := float64(buyRatings) / float64(totalRatings)
		sellPercentage := float64(sellRatings) / float64(totalRatings)
		score += (buyPercentage - sellPercentage) * 30
	}

	// Ensure score is within bounds
	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
	}

	return score
}

// calculateRiskLevel determines risk level based on analyst consensus
func (s *stockService) calculateRiskLevel(stocks []models.Stock) string {
	upgradeCount, downgradeCount := s.countUpgradesDowngrades(stocks)
	totalActions := upgradeCount + downgradeCount

	if totalActions == 0 {
		return "medium"
	}

	upgradeRatio := float64(upgradeCount) / float64(totalActions)
	if upgradeRatio > 0.7 {
		return "low"
	} else if upgradeRatio < 0.3 {
		return "high"
	}
	return "medium"
}

// calculateExpectedReturn estimates expected return based on target prices
func (s *stockService) calculateExpectedReturn(stocks []models.Stock) float64 {
	if len(stocks) == 0 {
		return 0.0
	}

	// Simple estimation based on analyst sentiment
	upgradeCount, downgradeCount := s.countUpgradesDowngrades(stocks)
	if upgradeCount > downgradeCount {
		return float64(upgradeCount-downgradeCount) * 2.5 // 2.5% per net upgrade
	} else if downgradeCount > upgradeCount {
		return float64(upgradeCount-downgradeCount) * 2.5 // Negative return
	}
	return 5.0 // Default 5% expected return
}

// generateReason creates a human-readable reason for the recommendation
func (s *stockService) generateReason(stocks []models.Stock, score float64) string {
	upgradeCount, downgradeCount := s.countUpgradesDowngrades(stocks)
	buyCount, sellCount, _ := s.countRatings(stocks)

	reasons := []string{}

	if upgradeCount > downgradeCount {
		reasons = append(reasons, fmt.Sprintf("%d upgrades vs %d downgrades", upgradeCount, downgradeCount))
	} else if downgradeCount > upgradeCount {
		reasons = append(reasons, fmt.Sprintf("%d downgrades vs %d upgrades", downgradeCount, upgradeCount))
	}

	if buyCount > sellCount {
		reasons = append(reasons, fmt.Sprintf("%d buy ratings vs %d sell ratings", buyCount, sellCount))
	} else if sellCount > buyCount {
		reasons = append(reasons, fmt.Sprintf("%d sell ratings vs %d buy ratings", sellCount, buyCount))
	}

	if score >= 70 {
		reasons = append(reasons, "Strong analyst consensus")
	} else if score <= 30 {
		reasons = append(reasons, "Weak analyst sentiment")
	}

	if len(reasons) == 0 {
		return "Mixed analyst opinions"
	}

	return strings.Join(reasons, "; ")
}

// calculateAnalystSentiment determines overall analyst sentiment
func (s *stockService) calculateAnalystSentiment(stocks []models.Stock) string {
	upgradeCount, downgradeCount := s.countUpgradesDowngrades(stocks)
	buyCount, sellCount, _ := s.countRatings(stocks)

	if upgradeCount > downgradeCount && buyCount > sellCount {
		return "bullish"
	} else if downgradeCount > upgradeCount && sellCount > buyCount {
		return "bearish"
	}
	return "neutral"
}

// countUpgradesDowngrades counts upgrade and downgrade actions
func (s *stockService) countUpgradesDowngrades(stocks []models.Stock) (int, int) {
	upgradeCount := 0
	downgradeCount := 0

	for _, stock := range stocks {
		action := strings.ToLower(stock.Action)
		if strings.Contains(action, "upgrade") {
			upgradeCount++
		} else if strings.Contains(action, "downgrade") {
			downgradeCount++
		}
	}

	return upgradeCount, downgradeCount
}

// countRatings counts buy, sell, and hold ratings
func (s *stockService) countRatings(stocks []models.Stock) (int, int, int) {
	buyCount := 0
	sellCount := 0
	holdCount := 0

	for _, stock := range stocks {
		ratingTo := strings.ToLower(stock.RatingTo)
		if strings.Contains(ratingTo, "buy") || strings.Contains(ratingTo, "outperform") || strings.Contains(ratingTo, "strong buy") {
			buyCount++
		} else if strings.Contains(ratingTo, "sell") || strings.Contains(ratingTo, "underperform") || strings.Contains(ratingTo, "strong sell") {
			sellCount++
		} else if strings.Contains(ratingTo, "hold") || strings.Contains(ratingTo, "neutral") {
			holdCount++
		}
	}

	return buyCount, sellCount, holdCount
}

// GetTopRecommendations retrieves top stock recommendations
func (s *stockService) GetTopRecommendations(limit int) ([]models.StockRecommendation, error) {
	return s.repo.GetTopRecommendations(limit)
}
