# Social Media API

A Go-based REST API for a social media platform with user management and post creation capabilities.

## Features

- User management (registration, authentication)
- Post creation and management
- PostgreSQL database
- Docker containerization
- Database migrations using Goose

## Prerequisites

- Go 1.23 or later
- Docker and Docker Compose
- [Goose](https://github.com/pressly/goose) for database migrations

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/batugkocak/social-go.git
cd social-go
```

### 2. Environment Setup

Create a `.env` file in the root directory:

```env
ENV = development
ADDR=:8080
DB_ADDR=postgres://admin:adminPassword@localhost/social?sslmode=disable
DB_MAX_OPEN_CONNS=30
DB_MAX_IDLE_CONNS=30
DB_MAX_IDLE_TIME=15min
```

### 3. Start Docker Services

```bash
docker-compose up -d
```

This will start the PostgreSQL database on port 5432.

### 4. Database Migrations

Install Goose:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Run migrations:

```bash
# Set your database connection string
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="postgres://admin:adminPassword@localhost/social?sslmode=disable"

# Run migrations
goose up
```

To create a new migration:

```bash
goose create name_of_your_migration sql
```

### 5. Run the Application

```bash
# Run directly
go run ./cmd/api

# Or use air for hot reload (if installed)
air
```

## API Endpoints

- `GET /v1/health` - Health check endpoint
- More endpoints documentation coming soon...

## Development

The project uses:

- Chi router for HTTP routing
- PostgreSQL for data storage
- Docker for containerization
- Goose for database migrations

## Project Structure

```
.
├── cmd/
│   ├── api/          # Main application
│   └── migrate/      # Database migrations
├── internal/
│   ├── db/          # Database connection
│   ├── env/         # Environment configuration
│   ├── scripts/     # Database scripts
│   └── store/       # Data access layer
└── docker-compose.yml
```

## License

[MIT License](LICENSE)
