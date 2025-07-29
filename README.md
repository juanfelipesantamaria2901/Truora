# Truora Stock Analysis Platform

A comprehensive full-stack application for hiring process in Truoara, the app makes stock analysis and investment recommendations. This platform combines a robust Go backend with a modern Vue.js frontend to provide professional-grade stock data management and AI-powered investment insights.

## 🚀 Project Overview

This App is designed for hiring process in Truora, the app makes stock analysis and investment recommendations.Use only for reference purposes, do not use for real investment decisions, the app is use false data for analysis.

## 🏗️ Architecture

```
truora/
├── app/
│   ├── backend/          # Go REST API (Port 8000)
│   └── frontend/         # Vue.js SPA (Port 5173)
├── infra/                # Terraform infrastructure
├── docker-compose.yaml   # Production deployment
├── docker-compose.dev.yaml # Development environment
└── nginx.conf           # Reverse proxy configuration
```

## 🎯 Key Features

### Backend (Go API)
- **Real-time Stock Data**: Integration with external financial APIs
- **AI Recommendations**: Sophisticated algorithm analyzing 6+ factors
- **RESTful Design**: Clean API with OpenAPI 3.1 specification
- **Database Management**: CockroachDB with GORM ORM
- **Security**: Parameterized queries, CORS, input validation
- **Performance**: Connection pooling, pagination, indexing

### Frontend (Vue.js)
- **Modern UI/UX**: Responsive design with Tailwind CSS
- **Real-time Updates**: Reactive state management with Pinia
- **Type Safety**: Full TypeScript support
- **Component Architecture**: Modular Vue 3 composition API
- **Routing**: Vue Router for SPA navigation
- **API Integration**: Axios-based service layer

## 🛠️ Technology Stack

### Backend Technologies
- **Language**: Go 1.23
- **Framework**: Gin Web Framework
- **Database**: CockroachDB (PostgreSQL-compatible)
- **ORM**: GORM v1.25.5
- **Configuration**: godotenv
- **Documentation**: OpenAPI 3.1 (Swagger)

### Frontend Technologies
- **Framework**: Vue.js 3.5.17
- **Language**: TypeScript 5.8
- **Build Tool**: Vite 7.0
- **State Management**: Pinia 3.0.3
- **Styling**: Tailwind CSS 4.1.11
- **Routing**: Vue Router 4.5.1

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Cloud**: Vercel (frontend), Railway/Render (backend)
- **Reverse Proxy**: Nginx
- **Infrastructure as Code**: Terraform

## 📁 Project Structure

### Root Level
```
truora/
├── app/                    # Main application code
├── infra/                  # Terraform infrastructure
├── docker-compose.yaml     # Production containers
├── docker-compose.dev.yaml # Development containers
├── nginx.conf             # Nginx reverse proxy
├── Dockerfile             # Production backend image
├── Dockerfile.dev         # Development backend image
├── start.sh              # Startup script
└── supervisord.conf      # Process manager
```

### Backend Structure
```
app/backend/
├── cmd/
│   ├── api/               # API server entry point
│   └── worker/            # Background worker (placeholder)
├── internal/              # Private application code
│   ├── app/
│   │   ├── handlers/      # HTTP request handlers
│   │   └── router/        # Route definitions
│   ├── pkg/
│   │   ├── models/        # Database models
│   │   ├── repository/    # Data access layer
│   │   └── service/       # Business logic
│   └── platform/
│       ├── cockroachdb/   # Database connection & migrations
│       └── kafka/         # Message queue (placeholder)
├── api/
│   └── openapi.yaml      # API specification
├── go.mod                # Go module dependencies
├── go.sum                # Dependency checksums
├── .env.example          # Environment template
└── README.md             # Backend documentation
```

### Frontend Structure
```
app/frontend/
├── src/
│   ├── components/        # Reusable Vue components
│   │   ├── Layout/        # Layout components (header, footer)
│   │   ├── Stocks/        # Stock-specific components
│   │   ├── UI/           # Generic UI components
│   │   └── icons/        # SVG icons
│   ├── views/            # Page components
│   ├── services/        # API service layer
│   ├── stores/          # Pinia state stores
│   ├── types/           # TypeScript type definitions
│   ├── router/          # Vue Router configuration
│   ├── assets/          # Static assets (CSS, images)
│   ├── main.ts          # Application entry point
│   └── App.vue          # Root component
├── public/              # Static files
├── package.json         # Node.js dependencies
├── vite.config.ts       # Vite configuration
├── tsconfig.json        # TypeScript configuration
├── tailwind.config.js   # Tailwind CSS configuration
├── postcss.config.js    # PostCSS configuration
└── README.md            # Frontend documentation
```

## 🚦 Quick Start

### Prerequisites
- **Docker & Docker Compose** (recommended)
- **Go 1.23+** (for local backend development)
- **Node.js 18+** (for local frontend development)
- **CockroachDB** (or PostgreSQL-compatible database)

### Development Setup

#### Option 1: Docker Development (Recommended)
```bash
# Clone repository
git clone <repository-url>
cd truora

# Start development environment
docker-compose -f docker-compose.dev.yaml up

# Services will be available at:
# - Frontend: http://localhost:5173
# - Backend API: http://localhost:8000
# - Health check: http://localhost:8000/health
```

#### Option 2: Local Development

**Backend Setup:**
```bash
cd app/backend
cp .env.example .env
# Edit .env with your configuration
go mod download

# Start CockroachDB
cockroach start-single-node --insecure --listen-addr=localhost:26257
cockroach sql --insecure -e "CREATE DATABASE truora_stocks;"

# Run backend
go run cmd/api/main.go
```

**Frontend Setup:**
```bash
cd app/frontend
npm install
npm run dev
# Frontend will be available at http://localhost:5173
```

### Production Deployment

#### Docker Production
```bash
# Build and start production containers
docker-compose up -d

# Services will be available at:
# - Frontend: http://localhost:80
# - Backend API: http://localhost:8000
```

#### Vercel Deployment (Frontend)
```bash
cd app/frontend
vercel --prod
```

## 🔧 Environment Configuration

### Backend Environment Variables
```bash
# Server Configuration
PORT=8000

# Database Configuration
DB_HOST=localhost
DB_PORT=26257
DB_USER=root
DB_PASSWORD=
DB_NAME=truora_stocks
DB_SSLMODE=disable

# External API Configuration
STOCK_API_URL=https://api.karenai.click/swechallenge/list
STOCK_API_KEY=Bearer <your-token>
```

### Frontend Environment Variables
```bash
# API Configuration
VITE_API_BASE_URL=http://localhost:8000/api/v1
```

## 📊 API Documentation

The backend provides comprehensive API documentation through OpenAPI 3.1 specification:

- **OpenAPI Spec**: Available at `http://localhost:8000/api/openapi.yaml`
- **Interactive Docs**: Swagger UI available at `http://localhost:8000/swagger`

### Key Endpoints

#### Stock Management
- `GET /api/v1/stocks` - List stocks with pagination and search
- `GET /api/v1/stocks/{symbol}` - Get specific stock details
- `POST /api/v1/stocks/fetch` - Fetch latest stock data from external API

#### Recommendations
- `GET /api/v1/recommendations` - Get AI-powered stock recommendations
- `POST /api/v1/recommendations/generate` - Generate new recommendations

#### Health & Monitoring
- `GET /health` - Service health check

## 🧪 Testing

### Backend Testing
```bash
cd app/backend
go test ./...
```

### Frontend Testing
```bash
cd app/frontend
npm run type-check  # Type checking
npm run lint        # Linting
npm run build        # Production build test
```

## 🏗️ Development Workflow

### Git Workflow
1. Create feature branch: `git checkout -b feature/stock-analysis`
2. Make changes and commit: `git commit -m "feat: add stock analysis feature"`
3. Push branch: `git push origin feature/stock-analysis`
4. Create Pull Request

### Code Style
- **Backend**: Follow Go conventions with `gofmt`
- **Frontend**: ESLint + Prettier configuration included
- **Git**: Conventional commits format

## 📈 Performance & Monitoring

### Backend Optimizations
- **Database**: Connection pooling, strategic indexing
- **Caching**: Redis integration ready (placeholder)
- **Rate Limiting**: API rate limiting implemented
- **Logging**: Structured logging with log levels

### Frontend Optimizations
- **Bundle Size**: Tree-shaking and code splitting
- **Images**: Optimized asset delivery
- **Caching**: Service worker for offline support
- **Performance**: Lazy loading for routes and components

## 🔒 Security Features

### Backend Security
- **SQL Injection**: Parameterized queries with GORM
- **CORS**: Configured for production domains
- **Input Validation**: Request validation middleware
- **Rate Limiting**: API abuse prevention
- **HTTPS**: SSL/TLS ready for production

### Frontend Security
- **XSS Protection**: Vue.js built-in protections
- **CSP**: Content Security Policy headers
- **HTTPS**: Enforced in production
- **API Security**: Token-based authentication ready

## 🆘 Troubleshooting

### Common Issues

**Database Connection Issues:**
```bash
# Check CockroachDB status
cockroach node status --insecure --host=localhost:26257

# Reset database
cockroach sql --insecure -e "DROP DATABASE truora_stocks; CREATE DATABASE truora_stocks;"
```

**Port Conflicts:**
```bash
# Check port usage
lsof -i :8000  # Backend
lsof -i :5173  # Frontend
lsof -i :26257 # Database
```

**Docker Issues:**
```bash
# Reset Docker environment
docker-compose -f docker-compose.dev.yaml down
docker system prune -f
docker-compose -f docker-compose.dev.yaml up --build
```

## 🤝 Contributing

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Commit** changes: `git commit -m 'feat: add amazing feature'`
4. **Push** to branch: `git push origin feature/amazing-feature`
5. **Open** a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Financial Data**: Powered by KarenAI financial API
- **Icons**: Heroicons and custom SVG icons
- **UI Framework**: Tailwind CSS team for excellent styling utilities
- **Go Community**: For excellent libraries and best practices
- **Vue.js Team**: For the amazing Vue 3 framework

## 📞 Support

For support, email [support@truora.com](mailto:jfsg2901@gmail.com)

---

**Built by Juan Felipe Santamariá Guerrero**