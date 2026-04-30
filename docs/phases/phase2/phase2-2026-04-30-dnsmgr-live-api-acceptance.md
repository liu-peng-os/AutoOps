# Phase 2 dnsmgr Live API Acceptance Notes

Date: 2026-04-30

## Decision

dnsmgr is the source of truth for domain and DNS record data.

AutoOps must not maintain a second copy of dnsmgr domains or DNS records because:

- dnsmgr may contain a large amount of domain and record data.
- Local synchronization can be slow, stale, or partial.
- Duplicating data creates split ownership between AutoOps and dnsmgr.
- The operations team needs direct, current dnsmgr state when changing DNS records.

Therefore, AutoOps acts as a backend-signed live API operation panel for dnsmgr.

## Implemented Scope

The domain management page under `/integration/domain` now supports these dnsmgr-backed actions:

- Health check
- Live domain list
- Domain detail
- Domain login link
- Live DNS record list by domain
- Add DNS record
- Edit DNS record
- Delete DNS record
- Enable or pause DNS record
- Change DNS record remark
- Batch enable, pause, delete, or change remark

SSL certificate order information is intentionally excluded from this slice.

## Removed Scope

The earlier sync model has been removed.

Removed backend behavior:

- No `POST /api/v1/domain/sync`
- No `GET /api/v1/domain/sync/last`
- No local zone list from `domain_zone`
- No local record list from `domain_record`
- No sync run tracking in `domain_sync_run`

Removed local tables:

- `domain_zone`
- `domain_record`
- `domain_sync_run`

Migration:

- `api/migrations/0009_drop_domain_sync_tables.sql`

## AutoOps API Surface

AutoOps keeps dnsmgr credentials server-side and exposes stable internal APIs:

```text
GET    /api/v1/domain/health
GET    /api/v1/domain/domains
GET    /api/v1/domain/domains/:id
GET    /api/v1/domain/records?domainId=:id
POST   /api/v1/domain/domains/:domainId/records
PUT    /api/v1/domain/domains/:domainId/records/:recordId
DELETE /api/v1/domain/domains/:domainId/records/:recordId
PUT    /api/v1/domain/domains/:domainId/records/:recordId/status
PUT    /api/v1/domain/domains/:domainId/records/:recordId/remark
POST   /api/v1/domain/domains/:domainId/records/batch
```

## dnsmgr API Mapping

AutoOps proxies to dnsmgr public API endpoints:

```text
POST /api/domain
POST /api/domain/:id
POST /api/record/data/:id
POST /api/record/add/:id
POST /api/record/update/:id
POST /api/record/delete/:id
POST /api/record/status/:id
POST /api/record/remark/:id
POST /api/record/batch/:id
```

Important protocol detail:

- For record update/delete/status/remark, dnsmgr path `:id` is the domain ID.
- The record ID is sent in the form field `recordid`.
- Batch operations send selected records as `recordinfo` JSON.

## Configuration

Use `.env` or deployment environment variables:

```bash
DNSMGR_ENABLED=true
DNSMGR_BASE_URL=https://domain.example.com
DNSMGR_UID=1002
DNSMGR_API_KEY=<secret>
```

The secret must not be committed.

dnsmgr request signing:

```text
uid=<uid>
timestamp=<unix_seconds>
sign=md5(uid + timestamp + apiKey)
```

## Deployment Acceptance

After pulling the latest `dev` branch on the server:

1. Confirm `.env` contains the dnsmgr variables above.
2. Rebuild/restart the API and Web services.
3. Confirm migration `0009_drop_domain_sync_tables.sql` has run.
4. Open `/integration/domain`.
5. Click the health-check button; it should report dnsmgr reachable.
6. Click the refresh-domains button; the page should show live dnsmgr domains.
7. Click one domain row; the record table should load live records.
8. Test a low-risk DNS record operation first, preferably remark update or pause/enable on a test record.
9. Confirm the same change is visible in dnsmgr directly.

## Verification Already Run Locally

Local checks completed before this document:

- `go test -mod=mod ./...` passed before final doc-only update.
- Targeted `go test -mod=mod ./api/domain/... ./router/integration ./common/config` passed with workspace `GOCACHE`.
- Targeted `go vet ./api/domain/... ./router/integration ./common/config` passed with workspace `GOCACHE`.
- `DomainManagement.vue` parsed successfully with `@vue/compiler-sfc`.
- `web/src/api/integration.js` and `web/src/router/integration.js` parsed successfully with `@babel/parser`.

Known local verification limitation:

- Full `npm run build` still times out locally at `Building for production...`; this is the existing frontend build hang/OOM class of issue, not a dnsmgr page parse error.
