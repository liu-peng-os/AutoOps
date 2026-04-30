# Phase 2 dnsmgr Domain Integration

Date: 2026-04-29

## Scope

The first Phase 2 external-system slice is dnsmgr domain integration.

The initial implementation is intentionally read-only:

- AutoOps connects to a private dnsmgr deployment through a backend adapter.
- dnsmgr remains the source of truth for DNS zones and records.
- AutoOps stores normalized domain zones and DNS records for monitoring and review.
- DNS write-back is not enabled in this slice.

## Backend

New endpoints:

- `GET /api/v1/domain/health`
- `POST /api/v1/domain/sync`
- `GET /api/v1/domain/sync/last`
- `GET /api/v1/domain/zones`
- `GET /api/v1/domain/records`

New persistence:

- `domain_zone`
- `domain_record`
- `domain_sync_run`

The adapter reads dnsmgr connection settings from `integrations.systems`.
The expected system can use `provider: dnsmgr`, `key: dnsmgr`, or `key: domain-management`.

For server deployments, prefer environment variables so secrets are not committed:

- `DNSMGR_ENABLED=true`
- `DNSMGR_BASE_URL=https://domain.example.com`
- `DNSMGR_UID=1002`
- `DNSMGR_API_KEY=<secret>`

Useful metadata keys:

- `uid`
- `apiKey`
- `healthPath`
- `zonesPath`
- `recordsPath`

dnsmgr API authentication uses request parameters:

- `uid`
- `timestamp`
- `sign = md5(uid + timestamp + apiKey)`

Default paths:

- `zonesPath: /api/domain`
- `recordsPath: /api/record/data/{zone_id}`

Example configuration shape:

```yaml
integrations:
  systems:
    - key: domain-management
      displayName: dnsmgr
      category: domain
      provider: dnsmgr
      mode: read-only
      baseUrl: https://domain.example.com
      enabled: true
      capabilities:
        - domain:list
        - record:list
      metadata:
        uid: "1002"
        apiKey: "<set from secret/config>"
```

## Frontend

The `/integration/domain` page now shows:

- dnsmgr health status
- last sync status
- synced zone and record counts
- zone table
- DNS record table
- manual sync action

## Follow-Up

Before enabling production sync, verify the private dnsmgr deployment's actual API paths and response shape.
After read-only sync is stable, write-back can be added behind approval and audit logging.
