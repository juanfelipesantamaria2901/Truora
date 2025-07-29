package models

import (
	"time"

	"gorm.io/gorm"
)

// Stock represents a stock entity in the database
type Stock struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Ticker      string         `json:"ticker" gorm:"uniqueIndex;not null;size:10"`
	Company     string         `json:"company" gorm:"not null;size:255"`
	TargetFrom  string         `json:"target_from" gorm:"size:20"`
	TargetTo    string         `json:"target_to" gorm:"size:20"`
	Action      string         `json:"action" gorm:"size:50"`
	Brokerage   string         `json:"brokerage" gorm:"size:255"`
	RatingFrom  string         `json:"rating_from" gorm:"size:50"`
	RatingTo    string         `json:"rating_to" gorm:"size:50"`
	Time        time.Time      `json:"time"`
	LastUpdated time.Time      `json:"last_updated" gorm:"autoUpdateTime"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// StockRecommendation represents a stock recommendation
type StockRecommendation struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	StockID           uint           `json:"stock_id" gorm:"not null;index"`
	Stock             Stock          `json:"stock" gorm:"foreignKey:StockID"`
	RecommendationScore float64      `json:"recommendation_score" gorm:"not null;type:decimal(5,2);index"`
	RiskLevel         string         `json:"risk_level" gorm:"not null;size:20"`
	ExpectedReturn    float64        `json:"expected_return" gorm:"type:decimal(5,2)"`
	TimeHorizon       string         `json:"time_horizon" gorm:"size:20"`
	Reason            string         `json:"reason" gorm:"type:text"`
	AnalystSentiment  string         `json:"analyst_sentiment" gorm:"size:20"`
	UpgradeCount      int            `json:"upgrade_count" gorm:"default:0"`
	DowngradeCount    int            `json:"downgrade_count" gorm:"default:0"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName sets the table name for Stock
func (Stock) TableName() string {
	return "stocks"
}

// TableName sets the table name for StockRecommendation
func (StockRecommendation) TableName() string {
	return "stock_recommendations"
}