package repository

import (
	"fmt"
	"truora-backend/internal/pkg/models"

	"gorm.io/gorm"
)

type StockRepository interface {
	Create(stock *models.Stock) error
	GetByTicker(ticker string) (*models.Stock, error)
	GetAll(limit, offset int) ([]models.Stock, error)
	Update(stock *models.Stock) error
	Delete(id uint) error
	BulkCreate(stocks []models.Stock) error
	GetTopRecommendations(limit int) ([]models.StockRecommendation, error)
	CreateRecommendation(recommendation *models.StockRecommendation) error
	GetStockCount() (int64, error)
	SearchStocks(query string, limit, offset int) ([]models.Stock, error)
}

type stockRepository struct {
	db *gorm.DB
}

// NewStockRepository creates a new stock repository
func NewStockRepository(db *gorm.DB) StockRepository {
	return &stockRepository{db: db}
}

// Create creates a new stock record
func (r *stockRepository) Create(stock *models.Stock) error {
	if err := r.db.Create(stock).Error; err != nil {
		return fmt.Errorf("failed to create stock: %w", err)
	}
	return nil
}

// GetByTicker retrieves a stock by its ticker
func (r *stockRepository) GetByTicker(ticker string) (*models.Stock, error) {
	var stock models.Stock
	if err := r.db.Where("ticker = ?", ticker).First(&stock).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get stock by ticker: %w", err)
	}
	return &stock, nil
}

// GetAll retrieves all stocks with pagination
func (r *stockRepository) GetAll(limit, offset int) ([]models.Stock, error) {
	var stocks []models.Stock
	query := r.db.Limit(limit).Offset(offset)
	
	if err := query.Find(&stocks).Error; err != nil {
		return nil, fmt.Errorf("failed to get all stocks: %w", err)
	}
	return stocks, nil
}

// Update updates an existing stock record
func (r *stockRepository) Update(stock *models.Stock) error {
	if err := r.db.Save(stock).Error; err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}
	return nil
}

// Delete soft deletes a stock record
func (r *stockRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.Stock{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete stock: %w", err)
	}
	return nil
}

// BulkCreate creates multiple stock records in a single transaction
func (r *stockRepository) BulkCreate(stocks []models.Stock) error {
	if len(stocks) == 0 {
		return nil
	}

	// Use batch insert for better performance
	batchSize := 100
	for i := 0; i < len(stocks); i += batchSize {
		end := i + batchSize
		if end > len(stocks) {
			end = len(stocks)
		}
		
		batch := stocks[i:end]
		if err := r.db.CreateInBatches(batch, batchSize).Error; err != nil {
			return fmt.Errorf("failed to bulk create stocks: %w", err)
		}
	}
	return nil
}

// GetTopRecommendations retrieves top stock recommendations
func (r *stockRepository) GetTopRecommendations(limit int) ([]models.StockRecommendation, error) {
	var recommendations []models.StockRecommendation
	if err := r.db.Preload("Stock").Order("recommendation_score DESC").Limit(limit).Find(&recommendations).Error; err != nil {
		return nil, fmt.Errorf("failed to get top recommendations: %w", err)
	}
	return recommendations, nil
}

// CreateRecommendation creates a new stock recommendation
func (r *stockRepository) CreateRecommendation(recommendation *models.StockRecommendation) error {
	if err := r.db.Create(recommendation).Error; err != nil {
		return fmt.Errorf("failed to create recommendation: %w", err)
	}
	return nil
}

// GetStockCount returns the total number of stocks
func (r *stockRepository) GetStockCount() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Stock{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to get stock count: %w", err)
	}
	return count, nil
}

// SearchStocks searches stocks by symbol or company name
func (r *stockRepository) SearchStocks(query string, limit, offset int) ([]models.Stock, error) {
	var stocks []models.Stock
	searchQuery := "%" + query + "%"
	
	if err := r.db.Where("ticker ILIKE ? OR company ILIKE ?", searchQuery, searchQuery).
		Limit(limit).Offset(offset).Find(&stocks).Error; err != nil {
		return nil, fmt.Errorf("failed to search stocks: %w", err)
	}
	return stocks, nil
}