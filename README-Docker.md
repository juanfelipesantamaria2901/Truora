# Truora Investment App - Docker & Infrastructure Guide

## Overview

This project includes comprehensive Docker and Infrastructure as Code (IaC) configurations for local development and production deployment.

## Project Structure

```
.
├── Dockerfile                 # Production multi-stage build
├── Dockerfile.dev            # Development with hot reloading
├── docker-compose.yaml       # Production compose
├── docker-compose.dev.yaml   # Development compose
├── nginx.conf               # Nginx configuration
├── supervisord.conf         # Process management
├── start.sh                 # Container startup script
├── Makefile                 # Build and deployment commands
├── .dockerignore           # Docker build exclusions
└── infra/
    ├── main.tf             # Basic Terraform configuration
    └── main-enhanced.tf    # Enhanced production-ready Terraform
```

## Local Development

### Prerequisites

- Docker and Docker Compose
- Make (optional, for convenience commands)

### Quick Start

1. **Development with hot reloading:**
   ```bash
   # Using Make
   make dev
   
   # Or directly with docker-compose
   docker-compose -f docker-compose.dev.yaml up --build
   ```

2. **Access the application:**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8000
   - CockroachDB Admin: http://localhost:8080
   - CockroachDB SQL: localhost:26257

3. **Stop development environment:**
   ```bash
   make dev-down
   # or
   docker-compose -f docker-compose.dev.yaml down
   ```

### Development Features

- **Hot Reloading**: Both frontend (Vite) and backend (Air) support hot reloading
- **Volume Mounts**: Source code is mounted for real-time changes
- **Separate Services**: Frontend, backend, and database run in separate containers
- **Database Persistence**: CockroachDB data persists between restarts

## Production Deployment

### Local Production Testing

1. **Build and run production container:**
   ```bash
   make build
   make run
   
   # Or directly
   docker-compose up --build
   ```

2. **Access the application:**
   - Application: http://localhost:80
   - CockroachDB Admin: http://localhost:8080

### Production Features

- **Multi-stage Build**: Optimized Docker image with minimal size
- **Nginx**: Serves static files and proxies API requests
- **Supervisor**: Manages multiple processes (Nginx, Go backend, CockroachDB)
- **Health Checks**: Built-in health monitoring
- **Security**: Non-root user, minimal attack surface

## Available Make Commands

```bash
# Development
make dev              # Start development environment
make dev-down         # Stop development environment
make dev-logs         # View development logs
make dev-frontend     # Start only frontend service
make dev-backend      # Start only backend service
make dev-db           # Start only database service

# Production
make build            # Build production image
make run              # Run production container
make stop             # Stop all containers
make clean            # Remove containers and images
make logs             # View production logs

# Database
make db-shell         # Connect to CockroachDB shell
make db-logs          # View database logs

# Testing
make test             # Run tests in containers
```

## Infrastructure as Code (Terraform)

### Basic Configuration (`infra/main.tf`)

The basic Terraform configuration provides:
- AWS ECS Fargate deployment
- Application Load Balancer
- ECR repository
- Basic VPC setup
- CloudWatch logging

### Enhanced Configuration (`infra/main-enhanced.tf`)

The enhanced configuration adds:
- **Security**: Enhanced security groups, IAM roles
- **Scalability**: Auto-scaling policies, capacity providers
- **Monitoring**: Comprehensive CloudWatch logging, ALB access logs
- **High Availability**: Multi-AZ deployment, health checks
- **Cost Optimization**: Fargate Spot instances, lifecycle policies
- **Variables**: Configurable parameters for different environments
- **Best Practices**: Proper tagging, encryption, deletion protection

### Deployment Steps

1. **Prerequisites:**
   ```bash
   # Install Terraform
   brew install terraform  # macOS
   
   # Configure AWS credentials
   aws configure
   ```

2. **Deploy infrastructure:**
   ```bash
   cd infra
   
   # Initialize Terraform
   terraform init
   
   # Plan deployment (using enhanced config)
   terraform plan -var-file="terraform.tfvars"
   
   # Apply changes
   terraform apply
   ```

3. **Build and push Docker image:**
   ```bash
   # Get ECR login token
   aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
   
   # Build and tag image
   docker build -t truora-investment-app .
   docker tag truora-investment-app:latest <ecr-repo-url>:latest
   
   # Push to ECR
   docker push <ecr-repo-url>:latest
   ```

### Terraform Variables

Create a `terraform.tfvars` file:

```hcl
region         = "us-east-1"
app_name       = "truora-investment-app"
environment    = "prod"
image_tag      = "latest"
cpu            = 1024
memory         = 2048
desired_count  = 2
enable_https   = true
domain_name    = "your-domain.com"
```

## Docker Configuration Details

### Multi-stage Build Process

1. **Backend Stage**: Builds Go application with optimizations
2. **Frontend Stage**: Builds Vue.js application with Vite
3. **CockroachDB Stage**: Prepares database with initialization
4. **Runtime Stage**: Combines all components with Nginx and Supervisor

### Security Considerations

- Non-root user execution
- Minimal base images (Alpine Linux)
- No sensitive data in images
- Health checks for monitoring
- Proper file permissions

### Performance Optimizations

- Multi-stage builds for smaller images
- Nginx for static file serving
- Gzip compression enabled
- Efficient layer caching
- .dockerignore for faster builds

## Troubleshooting

### Common Issues

1. **Port conflicts:**
   ```bash
   # Check what's using the port
   lsof -i :8080
   
   # Kill the process or change ports in docker-compose
   ```

2. **Database connection issues:**
   ```bash
   # Check database logs
   make db-logs
   
   # Connect to database shell
   make db-shell
   ```

3. **Build failures:**
   ```bash
   # Clean Docker cache
   docker system prune -a
   
   # Rebuild without cache
   docker-compose build --no-cache
   ```

### Logs and Debugging

```bash
# View all logs
make logs

# View specific service logs
docker-compose logs frontend
docker-compose logs backend
docker-compose logs crdb

# Follow logs in real-time
docker-compose logs -f
```

## Environment Variables

### Development
- `NODE_ENV=development`
- `VITE_API_URL=http://localhost:8000`
- `DATABASE_URL=postgresql://root@crdb:26257/myapp?sslmode=disable`

### Production
- `NODE_ENV=production`
- `PORT=8000`
- `DATABASE_URL=postgresql://root@localhost:26257/myapp?sslmode=disable`

## Next Steps

1. **CI/CD Pipeline**: Set up GitHub Actions or similar for automated builds
2. **Monitoring**: Add Prometheus/Grafana for metrics
3. **Secrets Management**: Use AWS Secrets Manager or similar
4. **Database**: Consider managed database service (AWS RDS)
5. **CDN**: Add CloudFront for static asset delivery
6. **SSL/TLS**: Configure HTTPS with ACM certificates

## Support

For issues or questions:
1. Check the logs using the provided commands
2. Review the Docker and Terraform configurations
3. Ensure all prerequisites are installed
4. Verify network connectivity and port availability