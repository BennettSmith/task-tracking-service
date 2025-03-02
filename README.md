# Task Tracking Service

A RESTful web service built with Go that manages tasks, following Clean Architecture principles.

## Development Environment Setup

### Prerequisites

1. **Install Homebrew** (Package Manager for macOS)
   ```bash
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```

2. **Install Go**
   ```bash
   brew install go
   ```
   Verify installation:
   ```bash
   go version
   ```

3. **Install Git** (if not already installed)
   ```bash
   brew install git
   ```

4. **Install Visual Studio Code** (recommended editor)
   ```bash
   brew install --cask visual-studio-code
   ```

5. **Install VS Code Go Extension**
   - Open VS Code
   - Press `Cmd+Shift+X`
   - Search for "Go"
   - Install the official Go extension by Google

### Project Setup

1. **Clone the Repository**
   ```bash
   git clone <repository-url>
   cd task-tracking-service
   ```

2. **Install Go Dependencies**
   ```bash
   go mod init task-tracking-service
   go mod tidy
   ```

3. **Install Development Tools**
   ```bash
   # Install common Go tools
   go install golang.org/x/tools/cmd/godoc@latest
   go install github.com/golang/mock/mockgen@latest
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

### Running the Application

1. **Start the Server**
   ```bash
   go run cmd/server/main.go
   ```
   The server will start on `http://localhost:8080`

2. **Run Tests**
   ```bash
   go test ./...
   ```

3. **Generate API Documentation**
   ```bash
   swag init -g cmd/server/main.go
   ```
   Access Swagger UI at `http://localhost:8080/swagger/index.html`

### Development Tools

- **Air** (Live reload for Go apps)
  ```bash
  go install github.com/air-verse/air@latest
  ```
  Run the app with live reload:
  ```bash
  air
  ```

- **Database Setup** (if needed)
  ```bash
  brew install postgresql
  brew services start postgresql
  ```

### IDE Configuration

#### VS Code Settings

Add these settings to your VS Code workspace:

```json
{
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.formatTool": "gofmt",
    "editor.formatOnSave": true,
    "[go]": {
        "editor.defaultFormatter": "golang.go"
    }
}
```

### Useful Commands

- Format code:
  ```bash
  go fmt ./...
  ```
- Run linter:
  ```bash
  golangci-lint run
  ```
- Check test coverage:
  ```bash
  go test -cover ./...
  ```

## Project Structure

### Directory Details

- **cmd/**: Contains the main applications of the project. Each subdirectory represents an executable program.
  
- **internal/**: Contains the private application code that shouldn't be imported by other projects.
  - **core/**: Implements the core business logic following Clean Architecture principles
  - **adapters/**: Contains all the interface adapters that convert data between the core and external agencies
  - **config/**: Handles all configuration-related code

- **pkg/**: Contains code that can be used by external applications. Libraries here should be stable and well-tested.

- **scripts/**: Contains various scripts for building, testing, and deploying the application.

- **test/**: Contains test helpers, fixtures, and integration tests.

- **docs/**: Contains project documentation, including API specifications and architecture decisions.

### Clean Architecture Layers

This project follows Clean Architecture principles with the following layers:

1. **Domain Layer** (internal/core/domain)
   - Contains enterprise business rules
   - Defines core entities and interfaces
   - Has no dependencies on external frameworks

2. **Use Case Layer** (internal/core/services)
   - Contains application-specific business rules
   - Orchestrates the flow of data between entities
   - Implements core business logic

3. **Interface Adapters** (internal/adapters)
   - Converts data between use cases and external agencies
   - Implements repositories and controllers
   - Handles HTTP routing and database operations

4. **External Interfaces** (cmd/server)
   - Contains frameworks and drivers
   - Handles external communications
   - Manages dependency injection

### VS Code Testing Features

#### Running Tests
- Click the "Testing" icon in the left sidebar (looks like a flask/beaker)
- You'll see all test files and individual tests listed
- Click the play button next to any test to run it
- Use these icons above the test tree:
  - ▶️ Run All Tests
  - 🔄 Repeat Last Run
  - ⏯️ Run Failed Tests
  - 🐞 Debug Selected Test

#### Test Codelens
Above each test function, you'll see these clickable links:
- `run test` - Runs this specific test
- `debug test` - Starts the debugger for this test
- `go test` - Shows the exact test command being run

#### Keyboard Shortcuts
- `Cmd+Shift+P` (Mac) or `Ctrl+Shift+P` (Windows/Linux), then type:
  - `Go: Run Test` - Run test at cursor
  - `Go: Run Package Tests` - Run all tests in current package
  - `Go: Run All Tests` - Run all tests in workspace

#### Test Output
- Results appear in the "Test Results" panel
- Green ✓ for passing tests
- Red ✗ for failing tests
- Detailed output shown for failures

#### Settings
Add these to your VS Code settings.json for better testing experience:
```json
{
    "go.testOnSave": false,
    "go.testFlags": ["-v"],
    "go.coverOnSave": true,
    "go.coverageDecorator": {
        "type": "highlight",
        "coveredHighlightColor": "rgba(64,128,128,0.2)",
        "uncoveredHighlightColor": "rgba(128,64,64,0.2)",
        "coveredBorderColor": "rgba(64,128,128,0.4)",
        "uncoveredBorderColor": "rgba(128,64,64,0.4)"
    }
}
```

#### Coverage Visualization
- Run tests with coverage using the "Testing" sidebar
- Code coverage is highlighted in the editor:
  - Green: Covered code
  - Red: Uncovered code

## Docker Support

### Building and Running with Docker

1. **Build the Docker Image**
   ```bash
   docker build -t task-tracking-service .
   ```

2. **Run the Container**
   ```bash
   docker run -d -p 8080:8080 --name task-service task-tracking-service
   ```

3. **View Container Logs**
   ```bash
   docker logs -f task-service
   ```

### API Usage Examples

Here are some example curl commands to interact with the service:

1. **Create a New Task**
   ```bash
   curl -X POST http://localhost:8080/api/v1/tasks \
     -H "Content-Type: application/json" \
     -d '{
       "title": "Complete Project",
       "description": "Finish the task tracking service",
       "due_date": "2024-12-31T23:59:59Z"
     }'
   ```

2. **Get a Task by ID**
   ```bash
   curl http://localhost:8080/api/v1/tasks/{task_id}
   ```

3. **List All Tasks**
   ```bash
   curl http://localhost:8080/api/v1/tasks
   ```

4. **Update a Task**
   ```bash
   curl -X PUT http://localhost:8080/api/v1/tasks/{task_id} \
     -H "Content-Type: application/json" \
     -d '{
       "title": "Complete Project",
       "description": "Updated description",
       "status": "in_progress",
       "due_date": "2024-12-31T23:59:59Z"
     }'
   ```

5. **Delete a Task**
   ```bash
   curl -X DELETE http://localhost:8080/api/v1/tasks/{task_id}
   ```

### Docker Management Commands

- **Stop the Container**
  ```bash
  docker stop task-service
  ```

- **Remove the Container**
  ```bash
  docker rm task-service
  ```

- **View Container Status**
  ```bash
  docker ps -a | grep task-service
  ```

### Running with Docker Compose

1. **Start the Services**
   ```bash
   # Start in detached mode
   docker-compose up -d

   # View logs
   docker-compose logs -f
   ```

2. **Check Service Status**
   ```bash
   docker-compose ps
   ```

3. **Stop the Services**
   ```bash
   # Stop services but preserve volumes
   docker-compose down

   # Stop services and remove volumes
   docker-compose down -v
   ```

4. **Rebuild and Restart**
   ```bash
   # Rebuild images and restart containers
   docker-compose up -d --build
   ```

5. **View Service Logs**
   ```bash
   # View all logs
   docker-compose logs -f

   # View specific service logs
   docker-compose logs -f app
   docker-compose logs -f postgres
   ```

### Development with Docker Compose

- **Access PostgreSQL Database**
  ```bash
  # Using psql from host
  psql postgresql://taskuser:taskpass@localhost:5432/taskdb

  # Using docker exec
  docker exec -it task-service-db psql -U taskuser -d taskdb
  ```

- **Execute Commands in Containers**
  ```bash
  # Run commands in app container
  docker-compose exec app /bin/sh

  # Run commands in database container
  docker-compose exec postgres /bin/sh
  ```

### Environment Configuration

1. **Local Development**
   ```bash
   # Copy example environment file
   cp .env.example .env

   # Edit variables as needed
   vim .env
   ```

2. **Environment Variables**
   
   Key configurations that can be customized:
   - `SERVER_PORT`: HTTP server port (default: 8080)
   - `DB_HOST`: Database host (default: postgres)
   - `DB_PORT`: Database port (default: 5432)
   - `LOG_LEVEL`: Logging level (default: info)
   - See `.env.example` for all available options

3. **Docker Environment**
   ```bash
   # Run with specific environment file
   docker-compose --env-file .env.production up -d

   # Override specific variables
   SERVER_PORT=3000 docker-compose up -d
   ```

4. **Configuration Precedence**
   1. Command-line overrides
   2. Environment variables
   3. `.env` file
   4. Default values
