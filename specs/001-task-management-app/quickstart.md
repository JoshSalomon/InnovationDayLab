# Quick Start Guide: Task Management Web Application

**Date**: 2025-11-09  
**Feature**: Task Management Web Application

## Prerequisites

- Go 1.21 or later installed
- Git (for cloning repository)
- Web browser (for accessing the application)

## Installation

### 1. Clone Repository

```bash
git clone <repository-url>
cd InnovationDayLab
```

### 2. Install Dependencies

```bash
# Navigate to backend source directory
cd backend/src

# Download dependencies
go mod download
```

This will download:
- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/gorilla/sessions` - Session management
- `modernc.org/sqlite` - Pure Go SQLite driver
- `golang.org/x/crypto/bcrypt` - Password hashing
- `github.com/stretchr/testify` - Testing assertions (dev dependency)

### 3. Build Application

```bash
# From project root
go build -o taskmanager ./backend/src

# Or from backend/src directory
cd backend/src
go build -o taskmanager
```

This creates a single executable binary `taskmanager` that includes:
- Backend API server
- Embedded frontend static assets
- SQLite database driver

### 4. Initialize Database

On first run, the application will automatically:
- Create database file (default: `./data/tasks.db`)
- Run schema migrations
- Create initial admin user (if configured)

**Note**: Initial admin user setup method depends on implementation:
- Environment variable: `ADMIN_USERNAME`, `ADMIN_PASSWORD`, `ADMIN_EMAIL`
- Configuration file: `config.yaml` or `config.toml`
- First-run setup prompt

### 5. Run Application

```bash
./taskmanager
```

Default server address: `http://localhost:8080`

## Configuration

### Environment Variables

```bash
# Server configuration
PORT=8080                    # HTTP server port (default: 8080)
HOST=localhost              # Server host (default: localhost)

# Database configuration
DB_PATH=./data/tasks.db     # SQLite database file path

# Session configuration
SESSION_SECRET=<random-key> # Secret key for session encryption (required in production)

# Admin user (first run)
ADMIN_USERNAME=admin
ADMIN_PASSWORD=<secure-password>
ADMIN_EMAIL=admin@example.com
```

### Configuration File (Alternative)

Create `config.yaml`:

```yaml
server:
  port: 8080
  host: localhost

database:
  path: ./data/tasks.db

session:
  secret: <random-secret-key>

admin:
  username: admin
  password: <secure-password>
  email: admin@example.com
```

## First Steps

### 1. Access Application

Open browser: `http://localhost:8080`

### 2. Login as Admin

- Username: `admin` (or configured admin username)
- Password: (set during setup)

### 3. Create Regular Users

As admin:
1. Navigate to User Management
2. Click "Create User"
3. Fill in:
   - Username (required)
   - Password (required, min 8 characters)
   - Email (required)
   - Display Name (optional)
4. Click "Create"

### 4. Create Tasks

As regular user:
1. Login with regular user credentials
2. Click "New Task"
3. Fill in:
   - Description (required)
   - Due Date (optional, can be past date with warning)
   - Status (required: completed, in_progress, not_started, blocked, deferred)
   - Progress (required: 0-100)
4. Click "Save"

### 5. Manage Task Dependencies

1. Open a task
2. Navigate to "Dependencies" section
3. Click "Add Dependency"
4. Select task to depend on
5. System prevents circular dependencies automatically

## API Usage

### Authentication

```bash
# Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"johndoe","password":"password123"}' \
  -c cookies.txt

# Logout
curl -X POST http://localhost:8080/api/logout \
  -b cookies.txt
```

### Task Operations

```bash
# List tasks (with session cookie)
curl http://localhost:8080/api/tasks?status=all \
  -b cookies.txt

# Create task
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "description": "Complete project documentation",
    "due_date": "2025-12-31",
    "status": "in_progress",
    "progress": 50
  }'

# Update task
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "status": "completed",
    "progress": 100
  }'

# Delete task
curl -X DELETE http://localhost:8080/api/tasks/1 \
  -b cookies.txt
```

See `contracts/api.yaml` for full API documentation.

## Development

### Run Tests

**Important:** The Go module is located in `backend/src/`, so tests must be run from that directory.

**Option 1: Use the test runner script (recommended)**
```bash
# From project root
./run-tests.sh

# With coverage
./run-tests.sh -cover

# Verbose output
./run-tests.sh -v
```

**Option 2: Navigate to backend/src directory**
```bash
# Navigate to the backend source directory (where go.mod is located)
cd backend/src

# Verify you're in the right place (should show go.mod)
ls go.mod

# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./services

# Verbose output
go test -v ./...
```

**Option 3: From project root (one-liner)**
```bash
cd backend/src && go test ./...
```

**Troubleshooting:**
If you get the error `"pattern ./...: directory prefix . does not contain main module"`, it means you're not in the `backend/src` directory. Make sure you've changed to that directory first.

### Run Development Server

```bash
# From project root
go run ./backend/src/main.go

# Or from backend/src directory
cd backend/src
go run main.go
```

### Database Management

**View database**:
```bash
sqlite3 ./data/tasks.db
.tables
SELECT * FROM users;
SELECT * FROM tasks;
```

**Backup database**:
```bash
cp ./data/tasks.db ./data/tasks.db.backup
```

**Reset database** (⚠️ deletes all data):
```bash
rm ./data/tasks.db
# Restart application to recreate schema
```

## Troubleshooting

### Database File Not Created

- Check write permissions in current directory
- Ensure `data/` directory exists or application can create it
- Verify SQLite driver is working: `go test modernc.org/sqlite`

### Session Not Persisting

- Check `SESSION_SECRET` is set
- Verify cookies are enabled in browser
- Check server logs for session errors

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill process or change PORT environment variable
PORT=8081 ./taskmanager
```

### Admin User Not Created

- Check environment variables or config file
- Verify database migrations ran successfully
- Check application logs for errors
- Manually create admin user via SQL if needed:
  ```sql
  INSERT INTO users (username, password_hash, email, user_type)
  VALUES ('admin', '<bcrypt-hash>', 'admin@example.com', 'admin');
  ```

## License Compliance

Generate license report:

```bash
# Install go-licenses if not installed
go install github.com/google/go-licenses@latest

# Generate report (from backend/src directory)
cd backend/src
go-licenses report ./... > ../../LICENSE_REPORT.txt
```

## Next Steps

- Review `data-model.md` for database schema details
- Review `contracts/api.yaml` for API specification
- Review `plan.md` for implementation architecture
- Proceed to `/speckit.tasks` for task breakdown
