# 🏗️ BE-POS (Building Material Point of Sale API)

A robust, enterprise-grade Backend API for a Building Material Point of Sale (POS) system. Built with Go (Golang) using Clean Architecture principles and Uber-Fx for Dependency Injection.

## 🛠️ Tech Stack

* **Language:** Go (Golang) 1.20+
* **Framework:** [Gin](https://github.com/gin-gonic/gin)
* **Dependency Injection:** [Uber-Fx](https://github.com/uber-go/fx)
* **Database:** PostgreSQL
* **ORM:** [Bun](https://bun.uptrace.dev/)
* **Migrations:** [dbmate](https://github.com/amacneil/dbmate)
* **Logging:** [Zap](https://github.com/uber-go/zap)
* **Authentication:** JWT (JSON Web Tokens) with RBAC

---

## 🚀 Getting Started

### 1. Prerequisites
Ensure you have the following installed:
* Go
* PostgreSQL
* `dbmate` (for database migrations)

### 2. Environment Setup
Clone the repository and set up your environment variables:
```bash
cp .env.example .env
```

Fill in your .env file with your local database credentials and a strong JWT secret:
```bash
DATABASE_URL="postgres://username:password@127.0.0.1:5432/pos_db?sslmode=disable"
SERVER_PORT=3000
JWT_SECRET="your_super_secret_key"
```

### 3. Database Migration
Create the database and run all migrations using dbmate:

```bash
dbmate up
```

To create migration file :
```bash
dbmate new migration_name
```

### 4. Run the Server
Start the application. Uber-Fx will automatically resolve and inject all dependencies.
```bash
go run cmd/api/main.go
```

## 📁 Project Structure (Clean Architecture)
This project strictly follows Clean Architecture to ensure separation of concerns.
```
├── cmd/                 # Main applications for this project
│   ├── api/             # API entry point
│   │   └── main.go      # Uber-Fx module registration for API
│   └── worker/          # (Future) Background jobs / Cron workers entry point
├── config/              # Centralized configuration (reads .env)
├── db/                  
│   └── migrations/      # SQL migration files (dbmate)
├── internal/            # Core application code (private packages)
│   ├── domain/          # 1. Core Entities & Interfaces (No external dependencies)
│   ├── repository/      # 2. Data access layer (PostgreSQL / Bun queries)
│   ├── usecase/         # 3. Business logic & rules
│   └── delivery/        # 4. HTTP layer
│       ├── http/
│       │   ├── dto/         # Data Transfer Objects (Request/Response shapes)
│       │   ├── handler/     # Gin HTTP Handlers
│       │   ├── middleware/  # Auth, RBAC, Logging
│       │   ├── router.go    # Endpoint groupings
│       │   └── server.go    # Gin Engine & Graceful Shutdown setup
└── pkg/                 # Shared public utilities (JWT, Response formatter, Validation)
```