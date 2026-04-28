# 阶段 1 每日进展 - 2026-04-23

## 已完成

- 后端数据库初始化流程新增 PostgreSQL 驱动支持和数据库方言切换能力。
- 阶段 1 默认数据库配置从 MySQL 思路切换到 PostgreSQL 思路。
- 数据库配置新增 `sslMode`、`migrationPath`、`autoMigrate` 等字段。
- 新增轻量 SQL migration 机制，使用 `schema_migrations` 记录已执行迁移。
- 新增第一份 baseline migration 文件：`api/migrations/0001_stage1_baseline.sql`。
- 统一 migration 执行入口：`api/scripts/migrate.go` 复用后端 DB 初始化流程。
- 配置加载支持环境变量覆盖，使 Docker 和本地启动能够使用同一套配置模型。
- 补充配置环境变量覆盖相关测试。
- Docker Compose、`.env` 和启动脚本开始对齐 PostgreSQL 16 和新的配置模型。
- Redis 初始化失败时改为清晰降级，避免保留半初始化客户端引用。

## 已验证

- 执行 `go test ./common/config ./pkg/db ./pkg/redis ./scripts`。
- `common/config` 相关测试通过，包括环境变量覆盖行为。

## 当日仓库状态

- 阶段 1 底座从 MySQL-first 路径转向 PostgreSQL-first 路径。
- SQL migration 和 `AutoMigrate` 当时仍然共存，`AutoMigrate` 仍作为临时兜底。
- Docker 配置结构已经调整，但当日尚未完成完整容器实启验证。

## 结转事项

这些事项已在 2026-04-24 继续推进并基本关闭：

- PostgreSQL 16 真实连接和启动链路验证。
- 将占位 baseline migration 替换为正式 schema migration。
- 在正式 migration 覆盖后移除 `AutoMigrate` 关键路径。
- 修复历史编译阻塞问题。
- 处理环境磁盘空间对测试和 Docker 构建的影响。
