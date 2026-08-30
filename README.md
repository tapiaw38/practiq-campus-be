# practiq-campus-be

Practiq Campus backend — virtual campus API in Go with clean architecture.

## Requirements

- Go 1.24+
- PostgreSQL 14+
- [golang-migrate CLI](https://github.com/golang-migrate/migrate)
- `auth-api-be` running on port 8082

## Configuration

Copy `.env.example` to `.env` and adjust the values:

```bash
cp .env.example .env
```

| Variable | Description | Default |
| --- | --- | --- |
| `PORT` | HTTP port | `8084` |
| `DATABASE_URL` | PostgreSQL DSN | see `.env.example` |
| `JWT_SECRET` | Must match `auth-api-be` and `practiq-be` | — |
| `AUTH_API_URL` | Auth API | `http://localhost:8082` |
| `PRACTIQ_API_URL` | Practiq API | `http://localhost:8083` |
| `FRONTEND_URL` | Campus FE CORS origin | `http://localhost:5175` |

## Database

From the monorepo root:

```bash
docker compose up -d campus-postgres-db
```

Or use an existing PostgreSQL instance and set `DATABASE_URL`.

## Migrations

SQL files live in `migrations/`. Docker applies them in `entrypoint.sh`.

## Run

```bash
go run ./cmd/api
```

The server listens on `http://localhost:8084`.

## Docker Compose

From the root:

```bash
docker compose up -d --build practiq-campus-be
```

Exposed by Traefik at `campus-api.practiq.localhost`.

## Architecture

```
practiq-campus-be/
├── cmd/api/
├── internal/
│   ├── adapters/
│   │   ├── datasources/repositories/
│   │   └── web/
│   ├── domain/
│   ├── platform/
│   └── usecases/
└── migrations/
```

## Modules

- Profile and preferences
- Courses, sections, enrollments, and materials
- Assignments and submissions
- Forums
- Calendar
- Messaging
