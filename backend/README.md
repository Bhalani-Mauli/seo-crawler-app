# SEO Crawler Checker Backend

A Go-based web crawler service that analyzes websites for SEO metrics, following SOLID principles and Go industry standards.

## Project Structure

```
backend/
├── cmd/
│   └── server/           # Application entry point
│       └── main.go
├── internal/             # Private application code
│   ├── api/             # API layer
│   │   ├── handlers.go  # HTTP handlers
│   │   └── router.go    # Route definitions
│   ├── config/          # Configuration management
│   │   └── config.go
│   ├── database/        # Data access layer
│   │   ├── connection.go # Database connection
│   │   ├── repository.go # Repository pattern implementation
│   │   └── schema.go    # Database schema management
│   ├── middleware/      # HTTP middleware
│   │   └── auth.go      # Authentication middleware
│   ├── models/          # Data models
│   │   └── crawl.go     # Crawl data structures
│   └── services/        # Business logic layer
│       └── crawl_service.go
├── pkg/                 # Public packages
│   └── crawler/         # Crawler implementation
│       └── crawler.go
├── go.mod
├── go.sum
├── Dockerfile
├── README.md
└── (optional) docker-compose.yml
```

## SOLID Principles Implementation

### 1. Single Responsibility Principle (SRP)

- Each package has a single, well-defined responsibility
- `models/` - Data structures only
- `database/` - Data access only
- `services/` - Business logic only
- `api/` - HTTP handling only
- `crawler/` - Web crawling only

### 2. Open/Closed Principle (OCP)

- Interfaces allow for extension without modification
- `Repository` interface enables different database implementations
- `CrawlerService` interface allows different crawling strategies
- `Handler` interface supports different API implementations

### 3. Liskov Substitution Principle (LSP)

- All implementations can be substituted for their interfaces
- Repository implementations are interchangeable
- Service implementations follow the same contract

### 4. Interface Segregation Principle (ISP)

- Small, focused interfaces
- `Repository` interface has specific database operations
- `CrawlerService` interface has only crawling operations
- `Handler` interface has only HTTP operations

### 5. Dependency Inversion Principle (DIP)

- High-level modules depend on abstractions
- Services depend on repository interfaces, not concrete implementations
- Handlers depend on service interfaces
- Configuration is injected rather than hardcoded

## Architecture Layers

### 1. Presentation Layer (`internal/api/`)

- HTTP handlers and routing
- Request/response handling
- Input validation

### 2. Business Logic Layer (`internal/services/`)

- Application business rules
- Orchestration of operations
- Transaction management

### 3. Data Access Layer (`internal/database/`)

- Database operations
- Repository pattern implementation
- Schema management

### 4. Domain Layer (`internal/models/`)

- Data structures
- Domain entities
- Value objects

### 5. Infrastructure Layer (`pkg/crawler/`)

- External service integration
- Web crawling implementation
- Third-party library usage

## Configuration

The application uses environment variables for configuration:

```bash
# Database Configuration
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=seo_user
DB_PASSWORD=seo_password
DB_NAME=seo_crawler

# Server Configuration
SERVER_PORT=8080

# API Configuration
API_KEY=seo-crawler-api-key-2025
```

## Local Development Commands

### 1. Install dependencies

```bash
go mod tidy
```

### 2. Build the application

```bash
go build -o seo-crawler ./cmd/server/main.go
```

### 3. Run the application locally

```bash
# Ensure MySQL is running and configured
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=password
export DB_NAME=seo_crawler_check
export SERVER_PORT=8080
export API_KEY=seo-crawler-api-key-2025 #<-- example key use secret key in prod

./seo-crawler
```

Or simply:

```bash
go run ./cmd/server/main.go
```

### 4. Run tests

```bash
go test ./...
```

## Docker Commands

### 1. Build the Docker image

```bash
docker build -t seo-crawler-backend .
```

### 2. Run the Docker container

```bash
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=3306 \
  -e DB_USER=seo_user \
  -e DB_PASSWORD=seo_password \
  -e DB_NAME=seo_crawler \
  -e SERVER_PORT=8080 \
  -e API_KEY=seo-crawler-api-key-2025 \
  seo-crawler-backend
```
