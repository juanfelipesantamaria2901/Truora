# Truora Stock API

A professional Go API for fetching, storing, and analyzing stock data with intelligent investment recommendations.

## Features

- **Stock Data Management**: Fetch and store stock data from external APIs
- **CockroachDB Integration**: Robust database storage with GORM ORM
- **Investment Recommendations**: AI-powered stock analysis and recommendations
- **RESTful API**: Clean, well-documented REST endpoints
- **Security**: Parameterized queries to prevent SQL injection
- **Performance**: Optimized database queries with proper indexing

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP web framework)
- **Database**: CockroachDB with PostgreSQL driver
- **ORM**: GORM
- **Environment**: godotenv for configuration

## Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── app/
│   │   ├── handlers/            # HTTP handlers
│   │   └── router/              # Route configuration
│   ├── pkg/
│   │   ├── models/              # Database models
│   │   ├── repository/          # Data access layer
│   │   └── service/             # Business logic
│   └── platform/
│       └── cockroachdb/         # Database connection and migrations
├── api/
│   └── openapi.yaml            # API specification
├── go.mod                      # Go module dependencies
├── .env.example               # Environment variables template
└── README.md                  # This file
```

## Setup Instructions

### Prerequisites

1. **Go 1.21+** installed
2. **CockroachDB** running locally or accessible remotely
3. **Git** for version control

### Installation

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd truora/app/backend
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Setup environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. **Setup CockroachDB**:
   ```bash
   # Start CockroachDB (if running locally)
   cockroach start-single-node --insecure --listen-addr=localhost:26257
   
   # Create database
   cockroach sql --insecure -e "CREATE DATABASE truora_stocks;"
   ```

5. **Run the application**:
   ```bash
   go run cmd/api/main.go
   ```

The API will be available at `http://localhost:8080`

## API Endpoints

### Health Check
- **GET** `/health` - Service health status

### Stock Management
- **GET** `/api/v1/stocks` - List all stocks with pagination
  - Query params: `limit`, `offset`, `q` (search)
- **GET** `/api/v1/stocks/:symbol` - Get specific stock by symbol
- **POST** `/api/v1/stocks/fetch` - Fetch and store stocks from external API

### Recommendations
- **GET** `/api/v1/recommendations` - Get top stock recommendations
  - Query params: `limit`
- **POST** `/api/v1/recommendations/generate` - Generate new recommendations

## Usage Examples

### 1. Fetch Stock Data
```bash
# Fetch all stock data from external API
curl -X POST http://localhost:8080/api/v1/stocks/fetch
```

### 2. Get All Stocks
```bash
# Get first 20 stocks
curl http://localhost:8080/api/v1/stocks

# Get stocks with pagination
curl "http://localhost:8080/api/v1/stocks?limit=50&offset=100"

# Search stocks
curl "http://localhost:8080/api/v1/stocks?q=AAPL"
```

### 3. Get Specific Stock
```bash
curl http://localhost:8080/api/v1/stocks/AAPL
```

### 4. Generate Recommendations
```bash
# Generate new recommendations
curl -X POST http://localhost:8080/api/v1/recommendations/generate

# Get top 10 recommendations
curl http://localhost:8080/api/v1/recommendations?limit=10
```

## Stock Recommendation Algorithm

The recommendation system analyzes multiple factors:

1. **Price Momentum (30% weight)**
   - Recent price changes and trends
   - Positive momentum increases score

2. **Valuation Analysis (25% weight)**
   - P/E ratio evaluation
   - Attractive valuations (P/E < 15) get higher scores

3. **Dividend Yield (15% weight)**
   - Higher dividend yields increase attractiveness
   - Provides income stability indicator

4. **Market Capitalization (15% weight)**
   - Large cap: Stability (lower risk)
   - Mid cap: Growth potential
   - Small cap: High growth, higher risk

5. **52-Week Performance (15% weight)**
   - Position relative to 52-week high/low
   - Stocks near 52-week lows may represent value

6. **Additional Factors**
   - Trading volume (liquidity)
   - Sector preferences (Technology, Healthcare)

Scores are normalized to 0-100, with recommendations generated for stocks scoring ≥50.

## Database Schema

### Stocks Table
- Primary stock information
- Financial metrics (P/E, dividend yield, market cap)
- 52-week high/low data
- Sector and industry classification

### Stock Recommendations Table
- Generated recommendation scores
- Reasoning and risk assessment
- Time horizon indicators
- Links to stock records

## Security Features

- **Parameterized Queries**: All database queries use GORM's parameterized approach
- **Input Validation**: Request validation and sanitization
- **CORS Configuration**: Proper cross-origin resource sharing setup
- **Error Handling**: Comprehensive error handling without information leakage

## Performance Optimizations

- **Database Indexing**: Strategic indexes on frequently queried columns
- **Batch Operations**: Bulk inserts for large datasets
- **Connection Pooling**: Efficient database connection management
- **Pagination**: Limit result sets to prevent memory issues

## Environment Variables

| Variable | Description | Default |
|----------|-------------|----------|
| `PORT` | Server port | `8080` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `26257` |
| `DB_USER` | Database user | `root` |
| `DB_PASSWORD` | Database password | `` |
| `DB_NAME` | Database name | `truora_stocks` |
| `DB_SSLMODE` | SSL mode | `require` |
| `STOCK_API_URL` | External API URL | (provided) |
| `STOCK_API_KEY` | External API key | (provided) |

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o truora-api cmd/api/main.go
```

### Docker Support
```dockerfile
# Example Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License.