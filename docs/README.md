# AutoOps 文档入口

接手项目时，建议从本文件开始看。

---

## 快速开始

**如果你是第一次接手这个项目**，按这个顺序读：

1. [docs/roadmap.md](roadmap.md) — 当前路线图、阶段状态、技术决策
2. [docs/phases/phase1/README.md](phases/phase1/README.md) — 阶段1交接说明、本地运行方式
3. [docs/decisions/](decisions/) — 阶段0的决策记录（为什么独立演进、为什么选这些技术）

**如果你想了解项目背景和历史**：

- [docs/history/INDEX.md](history/INDEX.md) — 历史流水账索引
- [docs/analysis/](analysis/) — 阶段0的架构分析（已过时，结论以 roadmap.md 为准）

---

## 当前状态（一句话）

阶段1已完成（PostgreSQL + goose + Docker 本地环境已验证），准备进入阶段2（接入型能力闭环）。

---

## 目录结构

```text
docs/
  README.md                    # 本文件：文档入口
  roadmap.md                   # 当前路线图和阶段状态

  decisions/                   # 阶段0决策记录（ADR 风格，只增不改）
    phase0-project-positioning.md    # 为什么独立演进
    phase0-tech-stack.md             # 技术栈选型推导过程
    phase0-external-system-principle.md  # 可替换依赖原则

  phases/                      # 各阶段执行记录和交接说明
    phase1/                    # 阶段1：通用底座搭建
      README.md                # 阶段1交接说明（已完成/仍需收口/下一步）
      phase1-postgresql-migration.md
      phase1-2026-04-23-progress.md
      phase1-2026-04-24-progress.md
      phase1-2026-04-28-bugfix.md  # 前端 null 崩溃修复、Vue 2 遗留代码清理
    phase2/                    # 阶段2：接入型能力闭环
      README.md                # 阶段2前置条件和第一个模块选择

  analysis/                    # 历史架构分析（已过时，结论以 roadmap.md 为准）
    01-architecture-overview.md
    02-modules-deep-dive.md
    03-final-report.md
    04-tech-stack-recommendation.md
    05-independent-evolution-roadmap.md  # 已被 roadmap.md 替代

  history/                     # 原始流水账（只读）
    INDEX.md                   # 历史记录索引
    2025-10-11.txt             # 原作者原项目更新记录
    2025-12-10.txt
    2025-12-21.txt
    2025-12-24.txt
    2026-01-07.txt
    2026-03-27.txt             # 原作者最后一次更新
    2026-04-23.txt             # 二开后第一次记录（阶段1初始落地）

api/docs/                      # 生成类 API 文档（Swagger）
  docs.go
  swagger.json
  swagger.yaml
```

---

## 文档规则

1. **项目级决策放在 `docs/decisions/`**
   - ADR（Architecture Decision Record）风格
   - 只增不改，如果决策变化，新增一个文件说明为什么变化

2. **阶段执行总结和交接说明放在 `docs/phases/phaseN/`**
   - 每个阶段有自己的目录
   - README.md 是交接说明，记录已完成、仍需收口、下一步

3. **生成类 API 文档只放在 `api/docs/`**
   - Swagger 自动生成的文档
   - 不要把阶段进展、路线图、交接说明放进 `api/docs/`

4. **`docs/history/` 里的 txt 文件是原始记录，不作为当前状态来源**
   - 2026-03-27 之前是原作者原项目的更新记录
   - 2026-04-23 开始是二开后的记录
   - 这些文件保留用于追溯历史，但当前状态以 `roadmap.md` 为准

---

## 本地运行

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

详细说明见 [docs/phases/phase1/README.md](phases/phase1/README.md)。

---

## 下一步

阶段2需要先完成三个前置条件（鉴权收口、日志收口、适配层接口定义），然后从一个窄范围接入能力开始。

详细说明见 [docs/phases/phase2/README.md](phases/phase2/README.md)。
