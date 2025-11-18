# Feature Specification: Task Management Web Application

**Feature Branch**: `001-task-management-app`  
**Created**: 2025-11-09  
**Status**: Draft  
**Input**: User description: "I want to build web app that handles tasks. It shows the tasks (each one has fields such as description, due-date, status, progress etc.). Task have dependencies so the app should be able to edis and show dependencies. circular dependencies are not allowed. This is a multi user application, each user sees only his own tasks. the app starts with a single admin user which is the only user which can craete additional users to the app - this is the only thing admin can do. non admin users can create, edit and delete tasks (delete requires a dialog with option to cancel). Tasks can be filtered in multiple ways, but the most important fileters are 'all', 'completed' and 'in progress' (meaning all that are not completed and not deferred). there is no way for a user to see the tasks of another user. The UI should be simple and intuitive. Tasks should persist between invications of the server On the next phase I would like to add Rest API to the server, with 2 APIs, login and get_tasks. No need to add APIs for editing tasks."

## Clarifications

### Session 2025-11-09

- Q: When a user deletes a task that other tasks depend on, what should happen? → A: Allow deletion and automatically remove dependencies - Delete the task and remove all dependency references from dependent tasks
- Q: What information is required when an admin creates a new user account? → A: Username, password, email, and optional display name - Email required, name optional
- Q: What are the valid task status values that users can assign to tasks? → A: completed, in progress, not started, blocked, and deferred - Five status values total
- Q: Can users create or edit tasks with due dates in the past? → A: Allow but warn - Permit past due dates but display a warning message to the user
- Q: How should the system handle invalid progress values (e.g., negative numbers or values over 100%)? → A: Enforce 0-100% range, block invalid values - Show validation error and prevent saving if progress is outside 0-100%

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Admin Creates New Users (Priority: P1)

An admin user needs to create additional users in the system so that multiple people can use the application. The admin logs into the application and navigates to a user management section where they can add new users by providing required information (username, password, email) and optionally a display name. Once created, the new users can log in and start managing their own tasks.

**Why this priority**: This is foundational for the multi-user system. Without user creation, the application cannot support multiple users, making it a single-user system. This must be implemented first to enable all other user-facing features.

**Independent Test**: Can be fully tested by having an admin user create a new user account, then verifying that the new user can log in and access the application. This delivers the core multi-user capability.

**Acceptance Scenarios**:

1. **Given** an admin user is logged in, **When** they navigate to user management and provide valid username, password, and email (and optionally display name), **Then** a new user account is created and the new user can log in
2. **Given** an admin user is logged in, **When** they attempt to create a user without required fields (username, password, or email), **Then** appropriate validation errors are displayed
3. **Given** a non-admin user is logged in, **When** they attempt to access user management, **Then** they are denied access or the feature is not visible

---

### User Story 2 - User Views Their Tasks (Priority: P1)

A user logs into the application and sees a list of all their tasks. Each task displays key information including description, due date, status, and progress. The user can quickly understand what tasks they have and their current state.

**Why this priority**: Viewing tasks is the core functionality of the application. Users must be able to see their tasks before they can manage them. This is the primary value proposition of the application.

**Independent Test**: Can be fully tested by logging in as a user and verifying that their tasks are displayed correctly with all required fields. This delivers immediate value by showing users their task list.

**Acceptance Scenarios**:

1. **Given** a user is logged in, **When** they access the main task view, **Then** all their tasks are displayed with description, due date, status, and progress
2. **Given** a user has no tasks, **When** they access the main task view, **Then** an appropriate empty state message is shown
3. **Given** a user is logged in, **When** they view their tasks, **Then** they cannot see tasks belonging to other users

---

### User Story 3 - User Creates New Tasks (Priority: P1)

A user needs to add new tasks to track their work. The user navigates to create a new task, fills in the required information (description, due date, status, progress), and saves it. The task immediately appears in their task list.

**Why this priority**: Creating tasks is essential functionality. Without the ability to create tasks, users cannot use the application to track their work. This is a core user action that enables the application's primary purpose.

**Independent Test**: Can be fully tested by having a user create a new task with all required fields, then verifying it appears in their task list. This delivers value by allowing users to start tracking their work.

**Acceptance Scenarios**:

1. **Given** a user is logged in, **When** they create a new task with description, due date, status, and progress, **Then** the task is saved and appears in their task list
2. **Given** a user is creating a task, **When** they attempt to save without required fields, **Then** validation errors are displayed
3. **Given** a user is creating or editing a task, **When** they set a due date in the past, **Then** the system allows the date but displays a warning message
4. **Given** a user creates a task, **When** they log out and log back in, **Then** the task persists and is still visible

---

### User Story 4 - User Edits Existing Tasks (Priority: P2)

A user needs to update task information as work progresses or circumstances change. The user selects a task, modifies any of its fields (description, due date, status, progress), and saves the changes. The updated information is reflected immediately.

**Why this priority**: While viewing and creating tasks are essential, editing is crucial for maintaining accurate task information over time. Users need to update tasks as work progresses, making this a high-priority feature.

**Independent Test**: Can be fully tested by having a user edit an existing task's fields and verifying the changes are saved and displayed correctly. This delivers value by allowing users to keep their task information current.

**Acceptance Scenarios**:

1. **Given** a user has an existing task, **When** they edit its description, due date, status, or progress, **Then** the changes are saved and reflected in the task list
2. **Given** a user is editing a task, **When** they modify task dependencies, **Then** the dependency relationships are updated correctly
3. **Given** a user edits a task, **When** they attempt to create a circular dependency, **Then** the system prevents the change and displays an error
4. **Given** a user is creating or editing a task, **When** they attempt to set progress to a value outside 0-100%, **Then** the system displays a validation error and prevents saving

---

### User Story 5 - User Deletes Tasks (Priority: P2)

A user needs to remove tasks that are no longer relevant or were created in error. The user selects a task to delete, confirms the deletion in a dialog that provides a cancel option, and the task is removed from their list.

**Why this priority**: Deleting tasks is important for maintaining a clean and relevant task list. While not as critical as creating tasks, users need this capability to manage their workflow effectively.

**Independent Test**: Can be fully tested by having a user delete a task with confirmation, then verifying it no longer appears in their task list. This delivers value by allowing users to remove unwanted tasks.

**Acceptance Scenarios**:

1. **Given** a user has an existing task, **When** they select delete and confirm in the dialog, **Then** the task is removed from their task list
2. **Given** a user initiates task deletion, **When** they click cancel in the confirmation dialog, **Then** the task remains unchanged and the dialog closes
3. **Given** a user deletes a task, **When** they log out and log back in, **Then** the task remains deleted and does not reappear
4. **Given** a user has a task that other tasks depend on, **When** they delete that task, **Then** the task is deleted and all dependency references are automatically removed from dependent tasks

---

### User Story 6 - User Filters Tasks (Priority: P2)

A user needs to view subsets of their tasks based on status to focus on specific work. The user selects a filter option ('all', 'completed', 'in progress') and the task list updates to show only tasks matching that filter.

**Why this priority**: Filtering is important for usability, especially as users accumulate many tasks. It helps users focus on relevant work and improves the application's utility. While not essential for basic functionality, it significantly enhances user experience.

**Independent Test**: Can be fully tested by having a user with multiple tasks apply different filters and verifying that only matching tasks are displayed. This delivers value by helping users organize and focus their work.

**Acceptance Scenarios**:

1. **Given** a user has tasks with different statuses, **When** they select the 'all' filter, **Then** all their tasks are displayed
2. **Given** a user has completed and incomplete tasks, **When** they select the 'completed' filter, **Then** only completed tasks are displayed
3. **Given** a user has completed and incomplete tasks, **When** they select the 'in progress' filter, **Then** only non-completed and non-deferred tasks are displayed

---

### User Story 7 - User Manages Task Dependencies (Priority: P3)

A user needs to establish relationships between tasks to reflect work dependencies. The user can view which tasks depend on others and edit these relationships. When viewing a task, its dependencies are clearly displayed. A user can't mark a task as completed if any of the tasks it depends on is not completed.

**Why this priority**: Task dependencies are an advanced feature that adds significant value for complex workflows, but they are not essential for basic task management. Users can still effectively use the application without dependencies, making this a lower priority enhancement.

**Independent Test**: Can be fully tested by having a user create tasks with dependencies, view the dependency relationships, and edit them. This delivers value by enabling users to model complex work relationships.

**Acceptance Scenarios**:

1. **Given** a user has multiple tasks, **When** they set task A to depend on task B, **Then** the dependency relationship is saved and displayed
2. **Given** a user views a task with dependencies, **When** they examine the task details, **Then** all dependent tasks are clearly shown
3. **Given** a user edits task dependencies, **When** they attempt to create a circular dependency (A depends on B, B depends on A), **Then** the system prevents this and displays an error message
4. **Given** a user marks a task as completed, **When** any of the tasks this tasks depends on is not completed, **Then** the UI shows an error messags and the system prevents marking this task as completed.

---

### Edge Cases

- When a user deletes a task that other tasks depend on, the system automatically removes all dependency references from the dependent tasks (resolved: deletion allowed with automatic dependency cleanup)
- When a user creates or edits a task with a due date in the past, the system allows it but displays a warning message (resolved: past dates allowed with warning)
- When a user attempts to set a task's progress to an invalid value (negative or over 100%), the system displays a validation error and prevents saving (resolved: 0-100% range enforced)
- What happens when a user filters tasks but has no tasks matching the filter criteria?
- How does the system handle concurrent edits when multiple users are logged in (each editing their own tasks)?
- What happens when task data becomes corrupted or invalid during persistence?
- How does the system handle very long task descriptions or invalid date formats?
- How does the system handle dependency validation when tasks are deleted or modified?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST authenticate users and distinguish between admin and non-admin user types
- **FR-002**: System MUST allow admin users to create new user accounts
- **FR-003**: System MUST restrict user creation functionality to admin users only
- **FR-004**: System MUST allow non-admin users to create tasks with description, due date, status, and progress fields
- **FR-018**: System MUST allow users to set task due dates in the past, but MUST display a warning message when a past due date is entered
- **FR-019**: System MUST enforce that task progress values are within the range of 0-100% and MUST block saving with validation error if progress is outside this range
- **FR-005**: System MUST allow non-admin users to edit existing tasks they own
- **FR-006**: System MUST allow non-admin users to delete tasks they own, with a confirmation dialog that includes a cancel option
- **FR-007**: System MUST display all tasks belonging to a user when they log in
- **FR-008**: System MUST ensure users can only view and manage their own tasks, never tasks belonging to other users
- **FR-009**: System MUST support task filtering by 'all', 'completed', and 'in progress' statuses
- **FR-010**: System MUST allow users to view task dependencies
- **FR-011**: System MUST allow users to edit task dependencies
- **FR-012**: System MUST prevent circular dependencies between tasks
- **FR-013**: System MUST persist all task data between server invocations
- **FR-014**: System MUST persist all user data between server invocations
- **FR-015**: System MUST provide a simple and intuitive user interface
- **FR-016**: System MUST start with a single admin user account that can be used to create additional users

### Key Entities *(include if feature involves data)*

- **User**: Represents an application user. Has a user type (admin or regular), authentication credentials (username and password), email address (required), optional display name, and is associated with their own tasks. Admin users have the exclusive ability to create other users. Regular users can manage their own tasks but cannot create other users.

- **Task**: Represents a work item that a user needs to complete. Has attributes including description (text), due date (date), status (one of: completed, in progress, not started, blocked, deferred), and progress (percentage or similar measure). Tasks belong to exactly one user. Tasks can have dependencies on other tasks (one task can depend on multiple other tasks). Tasks cannot have circular dependencies.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can view their complete task list within 2 seconds of logging in
- **SC-002**: Users can create a new task with all required fields in under 30 seconds
- **SC-003**: Users can successfully filter their tasks by status and see accurate results immediately
- **SC-004**: 95% of users can complete their primary task (viewing, creating, or editing tasks) on their first attempt without assistance
- **SC-005**: Task data persists correctly across server restarts with zero data loss for all users
- **SC-006**: Users cannot access tasks belonging to other users (100% isolation between user data)
- **SC-007**: System prevents 100% of circular dependency attempts with clear error messaging
- **SC-008**: Admin users can create a new user account in under 1 minute
- **SC-009**: Task deletion confirmation dialog allows users to cancel deletion successfully 100% of the time

## Assumptions

- The application uses standard web-based authentication (session-based or similar) for user login
- Task status values are: completed, in progress, not started, blocked, and deferred
- Progress is represented as a percentage (0-100%) or similar measurable value
- The "in progress" filter shows all tasks that are not completed and not deferred (includes: in_progress, not_started, and blocked statuses)
- Task dependencies are directional (task A depends on task B means A cannot be completed until B is completed)
- The initial admin user account is pre-configured or created through a setup process
- Data persistence uses standard storage mechanisms (database, file system, etc.) appropriate for the deployment environment
- The UI follows common web application patterns and conventions for intuitive user experience
- Server invocations refer to server restarts or application restarts, not individual HTTP requests
- REST API features mentioned (login and get_tasks) are planned for a future phase and are out of scope for this specification
