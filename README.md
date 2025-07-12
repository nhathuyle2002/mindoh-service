# Mindoh Backend Service

This is the backend service for Mindoh, a personal growth application. It is built with Go, Gin, GORM, JWT authentication, and PostgreSQL. The service supports:

- English vocabulary learning (words, flashcards, quizzes, progress tracking)
- To-do task management (tasks, reminders)
- Expenses & income tracking (transactions, categories, summaries)
- Role-based access control (user/admin)

## Getting Started

### Prerequisites
- Go 1.21+
- Docker & Docker Compose

### Development

1. Install dependencies:
   ```sh
   go mod tidy
   ```
2. Run the service:
   ```sh
   go run main.go
   ```
3. Health check:
   Visit [http://localhost:8080/health](http://localhost:8080/health)

### Docker

1. Build and start services:
   ```sh
   docker-compose up --build
   ```

## Project Structure
- `main.go`: Entry point
- `internal/`: Feature modules (vocab, tasks, finance, auth, user)
- `config/`: Configuration files
- `pkg/`: Shared utilities

---

For more details, see inline comments and documentation in each package.
