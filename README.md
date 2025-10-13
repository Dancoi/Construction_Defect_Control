## Construction Defect Control — Monorepo

This repository contains a small full-stack prototype for a Construction Defect Control system.

- backend: Go (Gin, GORM) — API server, authentication (JWT), projects, defects, attachments, comments
- frontend: React (Vite, Tailwind) — SPA client
- db: PostgreSQL (via Docker Compose)

This README explains how to run the system locally, how Docker Compose is configured, environment variables, common troubleshooting, and development tips.

---

## Repo layout

```
/backend        # Go server (cmd/, internal/ ...)
/frontend       # React app (Vite + Tailwind)
/docker-compose.yml
/README.md
```

## Quick start (Docker Compose)

Recommended for most local testing. From repository root:

```powershell
# build images and start services (api, db, frontend)
docker compose up --build -d

# list running containers
docker compose ps

# follow logs
docker compose logs -f

# stop and remove containers
docker compose down
```

Services started by compose
- db: PostgreSQL
- app: backend (Go) on port 8080
- frontend: nginx serving built SPA on port 5173 (proxies `/api/` to backend)

Notes:
- The compose file declares named volumes: `db-data` and `uploads-data`. Compose will create them automatically.
- API is available at: `http://localhost:8080/api/v1`
- Frontend is available at: `http://localhost:5173`

### Bootstrap the first admin (developer flow)

To create the first admin automatically on an empty database, set the environment variable `AUTH_BOOTSTRAP_FIRST_ADMIN=true` for the `app` service before the first registration. This is implemented in the register flow: if the flag is true and the `users` table is empty, the first registered user will receive the `admin` role.

You can either add this variable to the `app` environment in `docker-compose.yml` or set it via an `.env` file used by compose.

After starting the stack, register a new user via the UI or via API POST `/api/v1/auth/register` and that user will become admin.

**Important:** Disable `AUTH_BOOTSTRAP_FIRST_ADMIN` after creating the first admin.

## Environment variables

Important variables used by services (configure in compose or `.env`):

- `DATABASE_URL` — Postgres connection string (e.g. `postgres://postgres:postgres@db:5432/defectdb?sslmode=disable`)
- `UPLOADS_PATH` — path where attachments are stored in the backend container (default `/app/uploads`)
- `JWT_SECRET` — secret used to sign JWTs (set to the same value across deployments)
- `AUTH_BOOTSTRAP_FIRST_ADMIN` — if `true`, first registered user on empty DB becomes admin

For local dev you may set them in `docker-compose.yml` or in a `.env` file.

## Running backend locally (without Docker)

Requirements: Go 1.20+, Postgres running

```powershell
cd backend
# configure configs/config.yml or env vars
go run ./cmd
```

The backend uses Viper; configuration file is `backend/configs/config.yml` and environment variables override values.

## Running frontend locally (dev)

Requirements: Node 18+, npm

```powershell
cd frontend
npm install
npm run dev
# open http://localhost:5173
```

When running frontend in dev mode, backend CORS is configured to allow `http://localhost:5173`. If you open the app at a different host (VM IP), add that origin to backend CORS in `backend/cmd/main.go`.

## Common troubleshooting

1. 401 / invalid token on `/api/v1/auth/me`

- Make sure the `Authorization: Bearer <token>` header is present (check devtools Network). If using the nginx frontend container, ensure nginx proxies the header (`proxy_set_header Authorization $http_authorization;`) — the image in this repo is configured.
- Check `JWT_SECRET` in the backend container matches the secret used to sign tokens. If you change the secret, re-login to obtain a new token.

2. nginx error `host not found in upstream "backend"`

- Frontend nginx must proxy to the service name declared in `docker-compose.yml` (this repo uses `app`). We already use `proxy_pass http://app:8080/api/` in `frontend/nginx.conf`.

3. Volumes (db-data / uploads-data)

- Named volumes declared in `docker-compose.yml` are created by Compose. To inspect them:
	- `docker volume ls`
	- `docker volume inspect <volume>`
- To bind to a local folder instead of a named volume, change the `volumes` entry to `./data/postgres:/var/lib/postgresql/data` or similar.

4. File upload/download permissions

- Backend saves uploads to the `uploads` volume. Ensure `UPLOADS_PATH` in compose matches the path the backend expects (default `/app/uploads`).

## Development tips

- Auth: `frontend/src/api/axios.js` attaches `Authorization` header from `localStorage.token`.
- Role-based UI: Header and pages use `AuthContext` and `user.role` to show/hide features.
- To debug token parsing on the server, see `backend/internal/middleware/jwt.go` (it logs parse errors to stdout).

## Useful commands

```powershell
# build backend image only
docker compose build app

# build frontend image only
docker compose build frontend

# rebuild everything
docker compose up --build -d

# follow logs
docker compose logs -f
```

## Next steps / TODOs

- Add tests and CI
- Harden production configuration (TLS, secrets, healthchecks)
- Implement project-scoped roles

---

If you want, I can:

- Add a small entrypoint script to create an admin from env vars (ADMIN_EMAIL / ADMIN_PASSWORD) on first run
- Add healthchecks and restart policies to compose
- Switch named volumes to bind mounts for easier local inspection

Tell me which of the above you'd like me to implement next.
