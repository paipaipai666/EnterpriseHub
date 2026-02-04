# AGENTS.md - EnterpriseHub Development Guide

This file provides guidelines for AI agents working on the EnterpriseHub codebase.

## Build Commands

### Run a Single Service
```bash
make user_service     # Run user-service
make auth_service     # Run auth-service
make order_service    # Run order-service
make payment_service  # Run payment-service
make gateway          # Run gateway-service
```

### Manual Commands
```bash
cd <service-name> && go run cmd/main.go
cd <service-name> && go build -o <service-name> ./cmd/
```

### All Services
```bash
# Run each service in separate terminals:
cd user-service && go run cmd/main.go
cd auth-service && go run cmd/main.go
cd order-service && go run cmd/main.go
cd payment-service && go run cmd/main.go
cd gateway-service && go run cmd/main.go
```

### Dependencies
```bash
go mod tidy
```

### Generate gRPC Code
```bash
# From service directory:
protoc --go_out=. --go-grpc_out=. ./internal/pb/*.proto
```

## Testing Commands

### Run All Tests
```bash
cd <service-name> && go test ./...
```

### Run Single Test
```bash
cd <service-name> && go test -v ./internal/service/
cd <service-name> && go test -run TestFunctionName ./...
cd <service-name> && go test -v ./internal/repository/ -run TestCreate
```

### Test Coverage
```bash
cd <service-name> && go test -cover ./...
cd <service-name> && go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
```

## Linting and Formatting

### Format Code
```bash
gofmt -w .
goimports -w .
```

### Vet Code
```bash
go vet ./...
```

### Static Analysis
```bash
golangci-lint run ./...
```

## Code Style Guidelines

### Naming Conventions

**Interfaces**: Use `XxxService`, `XxxRepository`, `XxxController`, `XxxHandler` pattern.
```go
type OrderService interface { ... }
type PaymentRepository interface { ... }
type AuthController interface { ... }
```

**Implementations**: Use `xxxServiceImpl`, `xxxRepositoryImpl` pattern (unexported).
```go
type orderServiceImpl struct { ... }
type paymentRepositoryImpl struct { ... }
```

**Factory Functions**: Use `NewXxxService`, `NewXxxRepository` pattern.
```go
func NewOrderService(repo repository.OrderRepository) OrderService { ... }
func NewPaymentRepository() PaymentRepository { ... }
```

**Files**: Use snake_case: `order_service.go`, `payment_repository.go`.

**Variables**: Use camelCase for local variables, PascalCase for exported.
```go
orderId := "123"
var PaymentAmount float64
```

**Constants**: Use PascalCase or snake_case for constants.
```go
const MaxRetries = 3
const db_max_open_conns = 25
```

### Project Structure

Each service follows this structure:
```
<service>/
├── cmd/main.go           # Entry point
├── internal/
│   ├── api/              # HTTP Handlers/Controllers
│   ├── service/          # Business logic
│   ├── repository/       # Data access
│   ├── domain/          # Domain models and business rules
│   ├── model/           # Data models (for user-service)
│   ├── handler/         # gRPC handlers
│   ├── client/          # gRPC clients to other services
│   ├── middleware/     # HTTP middleware
│   └── config/          # Configuration
├── initializers/        # DB, Redis, env setup
├── migration/          # Database migrations
├── proto/              # gRPC proto files (if custom)
├── internal/pb/        # Generated gRPC code
└── go.mod
```

### Imports Organization

Group imports in this order:
1. Standard library
2. Third-party packages
3. Internal imports (relative paths)
```go
import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/paipaipai666/EnterpriseHub/order-service/internal/domain"
    "github.com/paipaipai666/EnterpriseHub/order-service/internal/repository"
)
```

### Error Handling

**Return errors with context**:
```go
if err != nil {
    return "", fmt.Errorf("无法获取用户信息: %v", err)
}
```

**Use sentinel errors** when appropriate:
```go
var ErrNotFound = errors.New("record not found")
```

**Check errors before using values**:
```go
payment, err := pgh.repo.Find(req.PaymentId)
if err != nil {
    return nil, err
}
```

**Handle errors at appropriate layers** - don't silently ignore errors.

### Domain Models

**Use value objects for business rules**:
```go
type Order struct {
    Id        string
    UserID    string
    Amount    float64
    Status    OrderStatus
    CreatedAt time.Time
}

func (o *Order) Cancel() error {
    if o.Status != OrderStatusCreated {
        return errors.New("只能取消创建状态的订单")
    }
    o.Status = OrderStatusCancelled
    return nil
}
```

### gRPC Services

**Handler pattern**:
```go
type PaymentGrpcHandler struct {
    pb.UnimplementedPaymentServiceServer
    repo repository.PaymentRepository
}

func NewPaymentGrpcHandler(repo repository.PaymentRepository) *PaymentGrpcHandler {
    return &PaymentGrpcHandler{repo: repo}
}
```

### HTTP Controllers

**Use interfaces for testability**:
```go
type AuthController interface {
    LoginWithJWT(c *gin.Context)
    Logout(c *gin.Context)
}
```

**Return consistent JSON responses** (see API.md):
```json
{
    "code": 0,
    "message": "success",
    "data": {}
}
```

### Database Operations

**Repository pattern with interfaces**:
```go
type PaymentRepository interface {
    Create(payment *domain.Payment) error
    Find(id string) (*domain.Payment, error)
    Save(payment *domain.Payment) error
}
```

**Use GORM for ORM operations**:
```go
err := initializers.DB.Create(&payment).Error
return err
```

### Configuration

**Load environment variables**:
```go
import "github.com/joho/godotenv"
func LoadEnv() {
    godotenv.Load()
}
```

**Use initializers package** for shared setup:
```go
initializers.Database()
initializers.Redis()
```

### Communication Patterns

**Service-to-service via gRPC**:
- Define protobuf services in `proto/` directories
- Generate Go code with `protoc`
- Use client wrappers in `internal/client/`

**Async via RabbitMQ**:
- Publish events after state changes
- Consume in notify-service

## Tech Stack Summary

- **Language**: Go 1.25
- **Web Framework**: Gin
- **RPC**: gRPC + Protobuf
- **ORM**: GORM
- **Database**: MySQL 8
- **Cache**: Redis
- **Message Queue**: RabbitMQ
- **Authentication**: JWT (github.com/golang-jwt/jwt/v5)
- **Logging**: Zap (recommended)

## Key Files

- `API.md` - HTTP API specification
- `README.md` - Project overview
- `Makefile` - Service run commands
