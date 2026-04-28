-- +goose Up

CREATE TABLE IF NOT EXISTS sys_operation_log (
    id BIGSERIAL PRIMARY KEY,
    admin_id BIGINT NOT NULL,
    username VARCHAR(64) NOT NULL,
    method VARCHAR(64) NOT NULL,
    ip VARCHAR(64),
    url VARCHAR(500),
    description VARCHAR(255),
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sys_operation_log_admin_id ON sys_operation_log (admin_id);
CREATE INDEX IF NOT EXISTS idx_sys_operation_log_create_time ON sys_operation_log (create_time);

CREATE TABLE IF NOT EXISTS cmdb_sql (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type BIGINT NOT NULL,
    account_id BIGINT NOT NULL,
    group_id BIGINT NOT NULL,
    tags VARCHAR(255),
    description VARCHAR(500),
    created_at TIMESTAMP WITHOUT TIME ZONE,
    updated_at TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_cmdb_sql_type ON cmdb_sql (type);
CREATE INDEX IF NOT EXISTS idx_cmdb_sql_account_id ON cmdb_sql (account_id);
CREATE INDEX IF NOT EXISTS idx_cmdb_sql_group_id ON cmdb_sql (group_id);

CREATE TABLE IF NOT EXISTS cmdb_sql_log (
    id BIGSERIAL PRIMARY KEY,
    instance_id VARCHAR(64) NOT NULL,
    database VARCHAR(128) NOT NULL,
    operation_type VARCHAR(32) NOT NULL,
    sql_content TEXT NOT NULL,
    exec_user VARCHAR(64) NOT NULL,
    ip VARCHAR(64) NOT NULL,
    scanned_rows BIGINT DEFAULT 0,
    affected_rows BIGINT DEFAULT 0,
    execution_time BIGINT DEFAULT 0,
    returned_rows BIGINT DEFAULT 0,
    result VARCHAR(32) NOT NULL,
    query_time TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_cmdb_sql_log_query_time ON cmdb_sql_log (query_time);
CREATE INDEX IF NOT EXISTS idx_cmdb_sql_log_instance_id ON cmdb_sql_log (instance_id);
CREATE INDEX IF NOT EXISTS idx_cmdb_sql_log_result ON cmdb_sql_log (result);

CREATE TABLE IF NOT EXISTS k8s_cluster (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    version VARCHAR(50) NOT NULL,
    status BIGINT NOT NULL DEFAULT 1,
    credential TEXT,
    description TEXT,
    cluster_type BIGINT NOT NULL DEFAULT 1,
    node_count BIGINT DEFAULT 0,
    ready_nodes BIGINT DEFAULT 0,
    master_nodes BIGINT DEFAULT 0,
    worker_nodes BIGINT DEFAULT 0,
    last_sync_at TIMESTAMP WITHOUT TIME ZONE,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_k8s_cluster_name ON k8s_cluster (name);
CREATE INDEX IF NOT EXISTS idx_k8s_cluster_status ON k8s_cluster (status);
CREATE INDEX IF NOT EXISTS idx_k8s_cluster_cluster_type ON k8s_cluster (cluster_type);
