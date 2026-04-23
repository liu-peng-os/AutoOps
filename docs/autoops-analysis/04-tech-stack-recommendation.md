# 天枢 AutoOps — 技术栈选型评估与二次开发建议

> 文档版本：v1.0 | 分析日期：2026-04-23 | 基于实际源码分析

---

## 一、现有技术栈客观评价

### 总体评价

选型整体**务实且合理**，没有过度工程化，也没有明显的错误决策。适合中小企业运维团队的实际场景。但有几个地方值得讨论。

### 后端

**Go + Gin — 合适**

运维平台的核心场景是 SSH 终端（大量长连接）、任务并发执行、WebSocket 实时推送，这些都是 Go 的强项。Gin 是 Go 生态里最成熟的 HTTP 框架，没有争议。

**GORM — 可接受，但有代价**

GORM 的 `AutoMigrate` 在开发阶段很方便，但生产环境用 AutoMigrate 管理 schema 是有风险的——它只会加列不会删列，schema 漂移问题难以追踪。更规范的做法是用 `golang-migrate` 或 `goose` 管理版本化迁移脚本。这不是 GORM 本身的问题，是使用方式的问题。

**dgrijalva/jwt-go — 有问题**

该库 2020 年就停止维护，已有 CVE。迁移到 `golang-jwt/jwt` 的成本极低（API 几乎完全兼容），没有理由继续用废弃库。这是技术债务里最容易修的一个。

**robfig/cron — 合适**

成熟稳定，没有问题。

**Redis 任务队列 vs 专用消息队列**

用 Redis List 做任务队列在这个场景下是合理的——运维任务不是高吞吐场景，Redis 已经是基础设施依赖，不引入额外组件是正确的工程判断。如果未来任务量上去了，迁移到 `asynq`（基于 Redis 的专用任务队列库）比引入 RabbitMQ/Kafka 代价小得多，且有 Web UI、重试策略、任务去重等开箱即用的能力。

### 前端

**Vue3 + Element Plus — 合适**

企业内部工具的标准选择，Element Plus 组件覆盖面广，开发效率高。

**Vue CLI（Webpack）— 落后**

这是最明显的问题。Vue CLI 已经进入维护模式，Vue 官方推荐 Vite。Webpack 冷启动在这个体量的项目下可能需要 30-60 秒，Vite 可以做到 1-2 秒。对于需要频繁迭代的运维工具，这个差距每天都在消耗开发者时间。迁移成本中等，但收益明显。

**Vuex v4 — 可以，但 Pinia 更现代**

Vuex 4 是 Vue3 的兼容版本，能用但不是最优选。Pinia 是 Vue 官方现在推荐的状态管理方案，TypeScript 支持更好，API 更简洁。对于新项目来说选 Pinia 更合适，但已有 Vuex 代码迁移成本不低，不是必须改的。

**moment.js — 应该替换**

moment.js 已宣布停止新功能开发，包体积大（约 67KB gzip）。`dayjs`（API 兼容，2KB）是更好的选择。替换成本极低。

### 基础设施

**Docker Compose — 合适**

对于中小企业运维平台，Docker Compose 是正确的起点，简单可靠。项目同时提供了 K8s 部署清单，给了用户选择空间。

**Prometheus + Pushgateway — 合适**

Agent 主动推送到 Pushgateway 的模式适合主机监控场景（主机可能在防火墙后面，Prometheus 无法主动 scrape）。这个架构决策是对的。

---

## 二、PostgreSQL 替换 MySQL 专项分析

**结论：抛开改造难度，PostgreSQL 是更好的选择。**

### 代码中的 MySQL 耦合点（已扫描确认）

```
1. api/pkg/db/db.go
   → gorm.io/driver/mysql（MySQL 专用驱动）
   → DSN 中 sql_mode=''（禁用严格模式，MySQL 特有）

2. api/sql/autoops.sql（初始化脚本）
   → ENGINE=InnoDB（每张表）
   → AUTO_INCREMENT（每个主键）
   → DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
   → ROW_FORMAT=DYNAMIC

3. model 层 GORM tag
   → type:tinyint（大量状态字段）
   → type:json（Application 模型的 JSON 数组字段）
   → type:longtext（knowledge.go 的 Content 字段）

4. api/api/app/service/application.go（原生 SQL）
   → HAVING service_count > 0（引用 SELECT 别名，MySQL 非标准行为）
   → CONCAT() 函数（标准 SQL，PG 支持）

5. dao 层两处 TRUNCATE TABLE
   → sys_logininfo.go / sysOperationLog.go（标准 SQL，PG 支持）
```

耦合点数量不多，且集中在可控位置，迁移工作量可控。

### 支持换 PostgreSQL 的理由

**1. JSON 字段支持更强**

项目里 `Application` 模型大量使用 JSON 列存储数组数据：

```go
DevOwners   JSONArray `gorm:"type:json"`
TestOwners  JSONArray `gorm:"type:json"`
HostIDs     JSONArray `gorm:"type:json"`
Domains     JSONArray `gorm:"type:json"`
Databases   JSONArray `gorm:"type:json"`
```

MySQL 的 JSON 类型没有索引支持，只能全量读取后在应用层过滤。PostgreSQL 的 `jsonb` 支持 GIN 索引，可以直接 `WHERE hosts @> '[1,2,3]'`，随着数据量增长差距会越来越明显。

**2. `sql_mode=''` 在掩盖数据质量问题**

`pkg/db/db.go` 的 DSN 里有 `sql_mode=''`，这禁用了 MySQL 的严格模式，允许空字符串写入 NOT NULL 列、截断超长字符串等不规范操作静默通过。PostgreSQL 默认就是严格的，这类问题会在开发阶段暴露，而不是在生产环境悄悄吞掉。

**3. `HAVING` 别名引用是非标准行为**

`application.go` 的原生 SQL 用了 `HAVING service_count > 0`（引用 SELECT 别名），这在 MySQL 里能跑，在标准 SQL 和 PostgreSQL 里需要写成子查询。说明现有代码在写 SQL 时没有考虑可移植性，换 PG 会强制规范 SQL 写法。

**4. 长期维护性更好**

PostgreSQL 的 SQL 标准兼容性更好，社区更活跃，扩展生态更丰富（PostGIS、TimescaleDB 等）。国内云厂商（阿里云、腾讯云）均提供托管 PG 服务。

### MySQL 也说得过去的理由

- 这是运维内部工具，数据量不大，MySQL 性能完全够用
- 国内运维团队对 MySQL 更熟悉，DBA 资源更容易获取
- 现有耦合点数量不多，但迁移仍需投入时间

**如果是新项目，选 PostgreSQL。主要原因是 JSON 字段的使用方式——MySQL 的 JSON 支持在这个场景下是个半成品。**

---

## 三、二次开发技术栈建议

### 后端

| 项目 | 现状 | 建议 | 优先级 |
|------|------|------|--------|
| Go + Gin | ✅ | 保留 | — |
| GORM | 保留，改用法 | 换掉 AutoMigrate，引入 `golang-migrate` | 高 |
| dgrijalva/jwt-go | ❌ 废弃库 | 换 `golang-jwt/jwt`（API 兼容） | 高 |
| MySQL | 现状 | **换 PostgreSQL** | 高 |
| Redis 任务队列 | 保留 | 可选换 `asynq` | 低 |
| logrus | 保留 | 可选换 Go 1.21 标准库 `slog` | 低 |

### 前端

| 项目 | 现状 | 建议 | 优先级 |
|------|------|------|--------|
| Vue3 + Element Plus | ✅ | 保留 | — |
| Vue CLI / Webpack | ❌ 停止维护 | **换 Vite** | 高 |
| Vuex v4 | 保留 | 新模块用 Pinia，老模块不动 | 低 |
| moment.js | ❌ 停止维护 | 换 `dayjs`（API 兼容，2KB） | 高 |
| Axios v0.27 | 保留 | 可升级到 v1.x | 低 |

### 基础设施

| 项目 | 现状 | 建议 | 优先级 |
|------|------|------|--------|
| Docker Compose | ✅ | 保留 | — |
| MySQL → PostgreSQL | — | 替换容器镜像 + 初始化脚本 | 配合后端 |
| Prometheus + Pushgateway | ✅ | 保留 | — |
| Nginx | ✅ | 保留 | — |

---

## 四、与上游保持同步的分支策略

### 场景定义

```
upstream/main（原项目）
    ↓ git fetch upstream
your/main          ← 镜像 upstream，只做 sync，不做任何开发
    ↓
your/compat        ← 薄适配层（见下方说明）
    ↓
your/dev           ← 你的重构版本，日常开发在此
```

### compat 分支的作用

上游更新了新功能（比如新增了一个 K8s 资源类型的管理），你想借鉴它的业务逻辑。流程如下：

```
1. git fetch upstream
2. git cherry-pick <upstream commit> → your/compat
3. 在 compat 上做最小化适配：
   - MySQL 特定写法 → PG 兼容
   - gorm:"type:tinyint" → gorm:"type:smallint"
   - gorm:"type:json" → gorm:"type:jsonb"
   - HAVING 别名引用 → 改为子查询
   - Vuex 新模块 → 改为 Pinia（如果涉及）
4. git merge your/compat → your/dev
```

这样 dev 分支永远只接收"已适配"的内容，不会因为上游用了 MySQL 特性而污染你的 PG 代码库。

### 降低冲突的核心原则

**把技术栈改动集中在"可替换层"，不要散落在业务代码里。**

```
api/pkg/db/db.go          ← 数据库驱动切换，只改这一个文件
api/pkg/db/migrate.go     ← 迁移策略，只改这一个文件
docker/docker-compose.yml ← 基础设施，只改这一个文件
web/vite.config.js        ← 构建工具，新增替换 vue.config.js
```

上游 99% 的业务逻辑变更不会碰这些文件，所以你的技术栈改动和上游的功能更新基本不会产生冲突。

### 需要手动处理的冲突类型

上游更新时，以下两类变更需要在 compat 分支手动适配：

| 上游变更类型 | 需要适配的内容 |
|------------|--------------|
| 新增 model 文件（新表） | `gorm:"type:tinyint"` → `smallint`；`type:json` → `jsonb`；`type:longtext` → `text` |
| 新增原生 SQL（db.Raw/db.Exec） | 检查是否有 MySQL 方言，改为 PG 兼容写法 |
| 新增前端模块 | 检查是否引入了新的 moment.js 用法，替换为 dayjs |

这两类冲突无法完全避免，但数量不多，在 compat 分支上集中处理即可。

---

## 五、改造优先级排序

```
第一批（高收益、低风险，先做）
  1. Vue CLI → Vite          开发体验立竿见影，冷启动 30s → 1s
  2. jwt-go → golang-jwt     消除安全隐患，2 小时内完成，API 兼容
  3. moment → dayjs          10 分钟，无风险，API 兼容

第二批（架构级，需要规划）
  4. MySQL → PostgreSQL      迁移 schema + 适配 GORM tag + 改原生 SQL
  5. AutoMigrate → golang-migrate  配合 PG 迁移一起做，建立版本化迁移体系

第三批（可选，按需）
  6. Redis 队列 → asynq      任务量增长后再考虑，现有方案够用
  7. Vuex → Pinia            新模块用 Pinia，老模块不动，渐进式迁移
  8. logrus → slog           减少外部依赖，优先级最低
```

---

## 六、MySQL → PostgreSQL 迁移要点速查

| 改动点 | MySQL 写法 | PostgreSQL 写法 |
|--------|-----------|----------------|
| GORM 驱动 | `gorm.io/driver/mysql` | `gorm.io/driver/postgres` |
| DSN 格式 | `user:pass@tcp(host)/db?charset=utf8mb4` | `host=h user=u password=p dbname=d sslmode=disable` |
| 整数状态字段 | `gorm:"type:tinyint"` | `gorm:"type:smallint"` |
| JSON 数组字段 | `gorm:"type:json"` | `gorm:"type:jsonb"` |
| 长文本字段 | `gorm:"type:longtext"` | `gorm:"type:text"` |
| HAVING 别名 | `HAVING service_count > 0` | `HAVING COUNT(*) > 0`（改为表达式） |
| 禁用严格模式 | DSN `sql_mode=''` | 删除，PG 默认严格 |
| 初始化脚本 | `ENGINE=InnoDB AUTO_INCREMENT` | 删除，PG 用 `SERIAL` 或 `BIGSERIAL` |
| 字符集声明 | `CHARSET=utf8mb4` | 删除，PG 默认 UTF-8 |
| TRUNCATE | `TRUNCATE TABLE t` | `TRUNCATE TABLE t`（兼容） |

---

*本文档基于 2026-04-23 的源码快照分析生成。技术选型建议仅供参考，最终决策需结合团队实际情况。*
