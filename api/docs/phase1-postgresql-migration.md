# Phase 1 PostgreSQL Foundation

## Current baseline
- The backend can now switch DB dialects through config and defaults to PostgreSQL.
- Startup runs SQL migrations from `./migrations/*.up.sql` before the temporary `AutoMigrate` fallback.
- Config loading now supports `config.yaml` defaults plus environment variable overrides, so local runs and Docker runs use the same model.
- Redis is treated as a degradable dependency during bootstrap. If Redis is unavailable, the app can still start, but cache and task-queue capabilities may be limited.

## Key environment variables
- `DB_DIALECTS`, `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`
- `DB_SSLMODE`, `DB_MIGRATION_PATH`, `DB_AUTO_MIGRATE`
- `REDIS_ADDR`, `REDIS_PASSWORD`
- `SERVER_ADDRESS`, `SERVER_PUBLIC_URL`, `IMAGE_HOST`

## Local startup

```bash
cd api
go run scripts/migrate.go
go run main.go
```

You can override connection settings through environment variables instead of editing `config.yaml`.

## Docker startup

```bash
cd docker
docker compose up -d
```

The Docker stack is aligned to PostgreSQL 16 and injects runtime settings through `.env`.

## Known limitations
- `AutoMigrate` is still enabled as a phase1 safety net and should be removed after the first real schema migrations are complete.
- The repository still has unrelated historical compile/test issues outside the DB foundation work, including disk-space-sensitive test runs and a pre-existing format-string error in `common/agent/agent.go`.

## 2026-04-23 Progress Summary
- Added PostgreSQL dialect support in the DB bootstrap path and switched the phase1 default config from MySQL-oriented assumptions to PostgreSQL-friendly settings.
- Added a lightweight SQL migration runner with `schema_migrations` tracking and committed the first baseline migration file at `api/migrations/0001_stage1_baseline.up.sql`.
- Kept `AutoMigrate` as a temporary compatibility fallback so an empty PostgreSQL instance can still boot during phase1.
- Reworked config loading so runtime environment variables can override `config.yaml`, which unifies local startup and Docker startup behavior.
- Restored config fields that were already present in YAML, including `server.publicUrl` and `monitor.webhook.token`, so they are no longer silently ignored.
- Added focused tests for config environment overrides and invalid env fallback behavior.
- Updated the Docker stack from MySQL to PostgreSQL 16, aligned `.env` variables with backend config keys, and simplified `devops-start.sh` so it updates runtime env instead of patching YAML files in place.
- Tightened Redis bootstrap behavior so failure leaves the shared Redis client in a clean degraded state instead of a half-initialized reference.
- Verified the focused backend packages with `go test ./common/config ./pkg/db ./pkg/redis ./scripts`.
