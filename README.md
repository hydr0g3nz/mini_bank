# Mini Bank API

A simple banking API built with Go, Gin, PostgreSQL, and Redis using Clean Architecture principles.

## Prerequisites

- Docker & Docker Compose
- Go 1.23.4+ (for local development)
- Git

## Quick Start

### 1. Clone Repository
```bash
git clone <repository-url>
cd mini_bank
```

### 2. Environment Setup
```bash
# Copy environment file
cp .env.example .env

# Edit .env file with your configurations (optional)
vim .env
```

### 3. Run with Docker Compose
```bash
# Start all services
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f app
```

### 4. Verify Installation
```bash
# Health check
curl http://localhost:8080/health

# Expected response:
# {"service":"mini-bank-api","status":"ok"}
```

## Local Development

### Run without Docker
```bash
# Install dependencies
go mod tidy

# Setup local PostgreSQL and Redis
# Update .env with local database credentials

# Run application
go run cmd/main.go
```

### Run with local database but containerized dependencies
```bash
# Start only database services
docker-compose up -d postgres redis

# Run application locally
go run cmd/main.go
```

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint

### Account Management
- `POST /api/v1/accounts` - Create new account
- `GET /api/v1/accounts` - List all accounts (with pagination)
- `GET /api/v1/accounts/:id` - Get specific account
- `PUT /api/v1/accounts/:id` - Update account information
- `DELETE /api/v1/accounts/:id` - Delete account
- `PATCH /api/v1/accounts/:id/suspend` - Suspend account
- `PATCH /api/v1/accounts/:id/activate` - Activate account
- `GET /api/v1/accounts/:id/transactions` - Get transactions for specific account

### Transaction Management
- `POST /api/v1/transactions` - Create new transaction
- `GET /api/v1/transactions` - List all transactions (with pagination)
- `GET /api/v1/transactions/:id` - Get specific transaction
- `PATCH /api/v1/transactions/:id/confirm` - Confirm pending transaction
- `PATCH /api/v1/transactions/:id/cancel` - Cancel pending transaction
- `GET /api/v1/transactions/status/:status` - Get transactions by status

### Authentication
All API endpoints (except `/health`) require API key authentication via `x-api-key` header.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_HOST` | Database host | `localhost` |
| `DB_USER` | Database username | `minibank_user` |
| `DB_PASSWORD` | Database password | `minibank_pass` |
| `DB_NAME` | Database name | `mini_bank` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PASSWORD` | Redis password | `redis_pass` |
| `API_KEY` | API authentication key | `your-secret-api-key-change-in-production` |
| `LOG_LEVEL` | Logging level | `info` |

## Docker Commands

```bash
# Build and start services
docker-compose up -d --build

# Stop services
docker-compose down

# Remove volumes (caution: deletes data)
docker-compose down -v

# View service logs
docker-compose logs -f [service_name]

# Execute commands in container
docker-compose exec app sh
docker-compose exec postgres psql -U minibank_user -d mini_bank
```

## Database Migration

The application automatically runs database migrations on startup using GORM AutoMigrate.

## API Testing

Use the provided Postman collection for testing all endpoints. Import the collection and set up environment variables for the API key and base URL.

## Stopping Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (deletes all data)
docker-compose down -v
```