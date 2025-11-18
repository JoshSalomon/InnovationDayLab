<!--
Sync Impact Report:
Version change: (none) → 1.0.0
Modified principles: N/A (initial creation)
Added sections: Additional Constraints, Development Workflow
Removed sections: N/A
Templates requiring updates:
  ✅ .specify/templates/plan-template.md (Constitution Check section already present)
  ✅ .specify/templates/spec-template.md (no constitution-specific references found)
  ✅ .specify/templates/tasks-template.md (no constitution-specific references found)
  ✅ .cursor/commands/speckit.constitution.md (no updates needed)
Follow-up TODOs: None
-->

# InnovationDayLab Constitution

## Core Principles

### I. Single Executable Architecture
Every component MUST be compiled into a single executable binary. No external executables or runtime dependencies are permitted. The application MUST be self-contained and runnable as a standalone binary. This ensures portability, simplifies deployment, and eliminates external dependency management for end users.

### II. Open Source Compliance
All components and dependencies MUST be open source with compatible licenses. A comprehensive license report MUST be maintained documenting all third-party components with their respective licenses. This ensures transparency, legal compliance, and community trust.

### III. Security & Authentication
The application MUST implement authentication mechanisms to protect user data. Each user MUST have isolated access to their own data set. Authentication MUST be implemented before any personal data access is permitted. This protects user privacy and ensures data isolation.

### IV. Personal Data Handling
User data MUST be presented according to the authenticated user. Each user MUST only access their own data set. Data persistence MUST use local disk storage. For the initial phase, encryption at rest is not required, but this MUST be documented as a known limitation. Future phases MUST address encryption requirements.

### V. Simplicity & Local-First
The application MUST prioritize simplicity for local/personal use. Complexity MUST be justified. YAGNI (You Aren't Gonna Need It) principles apply. The application MUST be designed for single-user or small-scale personal use scenarios. Avoid over-engineering for enterprise-scale requirements.

## Additional Constraints

**Technology Stack**: The application MUST be implemented in Go (Golang). All dependencies MUST be Go packages that compile into the single executable.

**Storage Requirements**: Persistent storage MUST use local disk. If database functionality is required, it MUST be implemented as an in-memory database within the single executable. No external database processes or files are permitted.

**Deployment Model**: The application is designed for local execution. It MUST run as a single process without requiring external services, containers, or orchestration.

**License Reporting**: A mechanism MUST exist to generate and maintain a report of all third-party components and their licenses. This report MUST be accessible and kept up to date with dependency changes.

## Development Workflow

**Constitution Compliance**: All code changes MUST comply with these principles. PRs and reviews MUST verify compliance before merge.

**Complexity Justification**: Any deviation from simplicity principles MUST be documented with rationale in the Complexity Tracking section of implementation plans.

**Testing Requirements**: Core functionality MUST be tested. Authentication and data isolation MUST have explicit test coverage.

**Documentation**: License reports, authentication mechanisms, and data handling policies MUST be documented and kept current.

## Governance

This constitution supersedes all other development practices and guidelines. Amendments require:

1. Documentation of the proposed change and rationale
2. Impact assessment on existing code and practices
3. Version increment according to semantic versioning:
   - **MAJOR**: Backward incompatible principle removals or redefinitions
   - **MINOR**: New principles or materially expanded guidance
   - **PATCH**: Clarifications, wording improvements, typo fixes
4. Update of dependent templates and documentation
5. Sync Impact Report generation

All development work MUST align with these principles. Violations MUST be justified in Complexity Tracking sections or addressed through constitution amendments.

**Version**: 1.0.0 | **Ratified**: 2025-11-09 | **Last Amended**: 2025-11-09
