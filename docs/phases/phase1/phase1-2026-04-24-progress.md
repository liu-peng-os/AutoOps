# 阶段 1 每日进展 - 2026-04-24

## 已完成

- 修复仓库级编译阻塞：`api/common/agent/agent.go` 中的格式化参数问题。
- Docker 运行模型调整为：PostgreSQL 使用外部实例，本地 Docker 只负责 Redis、Pushgateway、Prometheus、API、Web。
- 更新 `docker/docker-compose.yml`，使 API 默认连接外部 PostgreSQL。
- 更新 `docker/.env`，默认指向 `05-independent-evolution-roadmap.md` 中已经验证的 PostgreSQL 实例。
- 更新 `docker/devops-start.sh`，去掉本地 PostgreSQL 假设，并增加 Docker Engine 健康预检。
- 将阶段文档统一放到根目录 `docs/phase1/`，不再放入 `api/docs/`。
- 将占位 baseline migration 替换为第一批正式 PostgreSQL DDL。
- 将 `AutoMigrate` 从核心表管理路径中退出来。
- 新增测试，防止已由 SQL 管理的表重新漂回 `AutoMigrate`。
- 新增系统基础表 migration：管理员、角色、菜单、部门、岗位、角色关联、菜单关联、登录记录。
- 新增平台基础表 migration：操作日志、CMDB SQL 实例、CMDB SQL 审计、Kubernetes 集群元数据。
- 新增交付流表 migration：任务模板、任务、任务执行、Ansible 任务、应用、Jenkins 环境、快速发布记录。
- 将 migration 引擎从临时自研 runner 切换为 `goose`。
- 保留旧 `schema_migrations` 到 `goose_db_version` 的兼容导入路径。

## 已验证

- 执行 `go test ./common/agent ./common/config ./pkg/db ./pkg/redis ./scripts`。
- 确认历史 `common/agent` 编译失败已修复。
- 验证 SQL 管理表清单与 migration 覆盖范围保持一致。
- 清空当前已登记 GORM 模型的 `AutoMigrate` 待迁移清单。
- 使用文档中的外部 PostgreSQL 连接信息执行 `go run scripts/migrate.go`，确认 `goose` 当前版本为 `4`。
- 使用同一 PostgreSQL 配置执行后端主进程冒烟，确认 API 能保持运行。
- 验证 Redis 不在宿主机运行时，API 能按预期降级启动。
- 验证 `docker compose config` 能解析到正确的外部 PostgreSQL 主机、数据库和账号。
- 更新 `docker/api/Dockerfile`，使用 Go 1.25、阿里云 Alpine 源、`goproxy.cn` 和 `sum.golang.google.cn`。
- 确认 `docker-devops-web` 镜像能完成前端生产构建。
- 释放 Docker Desktop 磁盘空间后，API 镜像构建成功。
- 启动完整本地 compose 栈：Redis、Pushgateway、Prometheus、API、Web。
- 最终 compose 状态验证通过：Redis、Pushgateway、Prometheus、Web 为 healthy，API 为 running。
- HTTP 冒烟检查通过：
  - `http://localhost:8088`
  - `http://localhost:8000/api/v1/integration/catalog`
  - `http://localhost:9090/-/healthy`
  - `http://localhost:9091/-/healthy`

## 当前状态

- 阶段 1 已在代码、数据库、本地 Docker 环境三个层面完成底座验证。
- PostgreSQL 使用外部实例，并已完成真实连接和 migration 验证。
- 本地 Docker 提供 Redis、Pushgateway、Prometheus、API、Web。
- Docker Desktop 磁盘空间是已知运行约束。构建 API/Web 镜像前需要预留足够空间，因为 Go 和前端镜像构建会消耗数 GB 镜像层和构建缓存。

## 下一步建议

- 阶段 1 可以视为基本完成。
- 阶段 2 开始前，需要决定本地 Docker 开发是继续源码构建，还是回到更省磁盘的运行时镜像模式。
- 阶段 2 应优先选择一个接入型能力切片，验证 adapter、同步、回写、审计和前端基础交互。
