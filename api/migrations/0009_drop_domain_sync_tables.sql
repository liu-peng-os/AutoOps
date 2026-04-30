-- +goose Up
DROP TABLE IF EXISTS domain_sync_run;
DROP TABLE IF EXISTS domain_record;
DROP TABLE IF EXISTS domain_zone;

-- +goose Down
CREATE TABLE IF NOT EXISTS domain_zone (
    id BIGSERIAL PRIMARY KEY,
    name varchar(255) NOT NULL,
    display_name varchar(255) DEFAULT '',
    provider varchar(64) DEFAULT '',
    status varchar(32) DEFAULT 'unknown',
    external_source varchar(64) NOT NULL,
    external_id varchar(255) NOT NULL,
    raw_data TEXT,
    last_synced_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_domain_zone_external UNIQUE (external_source, external_id)
);

CREATE INDEX IF NOT EXISTS idx_domain_zone_name ON domain_zone (name);
CREATE INDEX IF NOT EXISTS idx_domain_zone_status ON domain_zone (status);

CREATE TABLE IF NOT EXISTS domain_record (
    id BIGSERIAL PRIMARY KEY,
    zone_id BIGINT NOT NULL,
    zone_name varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    type varchar(32) NOT NULL,
    value TEXT NOT NULL,
    line varchar(128) DEFAULT '',
    ttl BIGINT DEFAULT 0,
    priority BIGINT DEFAULT 0,
    status varchar(32) DEFAULT 'unknown',
    external_source varchar(64) NOT NULL,
    external_id varchar(255) NOT NULL,
    raw_data TEXT,
    last_synced_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_domain_record_external UNIQUE (external_source, external_id)
);

CREATE INDEX IF NOT EXISTS idx_domain_record_zone_id ON domain_record (zone_id);
CREATE INDEX IF NOT EXISTS idx_domain_record_zone_name ON domain_record (zone_name);
CREATE INDEX IF NOT EXISTS idx_domain_record_name ON domain_record (name);
CREATE INDEX IF NOT EXISTS idx_domain_record_type ON domain_record (type);
CREATE INDEX IF NOT EXISTS idx_domain_record_status ON domain_record (status);

CREATE TABLE IF NOT EXISTS domain_sync_run (
    id BIGSERIAL PRIMARY KEY,
    external_source varchar(64) NOT NULL,
    status varchar(32) NOT NULL,
    zones_synced BIGINT DEFAULT 0,
    records_synced BIGINT DEFAULT 0,
    message TEXT,
    started_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    finished_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_domain_sync_run_source ON domain_sync_run (external_source);
CREATE INDEX IF NOT EXISTS idx_domain_sync_run_status ON domain_sync_run (status);
CREATE INDEX IF NOT EXISTS idx_domain_sync_run_started_at ON domain_sync_run (started_at);
