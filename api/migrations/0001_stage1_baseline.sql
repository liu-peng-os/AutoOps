-- +goose Up

CREATE TABLE IF NOT EXISTS config_account (
    id BIGSERIAL PRIMARY KEY,
    alias VARCHAR(128) NOT NULL,
    host VARCHAR(128) NOT NULL,
    port BIGINT NOT NULL,
    name VARCHAR(128) NOT NULL,
    password TEXT NOT NULL,
    type BIGINT NOT NULL,
    remark TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE,
    updated_at TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_config_account_host ON config_account (host);
CREATE INDEX IF NOT EXISTS idx_config_account_type ON config_account (type);

CREATE TABLE IF NOT EXISTS config_ecsauth (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    type BIGINT NOT NULL,
    username VARCHAR(64),
    password VARCHAR(256),
    public_key TEXT,
    port BIGINT DEFAULT 22,
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    remark VARCHAR(500)
);

CREATE INDEX IF NOT EXISTS idx_config_ecsauth_type ON config_ecsauth (type);
CREATE INDEX IF NOT EXISTS idx_config_ecsauth_name ON config_ecsauth (name);

CREATE TABLE IF NOT EXISTS config_keymanage (
    id BIGSERIAL PRIMARY KEY,
    key_type BIGINT NOT NULL,
    key_id TEXT NOT NULL,
    key_secret TEXT NOT NULL,
    remark TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE,
    updated_at TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_config_keymanage_key_type ON config_keymanage (key_type);

CREATE TABLE IF NOT EXISTS config_sync_schedule (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    cron_expr VARCHAR(100) NOT NULL,
    key_types TEXT NOT NULL,
    status BIGINT NOT NULL DEFAULT 1,
    last_run_time TIMESTAMP WITHOUT TIME ZONE,
    next_run_time TIMESTAMP WITHOUT TIME ZONE,
    sync_log TEXT,
    remark TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE,
    updated_at TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_config_sync_schedule_status ON config_sync_schedule (status);
CREATE INDEX IF NOT EXISTS idx_config_sync_schedule_next_run_time ON config_sync_schedule (next_run_time);

CREATE TABLE IF NOT EXISTS cmdb_group (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT NOT NULL DEFAULT 0,
    name VARCHAR(50) NOT NULL,
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_cmdb_group_parent_id ON cmdb_group (parent_id);
CREATE INDEX IF NOT EXISTS idx_cmdb_group_name ON cmdb_group (name);

CREATE TABLE IF NOT EXISTS cmdb_host (
    id BIGSERIAL PRIMARY KEY,
    host_name VARCHAR(64) NOT NULL,
    group_id BIGINT NOT NULL,
    private_ip VARCHAR(64),
    public_ip VARCHAR(64),
    ssh_ip VARCHAR(64) NOT NULL,
    ssh_name VARCHAR(64),
    ssh_key_id BIGINT,
    ssh_port BIGINT DEFAULT 22,
    remark VARCHAR(500),
    vendor BIGINT,
    region VARCHAR(64),
    instance_id VARCHAR(128),
    name VARCHAR(64) NOT NULL,
    os VARCHAR(128),
    status BIGINT,
    cpu VARCHAR(32),
    memory VARCHAR(32),
    disk VARCHAR(128),
    billing_type VARCHAR(32),
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expire_time TIMESTAMP WITHOUT TIME ZONE,
    update_time TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_cmdb_host_group_id ON cmdb_host (group_id);
CREATE INDEX IF NOT EXISTS idx_cmdb_host_ssh_key_id ON cmdb_host (ssh_key_id);
CREATE INDEX IF NOT EXISTS idx_cmdb_host_instance_id ON cmdb_host (instance_id);
CREATE INDEX IF NOT EXISTS idx_cmdb_host_vendor ON cmdb_host (vendor);
CREATE INDEX IF NOT EXISTS idx_cmdb_host_status ON cmdb_host (status);

CREATE TABLE IF NOT EXISTS monitor_agent (
    id BIGSERIAL PRIMARY KEY,
    host_id BIGINT NOT NULL,
    host_name VARCHAR(128),
    version VARCHAR(32) DEFAULT '1.0.0',
    status BIGINT,
    install_path VARCHAR(256),
    port BIGINT DEFAULT 9100,
    pid BIGINT,
    last_heartbeat TIMESTAMP WITHOUT TIME ZONE,
    update_time TIMESTAMP WITHOUT TIME ZONE,
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    error_msg TEXT,
    install_progress BIGINT DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_monitor_agent_host_id ON monitor_agent (host_id);
CREATE INDEX IF NOT EXISTS idx_monitor_agent_status ON monitor_agent (status);

CREATE TABLE IF NOT EXISTS tool_link (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    icon VARCHAR(500),
    link VARCHAR(500) NOT NULL,
    sort BIGINT DEFAULT 0,
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_tool_link_sort ON tool_link (sort);

CREATE TABLE IF NOT EXISTS tool_service_deploy (
    id BIGSERIAL PRIMARY KEY,
    service_name VARCHAR(64) NOT NULL,
    service_id VARCHAR(64) NOT NULL,
    version VARCHAR(64) NOT NULL,
    host_id BIGINT NOT NULL,
    host_ip VARCHAR(64) NOT NULL,
    install_dir VARCHAR(255) NOT NULL,
    container_name VARCHAR(128),
    ports VARCHAR(255),
    env_vars TEXT,
    status BIGINT DEFAULT 0,
    deploy_log TEXT,
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_tool_service_deploy_service_id ON tool_service_deploy (service_id);
CREATE INDEX IF NOT EXISTS idx_tool_service_deploy_host_id ON tool_service_deploy (host_id);
CREATE INDEX IF NOT EXISTS idx_tool_service_deploy_status ON tool_service_deploy (status);
