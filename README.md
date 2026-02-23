# mindoh-service

Backend REST API for the Mindoh expense tracker, built with Go + Gin + GORM.

## Tech Stack

- **Go 1.24** — language
- **Gin** — HTTP framework
- **GORM** — ORM (PostgreSQL)
- **Supabase** — managed PostgreSQL (production)
- **Railway** — deployment
- **JWT** — authentication (HS256)
- **Swagger** — API docs (`/swagger/index.html`)

## Project Structure

```
mindoh-service/
├── config/           Config loader (config.yaml + env vars)
├── internal/
│   ├── auth/         JWT middleware, role guard
│   ├── currency/     Exchange rate endpoints
│   ├── db/           GORM models
│   ├── expense/      Expense CRUD, summary, groups
│   └── user/         Registration, login, profile
├── common/utils/     Shared helpers
├── docs/             Swagger generated docs
├── Dockerfile
├── railway.json
└── main.go
```

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | /api/register | | Register new user |
| POST | /api/login | | Login, returns JWT |
| GET | /api/users/me | JWT | Current user profile |
| PUT | /api/users/:id | JWT | Update profile |
| GET | /api/expenses/ | JWT | List expenses (paginated, filtered) |
| POST | /api/expenses/ | JWT | Create expense |
| PUT | /api/expenses/:id | JWT | Update expense |
| DELETE | /api/expenses/:id | JWT | Delete expense |
| GET | /api/expenses/summary | JWT | Totals by type and currency |
| GET | /api/expenses/groups | JWT | Time-bucketed groups (day/week/month/year) |
| GET | /api/expenses/types | JWT | Distinct types for current user |
| GET | /api/currency/exchange-rates | | Latest exchange rates |
| GET | /api/currency/currencies | | Available currencies |

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
cp .env.example .env   # fill in DB values
go mod download
go run main.go
```

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
