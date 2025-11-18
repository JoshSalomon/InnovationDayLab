# Research: Task Management Web Application

**Date**: 2025-11-09  
**Feature**: Task Management Web Application  
**Purpose**: Resolve technical decisions marked as NEEDS CLARIFICATION in plan.md

## Research Tasks

### 1. Go Version Selection

**Decision**: Go 1.21 or later

**Rationale**: 
- Go 1.21+ provides stable support for embedded filesystems (embed package), which is essential for embedding frontend static assets
- Good balance between modern features and stability
- Wide compatibility with dependency ecosystem
- Standard library includes robust HTTP server capabilities

**Alternatives considered**:
- Go 1.20: Missing some newer standard library improvements
- Go 1.22+: Latest features but may have less ecosystem stability

### 2. Web Framework Selection

**Decision**: Standard library `net/http` with minimal routing helper (e.g., `gorilla/mux` or `chi`)

**Rationale**:
- Constitution requires single executable with minimal dependencies
- Standard library `net/http` is sufficient for basic routing and handlers
- Lightweight router libraries (gorilla/mux, chi) provide clean URL routing without heavy framework overhead
- Avoids large framework dependencies that increase binary size
- Simple middleware pattern for authentication

**Alternatives considered**:
- Gin/Echo: More features but heavier dependencies, violates simplicity principle
- Full-featured frameworks: Unnecessary complexity for this use case
- Pure standard library: Possible but routing becomes verbose; small router library improves maintainability

**Selected**: `chi` router (github.com/go-chi/chi/v5) - lightweight, standard library compatible, minimal dependencies

### 3. Authentication Library

**Decision**: Session-based authentication using `gorilla/sessions` with secure cookie storage

**Rationale**:
- Specification assumes session-based authentication (Assumptions section)
- `gorilla/sessions` is lightweight, well-maintained, and open source
- Provides secure cookie-based sessions suitable for web applications
- No external dependencies (uses standard library crypto)
- Supports session storage in-memory (sufficient for single-executable requirement)

**Alternatives considered**:
- JWT tokens: More complex, requires token management, not needed for session-based auth
- OAuth2: Overkill for local application, adds external dependencies
- Custom session implementation: Reinventing the wheel, security risks

**Selected**: `gorilla/sessions` (github.com/gorilla/sessions) - simple, secure, minimal dependencies

### 4. Storage/Database Solution

**Decision**: SQLite with `modernc.org/sqlite` (pure Go SQLite driver) or `go-sqlite3` with CGO, with file-based persistence

**Rationale**:
- Constitution requires in-memory database within executable OR local disk storage
- SQLite provides both: can run in-memory or file-based, embedded in application
- File-based SQLite meets "local disk storage" requirement
- `modernc.org/sqlite` is pure Go (no CGO), aligns with single-executable goal
- Provides ACID transactions, SQL interface, and relational data modeling
- Well-suited for task dependencies (foreign keys, joins)

**Alternatives considered**:
- In-memory only (map-based): Would lose data on restart, violates persistence requirement (FR-013, FR-014)
- JSON file storage: Simpler but lacks query capabilities, harder to enforce relationships
- BadgerDB/BoltDB: Key-value stores, less suitable for relational task dependencies
- External database: Violates constitution (no external processes)

**Selected**: `modernc.org/sqlite` (pure Go SQLite) - no CGO dependencies, file-based persistence, relational capabilities

### 5. Testing Framework

**Decision**: Standard library `testing` package with `testify/assert` for assertions

**Rationale**:
- Go standard library `testing` is sufficient for test execution
- `testify/assert` provides readable assertions without heavy framework overhead
- `testify/suite` optional for test organization
- Minimal dependency addition
- Standard library testing integrates with `go test` toolchain

**Alternatives considered**:
- Pure standard library: Possible but assertions become verbose
- Ginkgo/Gomega: BDD-style but adds significant complexity and dependencies
- Custom assertion helpers: Unnecessary when testify provides clean API

**Selected**: `testify` (github.com/stretchr/testify) - minimal, widely used, open source

### 6. Frontend Framework Selection

**Decision**: Vanilla JavaScript with minimal framework, or lightweight option like Alpine.js

**Rationale**:
- Constitution emphasizes simplicity (YAGNI principle)
- Vanilla JS reduces dependencies and binary size
- For single executable, frontend assets are embedded static files
- Alpine.js provides reactivity without build toolchain complexity
- No Node.js/npm build step required - can write directly or use simple bundler

**Alternatives considered**:
- React/Vue/Angular: Heavy dependencies, requires build toolchain, violates simplicity
- HTMX: Interesting but adds HTTP dependency complexity
- Server-side rendering (Go templates): Simpler but less interactive UX
- Vanilla JS: Most aligned with constitution, sufficient for task management UI

**Selected**: Vanilla JavaScript with Alpine.js for minimal reactivity - no build step, embedded in binary

### 7. Performance Goals & Scale

**Decision**: 
- Target: Support 50-100 concurrent users (reasonable for local/personal use)
- Response time: <2s for task list load (per SC-001)
- Task creation: <30s user time (per SC-002), <500ms server processing
- No explicit throughput target (not a high-traffic public service)

**Rationale**:
- Specification success criteria define user-facing performance (SC-001, SC-002)
- Constitution emphasizes "single-user or small-scale personal use scenarios"
- SQLite handles 50-100 concurrent connections effectively
- No need to optimize for enterprise scale (violates simplicity principle)

**Alternatives considered**:
- Higher concurrency targets: Unnecessary for local/personal use case
- Lower targets: Too restrictive, may not meet multi-user requirement
- No targets: Need some guidance for design decisions

### 8. Password Hashing

**Decision**: `golang.org/x/crypto/bcrypt` for password hashing

**Rationale**:
- Standard Go crypto package, well-maintained
- Industry-standard bcrypt algorithm for password security
- No external C dependencies (pure Go)
- Appropriate security for local application
- Open source, compatible license

**Alternatives considered**:
- Argon2: More modern but adds dependency complexity
- SHA-256: Insecure for passwords (no salting, fast hashing)
- External auth service: Violates single-executable requirement

**Selected**: `golang.org/x/crypto/bcrypt` - standard, secure, pure Go

### 9. Static Asset Embedding

**Decision**: Go `embed` package (standard library, Go 1.16+)

**Rationale**:
- Standard library feature, no dependencies
- Compiles frontend assets directly into binary
- Meets single-executable requirement
- Simple to use: `//go:embed static/*`

**Alternatives considered**:
- External file serving: Violates single-executable requirement
- Base64 encoding: Inefficient, harder to maintain
- Third-party embedding tools: Unnecessary when standard library provides solution

### 10. License Compliance

**Decision**: Use `go-licenses` tool or manual license file generation

**Rationale**:
- Constitution requires license report mechanism
- `go-licenses` (github.com/google/go-licenses) can scan dependencies
- Generate report as part of build process
- Maintain LICENSE file in repository root
- Document all third-party dependencies

**Alternatives considered**:
- Manual tracking: Error-prone, doesn't scale
- No tracking: Violates constitution requirement
- Commercial tools: Unnecessary for open source compliance

**Selected**: `go-licenses` tool for automated license reporting

## Summary of Technology Choices

| Component | Selected Technology | License | Rationale |
|-----------|-------------------|---------|-----------|
| Language | Go 1.21+ | BSD-3-Clause | Constitution requirement, standard library features |
| Web Framework | chi router | MIT | Lightweight, standard library compatible |
| Authentication | gorilla/sessions | BSD-3-Clause | Session-based, secure, minimal |
| Storage | modernc.org/sqlite | BSD-3-Clause | Pure Go, file-based, relational |
| Testing | testify | MIT | Minimal assertions, standard library testing |
| Frontend | Vanilla JS + Alpine.js | MIT | Simple, no build step, embedded |
| Password Hashing | golang.org/x/crypto/bcrypt | BSD-3-Clause | Standard, secure, pure Go |
| Asset Embedding | Go embed (stdlib) | BSD-3-Clause | Standard library, no dependencies |
| License Tool | go-licenses | Apache-2.0 | Automated compliance reporting |

## Open Source License Compatibility

All selected dependencies use permissive licenses (MIT, BSD-3-Clause, Apache-2.0) that are compatible with open source distribution. No GPL or copyleft licenses that would require special handling.

## Next Steps

Proceed to Phase 1: Design & Contracts
- Generate data-model.md with User and Task entities
- Create API contracts for HTTP endpoints
- Generate quickstart.md with setup instructions
