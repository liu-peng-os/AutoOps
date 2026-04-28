# 阶段 1 PostgreSQL Schema 覆盖清单

## 当前结论

截至 2026-04-27，旧 MySQL 初始化 SQL 中真正的建表语句共 79 张表，当前 PostgreSQL goose migration 已覆盖 79 张表。

当前开发库验证结果：

```text
goose_version=7
public_tables=80
```

其中 `public_tables=80` 包含 79 张业务表和 1 张 `goose_db_version`。

## 口径说明

本清单只确认 schema 覆盖，不代表旧业务数据已经全量导入 PostgreSQL。

已完成：

- 表结构迁移覆盖旧 MySQL 初始化 SQL 中的 79 张业务表。
- 系统基础 seed 已迁移，包括岗位、部门、角色、用户、菜单、角色菜单、用户角色绑定。
- `admin / 123456` 已可登录，并能返回完整菜单与权限点。

未包含：

- 历史业务数据全量导入。
- 旧 MySQL dump 中运行日志、测试数据、历史监控结果等数据迁移。
- 外键约束的逐项恢复。阶段 1 先保证 schema 可用和代码不因缺表失败，约束治理放到后续按模块收敛。

## Migration 分布

| Migration | 作用 | 表数量 |
| --- | --- | --- |
| `0001_stage1_baseline.sql` | 配置中心、CMDB 主机、监控 agent、工具部署等底座表 | 9 |
| `0002_system_core.sql` | 系统权限核心表 | 8 |
| `0003_platform_foundation.sql` | 操作日志、SQL 审计、K8s 集群等平台基础表 | 4 |
| `0004_delivery_flow.sql` | 任务、发布、快速部署等交付流程表 | 9 |
| `0005_system_seed.sql` | 最小可登录系统 seed | 0 |
| `0006_system_full_seed.sql` | 完整系统基础 seed | 0 |
| `0007_legacy_schema_catchup.sql` | 旧系统剩余 schema 补齐 | 49 |

## 已覆盖表

### 0001

- `config_account`
- `config_ecsauth`
- `config_keymanage`
- `config_sync_schedule`
- `cmdb_group`
- `cmdb_host`
- `monitor_agent`
- `tool_link`
- `tool_service_deploy`

### 0002

- `sys_post`
- `sys_dept`
- `sys_role`
- `sys_menu`
- `sys_admin`
- `sys_admin_role`
- `sys_role_menu`
- `sys_login_info`

### 0003

- `sys_operation_log`
- `cmdb_sql`
- `cmdb_sql_log`
- `k8s_cluster`

### 0004

- `task_template`
- `task_job`
- `task_work`
- `task_ansible`
- `task_ansiblework`
- `app_application`
- `app_jenkins_env`
- `quick_deployments`
- `quick_deployment_tasks`

### 0007

- `ai_agent_chat_history`
- `ai_agent_task`
- `ai_model`
- `app_service_release`
- `app_service_release_item`
- `app_sh_release`
- `bastion_host_authorization`
- `bastion_user_group`
- `cmdb_snmp_devices`
- `cmdb_sql_records`
- `db`
- `db_es_instance`
- `db_export_task`
- `db_instance`
- `db_mongo_instance`
- `db_redis_instance`
- `db_sql`
- `db_sql_exec`
- `knowledge_base`
- `monitor_alert_mute`
- `monitor_alert_policy`
- `monitor_alert_policy_rule`
- `monitor_alert_record`
- `monitor_alert_rule`
- `monitor_alert_rule_template`
- `monitor_alert_source`
- `monitor_aliyun_config`
- `monitor_api_business_function`
- `monitor_api_monitor`
- `monitor_api_monitor_result`
- `monitor_domain`
- `monitor_domain_schedule`
- `monitor_incident`
- `monitor_notify_robot`
- `monitor_ssl_cert`
- `monitor_ssl_cert_deploy_log`
- `monitor_webhook_log`
- `monitor_webhook_notify_log`
- `notify_robot`
- `prompt_template`
- `redis_instance`
- `sys_activity_log`
- `sys_blocking_policy`
- `sys_command_audit`
- `sys_command_blocking`
- `sys_config`
- `sys_session_recording`
- `system_migration_meta`
- `tool_feishu_user_token`

## 验证命令

用于比对旧 MySQL 建表数量与当前 migration 覆盖数量：

```bash
python scripts/check_schema_coverage.py
```

脚本口径：

- 从 `docker/mysql/devops001.sql` 中按行首 `CREATE TABLE` 统计旧表。
- 从 `api/migrations/*.sql` 中按 `CREATE TABLE IF NOT EXISTS` 统计 PostgreSQL 表。
- 对比结果应为 `missing=0`、`extra=0`。

## 后续注意

- 如果后续新增模块表，必须通过新的 goose migration 增量添加。
- 不要重新改已经执行过的 migration；如果要修正表结构，新增 migration。
- 旧业务数据导入应单独设计，不要混入 schema migration。
