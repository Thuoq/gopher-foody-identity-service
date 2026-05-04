# Gopher Foody Identity Service

A production-ready Authentication and Identity microservice built with Go, implementing Clean Architecture.
It handles user registration, secure sign-in with Argon2id, and JWT session management.

## 📂 Project Structure

```text
.
├── cmd/server/main.go       # Application entry point
├── internal/
│   ├── application/         # Business logic (Use cases)
│   ├── core/                # Domain entities & repository interfaces
│   ├── infrastructure/      # DB, Repositories, JWT implementation
│   └── presentation/        # HTTP & gRPC layers
├── migrations/              # SQL migration files
└── pkg/                     # Reusable utilities
```

## 🚀 Getting Started

### 1. Prerequisites
- Go 1.26+
- PostgreSQL
- Docker (for migrations)

### 2. Environment Setup
Create a `.env` file from `.env.example`:
```env
APP_HTTP_PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable
JWT_ACCESS_SECRET=your-access-secret
JWT_REFRESH_SECRET=your-refresh-secret
```

### 3. Database Migrations
We use `golang-migrate` via Docker. Make sure your database is running before executing:

- **Run all migrations up**:
  ```bash
  make migrate-up
  ```
- **Rollback last migration**:
  ```bash
  make migrate-down
  ```
- **Create a new migration**:
  ```bash
  make migrate-create name=create_users_table
  ```

### 4. Running the Service
```bash
go mod tidy
go run cmd/server/main.go
```

The service will be available at `http://localhost:8080`.
