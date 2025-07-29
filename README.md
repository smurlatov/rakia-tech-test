# Blog API

REST API for a blog platform built with Go, following Domain-Driven Design (DDD) principles with a clean architecture.

## Features

- **CRUD Operations**: Create, Read, Update, Delete blog posts
- **RESTful Design**: Follows REST API conventions
- **Concurrency Safe**: Thread-safe in-memory storage with read-write mutexes
- **Domain-Driven Design**: Clean architecture with separated concerns
- **Comprehensive Testing**: Unit and integration tests with high coverage
- **Containerized**: Multi-stage Docker build for optimal image size
- **Developer Friendly**: Makefile with common development tasks

## Architecture

```
blog-api/
├── cmd/                     # Application entry point
│   └── main.go             # Main application file
├── internal/               # Internal packages (unexported)
│   ├── domain/             # Domain layer (entities, interfaces)
│   │   ├── entities/       # Business entities
│   │   └── repositories/   # Repository interfaces
│   ├── infrastructure/     # Infrastructure layer
│   │   ├── repositories/   # Repository implementations
│   │   └── loader/         # Data loading utilities
│   ├── application/        # Application layer
│   │   └── services/       # Business logic services
│   └── interfaces/         # Interface layer
│       └── rest/           # REST API handlers and routing
│           ├── dto/        # Data transfer objects
│           ├── post_handler.go  # HTTP handlers
│           └── router.go   # Route configuration
└── test/                   # Test files
    ├── integration/        # Integration tests
    └── concurrency/        # Concurrency tests
```

## API Endpoints

| Method | Endpoint        | Description              |
|--------|-----------------|--------------------------|
| GET    | `/health`       | Health check             |
| GET    | `/api/v1/posts` | Get all blog posts       |
| GET    | `/api/v1/posts/{id}` | Get specific blog post |
| POST   | `/api/v1/posts` | Create new blog post     |
| PUT    | `/api/v1/posts/{id}` | Update existing post |
| DELETE | `/api/v1/posts/{id}` | Delete blog post    |

## API Examples

### Create a Post
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first blog post.",
    "author": "John Doe"
  }'
```

### Get All Posts
```bash
curl http://localhost:8080/api/v1/posts
```

### Get Specific Post
```bash
curl http://localhost:8080/api/v1/posts/1
```

### Update a Post
```bash
curl -X PUT http://localhost:8080/api/v1/posts/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Title",
    "content": "Updated content.",
    "author": "John Doe"
  }'
```

### Delete a Post
```bash
curl -X DELETE http://localhost:8080/api/v1/posts/1
```

## Quick Start

### Prerequisites

- Go 1.21 or later
- Docker
- Make 

### Using Docker

1. Clone the repository

2. Build the Docker image:
   ```bash
   make docker-build
   # or
   docker build -t blog-api:latest .
   ```

3. Run the container:
   ```bash
   make docker-run
   # or
   docker run -d -p 8080:8080 --name blog-api blog-api:latest
   ```
   
The API will be available at `http://localhost:8080`

### Using Docker-cpmpose

   ```bash
   docker-compose up
   ```

## Docker Security & Optimization

The Docker build follows **production security best practices**:

### 🔒 Security Features

- **Non-root user**: Application runs as `appuser` (not root)
- **Minimal base image**: Uses Alpine Linux for reduced attack surface
- **Multi-stage build**: Separates build dependencies from runtime
- **Health checks**: Automatic container health monitoring

### ⚡ Optimization Features

- **Cache cleanup**: Removes package manager cache (`/var/cache/apk/*`, `/tmp/*`)
- **Small image size**: Final image is only ~57MB
- **Efficient layering**: Optimized Docker layer structure

## Graceful Shutdown

The application implements **production-grade graceful shutdown** for zero-downtime deployments:

### 🔄 Shutdown Process

1. **Signal Detection**: Listens for `SIGINT` (Ctrl+C) and `SIGTERM` (Docker stop)
2. **Graceful Stop**: Allows ongoing requests to complete
3. **Timeout Protection**: 30-second timeout prevents hanging
4. **Clean Exit**: Proper resource cleanup and logging

## Development

### Available Make Commands

```bash
make help              # Show all available commands
make build             # Build the application (runs tests first)
make build-only        # Build the application without running tests
make run               # Run the application
make test              # Run all tests
make test-coverage     # Run tests with coverage report
make test-race         # Run tests with race detection
make fmt               # Format Go code
make lint              # Lint Go code (requires golangci-lint)
make clean             # Clean build artifacts
make deps              # Download dependencies
make docker-build      # Build Docker image (runs tests first)
make docker-build-only # Build Docker image without running tests
make docker-run        # Run in Docker container
make docker-stop       # Stop Docker container
```

### Test-First Build Philosophy

The build system follows a **test-first approach** to ensure code quality:

- **`make build`** - Runs all tests before building. **Build fails if tests fail.**
- **`make build-only`** - Fast build without tests (for development)
- **`make docker-build`** - Runs tests before Docker build. **Build fails if tests fail.**
- **`make docker-build-only`** - Fast Docker build without tests

## Concurrency Safety

The in-memory repository implementation uses `sync.RWMutex` to ensure thread safety:

- **Read operations** (GetByID, GetAll, Exists) use read locks for concurrent access
- **Write operations** (Create, CreatePost, Update, Delete, LoadData) use exclusive write locks
- **Data isolation** is maintained by returning copies of stored data
- **ID generation** is atomic and thread-safe within CreatePost method

## Error Handling

The API provides consistent error responses:

```json
{
  "error": "error_code",
  "message": "Human readable error message"
}
```

Common error codes:
- `validation_error`: Invalid request data or ID format
- `not_found`: Resource not found
- `creation_failed`: Failed to create resource
- `update_failed`: Failed to update resource
- `internal_error`: Internal server error

## Logging

The application uses structured logging with JSON format in production. Log levels can be configured via the `LOG_LEVEL` environment variable.

Key log events:
- Data loading on startup
- CRUD operations with post IDs
- Error conditions with context

## Sample Data

The `blog_data.json` file contains 100 sample blog posts that are automatically loaded when the application starts. This provides immediate data for testing and development without requiring manual post creation. 