# mindoh-service

Backend REST API for the Mindoh expense tracker, built with Go + Gin + GORM.

## Tech Stack

- **Go 1.24** — language
- **Gin** — HTTP framework
- **GORM** — ORM (PostgreSQL)
- **Supabase** — managed PostgreSQL (production)
- **Railway** — deployment
- **JWT (HS256)** — authentication; token payload carries `username` and `role` (no numeric user ID)
- **Brevo SMTP** — transactional email (verification, password reset)
- **Swagger** — API docs (`/swagger/index.html`)

## Project Structure

```
mindoh-service/
├── config/           Config loader (config.yaml + env vars)
├── internal/
│   ├── auth/         JWT generation, middleware, role guard
│   ├── currency/     Exchange rate endpoints
│   ├── db/           GORM models
│   ├── dto/          Request / response DTOs
│   ├── expense/      Expense CRUD, summary, groups
│   └── user/         Registration, login, email verification, profile
├── common/utils/     Shared helpers
├── docs/             Swagger generated docs
├── Dockerfile
├── start.sh          Dev runner (ensures correct cwd for .env loading)
└── main.go
```

## API Endpoints

### Auth (public)

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/register | Register new user |
| POST | /api/login | Login — returns JWT + user profile |
| GET | /api/verify-email?token=... | Verify email address |
| POST | /api/resend-verification | Resend verification email |
| POST | /api/forgot-password | Send password-reset email |
| POST | /api/reset-password | Complete password reset with token |

### Users (JWT required)

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/users/me | Current user profile |
| PUT | /api/users/me | Update current user profile |
| POST | /api/users/change-password | Change password (requires current password) |
| GET | /api/users/:id | Get user by ID |
| PUT | /api/users/:id | Update user by ID |
| DELETE | /api/users/:id | Delete user by ID |

### Expenses (JWT required)

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/expenses/ | List expenses (paginated, filtered) |
| POST | /api/expenses/ | Create expense |
| PUT | /api/expenses/:id | Update expense |
| DELETE | /api/expenses/:id | Delete expense |
| GET | /api/expenses/summary | Totals by type and currency |
| GET | /api/expenses/groups | Time-bucketed groups (day/week/month/year) |
| GET | /api/expenses/types | Distinct types for authenticated user |

### Currency (JWT required)

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/currency/exchange-rates | Latest exchange rates |
| GET | /api/currency/currencies | Available currencies |

### Admin (JWT + admin role)

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/admin/users | Create user with explicit role |

> **Note:** User responses never include a numeric `id`. The JWT payload stores `username` instead of a sequential user ID to avoid leaking enumerable identifiers.

## Test Account

A pre-seeded account is available for testing:

| Field | Value |
|---|---|
| Username | `test111` |
| Password | `nvmQF6F2scnn..u` |

## Local Development

### Prerequisites

- Go 1.24+
- PostgreSQL

### Run

```sh
git clone https://github.com/nhathuyle2002/mindoh-service
cd mindoh-service
cp .env.example .env   # fill in values
go mod download
./start.sh             # or: go run main.go
```

> Use `start.sh` so that `godotenv` loads `.env` from the correct working directory.

Server: http://localhost:8080  
Swagger: http://localhost:8080/swagger/index.html

### Environment Variables

| Variable | Description | Example |
|---|---|---|
| APP_ENV | `dev` or `prod` | dev |
| POSTGRES_HOST | DB host | localhost |
| POSTGRES_PORT | DB port | 5431 |
| POSTGRES_USER | DB user | mindoh |
| POSTGRES_PASSWORD | DB password | 1234 |
| POSTGRES_NAME | DB name | mindoh |
| JWT_SECRET | JWT signing secret | your-secret |
| ALLOWED_ORIGINS | CORS origins (comma-separated) | * |
| SMTP_HOST | SMTP server host | smtp-relay.brevo.com |
| SMTP_PORT | SMTP server port | 587 |
| SMTP_USER | SMTP login username | user@smtp-brevo.com |
| SMTP_PASSWORD | SMTP login password | your-smtp-key |
| SMTP_FROM | Verified sender address | you@example.com |
| APP_URL | Frontend base URL (for email links) | http://localhost:5173 |

## Docker

```sh
docker build -t mindoh-service .

docker run -d \
  --name mindoh-backend \
  -p 8080:8080 \
  --env-file .env.docker \
  --add-host=host.docker.internal:host-gateway \
  mindoh-service
```

`.env.docker` is the same as `.env` but with `POSTGRES_HOST=host.docker.internal`.

## Deployment (Railway)

1. Connect the repo to Railway — Dockerfile is auto-detected via `railway.json`
2. Set env vars in the Railway dashboard (no quotes around values)
3. Railway injects `PORT` automatically

Production URL: https://mindoh-service-production.up.railway.app
