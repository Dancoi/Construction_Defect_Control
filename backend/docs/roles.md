# Roles and Permissions

This document describes the simple role model used by the Defect Control System.

Roles:

- engineer (default):
  - Create defects in projects
  - Upload attachments they own
  - View defects and attachments for projects they are part of

- manager:
  - All engineer permissions
  - Create projects
  - View and download attachments across projects
  - Manage defects (change status)

- stakeholder:
  - View projects and defects
  - Download attachments

- admin:
  - Full access to all resources
  - Can manage users and roles

Bootstrap and default role behavior:

- By default, new users are created with role `engineer`.
- To create the first admin on an empty database, set the environment variable `AUTH_BOOTSTRAP_FIRST_ADMIN=true` before starting the service. The first registered user will receive the `admin` role.

Examples of permission checks used in handlers:

- RequireRole middleware: protects endpoints that only specific roles may call. Example: creating a project requires `manager` or `admin`.
- Ownership check: for downloads we allow the attachment uploader, or users with role `manager`, `stakeholder`, or `admin`.

Notes and future improvements:

- Current model is global (one role per user). For project-scoped roles (e.g., project manager), extend the database with a `user_project_roles` table and check project-scoped roles in handlers.
- For a richer permission model, consider a permission matrix (capabilities -> roles) or use an existing RBAC library.
