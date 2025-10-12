Attachments (uploads) — design

Goal

Provide a simple, secure and auditable mechanism to attach files (images, PDFs, etc.) to Defects.
The design below defines the API surface, storage rules, data model and security constraints so implementation can be developed and tested.

Principles

- Files are owned by the Defect and by the authenticated user who uploaded them.
- Files are stored on disk (configurable path) and referenced in DB by a stable path/URL.
- Enforce size and content-type restrictions.
- Sanitize filenames and avoid collisions (use safe generated names + preserve original filename metadata).
- Access to upload and download endpoints is protected by JWT. Ownership and authorization checks applied on sensitive operations.

Config

- uploads.path (string) — directory where files are stored. Default: ./uploads
- uploads.max_size (int) — max bytes per file. Default: 10_000_000 (10 MB)
- uploads.allowed_types (list) — content types allowed, e.g. ["image/jpeg","image/png","application/pdf"]
- uploads.serve_via (string) — "file" or "proxy". If "file", server serves files from path; if "proxy", files are expected to be served by external static server (nginx) and server returns direct URL.

Database model

Attachment
- id (uint, PK)
- defect_id (uint) — FK to defects.id
- uploader_id (uint, nullable) — user who uploaded
- path (string) — internal path on disk (relative to uploads.path)
- filename (string) — original filename
- content_type (string)
- size (int64)
- created_at (timestamp)

OpenAPI / endpoints (sketch)

POST /api/v1/defects/{id}/attachments
- Description: upload one or more files attached to defect {id}
- Security: Bearer JWT
- Path params: id (uint) — defect id
- Body: multipart/form-data: files: file[] (one or multiple)
- Responses:
  - 201 Created: { "status": "ok", "data": [{ "id": <id>, "filename": "orig.jpg", "url": "/uploads/<file>" }]} 
  - 400 Bad Request: validation error / disallowed type / file too large
  - 401 Unauthorized: missing/invalid token
  - 403 Forbidden: user not allowed to attach to this defect (optional ownership policy)
  - 404 Not Found: defect not found

GET /api/v1/attachments/{attachment_id}
- Description: returns file (stream) or redirects to storage URL
- Security: Bearer JWT (or public depending on policy)
- Responses:
  - 200 OK: stream with correct Content-Type
  - 404 Not Found
  - 401 Unauthorized / 403 Forbidden

Response shape (upload)
- 201
{
  "status": "ok",
  "data": [
    {"id": 101, "filename": "photo1.jpg", "content_type": "image/jpeg", "size": 234234, "url": "/uploads/2025/10/11/abcd1234.jpg"}
  ]
}

Storage contract (server-side)
- On upload: validate defect exists (service layer). Validate file size, content-type.
- Compute safe filename: <YYYY>/<MM>/<DD>/<random-hash>.<ext>
- Persist file to disk under uploads.path, ensure directories exist and permissions safe.
- Create Attachment DB record with path relative to uploads.path and metadata (original filename, content-type, size, uploader).
- Return JSON list of created attachments with URLs.

Serving files
- If uploads.serve_via == "file": implement GET /api/v1/attachments/{id} which reads DB, checks auth/ownership if needed, and streams file with correct Content-Type. Use http.ServeFile or Gin's File method.
- If uploads.serve_via == "proxy": return a URL pointing to static server: e.g. https://static.example.com/uploads/2025/10/11/abcd1234.jpg

Security & ownership
- Only authenticated users can upload files.
- Optionally require that the uploader belongs to project or has permission to modify defect.
- When serving files, check that requesting user has permission to view the defect; otherwise 403.
- Avoid exposing internal file system paths in responses; expose public URLs or relative paths.

Validation rules
- Max file size per single file: uploads.max_size (default 10MB).
- Max total files per request: configurable (default 5).
- Allowed content types: configurable; by default allow image/* and pdf.

Filename sanitization and collisions
- Do not use original filename as storage filename.
- Use generated random names (crypto/rand -> hex) with directory sharding by date.
- Keep original filename in DB for display/download name.

Idempotency & duplicate uploads
- If the same file is uploaded twice, it will create two Attachment records (acceptable for MVP).

Error handling
- Return clear JSON error messages with HTTP status codes.
- Log storage errors and return 500 for unexpected failures.

Notes on thumbnails & processing (optional)
- After successful upload, run an asynchronous worker to generate thumbnails (e.g., 200x200) for image types.
- Store thumbnails alongside original with predictable naming.
- Consider rate-limiting or background queue for heavy image processing.

Tests
- Unit tests for storage service: save file to temp dir, verify file exists, DB record created.
- Handler integration tests: use httptest server with temp uploads path and in-memory sqlite (or temp Postgres) to verify full flow.

Docker
- In docker-compose, mount host directory to container (e.g., ./uploads -> /app/uploads)
- Ensure container user has write permissions (or use volume with correct config)

Example curl (upload):

curl -X POST http://localhost:8080/api/v1/defects/1/attachments \
  -H "Authorization: Bearer <token>" \
  -F "files=@./photo1.jpg" \
  -F "files=@./photo2.png"

Example curl (download):

curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/attachments/101 --output photo1.jpg

Acceptance criteria for implementation (MVP)
- POST /api/v1/defects/{id}/attachments accepts files, stores them on disk, creates Attachment rows and returns 201 with metadata.
- GET /api/v1/attachments/{id} returns the file if user authorized.
- Tests covering storage service and handler flow.

---

Next step: implement storage service that follows this contract and add Attachment model + handlers.