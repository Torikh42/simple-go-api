# Simple Product API — Golang

A simple RESTful API for product management built with Go, following clean architecture principles.

## Tech Stack

- **Language**: Go 1.22+
- **Database**: PostgreSQL (via Docker)
- **Driver**: [pgx/v5](https://github.com/jackc/pgx)
- **Query Generator**: [sqlc](https://sqlc.dev)
- **Auth**: [golang-jwt](https://github.com/golang-jwt/jwt) & bcrypt
- **Hot Reload**: [air](https://github.com/air-verse/air)

## Architecture

```
Handler → Service → Repository → PostgreSQL
```

```
go-api/
├── cmd/api/main.go                 # Entry point
├── db/
│   ├── queries/product.sql         # SQL queries (sqlc input)
│   └── schema/products.sql         # Table schema (sqlc input)
├── internal/
│   ├── db/                         # Generated code by sqlc (DO NOT EDIT)
│   ├── handlers/product.go         # HTTP layer
│   ├── middleware/                  # Logger, Recoverer
│   ├── models/product.go           # Data structs
│   ├── repository/                 # Database layer
│   ├── routes/routes.go            # Route definitions
│   └── services/product.go         # Business logic
└── migrations/                     # SQL migration files
```

## Getting Started

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Docker](https://www.docker.com/)

### 1. Clone the repository

```bash
git clone https://github.com/Torikh42/simple-go-api.git
cd simple-go-api
```

### 2. Set up environment variables

```bash
cp .env.example .env
```

Or create a `.env` file manually:

```env
DB_HOST=localhost
DB_PORT=5434
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_api_db
APP_PORT=8080
```

### 3. Start the database

```bash
docker compose up -d
```

This will start a PostgreSQL container named `simple-go-api` on port `5434` and automatically run the migration in `migrations/`.

### 4. Run the server

```bash
go run ./cmd/api
```

Or with hot reload:

```bash
air
```

The server will be available at `http://localhost:8080`.

## API Endpoints

| Method | Endpoint           | Description          |
|--------|--------------------|----------------------|
| GET    | `/health`          | Health check         |
| POST   | `/register`        | Register user baru   |
| POST   | `/login`           | Login & dapat JWT    |
| GET    | `/products`        | Get all products (🔒)|
| POST   | `/products`        | Create a product (🔒)|
| GET    | `/products/{id}`   | Get product by ID(🔒)|
| PUT    | `/products/{id}`   | Update a product (🔒)|
| DELETE | `/products/{id}`   | Delete a product (🔒)|

### Example Request

**Create a product:**
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Kopi Susu", "price": 25000, "stock": 100}'
```

**Response:**
```json
{
  "id": 1,
  "name": "Kopi Susu",
  "price": 25000,
  "stock": 100
}
```

## Development

### Regenerate sqlc code

After modifying `db/queries/product.sql` or `db/schema/products.sql`, run:

```bash
sqlc generate
```

### Install sqlc CLI

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```
