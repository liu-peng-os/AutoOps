# 阶段 1 交接说明

## 当前状态

**阶段 1 已完成，正式关闭。** (2026-04-28)

底座已稳定，可以进入阶段 2。

## 已完成

- PostgreSQL 16 已选型，并通过文档中的外部数据库实例完成验证。
- 后端运行时数据库配置支持 PostgreSQL 和环境变量覆盖。
- SQL migration 已由 `goose v3.24.3` 接管。
- 当前已登记的核心 GORM 模型已由 SQL migration 覆盖。
- `AutoMigrate` 不再作为主要 schema 管理路径。
- Docker Compose 已能启动 Redis、Pushgateway、Prometheus、API、Web。
- API 能连接外部 PostgreSQL，并确认 `goose` 当前版本为 `4`。
- Web 和 API 冒烟检查通过。
- Docker 运行数据、上传目录和临时日志已加入 `.gitignore`。
- **JWT 已从 `dgrijalva/jwt-go`（有CVE）迁移到 `golang-jwt/jwt/v5`。**
- **Docker 构建改用 vendor 模式，不再依赖构建时联网下载依赖。**
- **前端 CMDB 主机管理页面 null 崩溃已修复（空分组时后端返回 `[]` 而非 `null`）。**
- **前端 Vue 2 遗留代码已清理（`this.$set`、`slot="footer"`）。**

## 已验证命令

```bash
cd docker
docker compose up -d redis pushgateway prometheus devops-api devops-web
docker compose ps
```

期望服务状态：

```text
devops-redis         healthy
devops-pushgateway   healthy
devops-prometheus    healthy
devops-api           running
devops-web           healthy
```

冒烟检查地址：

```text
http://localhost:8088
http://localhost:8000/api/v1/integration/catalog
http://localhost:9090/-/healthy
http://localhost:9091/-/healthy
```

以上地址在阶段 1 验证时均返回 HTTP 200。

## 遗留（不阻塞阶段 2，建议早期处理）

- **日志**：确认最终日志标准、必备字段和容器日志输出方式（目前 logrus 和 zap 混用）。
- **适配层**：定义最小 provider/adapter 接口（等第一个真实接入场景出现后再抽象）。
- **同步框架**：定义可复用的同步、回写、审计生命周期（同上）。
- **前端图标**：全量注册 Element Plus 图标（约 300 个），建议阶段 2 改为按需引入。

## 进展记录

| 日期 | 文件 | 内容 |
|------|------|------|
| 2026-04-23 | [phase1-2026-04-23-progress.md](phase1-2026-04-23-progress.md) | 初始落地记录 |
| 2026-04-24 | [phase1-2026-04-24-progress.md](phase1-2026-04-24-progress.md) | PostgreSQL验证、goose接管、Docker联调 |
| 2026-04-28 | [phase1-2026-04-28-bugfix.md](phase1-2026-04-28-bugfix.md) | JWT升级、前端null修复、构建稳定化 |

## 下一步

阶段 2 建议从一个窄范围接入能力开始，完整跑通：

```text
外部数据源 -> adapter -> 统一模型 -> 同步任务 -> UI -> 回写 -> 审计
```

候选方向：

- 基于飞书多维表格的运营工单。
- 基于飞书多维表格的站点管理。
- 基于 `dnsmgr` 方向的域名管理。
