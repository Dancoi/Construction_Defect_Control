# Docker Compose

This repository can be started locally using Docker Compose.

Services:
- db: PostgreSQL
- app: Go backend (Gin)
- frontend: React app served by nginx

Quick start:

```powershell
# build images and start services
docker-compose up --build

# stop
docker-compose down
```

Environment:
- Backend reads `DATABASE_URL`, `UPLOADS_PATH`, `JWT_SECRET` from environment. You can set these in `backend/.env` or in the compose file.

Notes:
- Frontend nginx proxies `/api/` to the backend service `app`.
- Uploads are stored in a named volume `uploads-data`.
