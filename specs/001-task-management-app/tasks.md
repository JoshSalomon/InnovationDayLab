# Tasks: Task Management Web Application

**Input**: Design documents from `/specs/001-task-management-app/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are OPTIONAL - not explicitly requested in feature specification, so test tasks are not included. Focus on implementation tasks.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., [US1], [US2], [US3])
- Include exact file paths in descriptions

## Path Conventions

- **Web app**: `backend/src/`, `frontend/src/`, `static/` at repository root
- Paths follow plan.md structure: backend/src/, frontend/src/, static/

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create project directory structure (backend/src/, frontend/src/, static/, backend/tests/, frontend/tests/) per plan.md
- [x] T002 Initialize Go module with go mod init in backend/src/
- [x] T003 [P] Add chi router dependency (github.com/go-chi/chi/v5) to go.mod
- [x] T004 [P] Add gorilla/sessions dependency (github.com/gorilla/sessions) to go.mod
- [x] T005 [P] Add modernc.org/sqlite dependency to go.mod
- [x] T006 [P] Add golang.org/x/crypto/bcrypt dependency to go.mod
- [x] T007 [P] Add testify dependency (github.com/stretchr/testify) to go.mod
- [x] T008 Create backend/src/main.go with basic HTTP server structure
- [x] T009 Create backend/src/config/config.go for environment configuration management
- [x] T010 [P] Create backend/src/storage/ directory for database layer
- [x] T011 [P] Create backend/src/models/ directory for data models
- [x] T012 [P] Create backend/src/services/ directory for business logic
- [x] T013 [P] Create backend/src/api/ directory for HTTP handlers
- [x] T014 [P] Create static/ directory for frontend assets (HTML, CSS, JS)
- [x] T015 Create frontend/src/ directory structure (components/, pages/, services/)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T016 Create backend/src/storage/database.go with SQLite connection initialization
- [x] T017 Create backend/src/storage/migrations.go with database migration framework
- [x] T018 Implement initial database migration (create users, tasks, task_dependencies tables) in backend/src/storage/migrations.go per data-model.md schema
- [x] T019 Create backend/src/models/user.go with User struct matching data-model.md
- [x] T020 Create backend/src/models/task.go with Task struct matching data-model.md
- [x] T021 Create backend/src/models/task_dependency.go with TaskDependency struct matching data-model.md
- [x] T022 Create backend/src/api/middleware/auth.go with session authentication middleware using gorilla/sessions
- [x] T023 Create backend/src/api/middleware/admin.go with admin authorization middleware
- [x] T024 Create backend/src/api/middleware/user_context.go to extract authenticated user from session
- [x] T025 Create backend/src/api/router.go with chi router setup and middleware chain
- [x] T026 Create backend/src/api/errors.go for standardized error response handling
- [x] T027 Create backend/src/services/password.go with bcrypt password hashing utilities
- [x] T028 Implement database initialization and admin user creation logic in backend/src/storage/setup.go
- [x] T029 Create static/index.html with basic HTML structure and Alpine.js integration
- [x] T030 Create static/css/style.css with basic styling
- [x] T031 Create frontend/src/services/api.js for API client with session cookie handling
- [x] T032 Update backend/src/main.go to embed static files using Go embed package
- [x] T033 Update backend/src/main.go to serve static files and API routes

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Admin Creates New Users (Priority: P1) 🎯 MVP

**Goal**: Admin user can create new user accounts with username, password, email, and optional display name. New users can log in and access the application.

**Independent Test**: Admin logs in, navigates to user management, creates a new user with valid information. New user can log in with created credentials. Non-admin users cannot access user management.

### Implementation for User Story 1

- [x] T034 [US1] Create backend/src/services/user_service.go with CreateUser method (username, password, email, display_name validation)
- [x] T035 [US1] Implement password hashing in user_service.go using password service
- [x] T036 [US1] Implement user_type validation (admin/regular) in user_service.go
- [x] T037 [US1] Create backend/src/api/handlers/user_handler.go with POST /api/users endpoint handler
- [x] T038 [US1] Add admin authorization check to POST /api/users handler
- [x] T039 [US1] Register POST /api/users route in backend/src/api/router.go
- [x] T040 [US1] Create frontend/src/pages/user-management.html with user creation form
- [x] T041 [US1] Create frontend/src/components/user-form.js with form validation and API integration
- [x] T042 [US1] Add user management page routing in frontend navigation

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently. Admin can create users, new users can log in.

---

## Phase 4: User Story 2 - User Views Their Tasks (Priority: P1)

**Goal**: Authenticated user can view a list of all their tasks with description, due date, status, and progress. Users cannot see tasks belonging to other users.

**Independent Test**: User logs in, accesses main task view, sees all their tasks displayed with required fields. User with no tasks sees empty state message. User cannot see other users' tasks.

### Implementation for User Story 2

- [x] T043 [US2] Create backend/src/services/task_service.go with GetTasksByUserID method
- [x] T044 [US2] Implement user_id filtering in GetTasksByUserID to ensure data isolation (FR-008)
- [x] T045 [US2] Create backend/src/api/handlers/task_handler.go with GET /api/tasks endpoint handler
- [x] T046 [US2] Add authentication middleware to GET /api/tasks route
- [x] T047 [US2] Register GET /api/tasks route in backend/src/api/router.go
- [x] T048 [US2] Create frontend/src/pages/task-list.html with task list display
- [x] T049 [US2] Create frontend/src/components/task-item.js for individual task display
- [x] T050 [US2] Create frontend/src/components/empty-state.js for no tasks message
- [x] T051 [US2] Implement task list API integration in frontend/src/services/task-service.js
- [x] T052 [US2] Add task list page routing in frontend navigation

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently. Users can view their tasks.

---

## Phase 5: User Story 3 - User Creates New Tasks (Priority: P1)

**Goal**: Authenticated user can create new tasks with description, due date, status, and progress. Tasks persist between server restarts and appear immediately in task list.

**Independent Test**: User creates a new task with all required fields, task appears in task list. User attempts to create task without required fields, sees validation errors. User creates task, logs out and back in, task persists.

### Implementation for User Story 3

- [x] T053 [US3] Implement CreateTask method in backend/src/services/task_service.go with validation
- [x] T054 [US3] Add progress validation (0-100%) in CreateTask method per FR-019
- [x] T055 [US3] Add past due date warning logic in CreateTask method per clarification (allow but warn)
- [x] T056 [US3] Implement POST /api/tasks endpoint handler in backend/src/api/handlers/task_handler.go
- [x] T057 [US3] Add authentication middleware to POST /api/tasks route
- [x] T058 [US3] Register POST /api/tasks route in backend/src/api/router.go
- [x] T059 [US3] Create frontend/src/pages/task-create.html with task creation form
- [x] T060 [US3] Create frontend/src/components/task-form.js with form validation (description, status, progress 0-100)
- [x] T061 [US3] Implement past due date warning display in task-form.js
- [x] T062 [US3] Add task creation API integration in frontend/src/services/task-service.js
- [x] T063 [US3] Add task creation page routing and navigation

**Checkpoint**: At this point, User Stories 1, 2, AND 3 should all work independently. Users can create and view tasks.

---

## Phase 6: User Story 4 - User Edits Existing Tasks (Priority: P2)

**Goal**: Authenticated user can update task fields (description, due date, status, progress). Changes are saved and reflected immediately. Users can modify task dependencies.

**Independent Test**: User edits an existing task's fields, changes are saved and displayed. User modifies task dependencies, relationships update correctly. User attempts circular dependency, system prevents and shows error.

### Implementation for User Story 4

- [x] T064 [US4] Implement UpdateTask method in backend/src/services/task_service.go
- [x] T065 [US4] Add task ownership validation in UpdateTask (user can only edit own tasks)
- [x] T066 [US4] Add progress validation (0-100%) in UpdateTask method
- [x] T067 [US4] Implement PUT /api/tasks/{taskId} endpoint handler in backend/src/api/handlers/task_handler.go
- [x] T068 [US4] Add task ID parameter parsing and validation in PUT handler
- [x] T069 [US4] Register PUT /api/tasks/{taskId} route in backend/src/api/router.go
- [x] T070 [US4] Create frontend/src/pages/task-edit.html with task edit form
- [x] T071 [US4] Update frontend/src/components/task-form.js to support edit mode
- [x] T072 [US4] Add task edit API integration in frontend/src/services/task-service.js
- [x] T073 [US4] Add task edit page routing (navigate from task list)

**Checkpoint**: At this point, User Stories 1-4 should all work independently. Users can create, view, and edit tasks.

---

## Phase 7: User Story 5 - User Deletes Tasks (Priority: P2)

**Goal**: Authenticated user can delete tasks they own. Deletion requires confirmation dialog with cancel option. When a task with dependencies is deleted, all dependency references are automatically removed.

**Independent Test**: User deletes a task with confirmation, task is removed from list. User cancels deletion, task remains. User deletes task with dependencies, task deleted and dependencies removed. Deleted task persists as deleted after logout/login.

### Implementation for User Story 5

- [x] T074 [US5] Implement DeleteTask method in backend/src/services/task_service.go
- [x] T075 [US5] Add task ownership validation in DeleteTask (user can only delete own tasks)
- [x] T076 [US5] Implement automatic dependency cleanup in DeleteTask (CASCADE DELETE handled by database, verify in service)
- [x] T077 [US5] Implement DELETE /api/tasks/{taskId} endpoint handler in backend/src/api/handlers/task_handler.go
- [x] T078 [US5] Register DELETE /api/tasks/{taskId} route in backend/src/api/router.go
- [x] T079 [US5] Create frontend/src/components/delete-confirmation.js with confirmation dialog component
- [x] T080 [US5] Add delete button to frontend/src/components/task-item.js
- [x] T081 [US5] Implement delete confirmation flow in task-item.js (show dialog, handle cancel/confirm)
- [x] T082 [US5] Add task deletion API integration in frontend/src/services/task-service.js
- [x] T083 [US5] Update task list to refresh after deletion

**Checkpoint**: At this point, User Stories 1-5 should all work independently. Users can create, view, edit, and delete tasks.

---

## Phase 8: User Story 6 - User Filters Tasks (Priority: P2)

**Goal**: Authenticated user can filter tasks by status ('all', 'completed', 'in_progress'). Filter results update immediately. Empty filter results show appropriate message.

**Independent Test**: User with multiple tasks selects 'all' filter, sees all tasks. User selects 'completed' filter, sees only completed tasks. User selects 'in_progress' filter, sees only non-completed tasks (in_progress, not_started, blocked, deferred).

### Implementation for User Story 6

- [ ] T084 [US6] Update GetTasksByUserID method in backend/src/services/task_service.go to support status filtering
- [ ] T085 [US6] Implement status filter logic ('all', 'completed', 'in_progress' meaning non-completed) per spec
- [ ] T086 [US6] Update GET /api/tasks handler to accept status query parameter
- [ ] T087 [US6] Add status query parameter parsing and validation in GET /api/tasks handler
- [ ] T088 [US6] Create frontend/src/components/task-filter.js with filter UI (all/completed/in_progress buttons)
- [ ] T089 [US6] Implement filter state management in frontend/src/pages/task-list.html
- [ ] T090 [US6] Update task list API call to include status filter parameter
- [ ] T091 [US6] Update empty state component to handle filtered empty results
- [ ] T092 [US6] Add filter persistence (maintain filter selection during session)

**Checkpoint**: At this point, User Stories 1-6 should all work independently. Users can create, view, edit, delete, and filter tasks.

---

## Phase 9: User Story 7 - User Manages Task Dependencies (Priority: P3)

**Goal**: Authenticated user can view task dependencies, add dependency relationships, and edit dependencies. System prevents circular dependencies. Dependencies are clearly displayed when viewing a task.

**Independent Test**: User sets task A to depend on task B, relationship is saved and displayed. User views task with dependencies, dependent tasks are shown. User attempts circular dependency, system prevents and shows error.

### Implementation for User Story 7

- [x] T093 [US7] Create backend/src/services/dependency_service.go with dependency management methods
- [x] T094 [US7] Implement AddDependency method with circular dependency detection (graph cycle detection)
- [x] T095 [US7] Implement GetDependencies method to retrieve tasks a task depends on
- [x] T096 [US7] Implement RemoveDependency method in dependency_service.go
- [x] T097 [US7] Add task ownership validation (both tasks must belong to same user) in dependency methods
- [x] T098 [US7] Implement GET /api/tasks/{taskId}/dependencies endpoint handler
- [x] T099 [US7] Implement POST /api/tasks/{taskId}/dependencies endpoint handler
- [x] T100 [US7] Implement DELETE /api/tasks/{taskId}/dependencies/{dependsOnTaskId} endpoint handler
- [x] T101 [US7] Register dependency routes in backend/src/api/router.go
- [x] T102 [US7] Create frontend/src/components/dependency-manager.js for viewing and managing dependencies
- [x] T103 [US7] Add dependency display to frontend/src/pages/task-edit.html
- [x] T104 [US7] Create frontend/src/components/dependency-selector.js for adding dependencies
- [x] T105 [US7] Implement circular dependency error display in frontend
- [x] T106 [US7] Add dependency API integration in frontend/src/services/task-service.js

**Checkpoint**: At this point, all user stories should be independently functional. Users can manage tasks with full dependency support.

---

## Phase 10: Authentication & Session Management

**Purpose**: Complete authentication flow (login/logout) required for all user stories

**Note**: This is foundational but separated for clarity. Login/logout endpoints are needed before users can access any features.

- [x] T107 Create backend/src/services/auth_service.go with Login and Logout methods
- [x] T108 Implement password verification in Login method using bcrypt
- [x] T109 Implement session creation in Login method using gorilla/sessions
- [x] T110 Create backend/src/api/handlers/auth_handler.go with POST /api/login endpoint handler
- [x] T111 Create backend/src/api/handlers/auth_handler.go with POST /api/logout endpoint handler
- [x] T112 Register POST /api/login and POST /api/logout routes in backend/src/api/router.go
- [x] T113 Create frontend/src/pages/login.html with login form
- [x] T114 Create frontend/src/components/login-form.js with form validation and API integration
- [x] T115 Implement session cookie handling in frontend/src/services/api.js
- [x] T116 Add login page routing and redirect logic (redirect authenticated users)
- [x] T117 Add logout functionality to frontend navigation

---

## Phase 11: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T118 [P] Add input validation error messages for all forms (user creation, task creation/edit)
- [ ] T119 [P] Implement consistent error handling and user-friendly error messages across all endpoints
- [ ] T120 [P] Add loading states and spinners for async operations in frontend
- [ ] T121 [P] Implement responsive design for mobile/tablet viewports
- [ ] T122 [P] Add date picker component for due date selection in task forms
- [ ] T123 [P] Implement progress bar visualization for task progress (0-100%)
- [ ] T124 [P] Add task status badge styling (completed, in_progress, etc.)
- [ ] T125 [P] Implement keyboard shortcuts for common actions (create task, filter, etc.)
- [ ] T126 [P] Add confirmation dialogs for destructive actions beyond task deletion
- [ ] T127 [P] Implement proper logging for all operations (user creation, task CRUD, dependencies)
- [ ] T128 [P] Add database backup/restore functionality for data persistence
- [ ] T129 [P] Create license report generation script using go-licenses tool
- [ ] T130 [P] Update documentation (README.md, quickstart.md validation)
- [ ] T131 [P] Add environment variable documentation and example config file
- [x] T132 [P] Implement graceful shutdown handling in backend/src/main.go
- [x] T133 [P] Add health check endpoint GET /api/health for monitoring

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **Authentication (Phase 10)**: Depends on Foundational completion - Required before user stories
- **User Stories (Phases 3-9)**: All depend on Foundational + Authentication completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Phase 11)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational + Authentication - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational + Authentication - No dependencies on other stories
- **User Story 3 (P1)**: Can start after Foundational + Authentication - Depends on User Story 2 (needs task viewing)
- **User Story 4 (P2)**: Can start after User Story 3 - Depends on task creation
- **User Story 5 (P2)**: Can start after User Story 3 - Depends on task creation
- **User Story 6 (P2)**: Can start after User Story 2 - Depends on task viewing
- **User Story 7 (P3)**: Can start after User Story 3 - Depends on task creation and editing

### Within Each User Story

- Models before services
- Services before endpoints/handlers
- Backend endpoints before frontend integration
- Core implementation before UI polish
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational + Authentication complete, User Stories 1 and 2 can start in parallel
- Models within a story marked [P] can run in parallel
- Frontend components marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members (after dependencies met)

---

## Parallel Example: User Story 1

```bash
# Launch all parallel tasks for User Story 1 together:
Task: T034 [US1] Create backend/src/services/user_service.go
Task: T040 [US1] Create frontend/src/pages/user-management.html
Task: T041 [US1] Create frontend/src/components/user-form.js
```

---

## Parallel Example: User Story 2

```bash
# Launch all parallel tasks for User Story 2 together:
Task: T043 [US2] Create backend/src/services/task_service.go
Task: T048 [US2] Create frontend/src/pages/task-list.html
Task: T049 [US2] Create frontend/src/components/task-item.js
Task: T050 [US2] Create frontend/src/components/empty-state.js
```

---

## Implementation Strategy

### MVP First (User Stories 1-3 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 10: Authentication (required for all stories)
4. Complete Phase 3: User Story 1 (Admin creates users)
5. Complete Phase 4: User Story 2 (User views tasks)
6. Complete Phase 5: User Story 3 (User creates tasks)
7. **STOP and VALIDATE**: Test User Stories 1-3 independently
8. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational + Authentication → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (Admin can create users)
3. Add User Story 2 → Test independently → Deploy/Demo (Users can view tasks)
4. Add User Story 3 → Test independently → Deploy/Demo (MVP: Users can create tasks)
5. Add User Story 4 → Test independently → Deploy/Demo (Users can edit tasks)
6. Add User Story 5 → Test independently → Deploy/Demo (Users can delete tasks)
7. Add User Story 6 → Test independently → Deploy/Demo (Users can filter tasks)
8. Add User Story 7 → Test independently → Deploy/Demo (Full feature set)
9. Add Polish phase → Final release

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational + Authentication together
2. Once foundation is done:
   - Developer A: User Story 1 (Admin creates users)
   - Developer B: User Story 2 (User views tasks) - can start in parallel with US1
   - Developer C: User Story 3 (User creates tasks) - after US2 complete
3. Stories complete and integrate independently
4. Polish phase can be parallelized across team

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- Authentication phase (Phase 10) is foundational but separated for clarity - complete before user stories
- All endpoints require authentication middleware except /api/login
- Admin endpoints require admin authorization middleware
- All task operations must enforce user_id filtering for data isolation (FR-008)
