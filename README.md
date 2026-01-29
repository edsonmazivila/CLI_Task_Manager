# Task Manager

A production-grade CLI task manager built with Go, featuring clean architecture, real persistence, and enterprise-level practices.

## Features

- **Full CRUD Operations**: Add, list, view, update, complete, and delete tasks
- **Advanced Filtering**: Filter tasks by status, priority, and date range
- **Real Persistence**: SQLite storage with automatic migrations
- **Clean Architecture**: Separation of concerns with clear boundaries
- **Structured Logging**: Built-in structured logging with `slog`
- **Configuration Management**: Environment variables and YAML config support
- **Production-Ready**: No mocks, stubs, or placeholders

## Prerequisites

- Go 1.21 or later
- GCC or compatible C compiler (required for SQLite CGO support)
- SQLite3 development libraries (usually pre-installed on most systems)

## Installation

### From Source

```bash
# Install GCC if not already installed (Ubuntu/Debian)
sudo apt-get install build-essential

# Clone the repository
git clone https://github.com/edson-mazvila/task-manager.git
cd task-manager

# Install dependencies
go mod download

# Build the binary (CGO is required for SQLite)
CGO_ENABLED=1 go build -o task ./cmd/task

# Or use the Makefile
make build

# Optional: Install to system path
sudo mv task /usr/local/bin/
# Or
make install
```

### Quick Install

```bash
# Build and install in one command
go install ./cmd/task@latest
```

## Configuration

Task Manager can be configured using environment variables, a YAML configuration file, or both.

### Environment Variables

Copy the example environment file and customize it:

```bash
cp .env.example .env
```

Available environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_TYPE` | `sqlite` | Database type (sqlite or postgres) |
| `DB_PATH` | `~/.task-manager/tasks.db` | SQLite database file path |
| `DB_HOST` | `localhost` | PostgreSQL host (if using postgres) |
| `DB_PORT` | `5432` | PostgreSQL port (if using postgres) |
| `DB_NAME` | `taskmanager` | PostgreSQL database name |
| `DB_USER` | - | PostgreSQL username |
| `DB_PASSWORD` | - | PostgreSQL password |
| `DB_SSL_MODE` | `disable` | PostgreSQL SSL mode |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |
| `LOG_FORMAT` | `text` | Log format (text or json) |
| `CONFIG_FILE` | `config.yaml` | Path to YAML config file |

### Configuration File

Alternatively, use a YAML configuration file:

```bash
cp config.yaml.example config.yaml
```

Example `config.yaml`:

```yaml
database:
  type: sqlite
  path: ~/.task-manager/tasks.db

logging:
  level: info
  format: text
```

### Configuration Priority

1. Environment variables (highest priority)
2. Configuration file
3. Default values (lowest priority)

## Usage

### Add a Task

```bash
# Add a task with default medium priority
task add "Buy groceries"

# Add a high-priority task
task add "Submit project report" --priority high

# Add a task with description
task add "Fix bug in login" --priority high --description "Users unable to login with email"

# Short flags
task add "Call dentist" -p low -d "Schedule annual checkup"
```

### List Tasks

```bash
# List all tasks
task list

# List only pending tasks
task list --status pending

# List completed tasks
task list --status completed

# List high-priority tasks
task list --priority high

# List tasks created after a specific date
task list --from 2026-01-01

# List tasks in a date range
task list --from 2026-01-01 --to 2026-01-31

# Combine filters
task list --status pending --priority high
```

### View Task Details

```bash
# Get detailed information about a specific task
task get <task-id>
```

### Update a Task

```bash
# Update task title
task update <task-id> --title "New title"

# Update task priority
task update <task-id> --priority high

# Update multiple fields
task update <task-id> --title "Updated title" --description "New description" --priority low
```

### Complete a Task

```bash
# Mark a task as completed
task complete <task-id>
```

### Delete a Task

```bash
# Delete a task permanently
task delete <task-id>
```

### Get Help

```bash
# General help
task --help

# Command-specific help
task add --help
task list --help
```

## Project Structure

```
.
├── cmd/
│   └── task/
│       └── main.go                 # Application entry point
├── internal/
│   ├── cli/
│   │   └── commands.go             # CLI command implementations
│   ├── config/
│   │   └── config.go               # Configuration loading and validation
│   ├── domain/
│   │   ├── task.go                 # Domain models and interfaces
│   │   └── errors.go               # Domain-specific errors
│   ├── repository/
│   │   └── sqlite_task_repository.go # Data access layer
│   ├── service/
│   │   └── task_service.go         # Business logic layer
│   └── storage/
│       └── sqlite.go               # Database initialization and migrations
├── migrations/
│   └── 001_create_tasks_table.sql  # Database schema
├── .env.example                     # Example environment configuration
├── config.yaml.example              # Example YAML configuration
├── .gitignore                       # Git ignore rules
├── go.mod                           # Go module definition
├── go.sum                           # Go module checksums
└── README.md                        # This file
```

## CI/CD

GitHub Actions workflow is configured to:
- Run tests on Go 1.21, 1.22, and 1.23
- Build binaries for multiple platforms
- Run linters (golangci-lint with comprehensive checks)
- Perform security scanning (gosec, govulncheck)
- Generate test coverage reports
- Upload artifacts

Workflow files:
- `.github/workflows/ci.yml` - Main CI pipeline
- `.golangci.yml` - Linter configuration
- `scripts/run-tests.sh` - Test automation
- `scripts/ci-verify.sh` - Local CI verification

## Performance Benchmarks

Current benchmark results (reference hardware):

```
BenchmarkTaskCreation-8     459      2.77 ms/op    1668 B/op    32 allocs/op
BenchmarkTaskQuery-8      83182     14.56 μs/op    1592 B/op    53 allocs/op
BenchmarkListTasks-8       3838      305.8 μs/op  48880 B/op  1444 allocs/op
```

Run benchmarks: `make bench`

## Architecture

The application follows **Clean Architecture** principles with clear separation of concerns:

### Layers

1. **Domain Layer** (`internal/domain/`)
   - Core business entities (Task)
   - Domain interfaces (TaskRepository)
   - Business rules and validation
   - No external dependencies

2. **Service Layer** (`internal/service/`)
   - Business logic implementation
   - Orchestrates operations
   - Uses domain interfaces
   - Implements use cases

3. **Repository Layer** (`internal/repository/`)
   - Data access implementation
   - Implements domain interfaces
   - Database operations
   - Query construction

4. **Storage Layer** (`internal/storage/`)
   - Database connection management
   - Migration execution
   - Connection pooling

5. **CLI Layer** (`internal/cli/`)
   - Command-line interface
   - User input/output
   - Command routing
   - Uses service layer

6. **Configuration Layer** (`internal/config/`)
   - Configuration loading
   - Environment variable handling
   - Validation

### Design Principles

- **Dependency Inversion**: High-level modules don't depend on low-level modules
- **Interface Segregation**: Interfaces are focused and minimal
- **Single Responsibility**: Each component has one reason to change
- **Clean Boundaries**: Clear separation between layers
- **Testability**: Components can be tested independently

## Database Schema

The application uses the following database schema:

```sql
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL CHECK (status IN ('pending', 'completed')),
    priority TEXT NOT NULL CHECK (priority IN ('low', 'medium', 'high')),
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    completed_at DATETIME
);

CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_created_at ON tasks(created_at);
```

## Error Handling

The application provides clear, actionable error messages:

- **Invalid Input**: Descriptive validation errors
- **Not Found**: Clear indication when a task doesn't exist
- **Configuration Errors**: Helpful messages for misconfiguration
- **Database Errors**: Informative error messages without exposing internals

## Logging

Structured logging is implemented using Go's `slog` package:

- **Levels**: debug, info, warn, error
- **Formats**: text (human-readable) or json (machine-parseable)
- **Context**: All logs include relevant context (task IDs, operations, etc.)

## Development

### Building

```bash
# Build for current platform (requires CGO)
CGO_ENABLED=1 go build -o task ./cmd/task

# Or use Makefile
make build

# Build for Linux
make build-linux

# Build for macOS
make build-darwin

# Build for Windows
make build-windows
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run integration tests only
make test-integration

# Run all tests (unit + integration)
make test-all

# Run benchmarks
make bench

# Run full test suite with report
./scripts/run-tests.sh

# Run CI verification locally
./scripts/ci-verify.sh
```

#### Integration Tests

The project includes comprehensive integration tests that verify:

- **Complete task lifecycle** (create, read, update, complete, delete)
- **Advanced filtering** (status, priority, date ranges, combined filters)
- **Concurrent operations** (thread safety with 50 goroutines)
- **Error handling** (all error scenarios)
- **Database persistence** (data survives reconnections)
- **Migration idempotency** (migrations run safely multiple times)
- **Transaction rollback** (data integrity)
- **Timestamp correctness** (created_at, updated_at, completed_at)
- **Task ordering** (descending by creation date)
- **Database indexes** (performance optimization)
- **Configuration loading** (env vars, YAML files, precedence)

**Test Results**: 16/16 tests passing, 77.8% code coverage

See [TEST_SUMMARY.md](TEST_SUMMARY.md) for detailed test report.

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Check for vulnerabilities
govulncheck ./...
```

## Production Considerations

### Security

- Database path is validated and sanitized
- SQL injection prevention through parameterized queries
- No hardcoded credentials
- Proper error handling without exposing internals

### Performance

- Connection pooling configured appropriately
- Indexed columns for efficient queries
- Minimal memory allocations
- Efficient query construction

### Reliability

- Transactional migrations
- Graceful error handling
- Proper resource cleanup
- Context-aware operations for cancellation

### Observability

- Structured logging throughout
- Operation tracking with context
- Error logging with relevant details

## Troubleshooting

### CGO Required Error

If you see an error about CGO being disabled:

```bash
# Ensure CGO is enabled during build
CGO_ENABLED=1 go build -o task ./cmd/task

# On Ubuntu/Debian, install build tools
sudo apt-get install build-essential

# On macOS, install Xcode Command Line Tools
xcode-select --install

# On Windows, install MinGW-w64 or TDM-GCC
```

### Database Locked Error

If you encounter a "database is locked" error:

1. Ensure no other instances are running
2. Check file permissions on the database directory
3. SQLite uses a single writer at a time

### Permission Denied

If you get permission errors:

```bash
# Ensure the database directory exists and is writable
mkdir -p ~/.task-manager
chmod 755 ~/.task-manager
```

### Binary Not Found After Install

Ensure your `$GOPATH/bin` is in your `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please ensure:

1. All code follows Go best practices
2. Tests are included for new features
3. Documentation is updated
4. Code is formatted with `go fmt`
5. No breaking changes without discussion

## Support

For issues, questions, or contributions, please open an issue on GitHub.
