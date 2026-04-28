-- +goose Up

INSERT INTO sys_post (id, post_code, post_name, post_status, create_time, remark)
VALUES (1, 'admin', '管理员', 1, CURRENT_TIMESTAMP, '阶段 1 默认管理员岗位')
ON CONFLICT (post_code) DO UPDATE SET
    post_name = EXCLUDED.post_name,
    post_status = EXCLUDED.post_status,
    remark = EXCLUDED.remark;

INSERT INTO sys_dept (id, parent_id, dept_type, dept_name, dept_status, create_time)
SELECT 1, 0, 1, '默认部门', 1, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM sys_dept WHERE id = 1);

INSERT INTO sys_role (id, role_name, role_key, status, description, create_time)
VALUES (1, '超级管理员', 'admin', 1, '最大权限', CURRENT_TIMESTAMP)
ON CONFLICT (role_key) DO UPDATE SET
    role_name = EXCLUDED.role_name,
    status = EXCLUDED.status,
    description = EXCLUDED.description;

INSERT INTO sys_admin (
    post_id,
    dept_id,
    username,
    password,
    nickname,
    status,
    icon,
    email,
    phone,
    note,
    create_time
)
VALUES (
    1,
    1,
    'admin',
    'e10adc3949ba59abbe56e057f20f883e',
    '管理员',
    1,
    '',
    '',
    '',
    '阶段 1 默认管理员账号',
    CURRENT_TIMESTAMP
)
ON CONFLICT (username) DO UPDATE SET
    password = EXCLUDED.password,
    nickname = EXCLUDED.nickname,
    status = EXCLUDED.status,
    post_id = EXCLUDED.post_id,
    dept_id = EXCLUDED.dept_id,
    note = EXCLUDED.note;

INSERT INTO sys_menu (parent_id, menu_name, icon, value, menu_type, url, menu_status, sort, create_time)
SELECT 0, '仪表盘', 'House', 'dashboard', 1, 'dashboard', 2, 1, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM sys_menu WHERE menu_name = '仪表盘');

INSERT INTO sys_menu (parent_id, menu_name, icon, value, menu_type, url, menu_status, sort, create_time)
SELECT 0, '系统管理', 'Setting', 'system', 1, 'system', 2, 90, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM sys_menu WHERE menu_name = '系统管理');

INSERT INTO sys_menu (parent_id, menu_name, icon, value, menu_type, url, menu_status, sort, create_time)
SELECT parent.id, '用户信息', '', 'base:admin:list', 2, 'system/admin', 2, 1, CURRENT_TIMESTAMP
FROM sys_menu parent
WHERE parent.menu_name = '系统管理'
  AND NOT EXISTS (SELECT 1 FROM sys_menu WHERE menu_name = '用户信息');

INSERT INTO sys_menu (parent_id, menu_name, icon, value, menu_type, url, menu_status, sort, create_time)
SELECT parent.id, '角色信息', '', 'base:role:list', 2, 'system/role', 2, 2, CURRENT_TIMESTAMP
FROM sys_menu parent
WHERE parent.menu_name = '系统管理'
  AND NOT EXISTS (SELECT 1 FROM sys_menu WHERE menu_name = '角色信息');

INSERT INTO sys_admin_role (role_id, admin_id)
SELECT r.id, a.id
FROM sys_role r
CROSS JOIN sys_admin a
WHERE r.role_key = 'admin'
  AND a.username = 'admin'
ON CONFLICT DO NOTHING;

INSERT INTO sys_role_menu (role_id, menu_id)
SELECT r.id, m.id
FROM sys_role r
CROSS JOIN sys_menu m
WHERE r.role_key = 'admin'
ON CONFLICT DO NOTHING;
