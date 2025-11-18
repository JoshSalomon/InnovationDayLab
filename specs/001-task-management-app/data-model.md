# Data Model: Task Management Web Application

**Date**: 2025-11-09  
**Feature**: Task Management Web Application  
**Database**: SQLite (file-based, local disk)

## Entities

### User

Represents an application user with authentication credentials and role-based access.

**Fields**:
- `id` (INTEGER PRIMARY KEY) - Unique user identifier
- `username` (TEXT UNIQUE NOT NULL) - Unique username for login
- `password_hash` (TEXT NOT NULL) - Bcrypt hashed password
- `email` (TEXT NOT NULL) - User email address (required)
- `display_name` (TEXT) - Optional display name for user
- `user_type` (TEXT NOT NULL) - User role: 'admin' or 'regular'
- `created_at` (TIMESTAMP DEFAULT CURRENT_TIMESTAMP) - Account creation timestamp
- `updated_at` (TIMESTAMP DEFAULT CURRENT_TIMESTAMP) - Last update timestamp

**Constraints**:
- Username must be unique (UNIQUE constraint)
- Email is required (NOT NULL)
- User type must be either 'admin' or 'regular' (CHECK constraint or application-level validation)
- Password must be hashed using bcrypt before storage

**Relationships**:
- One-to-many with Task (a user has many tasks)

**Indexes**:
- `idx_user_username` on `username` (for login lookups)
- `idx_user_email` on `email` (for email-based operations if needed)

**Validation Rules** (application-level):
- Username: non-empty, reasonable length (e.g., 3-50 characters)
- Email: valid email format
- Password: minimum strength requirements (e.g., 8+ characters, enforced at creation)

### Task

Represents a work item that belongs to a user and may have dependencies on other tasks.

**Fields**:
- `id` (INTEGER PRIMARY KEY) - Unique task identifier
- `user_id` (INTEGER NOT NULL) - Foreign key to User.id
- `description` (TEXT NOT NULL) - Task description text
- `due_date` (DATE) - Optional due date for the task
- `status` (TEXT NOT NULL) - Task status: 'completed', 'in_progress', 'not_started', 'blocked', 'deferred'
- `progress` (INTEGER NOT NULL DEFAULT 0) - Progress percentage (0-100)
- `created_at` (TIMESTAMP DEFAULT CURRENT_TIMESTAMP) - Task creation timestamp
- `updated_at` (TIMESTAMP DEFAULT CURRENT_TIMESTAMP) - Last update timestamp

**Constraints**:
- User ID is required (NOT NULL, foreign key to User.id)
- Description is required (NOT NULL)
- Status must be one of: 'completed', 'in_progress', 'not_started', 'blocked', 'deferred' (CHECK constraint or application-level validation)
- Progress must be between 0 and 100 inclusive (CHECK constraint: `progress >= 0 AND progress <= 100`)
- Due date can be NULL (optional field)
- Due date can be in the past (application allows with warning per clarification)

**Relationships**:
- Many-to-one with User (each task belongs to one user)
- Many-to-many with Task (task dependencies) via TaskDependency junction table

**Indexes**:
- `idx_task_user_id` on `user_id` (for user task queries)
- `idx_task_status` on `status` (for filtering by status)
- `idx_task_due_date` on `due_date` (for date-based queries/sorting)
- Composite index `idx_task_user_status` on `(user_id, status)` for filtered queries

**Validation Rules** (application-level):
- Description: non-empty, reasonable length (e.g., max 5000 characters)
- Progress: integer between 0-100 (enforced at database and application level)
- Due date: valid date format (application-level validation)
- Past due dates: allowed but warning displayed (per clarification)

### TaskDependency

Junction table representing dependency relationships between tasks.

**Fields**:
- `id` (INTEGER PRIMARY KEY) - Unique dependency relationship identifier
- `task_id` (INTEGER NOT NULL) - Foreign key to Task.id (the dependent task)
- `depends_on_task_id` (INTEGER NOT NULL) - Foreign key to Task.id (the task this depends on)
- `created_at` (TIMESTAMP DEFAULT CURRENT_TIMESTAMP) - Dependency creation timestamp

**Constraints**:
- Both task_id and depends_on_task_id are required (NOT NULL)
- Both tasks must belong to the same user (application-level validation)
- Cannot have self-dependency: `task_id != depends_on_task_id` (CHECK constraint)
- Unique dependency: `(task_id, depends_on_task_id)` must be unique (UNIQUE constraint)
- Foreign keys must reference existing tasks
- Circular dependencies prevented by application-level validation (graph cycle detection)

**Relationships**:
- Many-to-one with Task (task_id) - the task that has the dependency
- Many-to-one with Task (depends_on_task_id) - the task being depended upon

**Indexes**:
- `idx_dependency_task_id` on `task_id` (for finding dependencies of a task)
- `idx_dependency_depends_on` on `depends_on_task_id` (for finding tasks that depend on a task)
- Unique index on `(task_id, depends_on_task_id)` to prevent duplicate dependencies

**Cascade Behavior**:
- When a task is deleted, all TaskDependency records referencing it (either as task_id or depends_on_task_id) are automatically deleted (CASCADE DELETE)
- This implements the clarification: "Allow deletion and automatically remove dependencies"

## Database Schema

### SQL DDL

```sql
-- Users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    email TEXT NOT NULL,
    display_name TEXT,
    user_type TEXT NOT NULL CHECK(user_type IN ('admin', 'regular')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_username ON users(username);
CREATE INDEX idx_user_email ON users(email);

-- Tasks table
CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    description TEXT NOT NULL,
    due_date DATE,
    status TEXT NOT NULL CHECK(status IN ('completed', 'in_progress', 'not_started', 'blocked', 'deferred')),
    progress INTEGER NOT NULL DEFAULT 0 CHECK(progress >= 0 AND progress <= 100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_task_user_id ON tasks(user_id);
CREATE INDEX idx_task_status ON tasks(status);
CREATE INDEX idx_task_due_date ON tasks(due_date);
CREATE INDEX idx_task_user_status ON tasks(user_id, status);

-- Task dependencies junction table
CREATE TABLE task_dependencies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id INTEGER NOT NULL,
    depends_on_task_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (depends_on_task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    CHECK(task_id != depends_on_task_id),
    UNIQUE(task_id, depends_on_task_id)
);

CREATE INDEX idx_dependency_task_id ON task_dependencies(task_id);
CREATE INDEX idx_dependency_depends_on ON task_dependencies(depends_on_task_id);
```

## State Transitions

### Task Status Transitions

Tasks can transition between any status values:
- `not_started` → `in_progress` → `completed`
- `in_progress` → `blocked` → `in_progress`
- `completed` → `in_progress` (reopening completed tasks)
- Any status → `deferred`

**Note**: No restrictions on status transitions are specified in the requirements. Users can change status freely.

### Task Progress Updates

- Progress can be set to any integer value between 0-100
- Progress updates are independent of status changes
- Progress validation: values outside 0-100 are rejected with error (per clarification)

## Data Isolation

**Critical Requirement (FR-008)**: Users can only access their own tasks.

**Implementation**:
- All task queries MUST include `WHERE user_id = ?` with the authenticated user's ID
- Foreign key constraints ensure tasks belong to users
- Application-level authorization checks before any task access
- No cross-user task visibility in queries

## Initial Data

### Admin User

The application starts with a single admin user account (FR-016). This user must be created during initial setup/migration.

**Setup Options**:
1. Pre-configured admin user in database migration
2. First-run setup script that creates admin user
3. Environment variable or config file with initial admin credentials

**Admin User Default** (example):
- Username: `admin` (or configurable)
- Password: Must be set during setup (no default password for security)
- Email: Required, set during setup
- User Type: `admin`

## Data Persistence

**Requirement (FR-013, FR-014)**: All task and user data must persist between server invocations.

**Implementation**:
- SQLite database file stored on local disk (e.g., `./data/tasks.db`)
- Database file location configurable (default: current directory or `./data/`)
- Database file created automatically if it doesn't exist
- All writes are transactional (ACID compliance via SQLite)
- No data loss on server restart (file-based persistence)

## Migration Strategy

**Initial Migration**:
1. Create database file if it doesn't exist
2. Run schema DDL to create tables
3. Create initial admin user (if not exists)
4. Verify schema version

**Future Migrations**:
- Use migration versioning system (e.g., simple version table)
- Apply migrations in order on application startup
- Rollback support for failed migrations
