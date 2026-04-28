# 阶段 1 PostgreSQL 与 Migration 底座

## 当前基线

- 后端可以通过配置切换数据库方言，阶段 1 默认使用 PostgreSQL。
- 启动时通过 `goose` 执行 `api/migrations/*.sql`。
- 配置加载支持 `config.yaml + 环境变量覆盖`，本地启动和 Docker 启动使用同一套模型。
- Redis 被视为可降级依赖。Redis 不可用时，应用仍可启动，但缓存和任务队列能力会受限。

## 关键环境变量

- `DB_DIALECTS`
- `DB_HOST`
- `DB_PORT`
- `DB_NAME`
- `DB_USER`
- `DB_PASSWORD`
- `DB_SSLMODE`
- `DB_MIGRATION_PATH`
- `DB_AUTO_MIGRATE`
- `REDIS_ADDR`
- `REDIS_PASSWORD`
- `SERVER_ADDRESS`
- `SERVER_PUBLIC_URL`
- `IMAGE_HOST`

## 本地启动

```bash
cd api
go run scripts/migrate.go
go run main.go
```

本地启动时建议用环境变量覆盖数据库连接信息，而不是直接修改 `config.yaml`。

对于外部 PostgreSQL 实例，优先显式覆盖：

```text
DB_HOST
DB_PORT
DB_NAME
DB_USER
DB_PASSWORD
DB_SSLMODE
```

## Docker 启动

```bash
cd docker
docker compose up -d
```

当前 Docker 栈假设 PostgreSQL 由外部实例提供，并通过 `.env` 注入运行时配置。

默认 `.env` 指向 `docs/autoops-analysis/05-independent-evolution-roadmap.md` 中已验证的 PostgreSQL 实例。只有在明确切换目标数据库时，才修改 `docker/.env` 中的 `DB_HOST`、`DB_PORT`、`DB_NAME`、`DB_USER`、`DB_PASSWORD`。

## 当前运行模型

- `docker compose` 不启动 PostgreSQL。
- 本地 Docker 只负责 Redis、Pushgateway、Prometheus、API、Web。
- API 容器通过 `.env` 连接已验证的外部 PostgreSQL。
- Docker 挂载根项目中的 `api/config.yaml` 和 `api/migrations`，作为后端运行配置和 schema 初始化的单一来源。
- API 和 Web 当前在本地 Docker 启动时从仓库源码构建，不依赖私有仓库中的预构建镜像。

## Migration 状态

- `goose` 是当前正式 migration 工具。
- `goose` 使用 `goose_db_version` 管理版本状态。
- 旧的 `schema_migrations` 只作为兼容来源读取一次，用于从临时自研 runner 升级到 `goose`。
- `AutoMigrate` 已从核心启动路径退场；当前已登记的 GORM 模型均已由 SQL migration 覆盖。

当前 SQL 管理表包括：

```text
config_account
config_ecsauth
config_keymanage
config_sync_schedule
cmdb_group
cmdb_host
cmdb_sql
cmdb_sql_log
monitor_agent
tool_link
tool_service_deploy
sys_admin
sys_role
sys_menu
sys_dept
sys_post
sys_admin_role
sys_role_menu
sys_login_info
sys_operation_log
k8s_cluster
task_template
task_job
task_work
task_ansible
task_ansiblework
app_application
app_jenkins_env
quick_deployments
quick_deployment_tasks
```

## 已完成验证

- 2026-04-24，使用 `docs/autoops-analysis/05-independent-evolution-roadmap.md` 中的 PostgreSQL 实例完成真实连接验证。
- 低层探测确认该端点是 PostgreSQL 监听器，拒绝 TLS，支持明文 SCRAM 认证。
- 应用 migration 入口执行成功，`goose` 当前版本为 `4`。
- 后端主进程冒烟通过，API 能连接外部 PostgreSQL 并保持运行。
- 本地 Docker 栈在释放 Docker Desktop 磁盘空间后启动成功。
- Redis、Pushgateway、Prometheus、API、Web 联调通过。
- Web、API integration catalog、Prometheus、Pushgateway 均返回 HTTP 200。

## 已知注意事项

- Docker Desktop 需要预留足够磁盘空间。API/Web 源码构建会产生数 GB 镜像层和构建缓存。
- API Docker 构建已降低 Go 编译并发，减少临时磁盘压力。
- 如果 Docker 出现 `input/output error` 或 `no space left on device`，应优先检查 `docker system df` 和宿主磁盘空间。
- 阶段 2 前需要决定本地 Docker 是否继续源码构建，还是切回运行时镜像模式以节省磁盘。

## 与 2026-04-23 记录的关系

2026-04-23 的记录中，PostgreSQL 接入和 migration 实链验证仍是待办。

2026-04-24 已完成这些待办：

- PostgreSQL 真实连接验证。
- 正式 SQL migration 落地。
- `goose` 替换临时 runner。
- `AutoMigrate` 退出核心路径。
- 本地 Docker 完整联调。
