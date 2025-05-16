# Stay For Long

This project implements a hexagonal architecture (also known as ports and adapters) in Go. It provides a REST API for managing stays and accommodations.

## Project Structure

```
.
├── cmd/                  # Application entry point
│   ├── config/           # Application configurations
│   ├── di/               # Application dependency injection
├── internal/             # Internal application code
│   ├── application/      # Use cases and application logic
│   ├── domain/           # Business entities and rules
│   ├── infra/            # Concrete implementations (adapters)
│   └── ports/            # Ports (interfaces) for adapters
└── pkg/                  # Public code that can be used by other projects
```

## Architecture Layers

### Domain
Contains business entities and business rules. This layer has no external dependencies.

### Ports
Defines the interfaces that the application needs to function. These interfaces are implemented by the adapters.

### Application
Contains application logic and use cases. Depends only on domain and ports.

### Infrastructure
Implements concrete adapters for databases, external services, etc.

## Development

### Requirements
- Go 1.24 or higher
- Docker and Docker Compose
- golangci-lint (for linting)

### Environment Variables
The application can be configured using the following environment variables:

```
SERVER_PORT=8080            # Port where the server will listen
READ_TIMEOUT=15             # Server read timeout in seconds
WRITE_TIMEOUT=15            # Server write timeout in seconds
IDLE_TIMEOUT=60             # Server idle timeout in seconds
```

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd stay-for-long
```

2. Install dependencies:
```bash
make mod
```

### Running the Application

Start the application with all dependencies:
```bash
make run
```

This command will:
1. Build and start the application using Docker Compose

The server will start on the configured port (default: 8080).

### Testing

Run all tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
```

### Linting

Run the linter:
```bash
make lint
```

### Mocks

Create and refresh mocks:
```bash
make mock
```

## Available Commands

- `make mod` - Update dependencies
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage
- `make run` - Start application with Docker
- `make lint` - Run linter
- `make mock` - Create new mocks

## API Endpoints

### Calculate Stats
Calculates the average, minimum, and maximum nightly rates for a set of bookings.

```bash
curl -X POST http://localhost:8080/stats \
  -H "Content-Type: application/json" \
  -d '[
    {
      "request_id": "bookata_XY123",
      "check_in": "2020-01-01",
      "nights": 5,
      "selling_rate": 200,
      "margin": 20
    },
    {
      "request_id": "kayete_PP234",
      "check_in": "2020-01-04",
      "nights": 4,
      "selling_rate": 156,
      "margin": 22
    }
  ]'
```
Response:
```json
{
  "avg_night": 8.29,
  "min_night": 8,
  "max_night": 8.58
}
```

### Maximize Profit
Finds the optimal combination of bookings that maximizes profit while avoiding booking overlaps.

```bash
curl -X POST http://localhost:8080/maximize \
  -H "Content-Type: application/json" \
  -d '[
    {
      "request_id": "bookata_XY123",
      "check_in": "2020-01-01",
      "nights": 5,
      "selling_rate": 200,
      "margin": 20
    },
    {
      "request_id": "kayete_PP234",
      "check_in": "2020-01-04",
      "nights": 4,
      "selling_rate": 156,
      "margin": 5
    },
    {
      "request_id": "atropote_AA930",
      "check_in": "2020-01-04",
      "nights": 4,
      "selling_rate": 150,
      "margin": 6
    },
    {
      "request_id": "acme_AAAAA",
      "check_in": "2020-01-10",
      "nights": 4,
      "selling_rate": 160,
      "margin": 30
    }
  ]'
```
Response:
```json
{
  "request_ids": [
    "bookata_XY123",
    "acme_AAAAA"
  ],
  "total_profit": 88,
  "avg_night": 10,
  "min_night": 8,
  "max_night": 12
}
```

### Error Handling

The API uses standard HTTP status codes and returns error messages in JSON format:

- 200 OK: Successful operation
- 400 Bad Request: Invalid request parameters or JSON format
- 500 Internal Server Error: Server-side error

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request 