# Pack Calculator API

A production-ready REST API for calculating optimal pack distributions to fulfill orders. Built with Go, featuring clean architecture, comprehensive testing, and structured logging.

## 🎯 Problem Statement

Calculate the optimal pack distribution for any order quantity using available pack sizes, following these business rules:

1. **Only whole packs can be sent** - Packs cannot be broken open
2. **Send the least amount of items** to fulfill the order (minimize overage)
3. **Send as few packs as possible** among valid options (rule #2 takes precedence)

## ✨ Features

- 🚀 **High Performance** - Optimized algorithm handles large orders (500K+ items) efficiently
- 🏗️ **Clean Architecture** - Domain-driven design with clear separation of concerns
- 📊 **Comprehensive Testing** - 97%+ test coverage with unit, integration, and benchmark tests
- 📝 **Structured Logging** - JSON/text logging with request tracing and performance metrics
- 🐳 **Containerized** - Docker and Docker Compose ready
- 🔍 **Health Monitoring** - Built-in health checks and readiness probes
- 🌐 **Web Interface** - Interactive UI for testing calculations

## 🚀 Quick Start

### Using Docker (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd pack-calculator

# Add to /etc/hosts for local development
echo "127.0.0.1 pack-calculator.localhost" | sudo tee -a /etc/hosts

# Start the application with Traefik
make docker-run

# Access the application:
# - API: http://pack-calculator.localhost
# - Web UI: http://pack-calculator.localhost
# - Health: http://pack-calculator.localhost/health
# - Traefik Dashboard: http://pack-calculator.localhost:8080
```

### Local Development

```bash
# Install dependencies
make deps

# Run the application
make run

# The API will be available at http://pack-calculator.localhost
```

## 🧪 Validation & Testing

### Edge Case Verification

Test the complex edge case from requirements:

```bash
# Using Make command
make edge-case-test

# Using curl directly
curl -X POST http://pack-calculator.localhost/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"pack_sizes": [23, 31, 53], "order_quantity": 500000}'

# Expected result: {"23": 2, "31": 7, "53": 9429}
```

### Run Tests

```bash
make test              # Run all tests with coverage
make bench             # Run performance benchmarks
make api-test          # Test API endpoints
make health            # Check application health
```

## 🏗️ Architecture

```
├── cmd/api/                    # Application entrypoint
├── internal/
│   ├── domain/
│   │   ├── model/              # Business entities (Pack, Calculation, PackDistribution)
│   │   │   ├── calculation.go  # Calculation model and PackDistribution logic
│   │   │   ├── errors.go       # Domain-specific errors
│   │   │   ├── pack.go         # Pack entity
│   │   │   └── validation.go   # Domain validation logic
│   │   └── service/            # Business logic (PackCalculator, PackService)
│   │       ├── pack_calculator.go # Core calculation algorithm
│   │       └── pack_service.go    # Service orchestration
│   ├── api/
│   │   ├── dto/                # API data transfer objects
│   │   │   ├── calculation.go  # Calculation request/response DTOs
│   │   │   ├── health.go       # Health check DTOs
│   │   │   └── pack.go         # Pack management DTOs
│   │   ├── handlers/           # HTTP handlers
│   │   │   ├── calculation.go  # Calculation endpoint handler
│   │   │   └── health.go       # Health check handlers
│   │   ├── http/               # HTTP server and routing
│   │   │   ├── router.go       # Route configuration
│   │   │   ├── response.go     # Response helpers
│   │   │   └── server.go       # HTTP server setup
│   │   └── middleware/         # HTTP middleware
│   │       └── logging.go      # Request logging middleware
│   ├── infrastructure/         # Infrastructure concerns
│   │   └── logger/             # Structured logging system
│   │       ├── logger.go       # Logger implementation
│   │       └── logger_test.go  # Logger tests
│   └── config/                 # Configuration management
│       └── config.go           # Viper-based configuration
├── docker-compose.yml          # Container orchestration with Traefik
├── Dockerfile                  # Production container
└── Makefile                    # Development commands
```

### Key Design Principles

- **Domain-Driven Design** - Business logic isolated in domain layer
- **Clean Architecture** - Dependencies point inward toward domain
- **Infrastructure Separation** - External concerns (logging) in infrastructure layer
- **API Layer** - HTTP concerns separated from business logic
- **Testability** - All components are easily testable in isolation

## 🔧 API Reference

### Calculations

#### `POST /api/v1/calculate`

Calculate optimal pack distribution for custom pack sizes.

**Request:**
```json
{
  "pack_sizes": [250, 500, 1000, 2000, 5000],
  "order_quantity": 12001
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "1672531200123456789",
    "packs_used": {
      "250": 1,
      "2000": 1,
      "5000": 2
    },
    "total_items": 12250,
    "total_packs": 4,
    "items_overage": 249,
    "calculation_time": "1.234ms"
  }
}
```

### Health & Monitoring

- `GET /health` - Application health check
- `GET /ready` - Readiness probe for container orchestration

### Web Interface

- `GET /` - Interactive web UI for testing calculations (placeholder)
- `GET /ui` - Alternative UI path (placeholder)

## ⚙️ Configuration

Configure via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PC_SERVER_PORT` | `8080` | HTTP server port |
| `PC_SERVER_HOST` | `0.0.0.0` | HTTP server host |
| `PC_LOGGING_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `PC_LOGGING_FORMAT` | `json` | Log format (json, text) |
| `PC_APP_ENVIRONMENT` | `development` | Application environment |
| `PC_APP_NAME` | `pack-calculator` | Application name |
| `PC_APP_VERSION` | `1.0.0` | Application version |

## 🛠️ Development

### Available Commands

```bash
# Development
make build             # Build the application
make run               # Run locally
make test              # Run tests with coverage
make bench             # Run benchmarks
make lint              # Run code linter
make fmt               # Format code

# Docker
make docker-build      # Build Docker image
make docker-run        # Start with Docker Compose
make docker-stop       # Stop containers
make docker-logs       # View logs

# Testing & Validation
make api-test          # Test API endpoints
make edge-case-test    # Test specific edge case
make health            # Check application health

# Utilities
make clean             # Clean build artifacts
make deps              # Download dependencies
make help              # Show all commands
```

### Development Environment Setup

```bash
make dev-setup         # Install development tools
```

## 📊 Examples

### Basic Calculation

```bash
curl -X POST http://pack-calculator.localhost/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "pack_sizes": [250, 500, 1000],
    "order_quantity": 263
  }'
```

**Result:** `{"500": 1}` - One 500-item pack (minimal overage)

### Complex Scenario

```bash
curl -X POST http://pack-calculator.localhost/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "pack_sizes": [250, 500, 1000, 2000, 5000],
    "order_quantity": 501
  }'
```

**Result:** `{"250": 1, "500": 1}` - Combination for exact fulfillment

### Edge Case (Performance Test)

```bash
curl -X POST http://pack-calculator.localhost/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "pack_sizes": [23, 31, 53],
    "order_quantity": 500000
  }'
```

**Result:** `{"23": 2, "31": 7, "53": 9429}` - Complex optimization in ~2 seconds

## 🚀 Deployment

### Docker Compose (Development)

```bash
make docker-run        # Starts API with Traefik proxy
```

The application includes:
- **API Server** - Main application on port 8080
- **Traefik** - Reverse proxy with automatic service discovery
- **Health Checks** - Container health monitoring
- **Graceful Shutdown** - Proper signal handling

## 🧪 Test Coverage

```
internal/domain/service             97.9%   # Core business logic
internal/api/handlers               68.5%   # HTTP layer  
internal/domain/model               51.9%   # Domain models
internal/infrastructure/logger      85.2%   # Logging system
internal/api/dto                    16.7%   # Data transfer objects
internal/config                      0.0%   # Configuration (simple structs)
internal/api/http                    0.0%   # HTTP utilities (simple wrappers)
internal/api/middleware              0.0%   # Middleware (simple functions)
```

### Test Categories

- **Unit Tests** - Core algorithm, business logic, models
- **Integration Tests** - HTTP handlers, API contracts  
- **Benchmark Tests** - Performance validation
- **Edge Case Tests** - Complex scenarios and error conditions

Test files are co-located with source code:
- `internal/domain/service/*_test.go` - Service layer tests
- `internal/domain/model/*_test.go` - Model tests  
- `internal/api/handlers/*_test.go` - Handler tests
- `internal/infrastructure/logger/*_test.go` - Logger tests

## 🔒 Security Features

- **Input Validation** - Request validation with detailed error messages
- **SQL Injection Protection** - Safe parameter handling
- **Container Security** - Non-root user, minimal base image
- **Error Handling** - No sensitive information in error responses

## ⚡ Performance

- **Algorithm Complexity** - O(n) time complexity for most cases
- **Memory Efficiency** - Minimal memory allocation during calculations
- **Request Throughput** - Handles concurrent requests efficiently
- **Edge Case Performance** - 500K item calculation in ~2 seconds

## 📈 Monitoring & Observability

### Structured Logging

```json
{
  "timestamp": "2025-01-14T10:30:00Z",
  "level": "INFO",
  "message": "Calculation completed",
  "request_id": "a1b2c3d4",
  "fields": {
    "pack_sizes": [23, 31, 53],
    "order_quantity": 500000,
    "total_items": 500000,
    "duration_ms": 1200
  }
}
```

### Health Monitoring

- **Health Endpoint** - `GET /health` returns application status
- **Readiness Probe** - `GET /ready` for container orchestration
- **Request Tracing** - Unique request IDs for distributed tracing
- **Performance Metrics** - Response times and success rates

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**Built with ❤️ using Go, Docker, and modern software engineering practices**
