-- +goose Up

CREATE TABLE IF NOT EXISTS sys_post (
    id BIGSERIAL PRIMARY KEY,
    post_code VARCHAR(64) NOT NULL,
    post_name VARCHAR(50) NOT NULL,
    post_status BIGINT NOT NULL DEFAULT 1,
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    remark VARCHAR(500)
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_post_post_code ON sys_post (post_code);
CREATE INDEX IF NOT EXISTS idx_sys_post_status ON sys_post (post_status);

CREATE TABLE IF NOT EXISTS sys_dept (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT NOT NULL,
    dept_type BIGINT NOT NULL,
    dept_name VARCHAR(30) NOT NULL,
    dept_status BIGINT DEFAULT 1,
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sys_dept_parent_id ON sys_dept (parent_id);
CREATE INDEX IF NOT EXISTS idx_sys_dept_status ON sys_dept (dept_status);

CREATE TABLE IF NOT EXISTS sys_role (
    id BIGSERIAL PRIMARY KEY,
    role_name VARCHAR(64) NOT NULL,
    role_key VARCHAR(64) NOT NULL,
    status BIGINT NOT NULL DEFAULT 1,
    description VARCHAR(500),
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_role_role_key ON sys_role (role_key);
CREATE INDEX IF NOT EXISTS idx_sys_role_status ON sys_role (status);

CREATE TABLE IF NOT EXISTS sys_menu (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT DEFAULT 0,
    menu_name VARCHAR(100),
    icon VARCHAR(100),
    value VARCHAR(100),
    menu_type BIGINT,
    url VARCHAR(100),
    menu_status BIGINT,
    sort BIGINT,
    create_time TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sys_menu_parent_id ON sys_menu (parent_id);
CREATE INDEX IF NOT EXISTS idx_sys_menu_menu_status ON sys_menu (menu_status);
CREATE INDEX IF NOT EXISTS idx_sys_menu_sort ON sys_menu (sort);

CREATE TABLE IF NOT EXISTS sys_admin (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT,
    dept_id BIGINT,
    username VARCHAR(64) NOT NULL,
    password VARCHAR(64) NOT NULL,
    nickname VARCHAR(64),
    status BIGINT NOT NULL DEFAULT 1,
    icon VARCHAR(500),
    email VARCHAR(64),
    phone VARCHAR(64),
    note VARCHAR(500),
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_admin_username ON sys_admin (username);
CREATE INDEX IF NOT EXISTS idx_sys_admin_status ON sys_admin (status);
CREATE INDEX IF NOT EXISTS idx_sys_admin_dept_id ON sys_admin (dept_id);
CREATE INDEX IF NOT EXISTS idx_sys_admin_post_id ON sys_admin (post_id);

CREATE TABLE IF NOT EXISTS sys_admin_role (
    role_id BIGINT NOT NULL,
    admin_id BIGINT NOT NULL,
    PRIMARY KEY (role_id, admin_id)
);

CREATE INDEX IF NOT EXISTS idx_sys_admin_role_admin_id ON sys_admin_role (admin_id);

CREATE TABLE IF NOT EXISTS sys_role_menu (
    role_id BIGINT NOT NULL,
    menu_id BIGINT NOT NULL,
    PRIMARY KEY (role_id, menu_id)
);

CREATE INDEX IF NOT EXISTS idx_sys_role_menu_menu_id ON sys_role_menu (menu_id);

CREATE TABLE IF NOT EXISTS sys_login_info (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50),
    ip_address VARCHAR(128),
    login_location VARCHAR(255),
    browser VARCHAR(50),
    os VARCHAR(50),
    login_status BIGINT,
    message VARCHAR(255),
    login_time TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_sys_login_info_username ON sys_login_info (username);
CREATE INDEX IF NOT EXISTS idx_sys_login_info_status ON sys_login_info (login_status);
CREATE INDEX IF NOT EXISTS idx_sys_login_info_login_time ON sys_login_info (login_time);
