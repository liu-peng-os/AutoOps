# Phase 2 dnsmgr Domain Integration

Date: 2026-04-29

## Scope

The first Phase 2 external-system slice is dnsmgr domain integration.

AutoOps now treats dnsmgr as the source of truth and acts as a live API operation panel:

- AutoOps does not sync or persist dnsmgr domain and record data locally.
- Domain and DNS record lists are queried from dnsmgr on demand.
- DNS record create/update/delete/status/remark/batch actions are proxied to dnsmgr.
- dnsmgr credentials stay on the backend and are injected by environment variables.
- SSL certificate order APIs are intentionally not integrated in this slice.

## Backend

AutoOps exposes stable internal endpoints under `/api/v1/domain` and signs outbound dnsmgr requests server-side.

Current endpoints:

- `GET /api/v1/domain/health`
- `GET /api/v1/domain/domains`
- `GET /api/v1/domain/domains/:id`
- `GET /api/v1/domain/records?domainId=:id`
- `POST /api/v1/domain/domains/:domainId/records`
- `PUT /api/v1/domain/domains/:domainId/records/:recordId`
- `DELETE /api/v1/domain/domains/:domainId/records/:recordId`
- `PUT /api/v1/domain/domains/:domainId/records/:recordId/status`
- `PUT /api/v1/domain/domains/:domainId/records/:recordId/remark`
- `POST /api/v1/domain/domains/:domainId/records/batch`

Removed local persistence:

- `domain_zone`
- `domain_record`
- `domain_sync_run`

Migration `0009_drop_domain_sync_tables.sql` drops those tables because dnsmgr data should not be duplicated inside AutoOps.

## Configuration

For server deployments, use environment variables so secrets are not committed:

- `DNSMGR_ENABLED=true`
- `DNSMGR_BASE_URL=https://domain.example.com`
- `DNSMGR_UID=1002`
- `DNSMGR_API_KEY=<secret>`

dnsmgr API authentication uses request parameters:

- `uid`
- `timestamp`
- `sign = md5(uid + timestamp + apiKey)`

Example configuration shape:

```yaml
integrations:
  systems:
    - key: domain-management
      displayName: dnsmgr
      category: domain
      provider: dnsmgr
      mode: live-api
      baseUrl: https://domain.example.com
      enabled: true
      capabilities:
        - domain:list
        - domain:detail
        - record:list
        - record:create
        - record:update
        - record:delete
        - record:status
        - record:remark
        - record:batch
      metadata:
        uid: "1002"
        apiKey: "<set from secret/config>"
```

## Frontend

The `/integration/domain` page is a live dnsmgr operation panel:

- health check
- live domain list
- domain detail
- domain login link
- live DNS record list
- add/edit/delete record
- enable/pause record
- change record remark
- batch enable/pause/delete/change remark

The old Sync Now flow was removed to avoid stale data and split ownership.

## Follow-Up

The private dnsmgr deployment should be used as the final API compatibility source during server acceptance.
If dnsmgr returns provider-specific field names, normalize them in the backend proxy without adding local persistence.
