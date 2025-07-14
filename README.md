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

4. **API Documentation:**
   Visit [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) to view the interactive Swagger documentation.

### Docker

1. Build and start services:
   ```sh
   docker-compose up --build
   ```

## API Documentation

The API documentation is available via Swagger UI at `/swagger/index.html` when the service is running.

### Base URL
```
http://localhost:8080/api
```

### Authentication

Most endpoints require Bearer token authentication. To use protected endpoints:

1. Register a new user via `POST /api/register`
2. Login via `POST /api/login` to get a JWT token
3. Include the token in the Authorization header: `Bearer <your-jwt-token>`

### Available Endpoints

#### Health Check
- `GET /health` - Check service health (no auth required)

#### Authentication & Users

##### Register User
- **POST** `/api/register`
- **Description:** Create a new user account
- **Request Body:**
  ```json
  {
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123",
    "name": "John Doe",
    "birthdate": "1990-01-01",
    "phone": "+1234567890",
    "address": "123 Main St, City, State"
  }
  ```
- **Response:** `201 Created`
  ```json
  {
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user",
      "name": "John Doe",
      "birthdate": "1990-01-01",
      "phone": "+1234567890",
      "address": "123 Main St, City, State",
      "created_at": "2025-07-14T10:30:00Z",
      "updated_at": "2025-07-14T10:30:00Z"
    }
  }
  ```

##### Login User
- **POST** `/api/login`
- **Description:** Authenticate user and get JWT token
- **Request Body:**
  ```json
  {
    "username": "john_doe",
    "password": "password123"
  }
  ```
- **Response:** `200 OK`
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user",
      "name": "John Doe"
    }
  }
  ```

##### Get User
- **GET** `/api/users/{id}` 🔒
- **Description:** Get user information by ID
- **Headers:** `Authorization: Bearer <token>`
- **Response:** `200 OK`
  ```json
  {
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user",
      "name": "John Doe",
      "birthdate": "1990-01-01",
      "phone": "+1234567890",
      "address": "123 Main St, City, State"
    }
  }
  ```

##### Update User
- **PUT** `/api/users/{id}` 🔒
- **Description:** Update user information
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
  ```json
  {
    "email": "newemail@example.com",
    "name": "Updated Name",
    "phone": "+0987654321"
  }
  ```
- **Response:** `200 OK`
  ```json
  {
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "newemail@example.com",
      "name": "Updated Name",
      "phone": "+0987654321"
    }
  }
  ```

##### Delete User
- **DELETE** `/api/users/{id}` 🔒
- **Description:** Delete user account
- **Headers:** `Authorization: Bearer <token>`
- **Response:** `200 OK`
  ```json
  {
    "message": "User deleted"
  }
  ```

#### Expenses

##### Create Expense
- **POST** `/api/expenses` 🔒
- **Description:** Create a new expense or income record
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
  ```json
  {
    "amount": 50.00,
    "currency": "USD",
    "kind": "expense",
    "type": "food",
    "description": "Lunch at restaurant",
    "date": "2025-07-14T12:00:00Z"
  }
  ```
- **Response:** `201 Created`
  ```json
  {
    "id": 1,
    "user_id": 1,
    "amount": 50.00,
    "currency": "USD",
    "kind": "expense",
    "type": "food",
    "description": "Lunch at restaurant",
    "date": "2025-07-14T12:00:00Z",
    "created_at": "2025-07-14T10:30:00Z",
    "updated_at": "2025-07-14T10:30:00Z"
  }
  ```

##### List Expenses
- **GET** `/api/expenses` 🔒
- **Description:** Get list of expenses with optional filtering
- **Headers:** `Authorization: Bearer <token>`
- **Query Parameters:**
  - `user_id` (int, optional): Filter by user ID
  - `kind` (string, optional): Filter by kind (`expense` or `income`)
  - `type` (string, optional): Filter by type (`food`, `salary`, `transport`, `entertainment`)
  - `from` (string, optional): Start date (YYYY-MM-DD)
  - `to` (string, optional): End date (YYYY-MM-DD)
- **Example:** `GET /api/expenses?kind=expense&type=food&from=2025-07-01&to=2025-07-14`
- **Response:** `200 OK`
  ```json
  [
    {
      "id": 1,
      "user_id": 1,
      "amount": 50.00,
      "currency": "USD",
      "kind": "expense",
      "type": "food",
      "description": "Lunch at restaurant",
      "date": "2025-07-14T12:00:00Z",
      "created_at": "2025-07-14T10:30:00Z",
      "updated_at": "2025-07-14T10:30:00Z"
    }
  ]
  ```

##### Daily Expense Summary
- **GET** `/api/expenses/summary/day` 🔒
- **Description:** Get total expenses for a specific day
- **Headers:** `Authorization: Bearer <token>`
- **Query Parameters:**
  - `date` (string, required): Date in YYYY-MM-DD format
  - `kind` (string, optional): Filter by kind (`expense` or `income`)
  - `type` (string, optional): Filter by type (`food`, `salary`, `transport`, `entertainment`)
- **Example:** `GET /api/expenses/summary/day?date=2025-07-14&kind=expense`
- **Response:** `200 OK`
  ```json
  {
    "total": 150.75
  }
  ```

### Data Types

#### Expense Kinds
- `expense` - Money spent
- `income` - Money received

#### Expense Types
- `food` - Food and dining
- `salary` - Salary and wages
- `transport` - Transportation costs
- `entertainment` - Entertainment expenses

#### User Roles
- `user` - Regular user (default)
- `admin` - Administrator

### Error Responses

All endpoints may return the following error responses:

- **400 Bad Request**
  ```json
  {
    "error": "Invalid request"
  }
  ```

- **401 Unauthorized**
  ```json
  {
    "error": "Invalid credentials"
  }
  ```

- **403 Forbidden**
  ```json
  {
    "error": "You can only access your own data"
  }
  ```

- **404 Not Found**
  ```json
  {
    "error": "Resource not found"
  }
  ```

- **500 Internal Server Error**
  ```json
  {
    "error": "Internal server error"
  }
  ```

### Regenerating API Documentation

If you make changes to the API endpoints, regenerate the Swagger documentation:

```sh
swag init
```

## Project Structure
- `main.go`: Entry point
- `internal/`: Feature modules (vocab, tasks, finance, auth, user)
- `config/`: Configuration files
- `pkg/`: Shared utilities
- `docs/`: Auto-generated Swagger documentation

---

For more details, see inline comments and documentation in each package.
