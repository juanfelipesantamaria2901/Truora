# Truora Backend - Detailed Documentation

## 🏗️ Architecture Overview

The Truora backend is a production-ready Go REST API built with enterprise-grade patterns and best practices. It follows a clean architecture approach with clear separation of concerns across handlers, services, repositories, and domain models.

## 📁 Detailed Project Structure

```
backend/
├── cmd/
│   ├── api/
│   │   └── main.go              # Production API server entry point
│   └── worker/
│       └── main.go              # Background worker (future: Kafka consumer)
├── internal/
│   ├── app/
│   │   ├── handlers/            # HTTP request handlers (controllers)
│   │   │   └── stock_handler.go # RESTful stock endpoints
│   │   └── router/              # Route definitions and middleware
│   │       └── router.go        # Gin router configuration
│   ├── pkg/
│   │   ├── models/              # Database models (entities)
│   │   │   └── stock.go         # Stock and recommendation models
│   │   ├── repository/          # Data access layer (DAL)
│   │   │   └── stock_repository.go # Database operations
│   │   └── service/             # Business logic layer
│   │       └── stock_service.go   # Stock analysis and recommendations
│   └── platform/
│       ├── cockroachdb/         # Database infrastructure
│       │   ├── connection.go    # Database connection management
│       │   └── migrations.go      # Schema migrations and versioning
│       └── kafka/               # Message queue integration (future)
├── api/
│   └── openapi.yaml            # OpenAPI 3.1 specification
├── configs/                    # Configuration files (future)
├── go.mod                      # Go module dependencies
├── go.sum                      # Dependency checksums
├── .env.example               # Environment variables template
├── .env                       # Local environment (git-ignored)
├── Dockerfile                 # Production container image
├── Dockerfile.dev            # Development container image
└── README.md                 # Basic backend documentation
```

## 🎯 Core Components Deep Dive

### 1. Entry Point (`cmd/api/main.go`)

The main entry point orchestrates the entire application lifecycle:

```go
// Key responsibilities:
// 1. Environment configuration loading
// 2. Database connection establishment
// 3. Migration execution
// 4. Dependency injection setup
// 5. Server startup with graceful shutdown
```

**Features implemented:**
- Environment variable loading with `.env` support
- Database connection pooling with GORM
- Automatic database migrations
- Dependency injection pattern
- Graceful error handling
- Production-ready logging

### 2. Database Models (`internal/pkg/models/`)

#### Stock Model (`stock.go`)
```go
// Complete database schema:
type Stock struct {
    gorm.Model
    Symbol           string  `gorm:"uniqueIndex;not null"`
    CompanyName      string  `gorm:"not null"`
    CurrentPrice     float64 `gorm:"type:decimal(10,2)"`
    MarketCap        float64 `gorm:"type:decimal(20,2)"`
    PE               float64 `gorm:"type:decimal(10,2)"`
    DividendYield      float64 `gorm:"type:decimal(5,2)"`
    Week52High       float64 `gorm:"type:decimal(10,2)"`
    Week52Low        float64 `gorm:"type:decimal(10,2)"`
    Volume           int64
    Sector           string
    Industry         string
    
    // Relationships
    Recommendations  []StockRecommendation `gorm:"foreignKey:StockID"`
}

type StockRecommendation struct {
    gorm.Model
    StockID     uint    `gorm:"index"`
    Score       float64 `gorm:"type:decimal(5,2)"`
    Reason      string  `gorm:"type:text"`
    RiskLevel   string  `gorm:"type:varchar(20)"`
    TimeHorizon string  `gorm:"type:varchar(20)"`
}
```

**Database indexes:**
- `idx_stocks_symbol` - Unique index on stock symbol
- `idx_stocks_sector` - Index for sector filtering
- `idx_recommendations_stock_id` - Foreign key index
- `idx_recommendations_score` - Index for sorting recommendations

### 3. Repository Layer (`internal/pkg/repository/`)

#### Stock Repository (`stock_repository.go`)

Implements the data access layer with repository pattern:

```go
// Interface definition
type StockRepository interface {
    Create(stock *models.Stock) error
    CreateBatch(stocks []*models.Stock) error
    FindBySymbol(symbol string) (*models.Stock, error)
    FindAll(limit, offset int, search string) ([]models.Stock, int64, error)
    Update(stock *models.Stock) error
    Delete(id uint) error
    
    // Recommendation methods
    SaveRecommendation(rec *models.StockRecommendation) error
    GetTopRecommendations(limit int) ([]models.StockRecommendation, error)
    ClearRecommendations() error
}
```

**Advanced features:**
- Batch insertions for performance
- Pagination with total count
- Search across symbol and company name
- Transaction support
- Connection pooling optimization

### 4. Service Layer (`internal/pkg/service/`)

#### Stock Service (`stock_service.go`)

Contains the business logic for stock analysis and recommendations:

```go
// Core services:
- FetchStocksFromAPI()    // External API integration
- AnalyzeStocks()         // AI recommendation algorithm
- CalculateRecommendationScore() // Multi-factor scoring
- GenerateRecommendations() // Batch recommendation generation
```

#### AI Recommendation Algorithm

The recommendation system uses a sophisticated weighted scoring algorithm:

```go
// Scoring factors and weights:
const (
    WeightPriceMomentum     = 0.30  // Recent price trends
    WeightValuation         = 0.25  // P/E ratio analysis
    WeightDividend          = 0.15  // Dividend yield
    WeightMarketCap         = 0.15  // Company size stability
    Weight52WeekPerformance = 0.15  // Long-term performance
)

// Scoring methodology:
// 1. Price Momentum: RSI calculation, trend analysis
// 2. Valuation: P/E ratio compared to sector averages
// 3. Dividend: Yield sustainability and growth
// 4. Market Cap: Liquidity and stability factors
// 5. 52-Week: Position relative to trading range
```

### 5. HTTP Handlers (`internal/app/handlers/`)

#### Stock Handler (`stock_handler.go`)

RESTful API handlers following clean architecture:

```go
// Handler methods:
- GetStocks()        // GET /api/v1/stocks
- GetStock()         // GET /api/v1/stocks/:symbol
- FetchStocks()      // POST /api/v1/stocks/fetch
- GetRecommendations() // GET /api/v1/recommendations
- GenerateRecommendations() // POST /api/v1/recommendations/generate
```

**Response formats:**
- Standardized JSON responses
- Consistent error handling
- Pagination metadata
- HATEOAS links for navigation

### 6. Router Configuration (`internal/app/router/`)

#### Router Setup (`router.go`)

Production-ready router with middleware stack:

```go
// Middleware chain:
1. CORS handler
2. Request logging
3. Rate limiting
4. Request validation
5. Error recovery
6. Gzip compression
```

**Route definitions:**
```go
// Health check
GET /health

// Stock management
GET    /api/v1/stocks
GET    /api/v1/stocks/:symbol
POST   /api/v1/stocks/fetch

// Recommendations
GET    /api/v1/recommendations
POST   /api/v1/recommendations/generate
```

### 7. Database Infrastructure (`internal/platform/cockroachdb/`)

#### Connection Management (`connection.go`)

Production database configuration:

```go
// Connection pool settings:
- Max open connections: 25
- Max idle connections: 25
- Connection lifetime: 5 minutes
- SSL mode: Configurable (require/disable)
- Connection timeout: 30 seconds
```

#### Migrations (`migrations.go`)

Database schema versioning:

```go
// Migration strategy:
// 1. Automatic schema creation
// 2. Index creation for performance
// 3. Foreign key constraints
// 4. Data type optimizations
```

## 🚀 API Specification

### OpenAPI 3.1 Documentation

The complete API specification is available in `api/openapi.yaml`:

#### Stock Endpoints

**GET /api/v1/stocks**
```yaml
parameters:
  - name: limit
    schema: { type: integer, default: 20, max: 100 }
  - name: offset
    schema: { type: integer, default: 0 }
  - name: q
    schema: { type: string }
    description: Search in symbol and company name
```

**GET /api/v1/stocks/{symbol}**
```yaml
parameters:
  - name: symbol
    required: true
    schema: { type: string }
    example: AAPL
```

**POST /api/v1/stocks/fetch**
```yaml
responses:
  200:
    description: Successfully fetched and stored stocks
    content:
      application/json:
        schema:
          type: object
          properties:
            message: { type: string }
            count: { type: integer }
```

## 🔧 Environment Configuration

### Required Environment Variables

```bash
# Server Configuration
PORT=8000                                    # Server port

# Database Configuration
DB_HOST=localhost                           # Database host
DB_PORT=26257                               # CockroachDB port
DB_USER=root                                # Database user
DB_PASSWORD=                                # Database password
DB_NAME=truora_stocks                       # Database name
DB_SSLMODE=disable                          # SSL mode (disable/require)

# External API Configuration
STOCK_API_URL=https://api.karenai.click/swechallenge/list
STOCK_API_KEY=Bearer <your-api-token>
```

### Optional Environment Variables

```bash
# Logging
LOG_LEVEL=info                              # debug, info, warn, error

# Performance
MAX_OPEN_CONNS=25                          # Database connection pool
CONN_MAX_LIFETIME=5m                       # Connection lifetime

# Security
CORS_ALLOWED_ORIGINS=http://localhost:5173 # Frontend URL
RATE_LIMIT_REQUESTS=100                    # Rate limit per minute
```

## 🧪 Testing Strategy

### Unit Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/pkg/service
```

### Integration Tests

```bash
# Test with real database
go test -tags=integration ./...

# Test API endpoints
go test -tags=api ./...
```

### Load Testing

```bash
# Using hey for load testing
hey -n 1000 -c 10 http://localhost:8000/api/v1/stocks

# Using Apache Bench
ab -n 1000 -c 10 http://localhost:8000/api/v1/stocks
```

## 📊 Performance Optimization

### Database Optimization

**Indexes:**
```sql
-- Primary indexes
CREATE INDEX idx_stocks_symbol ON stocks(symbol);
CREATE INDEX idx_stocks_sector ON stocks(sector);
CREATE INDEX idx_stocks_market_cap ON stocks(market_cap);
CREATE INDEX idx_stocks_pe ON stocks(pe);

-- Recommendation indexes
CREATE INDEX idx_recommendations_stock_id ON stock_recommendations(stock_id);
CREATE INDEX idx_recommendations_score ON stock_recommendations(score DESC);
```

**Query Optimization:**
- Use of prepared statements
- Batch insertions for bulk data
- Pagination with keyset pagination (future)
- Connection pooling configuration

### API Optimization

**Response Caching:**
- In-memory caching for recommendations (future)
- Redis integration for distributed caching
- ETag support for conditional requests

**Compression:**
- Gzip compression enabled
- Response size optimization
- JSON minification in production

## 🔒 Security Implementation

### SQL Injection Prevention
- **Parameterized queries**: All queries use GORM with parameter binding
- **Input validation**: Request validation middleware
- **Type safety**: Strong typing for all database fields

### API Security
- **Rate limiting**: Per-IP rate limiting
- **CORS protection**: Configured for specific origins
- **Input sanitization**: XSS prevention
- **Error handling**: No sensitive information in error messages

### Authentication (Future)
- **JWT tokens**: Stateless authentication
- **API keys**: For programmatic access
- **OAuth 2.0**: Third-party integration

## 🚀 Deployment Strategies

### Docker Deployment

#### Production Dockerfile
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

#### Development Dockerfile
```dockerfile
FROM golang:1.23-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/api/main.go
CMD ["./main"]
```

### Cloud Deployment

#### Railway Deployment
```bash
# Using Railway CLI
railway login
railway init
railway up
```

#### Render Deployment
```bash
# Using render.yaml
services:
  - type: web
    name: truora-backend
    env: go
    buildCommand: go build -o main cmd/api/main.go
    startCommand: ./main
```

## 📈 Monitoring & Observability

### Logging
- **Structured logging**: JSON format logs
- **Log levels**: Debug, Info, Warn, Error
- **Request tracking**: Request ID middleware
- **Performance metrics**: Database query timing

### Health Checks
- **Database connectivity**: Periodic health checks
- **External API status**: API dependency monitoring
- **Memory usage**: Runtime metrics
- **Response times**: Performance monitoring

### Metrics (Future)
- **Prometheus integration**: Custom metrics
- **Grafana dashboards**: Visual monitoring
- **Alerting**: Slack/email notifications
- **Distributed tracing**: Jaeger integration

## 🔄 CI/CD Pipeline

### GitHub Actions

#### Backend CI
```yaml
name: Backend CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: 1.23
      - run: go mod download
      - run: go test ./...
      - run: go build cmd/api/main.go
```

#### Deployment
```yaml
name: Deploy to Production
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: docker/build-push-action@v4
      - uses: railway/cli@latest
        with:
          command: up
```

## 🆘 Troubleshooting Guide

### Common Issues

#### Database Connection Issues
```bash
# Check CockroachDB status
cockroach node status --insecure --host=localhost:26257

# Test database connection
go run cmd/api/main.go --test-db

# Reset database
cockroach sql --insecure -e "DROP DATABASE IF EXISTS truora_stocks; CREATE DATABASE truora_stocks;"
```

#### Memory Issues
```bash
# Check memory usage
go tool pprof http://localhost:8000/debug/pprof/heap

# Memory profiling
go test -memprofile=mem.prof -bench=.
```

#### Performance Issues
```bash
# Database query analysis
EXPLAIN ANALYZE SELECT * FROM stocks WHERE symbol = 'AAPL';

# API response time analysis
curl -w "@curl-format.txt" -o /dev/null http://localhost:8000/api/v1/stocks
```

### Debug Mode

#### Enable Debug Logging
```bash
export LOG_LEVEL=debug
go run cmd/api/main.go
```

#### Database Debug
```bash
# Enable SQL logging
export GORM_DEBUG=true
go run cmd/api/main.go
```

## 📚 Additional Resources

### Learning Resources
- [Go Best Practices](https://golang.org/doc/effective_go.html)
- [GORM Documentation](https://gorm.io/docs/)
- [Gin Web Framework](https://gin-gonic.com/docs/)
- [CockroachDB Documentation](https://www.cockroachlabs.com/docs/)

### Contributing Guidelines
- [Code Review Checklist](docs/code-review.md)
- [API Design Guidelines](docs/api-design.md)
- [Database Migration Guide](docs/migrations.md)
- [Performance Tuning Guide](docs/performance.md)

---

**Next Steps**: Continue to the [Frontend Detailed Documentation](../frontend/DETAILED_README.md) for complete application understanding.