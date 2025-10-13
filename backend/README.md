# Defect Control System (prototype)

Run with Docker Compose:

```powershell
docker-compose up --build
```

Local run: set env vars per `.env.example` and `go run ./cmd`

Swagger/OpenAPI
----------------

This project uses `swag` (swaggo) to generate OpenAPI docs from handler annotations.

Install `swag` (recommended):

```powershell
go install github.com/swaggo/swag/cmd/swag@latest
```

Generate docs locally:

```powershell
swag init -g cmd/main.go -o docs
```

This will produce `docs/swagger.json` and `docs/docs.go` (the JSON is served by the app).

CI tip: to ensure the committed `docs/swagger.json` doesn't drift from annotations, add a job that runs `swag init -g cmd/main.go -o /tmp/generated` and fails if `diff -u docs/swagger.json /tmp/generated/swagger.json` is non-empty.

Example CI snippet (bash GNU coreutils required):

```bash
swag init -g cmd/main.go -o /tmp/generated
if ! diff -u docs/swagger.json /tmp/generated/swagger.json >/dev/null; then
	echo "Swagger docs out of date. Run: swag init -g cmd/main.go -o docs" >&2
	exit 1
fi
```

Security
--------

The API uses JWT Bearer tokens. Each token contains a `role` claim.
Protect endpoints with the `@Security BearerAuth` annotation so they appear in the generated OpenAPI spec.

Notes:
- This project's `Dockerfile` uses Go 1.25; the builder image is `golang:1.25-alpine`.
- The repository root contains `docker-compose.yml` which builds the backend from `./backend`.

