[![GitHub repo size](https://img.shields.io/github/repo-size/LucasHARosa/BE-Daily-Diet)](https://github.com/LucasHARosa/BE-Daily-Diet)
[![GitHub language count](https://img.shields.io/github/languages/count/LucasHARosa/BE-Daily-Diet)](https://github.com/LucasHARosa/BE-Daily-Diet)
[![GitHub top language](https://img.shields.io/github/languages/top/LucasHARosa/BE-Daily-Diet)](https://github.com/LucasHARosa/BE-Daily-Diet)
[![GitHub last commit](https://img.shields.io/github/last-commit/LucasHARosa/BE-Daily-Diet)](https://github.com/LucasHARosa/BE-Daily-Diet)

# Daily Diet API

REST API for the [Daily Diet](https://github.com/LucasHARosa/APP-Daily-Diet) mobile app — built with **Go**, **PostgreSQL** and a clean layered architecture. Handles authentication, meal tracking, metrics, food plans, user profiles and AI-powered calorie estimation via Google Gemini.

---

## Stack

| Layer | Library |
|---|---|
| Language | Go 1.24 |
| Router | chi v5 |
| Database | PostgreSQL 16 (Docker) |
| Driver | pgx v5 (pgxpool) |
| Query generation | sqlc |
| Migrations | goose |
| Auth | JWT (golang-jwt/jwt v5) + bcrypt |
| AI | Google Gemini API |
| Env | godotenv |

---

## Architecture

Every request follows this path:

```
HTTP Request
  → chi Router
    → Auth Middleware (validates JWT, injects user_id into context)
      → Handler (reads body/params, validates, extracts user_id from token)
        → Service (business rules)
          → Repository (database access via sqlc-generated code)
            → PostgreSQL
```

`user_id` **never** comes from the request body — always from the JWT. This ensures users can never access each other's data.

---

## Endpoints

### Auth
```
POST /users              → register
POST /sessions           → login  { access_token, refresh_token, user }
POST /sessions/refresh   → renew tokens
GET  /me                 → authenticated user info
```

### Meals
```
POST   /meals
GET    /meals?start=&end=&status=on_diet|off_diet
GET    /meals/:id
PUT    /meals/:id
DELETE /meals/:id
```

### Metrics
```
GET /metrics/summary
GET /metrics?start=2026-06-01&end=2026-06-30&groupBy=day|month
```

### Food Plans
```
POST   /food-plans
GET    /food-plans
GET    /food-plans/active
GET    /food-plans/:id
PUT    /food-plans/:id
PATCH  /food-plans/:id/active
DELETE /food-plans/:id

POST   /food-plans/:id/days/:weekday/meals
PUT    /food-plan-meals/:id
DELETE /food-plan-meals/:id
```

### Profile & AI
```
GET   /me/profile
PATCH /me/profile

POST /calorie-estimations
GET  /calorie-estimations
```

---

## Project Structure

```
BE-Daily-Diet/
├── cmd/api/
│   └── main.go                  # entry point, wires all layers
├── internal/
│   ├── config/                  # loads .env
│   ├── domain/                  # pure functions (CalculateBestStreak)
│   ├── application/
│   │   ├── auth/                # register, login, refresh token
│   │   ├── meals/               # meal CRUD
│   │   ├── metrics/             # summary + period grouping
│   │   ├── foodplans/           # food plan + days + meals
│   │   ├── profile/             # user profile upsert
│   │   └── calories/            # AI calorie estimation
│   ├── infra/
│   │   ├── auth/                # JWT service + bcrypt
│   │   ├── ai/                  # Google Gemini client
│   │   ├── database/            # pgxpool connection
│   │   └── postgres/
│   │       ├── sqlc/            # auto-generated query code
│   │       └── repositories/    # typed wrappers over sqlc
│   ├── http/
│   │   ├── handlers/            # one file per domain
│   │   ├── middleware/          # JWT auth middleware
│   │   └── routes/              # chi router setup
│   └── shared/errors/           # sentinel errors (NotFound, Unauthorized…)
├── db/
│   ├── migrations/              # goose SQL migration files
│   └── queries/                 # SQL queries (input for sqlc)
├── docker-compose.yml
├── sqlc.yaml
├── Makefile
├── .env.example
└── requests.http                # REST Client test file (VS Code)
```

---

## Getting Started

### Prerequisites

- Go 1.22+
- Docker + Docker Compose
- [goose](https://github.com/pressly/goose) — `go install github.com/pressly/goose/v3/cmd/goose@latest`
- [sqlc](https://sqlc.dev) — `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

### Setup

```bash
# 1. Clone and enter the directory
cd BE-Daily-Diet

# 2. Copy environment file
cp .env.example .env
# Fill in JWT_SECRET and GEMINI_API_KEY

# 3. Start PostgreSQL
docker compose up -d

# 4. Run migrations
goose -dir db/migrations postgres "postgres://docker:docker@localhost:5432/daily_diet?sslmode=disable" up

# 5. Start the API
go run ./cmd/api
```

The API will be available at `http://localhost:3333`.

### Environment Variables

```env
PORT=3333
DATABASE_URL=postgres://docker:docker@localhost:5432/daily_diet?sslmode=disable
JWT_SECRET=
JWT_ACCESS_TOKEN_EXPIRES_IN_MINUTES=30
JWT_REFRESH_TOKEN_EXPIRES_IN_DAYS=7
GEMINI_API_KEY=        # get yours at aistudio.google.com (free)
```

---

## Common Commands

```bash
# Start database
docker compose up -d

# Run API
go run ./cmd/api

# Create migration
goose -dir db/migrations create migration_name sql

# Apply migrations
goose -dir db/migrations postgres "$DATABASE_URL" up

# Check migration status
goose -dir db/migrations postgres "$DATABASE_URL" status

# Regenerate sqlc code (after editing .sql query files)
sqlc generate

# Run tests
go test ./...

# Build binary
go build -o bin/api ./cmd/api
```

---

<div align="center">

[![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com) &nbsp;
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com) &nbsp;
[![forthebadge](https://forthebadge.com/images/badges/open-source.svg)](https://forthebadge.com)

</div>
