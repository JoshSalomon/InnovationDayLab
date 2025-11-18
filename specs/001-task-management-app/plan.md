# Implementation Plan: Task Management Web Application

**Branch**: `001-task-management-app` | **Date**: 2025-11-09 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-task-management-app/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build a multi-user task management web application where users can create, edit, delete, and filter tasks with dependency relationships. The application must support user authentication with admin and regular user roles, where admins can create new users and regular users manage their own isolated task lists. Tasks persist between server restarts and include fields for description, due date, status (completed, in progress, not started, blocked, deferred), and progress (0-100%). The application must be implemented as a single Go executable with local disk storage, following the InnovationDayLab constitution requirements.

## Technical Context

**Language/Version**: Go 1.21+ (Golang)  
**Primary Dependencies**: chi router (github.com/go-chi/chi/v5), gorilla/sessions (github.com/gorilla/sessions), modernc.org/sqlite (pure Go SQLite), golang.org/x/crypto/bcrypt  
**Storage**: SQLite database (modernc.org/sqlite) with file-based persistence on local disk  
**Testing**: Go standard library testing package with testify/assert (github.com/stretchr/testify)  
**Target Platform**: Linux/Unix systems (single executable, cross-platform compilation possible)  
**Project Type**: web (backend + frontend in single executable)  
**Performance Goals**: Support 50-100 concurrent users; <2s task list load (SC-001); <500ms server processing for task creation (user time <30s per SC-002)  
**Constraints**: Single executable binary, local disk storage only, no external processes, all dependencies must compile into binary  
**Scale/Scope**: Multi-user application with user data isolation, initial admin user + additional users created by admin, tasks per user (no explicit limit specified)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Pre-Research Check (Phase 0)

**Single Executable**: ✅ PASS - Application will be compiled as single Go binary. No external executables required. Web server and storage will be embedded within the executable.

**Open Source Compliance**: ✅ PASS - All dependencies verified as open source with compatible licenses (MIT, BSD-3-Clause, Apache-2.0). License reporting will use go-licenses tool. See research.md for full dependency list.

**Security & Authentication**: ✅ PASS - Feature specification explicitly requires user authentication (FR-001). Admin and regular user roles must be distinguished. Authentication required before any task data access.

**Personal Data Handling**: ✅ PASS - Feature specification requires user data isolation (FR-008). Each user can only view and manage their own tasks. Data persistence uses local disk storage as required by constitution.

**Simplicity**: ✅ PASS - Application follows YAGNI principles. Core functionality is well-scoped: task CRUD, user management, filtering, dependencies. No unnecessary complexity identified. Future REST API features explicitly deferred to next phase.

**Storage**: ✅ PASS - Feature specification requires data persistence (FR-013, FR-014). Will use local disk storage with in-memory database within executable as required by constitution. No external database processes.

**Technology Stack**: ✅ PASS - Implementation must be in Go (Golang) per constitution. All dependencies must compile into single executable.

### Post-Design Check (Phase 1)

**Single Executable**: ✅ PASS - Design confirmed: Single Go binary with embedded frontend assets via `embed` package. SQLite database file-based (no external process). All dependencies compile into binary.

**Open Source Compliance**: ✅ PASS - All selected dependencies verified:
- chi router: MIT
- gorilla/sessions: BSD-3-Clause  
- modernc.org/sqlite: BSD-3-Clause
- golang.org/x/crypto/bcrypt: BSD-3-Clause
- testify: MIT
- go-licenses: Apache-2.0
All licenses compatible. License reporting mechanism defined in quickstart.md.

**Security & Authentication**: ✅ PASS - Design includes session-based authentication using gorilla/sessions. Password hashing with bcrypt. User roles (admin/regular) enforced at application level. All API endpoints require authentication (see contracts/api.yaml).

**Personal Data Handling**: ✅ PASS - Data model enforces user isolation via foreign keys and application-level authorization checks. All task queries include `WHERE user_id = ?` filter. Database schema includes user_id foreign key constraints. See data-model.md for implementation details.

**Simplicity**: ✅ PASS - Design maintains simplicity:
- Vanilla JS + Alpine.js frontend (no build toolchain)
- Standard library HTTP server with lightweight chi router
- File-based SQLite (no external database)
- Minimal dependencies
- No over-engineering for enterprise scale

**Storage**: ✅ PASS - Design uses SQLite (modernc.org/sqlite) with file-based persistence on local disk. Database file stored at `./data/tasks.db` (configurable). No external database processes. All data persists between server restarts.

**Technology Stack**: ✅ PASS - Implementation in Go 1.21+. All dependencies are Go packages that compile into single executable. Pure Go SQLite driver (no CGO). Frontend assets embedded via standard library `embed` package.

## Project Structure

### Documentation (this feature)

```text
specs/001-task-management-app/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
backend/
├── src/
│   ├── models/          # User, Task data models
│   ├── services/        # Business logic (auth, task management, dependency validation)
│   ├── api/            # HTTP handlers/routes
│   ├── storage/        # Local disk persistence layer
│   └── main.go         # Application entry point
└── tests/
    ├── contract/       # API contract tests
    ├── integration/    # End-to-end tests
    └── unit/           # Unit tests for models/services

frontend/
├── src/
│   ├── components/     # React/Vue/vanilla JS components
│   ├── pages/          # Page components (task list, user management, etc.)
│   └── services/       # API client services
└── tests/              # Frontend tests

# Static assets embedded in binary
static/                 # HTML, CSS, JS files (embedded at build time)
```

**Structure Decision**: Web application structure with separate backend and frontend directories. Backend handles API, authentication, and data persistence. Frontend provides user interface. Both compile into single executable with embedded static assets. This structure allows clear separation of concerns while maintaining single-executable requirement through Go's embedding capabilities.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None identified | - | - |
