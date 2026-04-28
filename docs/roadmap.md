# AutoOps 路线图

> 版本：3.0
> 更新时间：2026-04-24
> 当前阶段：阶段1已完成，准备进入阶段2

---

## 1. 项目定位

AutoOps 按独立运维平台演进，不再以低成本持续兼容 upstream 为首要目标。

核心原则：
- 保留原项目中已验证的产品理念和功能方向
- 允许在架构、数据模型、工程体系上做必要改造
- 所有外部系统都视为可替换依赖，不绑定到核心模型里
- 优先建设通用底座，再逐步做深业务能力

背景和推导过程见 `docs/decisions/`。

---

## 2. 阶段状态

| 阶段 | 名称 | 状态 | 进入下一阶段的标准 |
|---|---|---|---|
| 阶段0 | 方向与原则确认 | 已完成 | 独立演进方向、技术方向、接入原则已形成文档 |
| 阶段1 | 通用底座搭建 | 已完成 | 数据库、migration、配置、本地运行、前后端基础骨架可用 |
| 阶段2 | 接入型能力闭环 | 下一阶段 | 完成 1-2 个真实接入能力，跑通同步、回写、审计闭环 |
| 阶段3 | 模块化增强 | 未开始 | 根据真实使用情况增强重点模块 |
| 阶段4 | 治理与优化 | 未开始 | 完善审计、安全、可观测性、适配器治理和性能 |

---

## 3. 已确认的技术决策

| 领域 | 决策 | 状态 |
|---|---|---|
| 后端语言 | 保留 Go + Gin | 已确认 |
| ORM | 保留 GORM 处理常规 CRUD | 已确认 |
| Migration | 使用 goose，不再以 AutoMigrate 作为主 schema 机制 | 阶段1已落地 |
| 主数据库 | PostgreSQL 16 | 阶段1已落地并验证 |
| 缓存 | Redis 由本地 Docker 提供 | 阶段1已落地并验证 |
| 本地开发环境 | Compose 启动 Redis、Pushgateway、Prometheus、API、Web | 阶段1已落地并验证 |
| 前端方向 | Vue 3 + Element Plus，逐步引入 Vite/TS/Pinia | 当前代码可继续开发 |
| 外部系统接入 | adapter/provider 模式，核心模型不绑定外部系统 | 原则已确认，阶段2验证 |

详细推导过程见 `docs/decisions/`。

---

## 4. 当前运行模型

```text
外部 PostgreSQL 16（52.82.70.176:5432）
本地 Docker Redis
本地 Docker Pushgateway
本地 Docker Prometheus
本地 Docker API
本地 Docker Web
```

启动命令：

```bash
cd docker
docker compose up -d redis pushgateway prometheus devops-api devops-web
```

冒烟检查：

```text
Web:         http://localhost:8088
API catalog: http://localhost:8000/api/v1/integration/catalog
Prometheus:  http://localhost:9090/-/healthy
Pushgateway: http://localhost:9091/-/healthy
```

已知约束：
- Docker Desktop 需要预留足够磁盘空间，构建前检查 `docker system df`
- PostgreSQL 密码通过 `.env` 管理，不提交到 git

---

## 5. 阶段2前置条件（进入阶段2前必须完成）

以下是阶段1遗留的收口工作，**必须在阶段2第一个功能模块开发前完成**，否则会阻塞后续开发：

### 5.1 鉴权收口（P1，阻塞性）

- 将 `github.com/dgrijalva/jwt-go` 替换为 `github.com/golang-jwt/jwt`
- 确认 token 过期时间、刷新策略、签名算法和 middleware 边界
- 补一份最小鉴权安全说明，覆盖日志脱敏和失败行为

**为什么是阻塞性**：`dgrijalva/jwt-go` 已废弃，有已知安全漏洞，阶段2会新增需要鉴权的接口，不能继续用旧库。

### 5.2 日志收口（P1，阻塞性）

- 决定阶段2继续使用 `logrus`，还是开始逐步迁移到 `slog`
- 定义最小日志字段：时间、级别、request_id、user_id、模块、操作、错误
- 确认容器日志策略：运行日志优先 stdout/stderr

**为什么是阻塞性**：阶段2的同步任务、回写任务、审计日志都需要统一的日志格式，不定稿就没法做。

### 5.3 适配层接口定义（P1，阻塞性）

- 定义最小 integration adapter 接口（Provider 接口）
- 定义同步任务生命周期：pending → running → success/failed/skipped
- 定义外部对象映射规则和审计字段

**为什么是阻塞性**：阶段2的第一个接入模块需要实现这个接口，不定义就没法开始。

---

## 6. 阶段2入口建议

阶段2从一个窄范围真实接入能力开始，完整跑通：

```text
外部数据源 → adapter → 统一模型 → 同步任务 → UI 列表/详情 → 编辑/回写 → 审计日志
```

推荐首批候选（选一个开始）：

1. 基于飞书多维表格的运营工单
2. 基于飞书多维表格的站点管理
3. 基于 dnsmgr 的域名管理

第一个模块必须验证：
- adapter 接口是否可用
- 同步生命周期是否合理
- 回写行为是否正确
- 错误处理是否完善
- 审计和日志约定是否够用
- 前端列表、详情、编辑骨架是否可复用

---

## 7. 阶段2期间的非阻塞工作（P2）

以下工作不阻塞阶段2接入能力交付，但建议在阶段2期间逐步推进：

- 前端工程现代化：Vite/TypeScript/Pinia 按真实功能开发逐步迁移
- 引入 asynq 替代手工 Redis 队列（等第一个同步任务跑通后再引入）
- 模块目录重构（等阶段2有了真实模块后再重构）

---

## 8. 文档目录约定

```text
docs/
  README.md                    # 总入口，阅读顺序
  roadmap.md                   # 本文件：当前路线图和阶段状态
  decisions/                   # 阶段0决策记录（ADR 风格，只增不改）
    phase0-project-positioning.md
    phase0-tech-stack.md
    phase0-external-system-principle.md
  phases/
    phase1/                    # 阶段1执行记录和交接说明
    phase2/                    # 阶段2开始时建
  analysis/                    # 历史架构分析（已过时，结论以本文件为准）
  history/                     # 原始流水账（只读）
    INDEX.md
    YYYY-MM-DD.txt

api/docs/                      # 生成类 API 文档（Swagger）
```

规则：
- 项目级决策放在 `docs/decisions/`
- 阶段执行总结和交接说明放在 `docs/phases/phaseN/`
- 生成类 API 文档只放在 `api/docs/`
- `docs/history/` 里的 txt 文件是原始记录，不作为当前状态来源
