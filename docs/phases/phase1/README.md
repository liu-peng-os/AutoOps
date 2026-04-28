# 阶段 1 交接说明

## 当前状态

阶段 1 可以视为基本完成。

当前底座已经能够支撑阶段 2 的接入型能力开发。仍有少量横切方案需要在阶段 2 早期收口，但不阻塞先启动一个窄范围的真实接入闭环。

## 已完成

- PostgreSQL 16 已选型，并通过文档中的外部数据库实例完成验证。
- 后端运行时数据库配置支持 PostgreSQL 和环境变量覆盖。
- SQL migration 已由 `goose` 接管。
- 当前已登记的核心 GORM 模型已由 SQL migration 覆盖。
- `AutoMigrate` 不再作为主要 schema 管理路径。
- Docker Compose 已能启动 Redis、Pushgateway、Prometheus、API、Web。
- API 能连接外部 PostgreSQL，并确认 `goose` 当前版本为 `4`。
- Web 和 API 冒烟检查通过。
- Docker 运行数据、上传目录和临时日志已加入 `.gitignore`。

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

## 仍需收口

以下内容建议在阶段 2 早期处理：

- 鉴权：从 `dgrijalva/jwt-go` 迁移到 `golang-jwt/jwt`，并补充 token 策略说明。
- 日志：确认最终日志标准、必备字段和容器日志输出方式。
- 适配层：定义最小 provider/adapter 接口。
- 同步框架：定义可复用的同步、回写、审计生命周期。
- Docker 构建策略：确认本地开发继续源码构建，还是回到更省磁盘的运行时镜像模式。

## 与 2026-04-23 记录的关系

`docs/2026-04-23.txt` 是阶段 1 初始落地时的原始记录。当时 PostgreSQL 接入、migration 实链验证、本地 Docker 联调都还没有完成。

2026-04-24 的阶段 1 记录已经覆盖这些待办：

- PostgreSQL 已完成真实验证。
- `goose` 已替代临时 migration runner。
- 本地 Docker 栈已完成联调。
- 阶段 1 底座已达到进入阶段 2 的条件。

## 下一步

阶段 2 建议从一个窄范围接入能力开始，完整跑通：

```text
外部数据源 -> adapter -> 统一模型 -> 同步任务 -> UI -> 回写 -> 审计
```

候选方向：

- 基于飞书多维表格的运营工单。
- 基于飞书多维表格的站点管理。
- 基于 `dnsmgr` 方向的域名管理。
