# 📰 CMS Application

A **RESTful API CMS (Content Management System)** built with Go, featuring user authentication, news management, categories, comments, and custom pages.

![CMS Preview](https://github.com/user-attachments/assets/1f7c066c-052b-440b-979a-56b800df27bb)

---

## ⚙️ Tech Stack

- **Language**: Go 1.25
- **Framework**: Gin
- **Database**: PostgreSQL
- **Migration**: golang-migrate
- **Documentation**: Swagger
- **Testing**: Testify
- **Containerization**: Docker & Docker Compose

---

## 🧩 Prerequisites

Before running this project, make sure you have the following installed:

- Go **1.25** or higher
- PostgreSQL **12+**
- Docker & Docker Compose
- Make (optional, for Makefile commands)

---

## 🚀 Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd coding-test
```

### 2. Environment Setup

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` with your own configuration:

```env
APP_NAME=cms
APP_VERSION=0.0.1

HTTP_PORT=8080
LOG_LEVEL=debug

POSTGRES_USER=your_postgres_user
POSTGRES_HOST=localhost
POSTGRES_PASSWORD=your_postgres_password
POSTGRES_DB=cms_db
POSTGRES_PORT=5432
POSTGRES_POOL_MAX=20

ACCESS_TOKEN_SECRET_KEY=your_access_token_secret_key_here
REFRESH_TOKEN_SECRET_KEY=your_refresh_token_secret_key_here
ACCESS_TOKEN_TTL=5m
REFRESH_TOKEN_TTL=24h
```

---

## 🐳 Option A: Run with Docker (Recommended)

This is the easiest way to run the app with all dependencies.

```bash
# Start all services (application + database)
make up
# or
docker-compose up -d
```

The application will be available at 👉 **[http://localhost:8080](http://localhost:8080)**

To view logs:

```bash
docker-compose logs -f cms
```

To stop services:

```bash
docker-compose down
```

To stop and remove volumes (including DB data):

```bash
docker-compose down -v
```

---

## 💻 Option B: Local Development

### 1. Install Dependencies

```bash
go mod download
```

### 2. Setup Database

Make sure PostgreSQL is running and create a database:

```bash
createdb cms_db
```

### 3. Run Database Migrations

```bash
make migrate-up
# or
migrate -path migrations -database "postgres://your_user:your_password@localhost:5432/cms_db?sslmode=disable" up
```

### 4. Run the Application

```bash
make run
# or
go run cmd/app/main.go
```

---

## 📘 API Documentation

Once the app is running, visit:

👉 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

The Swagger UI provides:

- Full API endpoint documentation
- Request/response schemas
- Interactive API testing
- Authentication flow examples

---

## 📡 Available Endpoints

### 🔐 Authentication

| Method | Endpoint               | Description          |
| ------ | ---------------------- | -------------------- |
| POST   | `/api/v1/auth/login`   | User login           |
| POST   | `/api/v1/auth/refresh` | Refresh access token |

### 🗂 Categories

| Method | Endpoint                 | Description                     |
| ------ | ------------------------ | ------------------------------- |
| GET    | `/api/v1/categories`     | Get all categories (public)     |
| GET    | `/api/v1/categories/:id` | Get category by ID (public)     |
| POST   | `/api/v1/categories`     | Create category (auth required) |
| PUT    | `/api/v1/categories/:id` | Update category (auth required) |
| DELETE | `/api/v1/categories/:id` | Delete category (auth required) |

### 📰 News

| Method | Endpoint           | Description                 |
| ------ | ------------------ | --------------------------- |
| GET    | `/api/v1/news`     | Get all news (public)       |
| GET    | `/api/v1/news/:id` | Get news by ID (public)     |
| POST   | `/api/v1/news`     | Create news (auth required) |
| PUT    | `/api/v1/news/:id` | Update news (auth required) |
| DELETE | `/api/v1/news/:id` | Delete news (auth required) |

### 💬 Comments

| Method | Endpoint                    | Description             |
| ------ | --------------------------- | ----------------------- |
| POST   | `/api/v1/news{id}/comments` | Create comment (public) |

### 📄 Custom Pages

| Method | Endpoint                   | Description                        |
| ------ | -------------------------- | ---------------------------------- |
| GET    | `/api/v1/custom-pages`     | Get all custom pages (public)      |
| GET    | `/api/v1/custom-pages/:id` | Get custom page by ID (public)     |
| POST   | `/api/v1/custom-pages`     | Create custom page (auth required) |
| PUT    | `/api/v1/custom-pages/:id` | Update custom page (auth required) |
| DELETE | `/api/v1/custom-pages/:id` | Delete custom page (auth required) |

---

## 🧪 Development

### Run Tests

```bash
make test
# or
go test ./... -v
```

With coverage:

```bash
go test -v -race -covermode atomic -coverprofile=coverage.txt ./internal/...
```

### Run Linter

```bash
make linter-golangci
# or
golangci-lint run
```

### Regenerate Swagger Docs

```bash
make swag-v1
# or
swag init -g internal/controller/http/v1/router.go
```

---

## 🏗 Project Structure

```
.
├── cmd/
│   └── app/              # Application entry point
├── config/               # Configuration files
├── docs/                 # Swagger documentation
├── internal/
│   ├── app/              # Application initialization
│   ├── controller/       # HTTP controllers
│   │   └── http/v1/      # API v1 handlers
│   ├── dto/              # Data Transfer Objects
│   ├── entity/           # Domain entities
│   ├── repository/       # Data access layer
│   │   └── postgres/     # PostgreSQL implementations
│   └── usecase/          # Business logic
├── migrations/           # Database migrations
├── pkg/                  # Shared packages
│   ├── apperror/         # Application errors
│   ├── jwt/              # JWT utilities
│   ├── logger/           # Logger utilities
│   └── postgres/         # PostgreSQL utilities
├── docker-compose.yml    # Docker compose configuration
├── Dockerfile            # Docker image definition
├── Makefile              # Build commands
└── README.md             # Project documentation
```

---

## 🧱 Database Migrations

Migrations are located in the `migrations/` directory and are **automatically run** when using Docker.

To create a new migration:

```bash
migrate create -ext sql -dir migrations -seq migration_name
```

To run migrations manually:

```bash
make migrate-up
```

---

## 🧰 Makefile Commands (Quick Reference)

| Command                | Description                |
| ---------------------- | -------------------------- |
| `make up`              | Run Docker Compose         |
| `make down`            | Stop and remove containers |
| `make run`             | Run app locally            |
| `make test`            | Run tests                  |
| `make migrate-up`      | Apply migrations           |
| `make linter-golangci` | Run linter                 |
| `make swag-v1`         | Regenerate Swagger docs    |

---

---

```

```
