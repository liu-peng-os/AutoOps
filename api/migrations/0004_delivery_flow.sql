-- +goose Up

CREATE TABLE IF NOT EXISTS task_template (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type BIGINT NOT NULL,
    content TEXT NOT NULL,
    remark VARCHAR(500),
    created_by VARCHAR(50),
    updated_by VARCHAR(50),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_task_template_type ON task_template (type);

CREATE TABLE IF NOT EXISTS task_job (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255),
    type BIGINT,
    shell TEXT,
    host_ids TEXT,
    cron_expr VARCHAR(255),
    tasklog TEXT,
    status BIGINT,
    duration BIGINT,
    task_count BIGINT,
    execute_count BIGINT DEFAULT 0,
    next_run_time TIMESTAMP WITHOUT TIME ZONE,
    remark TEXT,
    start_time TIMESTAMP WITHOUT TIME ZONE,
    end_time TIMESTAMP WITHOUT TIME ZONE,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_task_job_status ON task_job (status);
CREATE INDEX IF NOT EXISTS idx_task_job_type ON task_job (type);

CREATE TABLE IF NOT EXISTS task_work (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT,
    template_id BIGINT,
    host_id BIGINT,
    type BIGINT,
    status BIGINT,
    log TEXT,
    log_path TEXT,
    start_time TIMESTAMP WITHOUT TIME ZONE,
    end_time TIMESTAMP WITHOUT TIME ZONE,
    duration BIGINT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    scheduled_time TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_task_work_task_id ON task_work (task_id);
CREATE INDEX IF NOT EXISTS idx_task_work_template_id ON task_work (template_id);
CREATE INDEX IF NOT EXISTS idx_task_work_host_id ON task_work (host_id);
CREATE INDEX IF NOT EXISTS idx_task_work_status ON task_work (status);

CREATE TABLE IF NOT EXISTS task_ansible (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type BIGINT NOT NULL DEFAULT 1,
    git_repo VARCHAR(255),
    host_groups TEXT NOT NULL,
    all_host_ids TEXT NOT NULL,
    global_vars TEXT,
    status BIGINT NOT NULL DEFAULT 1,
    error_msg TEXT,
    task_count BIGINT NOT NULL DEFAULT 0,
    total_duration BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_task_ansible_name ON task_ansible (name);
CREATE INDEX IF NOT EXISTS idx_task_ansible_status ON task_ansible (status);

CREATE TABLE IF NOT EXISTS task_ansiblework (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL,
    entry_file_name VARCHAR(255) NOT NULL,
    entry_file_path VARCHAR(255) NOT NULL,
    log_path VARCHAR(255),
    status BIGINT NOT NULL DEFAULT 1,
    start_time TIMESTAMP WITHOUT TIME ZONE,
    end_time TIMESTAMP WITHOUT TIME ZONE,
    duration BIGINT,
    exit_code BIGINT,
    error_msg TEXT
);

CREATE INDEX IF NOT EXISTS idx_task_ansiblework_task_id ON task_ansiblework (task_id);
CREATE INDEX IF NOT EXISTS idx_task_ansiblework_status ON task_ansiblework (status);

CREATE TABLE IF NOT EXISTS app_application (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL,
    business_group_id BIGINT NOT NULL,
    business_dept_id BIGINT NOT NULL,
    description TEXT,
    repo_url VARCHAR(500),
    dev_owners JSONB,
    test_owners JSONB,
    ops_owners JSONB,
    programming_lang VARCHAR(100),
    start_command TEXT,
    stop_command TEXT,
    health_api VARCHAR(500),
    domains JSONB,
    hosts JSONB,
    databases JSONB,
    other_res JSONB,
    status BIGINT DEFAULT 1,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_app_application_code ON app_application (code);
CREATE INDEX IF NOT EXISTS idx_app_application_group_id ON app_application (business_group_id);
CREATE INDEX IF NOT EXISTS idx_app_application_dept_id ON app_application (business_dept_id);
CREATE INDEX IF NOT EXISTS idx_app_application_status ON app_application (status);

CREATE TABLE IF NOT EXISTS app_jenkins_env (
    id BIGSERIAL PRIMARY KEY,
    app_id BIGINT NOT NULL,
    env_name VARCHAR(50) NOT NULL,
    jenkins_server_id BIGINT,
    job_name VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_app_jenkins_env_app_id ON app_jenkins_env (app_id);
CREATE INDEX IF NOT EXISTS idx_app_jenkins_env_server_id ON app_jenkins_env (jenkins_server_id);

CREATE TABLE IF NOT EXISTS quick_deployments (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    business_group_id BIGINT NOT NULL,
    business_dept_id BIGINT NOT NULL,
    description TEXT,
    status BIGINT DEFAULT 1,
    task_count BIGINT NOT NULL DEFAULT 0,
    execution_mode BIGINT DEFAULT 1,
    creator_id BIGINT NOT NULL,
    creator_name VARCHAR(100),
    start_time TIMESTAMP WITHOUT TIME ZONE,
    end_time TIMESTAMP WITHOUT TIME ZONE,
    duration BIGINT,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_quick_deployments_group_id ON quick_deployments (business_group_id);
CREATE INDEX IF NOT EXISTS idx_quick_deployments_dept_id ON quick_deployments (business_dept_id);
CREATE INDEX IF NOT EXISTS idx_quick_deployments_status ON quick_deployments (status);
CREATE INDEX IF NOT EXISTS idx_quick_deployments_creator_id ON quick_deployments (creator_id);

CREATE TABLE IF NOT EXISTS quick_deployment_tasks (
    id BIGSERIAL PRIMARY KEY,
    deployment_id BIGINT NOT NULL,
    app_id BIGINT NOT NULL,
    app_name VARCHAR(255),
    app_code VARCHAR(255),
    environment VARCHAR(50),
    jenkins_env_id BIGINT NOT NULL,
    jenkins_job_url VARCHAR(500),
    build_number BIGINT,
    status BIGINT DEFAULT 1,
    execute_order BIGINT NOT NULL,
    start_time TIMESTAMP WITHOUT TIME ZONE,
    end_time TIMESTAMP WITHOUT TIME ZONE,
    duration BIGINT,
    error_message TEXT,
    log_url VARCHAR(500),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_quick_deployment_tasks_deployment_id ON quick_deployment_tasks (deployment_id);
CREATE INDEX IF NOT EXISTS idx_quick_deployment_tasks_app_id ON quick_deployment_tasks (app_id);
CREATE INDEX IF NOT EXISTS idx_quick_deployment_tasks_jenkins_env_id ON quick_deployment_tasks (jenkins_env_id);
CREATE INDEX IF NOT EXISTS idx_quick_deployment_tasks_status ON quick_deployment_tasks (status);
