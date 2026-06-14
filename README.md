# go-user-api

A RESTful API built with **Go + Fiber** that manages users with a date-of-birth field and dynamically calculates their age on every fetch.

---

## Tech Stack

| Concern | Library |
|---------|---------|
| HTTP Framework | [GoFiber v2](https://gofiber.io/) |
| Database | PostgreSQL |
| DB Access Layer | [SQLC](https://sqlc.dev/) (generated code) |
| DB Driver | `lib/pq` |
| Logging | [Uber Zap](https://github.com/uber-go/zap) |
| Validation | [go-playground/validator v10](https://github.com/go-playground/validator) |
| Config | `godotenv` + env vars |

---

## Project Structure

```
.
├── cmd/server/main.go            # Entry point – wires everything together
├── config/config.go              # Config loader (env vars / .env)
├── db/
│   ├── migrations/
│   │   └── 001_create_users.sql  # Schema migration
│   └── sqlc/
│       ├── db.go                 # SQLC boilerplate
│       ├── models.go             # Generated User struct
│       ├── querier.go            # Generated interface + param types
│       ├── query.sql.go          # Generated query implementations
|       └── query.sql             # Raw SQL Queries for SQLC generation
|       
├── internal/
│   ├── handler/user_handler.go   # HTTP handlers (Fiber)
│   ├── logger/logger.go          # Singleton Zap logger
│   ├── middleware/middleware.go   # RequestID + RequestLogger
│   ├── models/
│   │   ├── user.go               # Request/response DTOs + CalculateAge()
│   │   └── user_test.go          # Unit tests for age calculation
│   ├── repository/
│   │   └── user_repository.go    # DB access abstraction
│   ├── routes/routes.go          # Route registration
│   └── service/
|        ├── user_service.go      # Business logic
|        ├── user_service_test.go # Service layer unit tests
|        └── age.go               # Exported CalculateAge() wrapper
├── sqlc.yaml                     # SQLC generation config
├── Dockerfile
├── docker-compose.yml
├── .env.example
└── README.md
```

---

## Prerequisites

- Go 1.21+
- PostgreSQL 14+ **or** Docker & Docker Compose

---

## Quick Start

### 1 – Clone & configure

```bash
git clone https://github.com/example/go-user-api.git
cd go-user-api
cp .env.example .env
# Edit .env with your DB credentials
```

### 2a – Run with Docker (recommended)

```bash
docker-compose up --build
```

This starts both PostgreSQL (with the migration applied automatically via `docker-entrypoint-initdb.d`) and the API on port **8080**.

### 2b – Run locally

```bash
# 1. Start PostgreSQL and create the database
psql -U postgres -c "CREATE DATABASE userdb;"

# 2. Apply the migration
psql -U postgres -d userdb -f db/migrations/001_create_users.sql

# 3. Download dependencies
go mod download

# 4. Run
go run ./cmd/server
```

---

## API Endpoints

Base URL: `http://localhost:8080`

### Health Check
```
GET /health
```
```json
{ "status": "ok" }
```

---

### Create User
```
POST /users
Content-Type: application/json
```
```json
{ "name": "Alice", "dob": "1990-05-10" }
```
**201 Created**
```json
{ "id": 1, "name": "Alice", "dob": "1990-05-10" }
```

---

### Get User by ID (includes age)
```
GET /users/:id
```
**200 OK**
```json
{ "id": 1, "name": "Alice", "dob": "1990-05-10", "age": 35 }
```

---

### Update User
```
PUT /users/:id
Content-Type: application/json
```
```json
{ "name": "Alice Updated", "dob": "1991-03-15" }
```
**200 OK**
```json
{ "id": 1, "name": "Alice Updated", "dob": "1991-03-15" }
```

---

### Delete User
```
DELETE /users/:id
```
**204 No Content**

---

### List Users (paginated, includes age)
```
GET /users?page=1&limit=10
```
**200 OK**
```json
{
  "data": [
    { "id": 1, "name": "Alice", "dob": "1990-05-10", "age": 35 }
  ],
  "total": 1,
  "page": 1,
  "limit": 10,
  "total_pages": 1
}
```

---

## Error Responses

All errors follow the same shape:
```json
{ "error": "human-readable message" }
```

| Status | Meaning |
|--------|---------|
| 400 | Bad request / invalid path param |
| 404 | User not found |
| 422 | Validation failed |
| 500 | Internal server error |

---

## Middleware

| Middleware | What it does |
|-----------|-------------|
| `RequestID` | Injects / echoes `X-Request-ID` response header |
| `RequestLogger` | Logs method, path, status code, and latency via Zap |

---

## Running Tests

```bash
go test ./...
```

The `internal/models` package contains unit tests for `CalculateAge()`:

```
=== RUN   TestCalculateAge/birthday_already_passed_this_year
=== RUN   TestCalculateAge/birthday_not_yet_this_year
=== RUN   TestCalculateAge/birthday_is_today
=== RUN   TestCalculateAge/newborn_(dob_=_today)
--- PASS: TestCalculateAge (0.00s)
```

---

## Regenerating SQLC Code

If you modify `db/sqlc/query.sql` or the migration:

```bash
# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Regenerate
sqlc generate
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_ENV` | `development` | `development` or `production` (affects log format) |
| `APP_PORT` | `8080` | Port the API listens on |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | PostgreSQL user |
| `DB_PASSWORD` | `postgres` | PostgreSQL password |
| `DB_NAME` | `userdb` | PostgreSQL database name |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL mode |
