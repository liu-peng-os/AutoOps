# 阶段0决策：技术栈选型

> 决策时间：2026-04-23
> 决策状态：已确认
> 相关文档：`docs/roadmap.md`

---

## 背景

确定独立演进方向后，需要明确技术栈。核心问题是：

**哪些技术保留，哪些升级，哪些重构？**

---

## 最终决策

### 保留（低风险高收益）

| 技术 | 理由 |
|---|---|
| Go 1.24+ | 团队熟悉，性能够用，生态成熟 |
| Gin | 轻量高效，不为"更现代"而换框架 |
| GORM | 日常 CRUD 够用，复杂查询可手写 SQL |
| Vue 3 | 团队熟悉，Element Plus 对后台场景高效 |
| Element Plus | 组件丰富，不需要重新造轮子 |
| Redis | 缓存和队列基础设施 |

### 升级（中等成本，长期收益）

| 技术 | 从 | 到 | 理由 |
|---|---|---|---|
| 数据库 | MySQL | PostgreSQL 16 | JSONB、复杂查询、资源标签更友好 |
| Migration | AutoMigrate | goose | Schema 演进可追溯、可回滚 |
| 前端构建 | Webpack | Vite | 开发体验更好 |
| 状态管理 | Vuex | Pinia | Vue 3 官方推荐 |
| JWT 库 | dgrijalva/jwt-go | golang-jwt/jwt | 原库已废弃 |

### 新增（填补空白）

| 技术 | 用途 |
|---|---|
| TypeScript | 前端类型安全 |
| asynq | 异步任务队列（替代手工 Redis 队列） |

### 不做（避免过度工程）

- 不换 React/Next（避免无意义重学）
- 不过早拆微服务（单体够用）
- 不引入复杂消息系统（Redis + asynq 够用）
- 不换 UI 框架（Element Plus 够用）

---

## 关键决策的推导过程

### 决策1：为什么切 PostgreSQL

**考虑的方案**：
- 方案A：继续用 MySQL
- 方案B：切 PostgreSQL
- 方案C：同时支持 MySQL 和 PostgreSQL

**选择方案B的理由**：
1. 平台型系统需要存储大量 JSON 配置、资源标签、策略规则，PostgreSQL 的 JSONB 更适合
2. 复杂查询（多表关联、聚合、窗口函数）PostgreSQL 性能更好
3. 既然独立演进，就不必被上游的 MySQL 偏好绑住
4. 团队有 PostgreSQL 经验，迁移成本可控

**否决方案A的理由**：
- MySQL 的 JSON 支持不如 PostgreSQL
- 复杂查询需要更多优化工作

**否决方案C的理由**：
- 同时支持两个数据库会增加维护成本
- 测试成本翻倍
- 某些特性（JSONB、窗口函数）无法充分利用

### 决策2：为什么用 goose 而不是 golang-migrate

**考虑的方案**：
- 方案A：继续用 AutoMigrate
- 方案B：自己实现 migration runner
- 方案C：用 golang-migrate
- 方案D：用 goose

**选择方案D的理由**：
1. goose 更简单，学习成本低
2. 支持 SQL 和 Go 两种格式
3. 社区活跃，文档完善
4. 支持 up/down migration

**否决方案A的理由**：
- AutoMigrate 不可追溯、不可回滚
- Schema 变更无法 review
- 生产环境风险高

**否决方案B的理由**：
- 自己实现容易有 bug
- 缺少 down migration 支持
- 维护成本高

**否决方案C的理由**：
- golang-migrate 功能更强大，但也更复杂
- 对于我们的场景，goose 够用

### 决策3：为什么保留 Gin 而不换 Echo/Fiber

**理由**：
1. Gin 性能够用，生态成熟
2. 团队熟悉，切换成本高
3. 真正的问题不在框架，而在模块边界和数据演进
4. 换框架不会解决核心问题

### 决策4：为什么保留 Vue 3 而不换 React

**理由**：
1. 团队熟悉 Vue
2. Element Plus 对后台场景高效
3. 换框架需要重写所有前端代码，成本太高
4. Vue 3 + Vite + TypeScript 已经足够现代

---

## 技术债清理优先级

### P0（阶段1必须做）

- ✅ 切换到 PostgreSQL
- ✅ 引入 goose migration
- ✅ 停止依赖 AutoMigrate

### P1（阶段2早期做）

- ⏳ jwt-go → golang-jwt/jwt
- ⏳ 统一日志标准（logrus 或 slog）
- ⏳ 前端 Vite + TypeScript + Pinia

### P2（阶段2期间按需做）

- ⏳ 引入 asynq 替代手工 Redis 队列
- ⏳ 模块目录重构

### P3（阶段3或更晚）

- ⏳ logrus → slog（如果阶段2选择继续用 logrus）
- ⏳ 更深的领域化目录结构

---

## 约束条件

这些决策基于以下约束：

1. **团队规模**：小团队，无法同时做多个大改造
2. **业务压力**：需要快速响应需求，不能长期停止开发
3. **技术债**：现有代码的技术债已经影响效率
4. **学习成本**：团队对 Go/Vue 熟悉，对 Rust/React 不熟悉

如果约束条件变化（比如团队扩大、业务压力减小），可以重新评估。

---

## 后续行动

基于这些决策，后续需要：

1. 阶段1完成 PostgreSQL + goose 迁移（已完成）
2. 阶段2早期完成 jwt-go 和日志标准收口
3. 阶段2期间逐步引入 asynq 和前端现代化
