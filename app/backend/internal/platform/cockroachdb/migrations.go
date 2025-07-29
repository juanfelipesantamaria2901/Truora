package cockroachdb

import (
	"fmt"
	"log"
	"truora-backend/internal/pkg/models"
)

// RunMigrations runs database migrations
func RunMigrations(db *Database) error {
	// Auto-migrate models
	if err := db.DB.AutoMigrate(&models.Stock{}, &models.StockRecommendation{}); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Create indexes for better performance
	if err := createIndexes(db); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// createIndexes creates database indexes for better performance
func createIndexes(db *Database) error {
	// Index on ticker for fast lookups
	if err := db.DB.Exec("CREATE INDEX IF NOT EXISTS idx_stocks_ticker ON stocks(ticker)").Error; err != nil {
		return fmt.Errorf("failed to create ticker index: %w", err)
	}

	// Index on company for search
	if err := db.DB.Exec("CREATE INDEX IF NOT EXISTS idx_stocks_company ON stocks(company)").Error; err != nil {
		return fmt.Errorf("failed to create company index: %w", err)
	}

	// Index on brokerage for filtering
	if err := db.DB.Exec("CREATE INDEX IF NOT EXISTS idx_stocks_brokerage ON stocks(brokerage)").Error; err != nil {
		return fmt.Errorf("failed to create brokerage index: %w", err)
	}

	// Index on action for filtering upgrades/downgrades
	if err := db.DB.Exec("CREATE INDEX IF NOT EXISTS idx_stocks_action ON stocks(action)").Error; err != nil {
		return fmt.Errorf("failed to create action index: %w", err)
	}

	// Index on rating_to for filtering by ratings
	if err := db.DB.Exec("CREATE INDEX IF NOT EXISTS idx_stocks_rating_to ON stocks(rating_to)").Error; err != nil {
		return fmt.Errorf("failed to create rating_to index: %w", err)
	}

	// Index on time for chronological queries
	if err := db.DB.Exec("CREATE INDEX IF NOT EXISTS idx_stocks_time ON stocks(time DESC)").Error; err != nil {
		return fmt.Errorf("failed to create time index: %w", err)
	}

	// Index on last_updated for data freshness queries
	if err := db.DB.Exec("CREATE INDEX IF NOT EXISTS idx_stocks_last_updated ON stocks(last_updated)").Error; err != nil {
		return fmt.Errorf("failed to create last_updated index: %w", err)
	}

	// Index on recommendation score for top recommendations
	if err := db.DB.Exec("CREATE INDEX IF NOT EXISTS idx_recommendations_score ON stock_recommendations(recommendation_score DESC)").Error; err != nil {
		return fmt.Errorf("failed to create recommendation_score index: %w", err)
	}

	// Composite index for text search
	if err := db.DB.Exec("CREATE INDEX IF NOT EXISTS idx_stocks_search ON stocks(company, ticker)").Error; err != nil {
		return fmt.Errorf("failed to create search index: %w", err)
	}

	// Composite index for ticker and time (for grouping latest analyst opinions)
	if err := db.DB.Exec("CREATE INDEX IF NOT EXISTS idx_stocks_ticker_time ON stocks(ticker, time DESC)").Error; err != nil {
		return fmt.Errorf("failed to create ticker_time index: %w", err)
	}

	return nil
}

// DropTables drops all tables (useful for testing)
func (d *Database) DropTables() error {
	log.Println("Dropping all tables...")

	err := d.DB.Migrator().DropTable(
		&models.StockRecommendation{},
		&models.Stock{},
	)

	if err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	log.Println("All tables dropped successfully")
	return nil
}

// ResetDatabase drops and recreates all tables
func (d *Database) ResetDatabase() error {
	log.Println("Resetting database...")

	if err := d.DropTables(); err != nil {
		return err
	}

	if err := RunMigrations(d); err != nil {
		return err
	}

	log.Println("Database reset completed successfully")
	return nil
}
