# Task Management Web Application

A multi-user task management web application built with Go, featuring authentication, task CRUD operations, task dependencies, and a modern web interface.

## Features

### Core Functionality
- ✅ **User Authentication** - Login/logout with session-based authentication
- ✅ **User Management** - Admin can create and manage users
- ✅ **Task Management** - Full CRUD operations for tasks
- ✅ **Task Dependencies** - Create and manage task dependencies with circular dependency detection
- ✅ **Task Status Tracking** - Track task status (not_started, in_progress, completed, blocked, deferred)
- ✅ **Progress Tracking** - Visual progress bars (0-100%)
- ✅ **Due Date Management** - Set and track task due dates
- ✅ **Task Filtering** - Filter tasks by status
- ✅ **User Roles** - Admin and regular user roles with appropriate permissions

### User Interface
- ✅ **Modern Web UI** - Clean, responsive interface built with Alpine.js
- ✅ **Loading States** - Visual feedback for async operations
- ✅ **Progress Visualization** - Visual progress bars with color coding
- ✅ **Status Badges** - Color-coded status indicators
- ✅ **Input Validation** - Client-side and server-side validation
- ✅ **Error Handling** - User-friendly error messages

### Technical Features
- ✅ **SQLite Database** - Embedded database with automatic migrations
- ✅ **RESTful API** - Clean API endpoints for all operations
- ✅ **Session Management** - Secure session-based authentication
- ✅ **Password Hashing** - Bcrypt password hashing
- ✅ **Static File Embedding** - Frontend assets embedded in binary
- ✅ **Graceful Shutdown** - Proper server shutdown handling
- ✅ **Health Check Endpoint** - `/api/health` for monitoring

## Prerequisites

- **Go 1.23.9** or compatible version
- **SQLite** (included via `modernc.org/sqlite` - pure Go, no CGO required)
- **Make** (optional, for convenience commands)

## Project Structure

```
backend/src/
├── main.go              # Application entry point
├── config/              # Configuration management
├── storage/             # Database layer (migrations, initialization)
├── models/              # Data models (User, Task, TaskDependency)
├── services/            # Business logic (auth, tasks, users, dependencies)
├── api/                 # HTTP handlers and routing
│   ├── handlers/       # Request handlers (auth, tasks, users, dependencies, health)
│   ├── middleware/     # Auth and admin middleware
│   └── response/        # Response utilities
└── static/              # Frontend assets (embedded in binary)
    ├── index.html       # Homepage
    ├── login.html       # Login page
    ├── tasks.html       # Task management page
    ├── user-management.html  # User management page (admin only)
    └── css/
        └── style.css    # Stylesheet
```

## Building the Application

### Build from Source

```bash
cd backend/src
go build -o task-management-app
```

This creates a single executable binary (`task-management-app`) that includes:
- All Go dependencies compiled in
- Frontend static files embedded via `//go:embed`
- SQLite driver (pure Go, no external dependencies)

### Build for Different Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o task-management-app-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o task-management-app-macos

# Windows
GOOS=windows GOARCH=amd64 go build -o task-management-app.exe
```

## Configuration

The application can be configured via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DATABASE_PATH` | `./data/tasks.db` | SQLite database file path |
| `SESSION_KEY` | `change-me-in-production-secret-key-min-32-chars` | Secret key for session cookies (must be at least 32 characters) |
| `ADMIN_USERNAME` | `admin` | Initial admin username |
| `ADMIN_PASSWORD` | `admin123` | Initial admin password (⚠️ change in production) |
| `ADMIN_EMAIL` | `admin@example.com` | Initial admin email |

### Example Configuration

```bash
export PORT=3000
export DATABASE_PATH=/var/lib/task-app/tasks.db
export SESSION_KEY="your-secret-key-minimum-32-characters-long"
export ADMIN_PASSWORD="secure-password-here"
```

## Running the Application

### Local Development

1. **Build the application:**
   ```bash
   cd backend/src
   go build -o task-management-app
   ```

2. **Run the server:**
   ```bash
   ./task-management-app
   ```
   
   Or with custom configuration:
   ```bash
   PORT=3000 DATABASE_PATH=./mydata.db ./task-management-app
   ```

3. **Access the application:**
   - Web UI: http://localhost:8080
   - Login page: http://localhost:8080/login.html
   - Tasks page: http://localhost:8080/tasks.html
   - API: http://localhost:8080/api/*
   - Health check: http://localhost:8080/api/health

### Default Credentials

After first startup, an admin user is automatically created:
- **Username:** `admin` (or value of `ADMIN_USERNAME`)
- **Password:** `admin123` (or value of `ADMIN_PASSWORD`)
- **Email:** `admin@example.com` (or value of `ADMIN_EMAIL`)

⚠️ **Important:** Change the default password in production!

## API Endpoints

### Authentication
- `POST /api/login` - Login with username and password
- `POST /api/logout` - Logout current user
- `GET /api/me` - Get current user information

### Tasks (requires authentication)
- `GET /api/tasks` - List all tasks (supports `?status=` filter)
- `POST /api/tasks` - Create a new task
- `GET /api/tasks/:id` - Get a specific task
- `PUT /api/tasks/:id` - Update a task
- `DELETE /api/tasks/:id` - Delete a task

### Task Dependencies (requires authentication)
- `GET /api/tasks/:id/dependencies` - Get task dependencies
- `POST /api/tasks/:id/dependencies` - Add a dependency
- `DELETE /api/tasks/:id/dependencies/:depId` - Remove a dependency

### User Management (admin only)
- `GET /api/users` - List all users
- `POST /api/users` - Create a new user
- `GET /api/users/:id` - Get a specific user
- `PUT /api/users/:id` - Update a user
- `DELETE /api/users/:id` - Delete a user

### System
- `GET /api/health` - Health check endpoint

## Usage Examples

### Creating a Task

```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -H "Cookie: session=your-session-cookie" \
  -d '{
    "description": "Complete project documentation",
    "due_date": "2024-12-31",
    "status": "in_progress",
    "progress": 50
  }'
```

### Adding a Task Dependency

```bash
curl -X POST http://localhost:8080/api/tasks/2/dependencies \
  -H "Content-Type: application/json" \
  -H "Cookie: session=your-session-cookie" \
  -d '{"dependency_id": 1}'
```

### Filtering Tasks by Status

```bash
curl http://localhost:8080/api/tasks?status=completed \
  -H "Cookie: session=your-session-cookie"
```

## Production Deployment

### 1. Set Secure Environment Variables

```bash
export SESSION_KEY="$(openssl rand -base64 32)"
export ADMIN_PASSWORD="$(openssl rand -base64 16)"
export DATABASE_PATH="/var/lib/task-app/tasks.db"
export PORT=8080
```

### 2. Create Data Directory

```bash
sudo mkdir -p /var/lib/task-app
sudo chown $USER:$USER /var/lib/task-app
```

### 3. Build and Deploy

```bash
cd backend/src
go build -o /usr/local/bin/task-management-app
```

### 4. Run as a Service (systemd example)

Create `/etc/systemd/system/task-app.service`:

```ini
[Unit]
Description=Task Management Application
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/usr/local/bin
ExecStart=/usr/local/bin/task-management-app
Environment="PORT=8080"
Environment="DATABASE_PATH=/var/lib/task-app/tasks.db"
Environment="SESSION_KEY=your-secret-key-here"
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable task-app
sudo systemctl start task-app
sudo systemctl status task-app
```

### Docker Deployment (Optional)

Create a `Dockerfile`:

```dockerfile
FROM golang:1.23.9-alpine AS builder
WORKDIR /app
COPY backend/src/ .
RUN go build -o task-management-app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/task-management-app .
EXPOSE 8080
CMD ["./task-management-app"]
```

Build and run:

```bash
docker build -t task-management-app .
docker run -p 8080:8080 \
  -e SESSION_KEY="your-secret-key" \
  -e ADMIN_PASSWORD="secure-password" \
  -v $(pwd)/data:/root/data \
  task-management-app
```

## Testing

### Manual Testing

#### 1. Test Login

```bash
# Successful login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  -c cookies.txt

# Expected response:
# {"id":1,"username":"admin","email":"admin@example.com","user_type":"admin",...}
```

#### 2. Test Task Creation

```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "description": "Test task",
    "status": "not_started",
    "progress": 0
  }'
```

#### 3. Test Health Check

```bash
curl http://localhost:8080/api/health

# Expected response:
# {"status":"ok"}
```

## Troubleshooting

### Server won't start

**Error:** `bind: address already in use`
- **Solution:** Another process is using port 8080. Change the port:
  ```bash
  PORT=3000 ./task-management-app
  ```

**Error:** `unable to open database file: out of memory`
- **Solution:** Ensure the database directory exists and is writable:
  ```bash
  mkdir -p data
  chmod 755 data
  ```

### Authentication issues

**Problem:** Login works but authenticated endpoints return "Not authenticated"
- **Check:** Ensure cookies are being sent with requests:
  ```bash
  curl -v -X GET http://localhost:8080/api/me -b cookies.txt
  ```

### Database issues

**Problem:** Database file not created
- **Solution:** Check directory permissions and ensure the path is writable:
  ```bash
  ls -ld data/
  chmod 755 data/
  ```

**Problem:** Migrations fail
- **Solution:** Delete the database file and restart:
  ```bash
  rm data/tasks.db
  ./task-management-app
  ```

## Development

### Code Structure

- **Models:** Data structures (`models/user.go`, `models/task.go`, `models/task_dependency.go`)
- **Services:** Business logic (`services/auth_service.go`, `services/task_service.go`, `services/user_service.go`, `services/dependency_service.go`)
- **Handlers:** HTTP request handlers (`api/handlers/*.go`)
- **Middleware:** Request middleware (`api/middleware/auth.go`, `api/middleware/admin.go`)
- **Storage:** Database operations (`storage/database.go`, `storage/migrations.go`)

### Adding New Features

1. Define data models in `models/`
2. Implement business logic in `services/`
3. Create HTTP handlers in `api/handlers/`
4. Register routes in `api/router.go`
5. Update frontend in `static/` if needed

### Dependencies

Key dependencies:
- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/gorilla/sessions` - Session management
- `modernc.org/sqlite` - SQLite driver (pure Go)
- `golang.org/x/crypto/bcrypt` - Password hashing
- Alpine.js (via CDN) - Frontend framework

## License

See [LICENSE](LICENSE) file for details.

## Status

**Current Progress:** Core features implemented and functional

- ✅ Phase 1: Setup
- ✅ Phase 2: Foundational
- ✅ Phase 3-9: User Stories (Core functionality)
- 🔄 Phase 11: Polish & Cross-Cutting Concerns (Partially complete)

### Completed Features
- User authentication and session management
- Task CRUD operations
- Task dependencies with validation
- User management (admin only)
- Modern web UI with Alpine.js
- Loading states and visual feedback
- Progress bars and status badges
- Input validation
- Graceful shutdown
- Health check endpoint

### In Progress / Planned
- Responsive design for mobile/tablet
- Date picker component
- Keyboard shortcuts
- Enhanced logging
- Database backup/restore
- Additional polish features
