-- +goose Up

-- Remaining schema migrated from docker/mysql/devops001.sql.
-- This migration creates PostgreSQL tables that were not covered by 0001-0004.

CREATE TABLE IF NOT EXISTS ai_agent_chat_history (
    id BIGSERIAL NOT NULL,
    session_id varchar(64) NOT NULL,
    user_id BIGINT NOT NULL,
    role varchar(20) NOT NULL,
    message TEXT NOT NULL,
    intent varchar(50) DEFAULT NULL,
    intent_conf decimal(3,2) DEFAULT NULL,
    entities JSONB DEFAULT NULL,
    task_id BIGINT DEFAULT NULL,
    task_type varchar(50) DEFAULT NULL,
    status SMALLINT DEFAULT 1,
    error_msg TEXT,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_ai_agent_chat_history_session_id ON ai_agent_chat_history (session_id);
CREATE INDEX IF NOT EXISTS idx_ai_agent_chat_history_user_id ON ai_agent_chat_history (user_id);
CREATE INDEX IF NOT EXISTS idx_ai_agent_chat_history_task_id ON ai_agent_chat_history (task_id);

CREATE TABLE IF NOT EXISTS ai_agent_task (
    id BIGSERIAL NOT NULL,
    type varchar(50) NOT NULL,
    name varchar(200) NOT NULL,
    description varchar(500) DEFAULT NULL,
    priority BIGINT DEFAULT 5,
    status BIGINT NOT NULL,
    params TEXT,
    result TEXT,
    error_msg TEXT,
    retry BIGINT DEFAULT 0,
    max_retry BIGINT DEFAULT 3,
    start_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    end_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_ai_agent_task_type ON ai_agent_task (type);
CREATE INDEX IF NOT EXISTS idx_ai_agent_task_status ON ai_agent_task (status);

CREATE TABLE IF NOT EXISTS ai_model (
    id BIGSERIAL NOT NULL,
    name varchar(100) NOT NULL,
    type varchar(50) NOT NULL,
    url varchar(500) NOT NULL,
    api_key varchar(200) NOT NULL,
    model varchar(100) NOT NULL,
    status BIGINT NOT NULL DEFAULT 1,
    remark varchar(500) DEFAULT NULL,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS app_service_release (
    id BIGSERIAL NOT NULL,
    title varchar(255) NOT NULL,
    business_group_id BIGINT NOT NULL,
    impact_feature TEXT,
    applicant_id BIGINT NOT NULL,
    applicant_name varchar(100) NOT NULL,
    owner_approver_id BIGINT DEFAULT NULL,
    owner_approver_name varchar(100) DEFAULT NULL,
    security_approver_id BIGINT DEFAULT NULL,
    security_approver_name varchar(100) DEFAULT NULL,
    test_approver_id BIGINT DEFAULT NULL,
    test_approver_name varchar(100) DEFAULT NULL,
    owner_approval_status BIGINT DEFAULT 1,
    security_approval_status BIGINT DEFAULT 1,
    test_approval_status BIGINT DEFAULT 1,
    owner_approval_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    security_approval_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    test_approval_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    owner_approval_remark TEXT,
    security_approval_remark TEXT,
    test_approval_remark TEXT,
    deploy_status BIGINT DEFAULT 1,
    regression_test_status BIGINT DEFAULT 1,
    status BIGINT DEFAULT 1,
    start_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    end_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    duration BIGINT DEFAULT 0,
    service_count BIGINT DEFAULT 0,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    deleted_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_app_service_release_deleted_at ON app_service_release (deleted_at);

CREATE TABLE IF NOT EXISTS app_service_release_item (
    id BIGSERIAL NOT NULL,
    release_id BIGINT NOT NULL,
    app_id BIGINT NOT NULL,
    app_name varchar(255) NOT NULL,
    app_code varchar(100) NOT NULL,
    project_name varchar(255) NOT NULL,
    repo_url varchar(500) NOT NULL,
    branch varchar(100) DEFAULT 'master',
    commit_id varchar(100) NOT NULL,
    impact_feature TEXT,
    function_module TEXT,
    db_change TEXT,
    config_change TEXT,
    remark TEXT,
    jenkins_env_id BIGINT DEFAULT NULL,
    jenkins_job_url varchar(500) DEFAULT NULL,
    parameters TEXT,
    build_number BIGINT DEFAULT 0,
    log_url varchar(500) DEFAULT NULL,
    status BIGINT DEFAULT 1,
    start_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    end_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    duration BIGINT DEFAULT 0,
    error_message TEXT,
    execute_order BIGINT DEFAULT 0,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    deleted_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_app_service_release_item_release_id ON app_service_release_item (release_id);
CREATE INDEX IF NOT EXISTS idx_app_service_release_item_deleted_at ON app_service_release_item (deleted_at);

CREATE TABLE IF NOT EXISTS app_sh_release (
    id BIGSERIAL NOT NULL,
    title varchar(255) NOT NULL,
    reason TEXT NOT NULL,
    business_group_id BIGINT NOT NULL,
    app_id BIGINT NOT NULL,
    app_name varchar(255) NOT NULL,
    app_code varchar(100) NOT NULL,
    applicant_id BIGINT NOT NULL,
    applicant_name varchar(100) NOT NULL,
    approver_id BIGINT DEFAULT NULL,
    approver_name varchar(100) DEFAULT NULL,
    executor_id BIGINT DEFAULT NULL,
    executor_name varchar(100) DEFAULT NULL,
    execute_dir varchar(500) NOT NULL,
    script_content TEXT NOT NULL,
    approval_status BIGINT DEFAULT 1,
    approval_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    approval_remark TEXT,
    execute_status BIGINT DEFAULT 1,
    status BIGINT DEFAULT 1,
    start_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    end_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    duration BIGINT DEFAULT 0,
    jenkins_env_id BIGINT DEFAULT NULL,
    build_number BIGINT DEFAULT 0,
    log_url varchar(500) DEFAULT NULL,
    error_message TEXT,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    deleted_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    parameters TEXT,
    server_host_id BIGINT NOT NULL,
    pull_code_start_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    pull_code_end_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    script_output TEXT,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_app_sh_release_deleted_at ON app_sh_release (deleted_at);

CREATE TABLE IF NOT EXISTS bastion_host_authorization (
    id BIGSERIAL NOT NULL,
    subject_type varchar(191) NOT NULL,
    subject_id BIGINT NOT NULL,
    host_id BIGINT NOT NULL,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_subject_host ON bastion_host_authorization (subject_type, subject_id, host_id);

CREATE TABLE IF NOT EXISTS bastion_user_group (
    id BIGSERIAL NOT NULL,
    name varchar(64) NOT NULL,
    description varchar(500) DEFAULT NULL,
    member_ids JSONB DEFAULT NULL,
    status BIGINT NOT NULL DEFAULT 1,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_bastion_user_group_name ON bastion_user_group (name);

CREATE TABLE IF NOT EXISTS cmdb_snmp_devices (
    id BIGSERIAL NOT NULL,
    device_name varchar(100) NOT NULL,
    device_type varchar(20) NOT NULL,
    ip_address varchar(50) NOT NULL,
    snmp_version varchar(10) DEFAULT 'v2c',
    snmp_community varchar(100) DEFAULT 'public',
    snmp_port BIGINT DEFAULT 161,
    remote_port BIGINT DEFAULT NULL,
    remote_username varchar(100) DEFAULT NULL,
    remote_password varchar(255) DEFAULT NULL,
    remote_domain varchar(100) DEFAULT NULL,
    monitor_enabled BIGINT DEFAULT 1,
    location varchar(200) DEFAULT NULL,
    owner varchar(100) DEFAULT NULL,
    remark TEXT,
    status varchar(20) DEFAULT 'offline',
    last_seen TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    ssh_port BIGINT DEFAULT 22,
    ssh_username varchar(100) DEFAULT NULL,
    ssh_password varchar(255) DEFAULT NULL,
    rdp_port BIGINT DEFAULT 3389,
    rdp_username varchar(100) DEFAULT NULL,
    rdp_password varchar(255) DEFAULT NULL,
    rdp_domain varchar(100) DEFAULT NULL,
    exporter_installed BIGINT DEFAULT 0,
    exporter_port BIGINT DEFAULT 9182,
    exporter_version varchar(50) DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_cmdb_snmp_devices_ip_address ON cmdb_snmp_devices (ip_address);

CREATE TABLE IF NOT EXISTS cmdb_sql_records (
    id BIGSERIAL NOT NULL,
    instance_id varchar(64) NOT NULL,
    database varchar(128) NOT NULL,
    operation_type varchar(32) NOT NULL,
    sql_content TEXT NOT NULL,
    exec_user varchar(64) NOT NULL,
    scanned_rows BIGINT DEFAULT 0,
    affected_rows BIGINT DEFAULT 0,
    execution_time BIGINT DEFAULT 0,
    returned_rows BIGINT DEFAULT 0,
    result varchar(32) NOT NULL,
    query_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    name varchar(64) NOT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_cmdb_sql_records_query_time ON cmdb_sql_records (query_time);

CREATE TABLE IF NOT EXISTS db (
    id BIGSERIAL NOT NULL,
    code varchar(36) NOT NULL,
    name varchar(100) NOT NULL,
    database varchar(500) NOT NULL,
    remark varchar(500) DEFAULT NULL,
    instance_id BIGINT NOT NULL,
    instance_code varchar(36) DEFAULT NULL,
    status BIGINT DEFAULT 1,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    creator varchar(64) DEFAULT NULL,
    creator_id BIGINT DEFAULT NULL,
    modifier varchar(64) DEFAULT NULL,
    modifier_id BIGINT DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_db_code ON db (code);
CREATE INDEX IF NOT EXISTS idx_db_instance_id ON db (instance_id);

CREATE TABLE IF NOT EXISTS db_es_instance (
    id BIGSERIAL NOT NULL,
    code varchar(36) NOT NULL,
    name varchar(100) NOT NULL,
    protocol varchar(10) DEFAULT 'http',
    host varchar(255) NOT NULL,
    port BIGINT DEFAULT 9200,
    username varchar(100) DEFAULT NULL,
    password varchar(500) DEFAULT '',
    remark varchar(500) DEFAULT NULL,
    ssh_tunnel_machine_id BIGINT DEFAULT 0,
    status BIGINT DEFAULT 1,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    creator varchar(64) DEFAULT NULL,
    creator_id BIGINT DEFAULT NULL,
    modifier varchar(64) DEFAULT NULL,
    modifier_id BIGINT DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_db_es_instance_code ON db_es_instance (code);

CREATE TABLE IF NOT EXISTS db_export_task (
    id BIGSERIAL NOT NULL,
    task_id varchar(36) NOT NULL,
    db_id BIGINT NOT NULL,
    db_name varchar(100) NOT NULL,
    export_type varchar(20) NOT NULL,
    status varchar(20) NOT NULL,
    file_path varchar(500) DEFAULT NULL,
    file_size BIGINT DEFAULT NULL,
    error_message TEXT,
    start_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    end_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    creator varchar(64) DEFAULT NULL,
    creator_id BIGINT DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_db_export_task_task_id ON db_export_task (task_id);

CREATE TABLE IF NOT EXISTS db_instance (
    id BIGSERIAL NOT NULL,
    code varchar(36) NOT NULL,
    name varchar(100) NOT NULL,
    type varchar(20) NOT NULL,
    host varchar(100) NOT NULL,
    port BIGINT NOT NULL,
    network varchar(20) DEFAULT 'tcp',
    params varchar(500) DEFAULT NULL,
    username varchar(100) NOT NULL,
    password varchar(500) NOT NULL,
    remark varchar(500) DEFAULT NULL,
    ssh_tunnel_machine_id BIGINT DEFAULT 0,
    status BIGINT DEFAULT 1,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    creator varchar(64) DEFAULT NULL,
    creator_id BIGINT DEFAULT NULL,
    modifier varchar(64) DEFAULT NULL,
    modifier_id BIGINT DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_db_instance_code ON db_instance (code);

CREATE TABLE IF NOT EXISTS db_mongo_instance (
    id BIGSERIAL NOT NULL,
    code varchar(36) NOT NULL,
    name varchar(100) NOT NULL,
    uri varchar(500) NOT NULL,
    ssh_tunnel_machine_id BIGINT DEFAULT 0,
    remark varchar(500) DEFAULT NULL,
    status BIGINT DEFAULT 1,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    creator varchar(64) DEFAULT NULL,
    creator_id BIGINT DEFAULT NULL,
    modifier varchar(64) DEFAULT NULL,
    modifier_id BIGINT DEFAULT NULL,
    type varchar(20) DEFAULT 'mongodb',
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_db_mongo_instance_code ON db_mongo_instance (code);

CREATE TABLE IF NOT EXISTS db_redis_instance (
    id BIGSERIAL NOT NULL,
    code varchar(36) NOT NULL,
    name varchar(100) NOT NULL,
    mode varchar(20) NOT NULL,
    host varchar(300) NOT NULL,
    port BIGINT DEFAULT 0,
    db BIGINT DEFAULT 0,
    username varchar(100) DEFAULT NULL,
    password varchar(500) NOT NULL,
    redis_node_password varchar(500) DEFAULT '',
    remark varchar(500) DEFAULT NULL,
    ssh_tunnel_machine_id BIGINT DEFAULT 0,
    status BIGINT DEFAULT 1,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    creator varchar(64) DEFAULT NULL,
    creator_id BIGINT DEFAULT NULL,
    modifier varchar(64) DEFAULT NULL,
    modifier_id BIGINT DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_db_redis_instance_code ON db_redis_instance (code);

CREATE TABLE IF NOT EXISTS db_sql (
    id BIGSERIAL NOT NULL,
    db_id BIGINT NOT NULL,
    db varchar(100) NOT NULL,
    name varchar(100) NOT NULL,
    type BIGINT DEFAULT 1,
    sql TEXT NOT NULL,
    remark varchar(500) DEFAULT NULL,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    creator varchar(64) DEFAULT NULL,
    creator_id BIGINT DEFAULT NULL,
    modifier varchar(64) DEFAULT NULL,
    modifier_id BIGINT DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_db_sql_db_id ON db_sql (db_id);

CREATE TABLE IF NOT EXISTS db_sql_exec (
    id BIGSERIAL NOT NULL,
    db_id BIGINT NOT NULL,
    db_name varchar(100) NOT NULL,
    table_name varchar(100) DEFAULT NULL,
    type SMALLINT NOT NULL,
    sql TEXT NOT NULL,
    old_value TEXT,
    remark varchar(500) DEFAULT NULL,
    status SMALLINT NOT NULL,
    res TEXT,
    exec_time BIGINT DEFAULT NULL,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    creator varchar(64) DEFAULT NULL,
    creator_id BIGINT DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_db_sql_exec_db_id ON db_sql_exec (db_id);
CREATE INDEX IF NOT EXISTS idx_db_sql_exec_create_time ON db_sql_exec (create_time);

CREATE TABLE IF NOT EXISTS knowledge_base (
    id BIGSERIAL NOT NULL,
    type varchar(50) NOT NULL,
    category varchar(50) NOT NULL,
    title varchar(500) NOT NULL,
    content TEXT NOT NULL,
    keywords varchar(1000) DEFAULT NULL,
    tags varchar(500) DEFAULT NULL,
    score decimal(3,2) DEFAULT 0.50,
    use_count BIGINT DEFAULT 0,
    enabled BIGINT DEFAULT 1,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_knowledge_base_type ON knowledge_base (type);
CREATE INDEX IF NOT EXISTS idx_knowledge_base_category ON knowledge_base (category);

CREATE TABLE IF NOT EXISTS monitor_alert_mute (
    id BIGSERIAL NOT NULL,
    name varchar(200) NOT NULL,
    rule_ids TEXT,
    target_type varchar(20) DEFAULT NULL,
    target_ids TEXT,
    label_matchers TEXT,
    start_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    end_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    periodic SMALLINT DEFAULT 0,
    period_days varchar(50) DEFAULT NULL,
    period_start varchar(10) DEFAULT NULL,
    period_end varchar(10) DEFAULT NULL,
    status SMALLINT DEFAULT 1,
    remark varchar(500) DEFAULT NULL,
    created_by varchar(100) DEFAULT NULL,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_mute_end_time ON monitor_alert_mute (end_time);

CREATE TABLE IF NOT EXISTS monitor_alert_policy (
    id BIGSERIAL NOT NULL,
    name varchar(200) NOT NULL,
    description varchar(500) DEFAULT NULL,
    group_id BIGINT NOT NULL,
    group_name varchar(100) DEFAULT NULL,
    template_ids TEXT NOT NULL,
    rule_count BIGINT DEFAULT 0,
    robot_ids TEXT,
    silence_period BIGINT DEFAULT 300,
    repeat_interval BIGINT DEFAULT 600,
    max_notify_count BIGINT DEFAULT 10,
    send_recovery SMALLINT DEFAULT 1,
    cron_pattern varchar(100) DEFAULT '@every 30s',
    status SMALLINT DEFAULT 0,
    append_tags TEXT,
    created_by varchar(100) DEFAULT NULL,
    updated_by varchar(100) DEFAULT NULL,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    health_status varchar(20) DEFAULT 'healthy',
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_name ON monitor_alert_policy (name);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_policy_group_id ON monitor_alert_policy (group_id);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_policy_status ON monitor_alert_policy (status);

CREATE TABLE IF NOT EXISTS monitor_alert_policy_rule (
    id BIGSERIAL NOT NULL,
    policy_id BIGINT NOT NULL,
    template_id BIGINT NOT NULL,
    rule_name varchar(200) NOT NULL,
    category varchar(50) DEFAULT NULL,
    metric_key varchar(100) DEFAULT NULL,
    unit varchar(20) DEFAULT NULL,
    operator varchar(10) NOT NULL,
    threshold decimal(10,2) NOT NULL,
    prom_ql TEXT NOT NULL,
    severity SMALLINT DEFAULT 2,
    duration BIGINT DEFAULT 60,
    continuous_times BIGINT DEFAULT 1,
    recover_duration BIGINT DEFAULT 60,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    description varchar(500) DEFAULT NULL,
    health_status varchar(20) DEFAULT 'healthy',
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_policy_rule_policy_id ON monitor_alert_policy_rule (policy_id);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_policy_rule_template_id ON monitor_alert_policy_rule (template_id);

CREATE TABLE IF NOT EXISTS monitor_alert_record (
    id BIGSERIAL NOT NULL,
    rule_id BIGINT NOT NULL,
    rule_name varchar(200) DEFAULT NULL,
    target_type varchar(20) DEFAULT NULL,
    target_id BIGINT DEFAULT NULL,
    target_name varchar(200) DEFAULT NULL,
    metric_name varchar(100) DEFAULT NULL,
    current_value decimal(10,2) DEFAULT NULL,
    threshold decimal(10,2) DEFAULT NULL,
    severity SMALLINT DEFAULT NULL,
    message TEXT,
    labels TEXT,
    status varchar(20) DEFAULT 'firing',
    start_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    end_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    duration BIGINT DEFAULT NULL,
    notify_sent SMALLINT DEFAULT 0,
    notify_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    notify_count BIGINT DEFAULT 0,
    notify_error TEXT,
    group_id BIGINT DEFAULT 0,
    group_name varchar(100) DEFAULT NULL,
    fingerprint varchar(64) DEFAULT NULL,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    policy_id BIGINT DEFAULT NULL,
    policy_name varchar(200) DEFAULT NULL,
    template_id BIGINT DEFAULT NULL,
    summary varchar(500) DEFAULT NULL,
    handler varchar(100) DEFAULT NULL,
    handle_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    handle_note TEXT,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_monitor_alert_record_fingerprint ON monitor_alert_record (fingerprint);
CREATE INDEX IF NOT EXISTS idx_rule_target ON monitor_alert_record (rule_id, target_id);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_record_target_type ON monitor_alert_record (target_type);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_record_metric_name ON monitor_alert_record (metric_name);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_record_severity ON monitor_alert_record (severity);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_record_status ON monitor_alert_record (status);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_record_start_time ON monitor_alert_record (start_time);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_record_group_id ON monitor_alert_record (group_id);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_record_policy_id ON monitor_alert_record (policy_id);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_record_template_id ON monitor_alert_record (template_id);

CREATE TABLE IF NOT EXISTS monitor_alert_rule (
    id BIGSERIAL NOT NULL,
    name varchar(200) NOT NULL,
    description varchar(500) DEFAULT NULL,
    target_type varchar(20) NOT NULL,
    target_ids TEXT,
    prom_ql TEXT NOT NULL,
    duration BIGINT DEFAULT 60,
    metric_name varchar(100) DEFAULT NULL,
    operator varchar(10) DEFAULT NULL,
    threshold decimal(10,2) DEFAULT NULL,
    severity SMALLINT DEFAULT 2,
    robot_ids TEXT,
    silence_period BIGINT DEFAULT 300,
    repeat_interval BIGINT DEFAULT 3600,
    max_notify_count BIGINT DEFAULT 3,
    send_recovery SMALLINT DEFAULT 1,
    recover_duration BIGINT DEFAULT 60,
    cron_pattern varchar(100) DEFAULT '@every 30s',
    append_tags TEXT,
    annotations TEXT,
    status SMALLINT DEFAULT 0,
    group_id BIGINT DEFAULT 0,
    group_name varchar(100) DEFAULT NULL,
    created_by varchar(100) DEFAULT NULL,
    updated_by varchar(100) DEFAULT NULL,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_name ON monitor_alert_rule (name);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_rule_target_type ON monitor_alert_rule (target_type);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_rule_status ON monitor_alert_rule (status);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_rule_group_id ON monitor_alert_rule (group_id);

CREATE TABLE IF NOT EXISTS monitor_alert_rule_template (
    id BIGSERIAL NOT NULL,
    name varchar(200) NOT NULL,
    description varchar(500) DEFAULT NULL,
    category varchar(50) NOT NULL,
    prom_ql TEXT NOT NULL,
    duration BIGINT DEFAULT 60,
    metric_name varchar(100) DEFAULT NULL,
    operator varchar(10) DEFAULT NULL,
    threshold decimal(10,2) DEFAULT NULL,
    unit varchar(20) DEFAULT NULL,
    severity SMALLINT DEFAULT 2,
    recover_duration BIGINT DEFAULT 60,
    tags TEXT,
    status SMALLINT DEFAULT 1,
    is_builtin SMALLINT DEFAULT 0,
    created_by varchar(100) DEFAULT NULL,
    updated_by varchar(100) DEFAULT NULL,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    metric_key varchar(100) NOT NULL,
    value_type varchar(20) DEFAULT 'number',
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_monitor_alert_rule_template_metric_key ON monitor_alert_rule_template (metric_key);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_rule_template_category ON monitor_alert_rule_template (category);

CREATE TABLE IF NOT EXISTS monitor_alert_source (
    id BIGSERIAL NOT NULL,
    name TEXT NOT NULL,
    type BIGINT NOT NULL,
    app_key TEXT NOT NULL,
    api_base_url TEXT,
    status BIGINT DEFAULT 1,
    remark TEXT,
    create_time BIGINT DEFAULT NULL,
    update_time BIGINT DEFAULT NULL,
    key_id BIGINT DEFAULT 0,
    host_id BIGINT DEFAULT 0,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS monitor_aliyun_config (
    id BIGSERIAL NOT NULL,
    name TEXT NOT NULL,
    access_key TEXT NOT NULL,
    access_secret TEXT NOT NULL,
    region varchar(191) DEFAULT 'cn-hangzhou',
    email TEXT,
    phone TEXT,
    username TEXT,
    status BIGINT DEFAULT 1,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    eab_kid TEXT,
    eab_hmac_key TEXT,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS monitor_api_business_function (
    id BIGSERIAL NOT NULL,
    name varchar(200) NOT NULL,
    description varchar(500) DEFAULT NULL,
    group_id BIGINT NOT NULL,
    group_name varchar(100) DEFAULT NULL,
    sort BIGINT DEFAULT 0,
    status SMALLINT DEFAULT 1,
    created_by varchar(100) DEFAULT NULL,
    updated_by varchar(100) DEFAULT NULL,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_monitor_api_business_function_group_id ON monitor_api_business_function (group_id);
CREATE INDEX IF NOT EXISTS idx_monitor_api_business_function_status ON monitor_api_business_function (status);

CREATE TABLE IF NOT EXISTS monitor_api_monitor (
    id BIGSERIAL NOT NULL,
    name varchar(200) NOT NULL,
    description varchar(500) DEFAULT NULL,
    group_id BIGINT DEFAULT NULL,
    group_name varchar(100) DEFAULT NULL,
    service_name varchar(100) DEFAULT NULL,
    url varchar(1000) NOT NULL,
    method varchar(10) DEFAULT 'GET',
    headers TEXT,
    body TEXT,
    timeout BIGINT DEFAULT 5000,
    expected_status BIGINT DEFAULT 200,
    expected_keyword varchar(500) DEFAULT NULL,
    validate_json TEXT,
    check_interval BIGINT DEFAULT 60,
    retry_count BIGINT DEFAULT 3,
    alert_enabled SMALLINT DEFAULT 1,
    policy_id BIGINT DEFAULT NULL,
    status SMALLINT DEFAULT 1,
    health_status varchar(20) DEFAULT 'unknown',
    last_check_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    last_response_time BIGINT DEFAULT NULL,
    last_status_code BIGINT DEFAULT NULL,
    last_error_message varchar(1000) DEFAULT NULL,
    continuous_failures BIGINT DEFAULT 0,
    tags TEXT,
    created_by varchar(100) DEFAULT NULL,
    updated_by varchar(100) DEFAULT NULL,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    business_function_id BIGINT DEFAULT NULL,
    business_function_name varchar(200) DEFAULT NULL,
    parent_id BIGINT DEFAULT NULL,
    hierarchy_level BIGINT DEFAULT 0,
    is_aggregate SMALLINT DEFAULT 0,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_monitor_api_monitor_group_id ON monitor_api_monitor (group_id);
CREATE INDEX IF NOT EXISTS idx_monitor_api_monitor_service_name ON monitor_api_monitor (service_name);
CREATE INDEX IF NOT EXISTS idx_monitor_api_monitor_status ON monitor_api_monitor (status);
CREATE INDEX IF NOT EXISTS idx_monitor_api_monitor_business_function_id ON monitor_api_monitor (business_function_id);
CREATE INDEX IF NOT EXISTS idx_monitor_api_monitor_parent_id ON monitor_api_monitor (parent_id);

CREATE TABLE IF NOT EXISTS monitor_api_monitor_result (
    id BIGSERIAL NOT NULL,
    monitor_id BIGINT NOT NULL,
    check_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    success SMALLINT NOT NULL,
    status_code BIGINT DEFAULT NULL,
    response_time BIGINT DEFAULT NULL,
    error_message varchar(1000) DEFAULT NULL,
    response_body TEXT,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_monitor_api_monitor_result_monitor_id ON monitor_api_monitor_result (monitor_id);
CREATE INDEX IF NOT EXISTS idx_monitor_api_monitor_result_check_time ON monitor_api_monitor_result (check_time);

CREATE TABLE IF NOT EXISTS monitor_domain (
    id BIGSERIAL NOT NULL,
    domain varchar(255) NOT NULL,
    tags varchar(500) DEFAULT NULL,
    remark TEXT,
    status BIGINT DEFAULT 1,
    is_alive BIGINT DEFAULT 0,
    status_code BIGINT DEFAULT NULL,
    response_time BIGINT DEFAULT NULL,
    ssl_expire_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    ssl_days_left BIGINT DEFAULT NULL,
    ssl_issuer varchar(255) DEFAULT NULL,
    last_check_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    error_msg TEXT,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_monitor_domain_domain ON monitor_domain (domain);

CREATE TABLE IF NOT EXISTS monitor_domain_schedule (
    id BIGSERIAL NOT NULL,
    enabled SMALLINT DEFAULT 0,
    cron_expr varchar(100) DEFAULT NULL,
    next_run_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    last_run_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    status varchar(50) DEFAULT NULL,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    notify_enabled SMALLINT DEFAULT 0,
    notify_robot_id BIGINT DEFAULT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS monitor_incident (
    id BIGSERIAL NOT NULL,
    alert_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    business_line TEXT,
    frequency TEXT,
    alert_desc TEXT,
    alert_level varchar(191) DEFAULT 'P4',
    incident_cause TEXT,
    department TEXT,
    solution TEXT,
    detail_url TEXT,
    handler TEXT,
    handler_id BIGINT DEFAULT NULL,
    status BIGINT DEFAULT 1,
    remark TEXT,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    business_line_id BIGINT DEFAULT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS monitor_notify_robot (
    id BIGSERIAL NOT NULL,
    name varchar(100) NOT NULL,
    type varchar(20) NOT NULL,
    webhook varchar(500) DEFAULT NULL,
    secret varchar(200) DEFAULT NULL,
    status SMALLINT DEFAULT 1,
    remark varchar(500) DEFAULT NULL,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    server varchar(100) DEFAULT NULL,
    port BIGINT DEFAULT NULL,
    username varchar(100) DEFAULT NULL,
    password varchar(200) DEFAULT NULL,
    nickname varchar(100) DEFAULT NULL,
    headers TEXT,
    method varchar(10) DEFAULT 'POST',
    template TEXT,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS monitor_ssl_cert (
    id BIGSERIAL NOT NULL,
    domain varchar(191) NOT NULL,
    aliyun_config_id BIGINT DEFAULT NULL,
    order_id varchar(191) DEFAULT NULL,
    cert_id BIGINT DEFAULT NULL,
    cert_name TEXT,
    product_code varchar(191) DEFAULT 'digicert-free-1-free',
    status BIGINT DEFAULT NULL,
    validate_type varchar(191) DEFAULT 'DNS',
    validate_info TEXT,
    cert TEXT,
    private_key TEXT,
    issue_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    expire_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    days_left BIGINT DEFAULT NULL,
    error_msg TEXT,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    cert_source varchar(191) DEFAULT 'aliyun_cas',
    ca_provider varchar(191) DEFAULT 'DigiCert',
    issuer_cert TEXT,
    algorithm varchar(191) DEFAULT 'RSA2048',
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_monitor_ssl_cert_domain ON monitor_ssl_cert (domain);
CREATE INDEX IF NOT EXISTS idx_monitor_ssl_cert_aliyun_config_id ON monitor_ssl_cert (aliyun_config_id);
CREATE INDEX IF NOT EXISTS idx_monitor_ssl_cert_order_id ON monitor_ssl_cert (order_id);
CREATE INDEX IF NOT EXISTS idx_monitor_ssl_cert_cert_source ON monitor_ssl_cert (cert_source);

CREATE TABLE IF NOT EXISTS monitor_ssl_cert_deploy_log (
    id BIGSERIAL NOT NULL,
    cert_id BIGINT NOT NULL,
    domain varchar(255) DEFAULT NULL,
    host_id BIGINT NOT NULL,
    host_name varchar(255) DEFAULT NULL,
    deploy_path varchar(500) DEFAULT NULL,
    status BIGINT DEFAULT 1,
    backup_files TEXT,
    deploy_files TEXT,
    logs TEXT,
    error_msg TEXT,
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_cert_id ON monitor_ssl_cert_deploy_log (cert_id);
CREATE INDEX IF NOT EXISTS idx_host_id ON monitor_ssl_cert_deploy_log (host_id);

CREATE TABLE IF NOT EXISTS monitor_webhook_log (
    id BIGSERIAL NOT NULL,
    source varchar(50) NOT NULL,
    title varchar(200) NOT NULL,
    content TEXT NOT NULL,
    level varchar(20) DEFAULT NULL,
    tags varchar(500) DEFAULT NULL,
    extra TEXT,
    notify_robot_ids varchar(200) DEFAULT NULL,
    status varchar(20) DEFAULT 'success',
    error_msg TEXT,
    notify_count BIGINT DEFAULT 0,
    success_count BIGINT DEFAULT 0,
    failed_count BIGINT DEFAULT 0,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS monitor_webhook_notify_log (
    id BIGSERIAL NOT NULL,
    webhook_log_id BIGINT NOT NULL,
    robot_id BIGINT NOT NULL,
    robot_name varchar(100) DEFAULT NULL,
    robot_type varchar(20) DEFAULT NULL,
    status varchar(20) NOT NULL,
    error_msg TEXT,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_monitor_webhook_notify_log_webhook_log_id ON monitor_webhook_notify_log (webhook_log_id);

CREATE TABLE IF NOT EXISTS notify_robot (
    id BIGSERIAL NOT NULL,
    name varchar(100) NOT NULL,
    type varchar(20) NOT NULL,
    webhook varchar(500) NOT NULL,
    secret varchar(200) DEFAULT NULL,
    status SMALLINT DEFAULT 1,
    remark varchar(500) DEFAULT NULL,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS prompt_template (
    id BIGSERIAL NOT NULL,
    name varchar(200) NOT NULL,
    category varchar(50) NOT NULL,
    scene varchar(50) NOT NULL,
    template TEXT NOT NULL,
    variables varchar(1000) DEFAULT NULL,
    system_prompt TEXT,
    temperature decimal(3,2) DEFAULT 0.70,
    max_tokens BIGINT DEFAULT 2000,
    model_id BIGINT DEFAULT NULL,
    enabled BIGINT DEFAULT 1,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_prompt_template_name ON prompt_template (name);
CREATE INDEX IF NOT EXISTS idx_prompt_template_category ON prompt_template (category);
CREATE INDEX IF NOT EXISTS idx_prompt_template_scene ON prompt_template (scene);

CREATE TABLE IF NOT EXISTS redis_instance (
    id BIGSERIAL NOT NULL,
    code varchar(36) NOT NULL,
    name varchar(100) NOT NULL,
    mode varchar(20) NOT NULL,
    host varchar(300) NOT NULL,
    port BIGINT DEFAULT 0,
    db BIGINT DEFAULT 0,
    username varchar(100) DEFAULT NULL,
    password varchar(500) NOT NULL,
    redis_node_password varchar(500) DEFAULT '',
    remark varchar(500) DEFAULT NULL,
    ssh_tunnel_machine_id BIGINT DEFAULT 0,
    status BIGINT DEFAULT 1,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    creator varchar(64) DEFAULT NULL,
    creator_id BIGINT DEFAULT NULL,
    modifier varchar(64) DEFAULT NULL,
    modifier_id BIGINT DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_redis_instance_code ON redis_instance (code);

CREATE TABLE IF NOT EXISTS sys_activity_log (
    id BIGSERIAL NOT NULL,
    activity_type BIGINT NOT NULL,
    title varchar(200) NOT NULL,
    content TEXT,
    status BIGINT NOT NULL DEFAULT 1,
    related_id BIGINT DEFAULT NULL,
    summary varchar(500) DEFAULT NULL,
    duration BIGINT DEFAULT NULL,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS sys_blocking_policy (
    id BIGSERIAL NOT NULL,
    name varchar(100) NOT NULL,
    description varchar(500) DEFAULT NULL,
    enabled SMALLINT DEFAULT 1,
    priority BIGINT DEFAULT 0,
    block_mode BIGINT DEFAULT 1,
    enable_alias_resolution SMALLINT DEFAULT 1,
    enable_script_scanning SMALLINT DEFAULT 1,
    custom_rules TEXT,
    whitelist_cmds TEXT,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS sys_command_audit (
    id BIGSERIAL NOT NULL,
    recording_id BIGINT NOT NULL,
    session_id varchar(64) NOT NULL,
    command TEXT NOT NULL,
    timestamp DOUBLE PRECISION NOT NULL,
    sequence BIGINT NOT NULL,
    is_sensitive SMALLINT DEFAULT 0,
    risk_level BIGINT DEFAULT 0,
    risk_reason varchar(512) DEFAULT NULL,
    execute_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_sys_command_audit_recording_id ON sys_command_audit (recording_id);
CREATE INDEX IF NOT EXISTS idx_sys_command_audit_session_id ON sys_command_audit (session_id);
CREATE INDEX IF NOT EXISTS idx_sys_command_audit_is_sensitive ON sys_command_audit (is_sensitive);
CREATE INDEX IF NOT EXISTS idx_sys_command_audit_risk_level ON sys_command_audit (risk_level);
CREATE INDEX IF NOT EXISTS idx_sys_command_audit_execute_time ON sys_command_audit (execute_time);

CREATE TABLE IF NOT EXISTS sys_command_blocking (
    id BIGSERIAL NOT NULL,
    session_id varchar(64) DEFAULT NULL,
    recording_id BIGINT DEFAULT NULL,
    original_cmd TEXT,
    resolved_cmd TEXT,
    blocked SMALLINT DEFAULT NULL,
    block_reason TEXT,
    risk_level BIGINT DEFAULT NULL,
    is_alias SMALLINT DEFAULT 0,
    is_script SMALLINT DEFAULT 0,
    detected_issues TEXT,
    policy_id BIGINT DEFAULT NULL,
    policy_name varchar(100) DEFAULT NULL,
    user_id BIGINT DEFAULT NULL,
    username varchar(50) DEFAULT NULL,
    host_id BIGINT DEFAULT NULL,
    host_ip varchar(50) DEFAULT NULL,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_sys_command_blocking_session_id ON sys_command_blocking (session_id);
CREATE INDEX IF NOT EXISTS idx_sys_command_blocking_recording_id ON sys_command_blocking (recording_id);
CREATE INDEX IF NOT EXISTS idx_sys_command_blocking_blocked ON sys_command_blocking (blocked);
CREATE INDEX IF NOT EXISTS idx_sys_command_blocking_risk_level ON sys_command_blocking (risk_level);
CREATE INDEX IF NOT EXISTS idx_sys_command_blocking_user_id ON sys_command_blocking (user_id);
CREATE INDEX IF NOT EXISTS idx_sys_command_blocking_host_id ON sys_command_blocking (host_id);
CREATE INDEX IF NOT EXISTS idx_sys_command_blocking_created_at ON sys_command_blocking (created_at);

CREATE TABLE IF NOT EXISTS sys_config (
    id BIGSERIAL NOT NULL,
    config_key varchar(100) NOT NULL,
    config_type varchar(50) NOT NULL,
    config_data TEXT NOT NULL,
    status BIGINT NOT NULL DEFAULT 1,
    remark varchar(500) DEFAULT NULL,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_config_config_key ON sys_config (config_key);

CREATE TABLE IF NOT EXISTS sys_session_recording (
    id BIGSERIAL NOT NULL,
    session_id varchar(64) NOT NULL,
    admin_id BIGINT NOT NULL,
    username varchar(64) NOT NULL,
    host_id BIGINT NOT NULL,
    host_name varchar(128) NOT NULL,
    host_ip varchar(64) NOT NULL,
    ssh_user varchar(64) NOT NULL,
    start_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    end_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    duration BIGINT DEFAULT NULL,
    terminal_width BIGINT DEFAULT 80,
    terminal_height BIGINT DEFAULT 24,
    file_path varchar(512) NOT NULL,
    file_size BIGINT DEFAULT NULL,
    storage_type BIGINT DEFAULT 1,
    oss_key varchar(512) DEFAULT NULL,
    input_count BIGINT DEFAULT 0,
    output_count BIGINT DEFAULT 0,
    resize_count BIGINT DEFAULT 0,
    command_count BIGINT DEFAULT 0,
    client_ip varchar(64) DEFAULT NULL,
    user_agent varchar(512) DEFAULT NULL,
    risk_level BIGINT DEFAULT 0,
    has_sensitive_cmd SMALLINT DEFAULT 0,
    status BIGINT DEFAULT 1,
    error_msg varchar(512) DEFAULT NULL,
    create_time TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL,
    update_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    delete_time TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_session_recording_session_id ON sys_session_recording (session_id);
CREATE INDEX IF NOT EXISTS idx_sys_session_recording_host_id ON sys_session_recording (host_id);
CREATE INDEX IF NOT EXISTS idx_sys_session_recording_host_ip ON sys_session_recording (host_ip);
CREATE INDEX IF NOT EXISTS idx_sys_session_recording_start_time ON sys_session_recording (start_time);
CREATE INDEX IF NOT EXISTS idx_sys_session_recording_risk_level ON sys_session_recording (risk_level);
CREATE INDEX IF NOT EXISTS idx_sys_session_recording_status ON sys_session_recording (status);

CREATE TABLE IF NOT EXISTS system_migration_meta (
    meta_key varchar(128) NOT NULL,
    meta_value varchar(128) NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (meta_key)
);

CREATE TABLE IF NOT EXISTS tool_feishu_user_token (
    id BIGSERIAL NOT NULL,
    key_manage_id BIGINT NOT NULL,
    open_id varchar(100) DEFAULT NULL,
    union_id varchar(100) DEFAULT NULL,
    user_name varchar(100) DEFAULT NULL,
    avatar_url varchar(500) DEFAULT NULL,
    access_token TEXT,
    refresh_token TEXT,
    expires_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    refresh_expires_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_feishu_user ON tool_feishu_user_token (open_id);
CREATE INDEX IF NOT EXISTS idx_tool_feishu_user_token_key_manage_id ON tool_feishu_user_token (key_manage_id);
