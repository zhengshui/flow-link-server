# Flow Link Server (FitEasy Backend)

## Project Overview

**Flow Link Server** is the backend API for the FitEasy fitness tracking application. It is built using **Go (Golang)** and follows strict **Clean Architecture** principles to ensure separation of concerns, testability, and maintainability.

**Key Technologies:**
*   **Language:** Go
*   **Web Framework:** Gin
*   **Database:** MongoDB (Official Go Driver)
*   **Authentication:** JWT (JSON Web Tokens) with Access and Refresh tokens
*   **Configuration:** Viper
*   **Containerization:** Docker & Docker Compose

## Architecture

The project is structured into layers:

1.  **Domain (`domain/`)**: Defines the core business entities (e.g., `User`, `TrainingRecord`, `FitnessPlan`) and interfaces for Repositories and Usecases. This layer is independent of other layers.
2.  **Repository (`repository/`)**: Implements the interfaces defined in the Domain layer to interact with the database (MongoDB).
3.  **Usecase (`usecase/`)**: Contains the business logic. It depends on the Domain layer and Repository interfaces.
4.  **API (`api/`)**:
    *   **Controller (`api/controller/`)**: Handles HTTP requests, calls Usecases, and returns responses.
    *   **Route (`api/route/`)**: Defines API endpoints and maps them to Controllers.
    *   **Middleware (`api/middleware/`)**: Handles cross-cutting concerns like JWT authentication.
5.  **Bootstrap (`bootstrap/`)**: Handles application initialization (Env loading, DB connection).

**Data Flow:**
`Route` -> `Controller` -> `Usecase` -> `Repository` -> `Database`

## Key Entities & Features

*   **User:** Authentication (Signup/Login/Refresh), Profile management.
*   **Training Record:** Log exercises, sets, reps, and weight.
*   **Fitness Plan:** Manage workout plans (custom or from templates).
*   **Plan Template:** Pre-defined workout templates.
*   **Stats:** diverse statistics including training volume, muscle group usage, and personal records.

## Building and Running

### Prerequisites
*   Go 1.19+
*   MongoDB 6.0+
*   Docker & Docker Compose (Optional)

### Environment Setup
1.  Copy `.env.example` to `.env`.
2.  Configure `DB_HOST`, `DB_PORT`, `DB_NAME`, and JWT secrets in `.env`.
    *   For local run without Docker: `DB_HOST=localhost`
    *   For Docker run: `DB_HOST=mongodb`

### Running Locally
```bash
# Install dependencies
go mod download

# Run the application
go run cmd/main.go
```
The API will be available at `http://localhost:8080`.

### Running with Docker
```bash
docker-compose up -d
```

### Building for Production
```bash
go build -o bin/app cmd/main.go
```

### Testing
```bash
# Run all tests
go test ./...
```

### Generating Mocks
The project uses `mockery` for generating mocks for testing.
```bash
# Generate mocks for domain interfaces
mockery --dir=domain --output=domain/mocks --outpkg=mocks --all

# Generate mocks for mongo
mockery --dir=mongo --output=mongo/mocks --outpkg=mocks --all
```

## Development Conventions

*   **Response Format:** All API responses must follow the unified structure defined in `domain/response.go`:
    ```json
    {
      "code": 200,
      "message": "Success",
      "data": { ... }
    }
    ```
*   **Error Handling:** Use `domain/error_response.go` for consistent error returns.
*   **Naming:** Follow Go conventions (CamelCase).
*   **Authentication:** Protected routes must use the `JwtAuthMiddleware`. Access the user ID from the context using `ctx.GetString("x-user-id")`.
*   **Context:** Pass `context.Context` through all layers (Controller -> Usecase -> Repository) to support cancellation and timeouts.

## API Documentation

Refer to `API_DOCUMENTATION.md` for a complete list of endpoints, request/response examples, and error codes.
Refer to `IMPLEMENTATION_SUMMARY.md` for a detailed summary of the implementation status.
