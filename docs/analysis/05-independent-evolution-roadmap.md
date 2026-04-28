# AutoOps 独立演进路线图

> 版本：2.0
> 更新时间：2026-04-24
> 状态来源：已对齐 `docs/2026-04-23.txt` 与 `docs/phase1/`。

## 1. 项目定位

AutoOps 后续按独立运维平台演进，不再以低成本持续兼容 upstream 为首要目标。上游项目仍可作为参考，但本项目优先追求长期可维护、模块边界清晰、数据库演进可控、外部系统可替换。

核心方向：

- 保留原项目中已经验证过的产品理念和功能方向。
- 允许在架构、数据模型、工程体系上做必要改造。
- 所有外部系统都视为可替换依赖，不绑定到核心模型里。
- 优先建设通用底座，再逐步做深业务能力。

## 2. 已确认决策

| 领域 | 决策 | 当前状态 |
| --- | --- | --- |
| 后端语言 | 保留 Go + Gin | 已确认 |
| ORM | 保留 GORM 处理常规 CRUD | 已确认 |
| 数据库 migration | 使用 `goose`，不再以 `AutoMigrate` 作为主 schema 机制 | 阶段 1 已落地 |
| 主数据库 | PostgreSQL 16 | 已落地并验证 |
| PostgreSQL 运行模型 | 使用外部 PostgreSQL，不由本地 compose 启动 | 已落地并验证 |
| 缓存/本地服务 | Redis 由本地 Docker 提供 | 已落地并验证 |
| 监控本地服务 | Prometheus + Pushgateway 由本地 Docker 提供 | 已落地并验证 |
| 前端方向 | 当前继续使用 Vue 3 + Element Plus 体系 | 当前代码可继续开发 |
| 本地开发环境 | Compose 启动 Redis、Pushgateway、Prometheus、API、Web | 已落地并验证 |
| 文档位置 | 项目级文档放根目录 `docs/`，生成类 API 文档保留在 `api/docs/` | 已确认 |

已验证的开发 PostgreSQL 信息：

```text
Host: 52.82.70.176
Port: 5432
Database: mydb
Username: admin
SSL mode: disable
```

不要在新增文档或示例里提交真实密码。运行时凭据应通过本地 `.env` 或环境变量管理。

## 3. 当前阶段状态

项目当前处于阶段 1 收尾完成状态。

| 阶段 | 名称 | 状态 | 进入下一阶段的标准 |
| --- | --- | --- | --- |
| 阶段 0 | 方向与原则确认 | 已完成 | 独立演进方向、技术方向、接入原则已形成文档 |
| 阶段 1 | 通用底座搭建 | 基本完成 | 数据库、migration、配置、本地运行、前后端基础骨架可用 |
| 阶段 2 | 接入型能力闭环 | 下一阶段 | 完成 1-2 个真实接入能力，跑通同步、回写、审计闭环 |
| 阶段 3 | 模块化增强 | 未开始 | 根据真实使用情况增强重点模块 |
| 阶段 4 | 治理与优化 | 未开始 | 完善审计、安全、可观测性、适配器治理和性能 |

## 4. 阶段 1 对齐说明

本节用于对齐原始阶段 1 目标、`docs/2026-04-23.txt` 原始记录，以及当前 `docs/phase1/` 的实际状态。

### 4.1 原始阶段 1 目标

阶段 1 的核心目标是：底座稳定，可支撑接入型能力开发。

当时明确要回答的问题：

- 新底座的目录组织和工程骨架是否稳定。
- migration、日志、配置、鉴权方案是否正式定稿。
- 本地开发环境是否稳定可复用。
- 前后端基础骨架是否达到可持续开发状态。
- 外部适配层和同步框架是否达到可复用状态。

### 4.2 2026-04-23 的状态

`docs/2026-04-23.txt` 记录的是阶段 1 初始落地状态：

- 独立演进路线和阶段 1 方向已写入文档。
- 接入目录 API 骨架和首批前端入口占位已加入。
- 修复了鉴权失败日志泄露敏感信息的问题。
- 已选择 PostgreSQL 16，但 PostgreSQL 尚未完成真实接入验证。
- migration 和完整后端验证仍在待办中。

这份记录作为历史快照仍然有效，但不能作为当前状态来读。

### 4.3 2026-04-24 的更新

2026-04-24 继续推进并关闭了 2026-04-23 的主要待办：

- 后端支持 PostgreSQL 驱动和方言切换。
- 配置支持 `config.yaml + 环境变量覆盖`。
- 正式 SQL migration 已落地，并转换为 `goose` 格式。
- 临时自研 migration runner 已被 `goose` 替换。
- 当前核心 GORM 模型已由 SQL migration 覆盖。
- `AutoMigrate` 不再作为主 schema 路径。
- 文档中的 PostgreSQL 16 实例已完成端到端验证。
- API 主进程已使用外部 PostgreSQL 完成冒烟验证。
- 本地 Docker Compose 已能启动 Redis、Pushgateway、Prometheus、API、Web。
- Web、API integration catalog、Prometheus、Pushgateway 的 HTTP 冒烟检查均通过。

详细执行记录见：

- `docs/phase1/phase1-2026-04-23-progress.md`
- `docs/phase1/phase1-2026-04-24-progress.md`
- `docs/phase1/phase1-postgresql-migration.md`

## 5. 阶段 1 验收对照

| 阶段 1 问题 | 当前结论 | 说明 |
| --- | --- | --- |
| 目录组织和工程骨架是否稳定 | 基本稳定 | 当前仍保持 `api/`、`web/`、`docker/`、`docs/` 结构。更深的领域化目录重构延后。 |
| migration 是否正式定稿 | 是 | `goose` 已作为正式 migration 工具，迁移文件位于 `api/migrations/`。 |
| 配置方案是否正式定稿 | 基本是 | 当前方案是 `config.yaml + 环境变量覆盖`，Docker 通过 `.env` 注入运行时配置。 |
| 日志方案是否正式定稿 | 否 | 现有日志可用，但结构化日志标准、必备字段、容器日志约定尚未定稿。 |
| 鉴权方案是否正式定稿 | 否 | 当前鉴权可运行，并已修复一处敏感日志问题；`jwt-go` 到 `golang-jwt/jwt` 的迁移仍未完成。 |
| 本地开发环境是否稳定可复用 | 是，但有磁盘约束 | Compose 联调通过。Docker Desktop 需要预留足够磁盘空间。 |
| 前后端基础骨架是否可持续开发 | 基本是 | API/Web 能构建和启动。Vite/TS/Pinia 前端现代化不是阶段 1 阻塞项。 |
| 外部适配层是否可复用 | 部分完成 | integration catalog 已有，但 adapter/provider 接口还未正式抽象。 |
| 同步框架是否可复用 | 部分完成 | 调度器和任务队列能运行，但统一同步、回写、审计生命周期仍需在阶段 2 验证。 |

结论：阶段 1 可以视为达到进入阶段 2 的条件。但日志、鉴权、adapter/sync 抽象应在阶段 2 早期收口。

## 6. 当前运行模型

本地开发运行模型：

```text
外部 PostgreSQL 16
本地 Docker Redis
本地 Docker Pushgateway
本地 Docker Prometheus
本地 Docker API
本地 Docker Web
```

常用命令：

```bash
cd docker
docker compose up -d redis pushgateway prometheus devops-api devops-web
docker compose ps
```

冒烟检查：

```text
Web:         http://localhost:8088
API catalog: http://localhost:8000/api/v1/integration/catalog
Prometheus:  http://localhost:9090/-/healthy
Pushgateway: http://localhost:9091/-/healthy
```

已知运行约束：

- 构建 API/Web 镜像前，Docker Desktop 必须预留足够磁盘空间。
- API Docker 构建已降低 Go 编译并发，减少临时磁盘压力。
- 如果 Docker 出现 `input/output error` 或 `no space left on device`，优先检查 `docker info`、`docker system df` 和宿主磁盘空间。

## 7. 仍需补充的工作

以下是进入阶段 2 前后建议补齐的底座事项。

### P0：文档与交接

- 保持 `docs/README.md` 作为文档总入口。
- 阶段执行记录统一放在 `docs/phase1/`。
- 生成类 Swagger 文档保留在 `api/docs/`。
- 不再把阶段进展日志混入 `api/docs/`。

### P1：鉴权收口

- 将 `github.com/dgrijalva/jwt-go` 替换为 `github.com/golang-jwt/jwt`。
- 确认 token 过期时间、刷新策略、签名算法和 middleware 边界。
- 补一份最小鉴权安全说明，覆盖日志脱敏和失败行为。

### P1：日志收口

- 决定阶段 2 继续使用 `logrus`，还是开始逐步迁移到 `slog`。
- 定义最小日志字段：时间、级别、request id、用户 id、模块、操作、错误。
- 确认容器日志策略：运行日志优先 stdout/stderr，文件挂载只保留现有必要场景。

### P1：适配层与同步骨架

- 定义最小 integration adapter 接口。
- 定义同步任务生命周期：pending、running、success、failed、skipped。
- 定义外部对象映射规则和审计字段。
- 用阶段 2 的第一个真实能力验证抽象，不提前过度设计。

### P2：前端工程现代化

- 保持当前 Vue 3 + Element Plus 代码可继续使用。
- Vite/TypeScript/Pinia 按真实功能开发逐步迁移。
- 不让完整前端重构阻塞阶段 2 接入能力交付。

## 8. 阶段 2 入口建议

阶段 2 建议从一个窄范围真实接入能力开始，完整跑通：

```text
外部数据源 -> adapter -> 统一模型 -> 同步任务 -> UI 列表/详情 -> 编辑/回写 -> 审计日志
```

推荐首批候选：

1. 基于飞书多维表格的运营工单。
2. 基于飞书多维表格的站点管理。
3. 基于现有 `dnsmgr` 方向的域名管理。

第一个阶段 2 模块应足够小，但必须真实验证：

- adapter 接口
- 同步生命周期
- 回写行为
- 错误处理
- 审计和日志约定
- 前端列表、详情、编辑骨架

## 9. 文档目录约定

推荐目录含义：

```text
docs/
  README.md                         # 文档总入口和阅读顺序
  autoops-analysis/                 # 架构、路线图、决策、分析
  phase1/                           # 阶段 1 执行记录和交接说明
  YYYY-MM-DD.txt                    # 原始历史记录

api/docs/
  docs.go
  swagger.json
  swagger.yaml
```

规则：

- 项目级决策放在 `docs/autoops-analysis/`。
- 阶段执行总结和交接说明放在 `docs/phase1/`。
- 生成类 API 文档只放在 `api/docs/`。
- 旧的 `docs/YYYY-MM-DD.txt` 视为原始历史记录，不作为当前状态唯一来源。

## 10. 一句话结论

AutoOps 已完成进入接入型能力开发所需的阶段 1 底座：PostgreSQL、`goose`、配置、本地 Docker 运行链路均已验证；下一步应在阶段 2 早期正式收口鉴权/日志，并用第一个真实接入模块验证 adapter/sync 模型。
