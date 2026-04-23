# 天枢 AutoOps 运维管理系统 — 架构全景解析

> 文档版本：v1.0 | 分析日期：2026-04-23 | 基于实际代码分析

---

## 一、项目定位

### 项目是什么

**天枢 AutoOps**（全称：天枢 AutoOps 运维管理系统）是一套基于 **Go + Vue3** 开发的企业级运维自动化平台。"枢"取枢纽、中心之意，寓意将分散的运维工具整合为统一入口。

### 解决什么问题

| 痛点 | 现状 | 平台方案 |
|------|------|---------|
| 数据孤岛 | 服务器台账在飞书多维表格、工单手动流转、成本账单散落各处 | 统一 CMDB + 飞书自动同步，数据打通 |
| 工具碎片化 | Jenkins、JumpServer、Nacos、Grafana、Nightingale 各自独立，运维人员多系统切换 | 统一门户集成，一个界面管全部 |
| 可观测性不足 | 告警分散在多个系统，无统一仪表盘 | 内置 Prometheus 监控体系，统一告警速览 |
| 合规审计缺失 | 无统一操作审计入口 | 登录日志、操作日志、数据库操作、终端录像全面审计 |

### 目标用户

**仅面向运维团队内部使用**，覆盖以下角色：超级管理员、运维主管、运维工程师、SRE 工程师、只读用户。

核心适用场景：中小企业混合环境（VM + K8s）、需要传统运维与容器化双轨并行的团队。

---

## 二、技术栈清单

### 后端技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| **Go** | 1.24 | 主语言，高并发、低资源占用 |
| **Gin** | v1.11.0 | HTTP Web 框架，处理 REST API 路由 |
| **GORM** | v1.31.0 | ORM 框架，数据库操作抽象层 |
| **MySQL** | 8.0.33 | 主关系型数据库，存储所有业务数据 |
| **Redis** | 7.0 | 缓存层，存储 Session/Token、临时数据 |
| **go-redis** | v8.11.5 | Go Redis 客户端 |
| **JWT (dgrijalva/jwt-go)** | v3.2.0 | 用户认证令牌生成与校验 |
| **robfig/cron** | v3.0.1 | 定时任务调度器（Cron 表达式） |
| **gorilla/websocket** | v1.5.4 | WebSocket 支持（终端、任务日志实时推送） |
| **gobwas/ws** | v1.4.0 | 轻量 WebSocket 实现（WebSocket 握手鉴权） |
| **pkg/sftp** | v1.13.10 | SFTP 文件传输（主机文件操作） |
| **golang.org/x/crypto** | v0.42.0 | SSH 加密支持（主机 SSH 连接） |
| **k8s.io/client-go** | v0.34.1 | Kubernetes 官方 Go 客户端 |
| **k8s.io/api + apimachinery** | v0.34.1 | K8s API 对象模型 |
| **k8s.io/metrics** | v0.34.1 | K8s Metrics API（资源监控） |
| **prometheus/client_golang** | v1.23.2 | Prometheus 指标采集与上报 |
| **bce-sdk-go (百度云)** | v0.9.245 | 百度云 BCC 主机资产对接 |
| **tencentcloud-sdk-go** | v1.1.35 | 腾讯云 CVM 主机资产对接 |
| **go-ldap** | v3.4.1 | LDAP 认证集成 |
| **mojocn/base64Captcha** | v1.3.8 | 图形验证码生成 |
| **swaggo/swag** | v1.16.6 | Swagger API 文档自动生成 |
| **sirupsen/logrus** | v1.9.3 | 结构化日志框架 |
| **lestrrat-go/file-rotatelogs** | v2.4.0 | 日志文件滚动切割 |
| **shirou/gopsutil** | v3.24.5 | 主机系统信息采集（CPU/内存/磁盘） |
| **xuri/excelize** | v2.9.1 | Excel 文件读写（主机批量导入） |
| **spf13/viper** | v1.8.1 | 配置文件解析（YAML） |
| **dnsjia/luban** | v1.0.2 | K8s 辅助工具库 |

### 前端技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| **Vue 3** | v3.5.17 | 前端核心框架，Composition API |
| **Vue Router** | v4.5.1 | SPA 前端路由管理 |
| **Vuex** | v4.0.2 | 全局状态管理 |
| **Element Plus** | v2.10.2 | 企业级 UI 组件库（主要视觉框架） |
| **@element-plus/icons-vue** | v2.3.1 | 图标库 |
| **Axios** | v0.27.2 | HTTP 请求封装，对接后端 REST API |
| **ECharts** | v5.1.2 | 数据可视化图表（监控仪表盘） |
| **@xterm/xterm** | v5.5.0 | Web 终端模拟器（SSH Terminal） |
| **@xterm/addon-attach** | v0.11.0 | xterm WebSocket 附加器（连接后端 WebSocket） |
| **@xterm/addon-fit** | v0.10.0 | xterm 自适应尺寸 |
| **@antv/x6** | v2.18.1 | 流程图/拓扑图引擎（流程可视化） |
| **@logicflow/core** | v2.1.0 | 逻辑流程图引擎 |
| **@guolao/vue-monaco-editor** | v1.6.0 | Monaco 代码编辑器（配置/脚本编辑） |
| **highlight.js** | v11.11.1 | 代码语法高亮 |
| **marked** | v17.0.5 | Markdown 渲染 |
| **js-yaml** | v4.1.0 | YAML 解析（K8s 资源配置编辑） |
| **moment** | v2.30.1 | 时间格式化 |
| **vue3-treeselect** | v0.1.10 | 树形下拉选择器（部门/菜单树） |
| **Vue CLI** | ~5.0.0 | 构建工具链（webpack） |
| **Less** | v4.1.1 | CSS 预处理器 |

### 基础设施

| 组件 | 版本 | 用途 |
|------|------|------|
| **Docker + Docker Compose** | — | 容器化部署，一键启动所有服务 |
| **Nginx** | — | 前端静态资源服务 + 反向代理（将 `/api` 请求转发到后端） |
| **MySQL** | 8.0.33 | 持久化存储 |
| **Redis** | 7.0 Alpine | 缓存、分布式锁 |
| **Prometheus** | v2.47.0 | 时序指标采集与存储 |
| **Pushgateway** | v1.9.0 | Agent 主动推送指标到 Prometheus |
| **Kubernetes（可选）** | — | k8s 目录提供 K8s 部署清单，支持云原生部署 |

---

## 三、系统分层架构

### 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                       浏览器客户端                                │
│              Vue3 + Element Plus + Vuex + Vue Router             │
└──────────────────────────┬──────────────────────────────────────┘
                           │ HTTP / WebSocket
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Nginx（前端静态托管 + 反向代理）                 │
│   :80 → 静态 HTML/JS/CSS     /api/* → devops-api:8000           │
└──────────────────────────┬──────────────────────────────────────┘
                           │ REST API / WebSocket
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                  Go 后端（Gin）:8000                              │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────────────┐ │
│  │  Middleware  │  │    Router    │  │     Scheduler          │ │
│  │  · Cors      │  │  /api/v1/*  │  │  · robfig/cron         │ │
│  │  · JWT Auth  │  │  /ws/*      │  │  · 全局调度器           │ │
│  │  · OpLog     │  │  /swagger/* │  │  · 任务队列             │ │
│  └─────────────┘  └──────┬───────┘  └────────────────────────┘ │
│                           │                                      │
│  ┌────────────────────────▼─────────────────────────────────┐   │
│  │                   Controller 层                           │   │
│  │  system | cmdb | k8s | monitor | task | app |            │   │
│  │  configcenter | dashboard | tool                         │   │
│  └────────────────────────┬─────────────────────────────────┘   │
│                           │ 调用                                  │
│  ┌────────────────────────▼─────────────────────────────────┐   │
│  │                   Service 层                              │   │
│  │  业务逻辑、外部系统对接（Jenkins/云SDK/K8s/LDAP/SSH）      │   │
│  └────────────────────────┬─────────────────────────────────┘   │
│                           │ 调用                                  │
│  ┌────────────────────────▼─────────────────────────────────┐   │
│  │                    DAO 层                                 │   │
│  │  GORM 封装的数据库读写操作                                  │   │
│  └────────────────────────┬─────────────────────────────────┘   │
│                           │                                      │
│  ┌────────────────────────▼─────────────────────────────────┐   │
│  │                    Model 层                               │   │
│  │  GORM 结构体定义（对应数据库表）                             │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │             pkg（基础设施包）                              │   │
│  │   db | redis | jwt | log | cache | util | valid          │   │
│  └──────────────────────────────────────────────────────────┘   │
└───────────────────┬───────────────────────────────┬─────────────┘
                    │                               │
          ┌─────────▼──────────┐       ┌────────────▼───────────┐
          │   MySQL :3306      │       │   Redis :6379           │
          │   业务数据持久化     │       │   Token/缓存/会话        │
          └────────────────────┘       └────────────────────────┘
                    │
          ┌─────────▼──────────┐
          │ Prometheus :9090   │
          │ Pushgateway :9091  │
          │ 监控指标存储与采集   │
          └────────────────────┘
```

### 各层职责详述

#### 1. 路由层（`router/`）
- **职责**：URL 路由注册与分组，将请求分发到对应 Controller。
- **实现**：按业务域拆分子路由文件（`router/cmdb/`、`router/k8s/` 等），统一在 `router/router.go` 的 `InitRouter()` 中组装。
- **特点**：所有业务路由统一挂载在 `/api/v1` 前缀下，WebSocket 路由单独注册到根路径（`/ws/*`）。

#### 2. 中间件层（`middleware/`）
- **职责**：横切关注点处理。
- **组成**：
  - `cors.go` — 跨域处理（CORS），支持前后端分离部署
  - `authMiddleware.go` — JWT 鉴权，同时兼容 HTTP Bearer Token、WebSocket URL Token、SSE Query Token 三种认证方式
  - `logMiddleware.go` — 操作日志记录，将用户操作异步写入数据库审计表

#### 3. Controller 层（`api/*/controller/`）
- **职责**：接收 HTTP 请求、参数绑定与校验、调用 Service、封装 HTTP 响应。
- **规范**：使用 `common/result` 包统一响应格式（`{code, message, data}`）。
- **WebSocket**：终端（SSH、K8s Terminal）、任务日志实时推送均在 Controller 层建立 WebSocket 连接并维护会话。

#### 4. Service 层（`api/*/service/`）
- **职责**：核心业务逻辑、外部系统调用、数据组合处理。
- **关键能力**：
  - SSH/SFTP 连接管理（主机终端、远程执行）
  - K8s client-go API 调用（集群资源管理）
  - 云厂商 SDK 调用（阿里云、腾讯云、百度云主机同步）
  - LDAP 认证对接
  - Jenkins API 对接（应用发布）
  - Prometheus HTTP API 查询（监控数据）

#### 5. DAO 层（`api/*/dao/`）
- **职责**：数据访问对象，封装所有数据库 CRUD 操作。
- **实现**：基于 GORM，每个业务域有独立的 DAO 文件，与 Service 层解耦。

#### 6. Model 层（`api/*/model/`）
- **职责**：定义数据库表结构对应的 Go 结构体，GORM 标签声明字段约束。

#### 7. 调度器（`scheduler/`）
- **职责**：定时任务的生命周期管理。
- **组成**：
  - `manager.go` — 调度器管理器，统一启动/停止
  - `syncScheduler.go` — 数据同步定时任务（如云主机资产定时同步）
  - `globalScheduler.go`（在 task/service 中）— 全局 Cron 调度器
  - `taskQueue.go` — 任务队列，管理并发任务执行
  - `taskJob.go` — 单个任务 Job 封装

#### 8. pkg（公共基础包）
- `pkg/db/` — GORM 初始化与连接池配置
- `pkg/redis/` — Redis 客户端初始化
- `pkg/jwt/` — JWT Token 生成与验证
- `pkg/log/` — Logrus 日志初始化、自定义 Gin Logger
- `pkg/cache/` — Redis 缓存操作封装
- `pkg/util/` — 通用工具函数

#### 9. 前端层（`web/src/`）

```
web/src/
├── api/          # 按业务域拆分的 Axios 请求封装
│   ├── index.js  # Axios 实例（设置 baseURL、Token 拦截器）
│   ├── cmdb.js   # CMDB 相关接口
│   ├── k8s.js    # K8s 相关接口
│   ├── system.js # 系统管理接口
│   └── ...
├── views/        # 页面组件（与后端模块一一对应）
│   ├── dashboard/   仪表盘
│   ├── cmdb/        主机/数据库资产
│   ├── K8s/         K8s集群管理
│   ├── app/         应用服务管理
│   ├── monitor/     监控告警
│   ├── task/        任务中心
│   ├── configcenter/配置中心
│   ├── system/      系统管理
│   ├── Tools/       运维工具
│   └── work/        工单
├── components/   # 公共组件（CodeEditor、HostSelector、Tags等）
├── router/       # Vue Router 路由配置
├── store/        # Vuex 状态管理（用户信息、权限菜单）
├── permission/   # 路由守卫（前端权限控制）
└── utils/        # 工具函数（request封装、权限判断、SSE管理等）
```

**前端数据流**：用户操作 → Vue 组件调用 `api/*.js` → Axios 请求携带 JWT Token → Nginx 反代到后端 → 响应数据更新 Vuex Store → 组件响应式渲染。

---

## 四、核心模块清单

| 模块名 | 后端路径 | 前端路径 | 职责说明 |
|--------|---------|---------|---------|
| **系统管理（system）** | `api/api/system/` | `views/system/` | 用户、角色（RBAC）、菜单、部门、岗位、LDAP 认证管理，是整个平台的权限基础 |
| **CMDB 资产管理（cmdb）** | `api/api/cmdb/` | `views/cmdb/` | 主机资产管理（含阿里云/腾讯云/百度云同步）、主机分组、主机 SSH 终端、五类数据库（MySQL/PgSQL/Redis/ES/MongoDB）管理与 SQL 执行审计 |
| **K8s 集群管理（k8s）** | `api/api/k8s/` | `views/K8s/` | 多集群注册管理、节点管理、工作负载（Pod/Deployment 等）、配置（ConfigMap/Secret）、存储（PV/PVC）、网络（Service/Ingress）管理，以及 Web 终端直连 Pod |
| **应用服务管理（app）** | `api/api/app/` | `views/app/` | 应用发布（对接 Jenkins）、工单式上线流程、应用列表与发布历史管理 |
| **监控中心（monitor）** | `api/api/monitor/` | `views/monitor/` | 域名监控、主机基础资源监控（Prometheus 数据查询）、Agent 心跳管理、故障管理、告警 Webhook 接收 |
| **任务中心（task）** | `api/api/task/` | `views/task/` | 定时任务调度（Shell/Python 脚本）、Ansible Playbook 任务执行、任务日志实时推送（WebSocket/SSE） |
| **配置中心（configcenter）** | `api/api/configcenter/` | `views/configcenter/` | 主机 SSH 凭据管理、云资源密钥（AK/SK）管理、通用账号管理 |
| **仪表盘（dashboard）** | `api/api/dashboard/` | `views/dashboard/` | 平台总览：资产统计、告警速览、发布概况、快捷导航 |
| **运维工具（tool）** | `api/api/tool/` | `views/Tools/` | 常用运维资源安装向导、Agent 监控工具一键安装 |
| **工单（work）** | —（复用 app 模块） | `views/work/` | 运营工单流转（发起、处理、关闭） |
| **调度器（scheduler）** | `api/scheduler/` | — | 后台定时任务的注册、触发与生命周期管理，不对外暴露 HTTP 接口 |
| **公共包（pkg）** | `api/pkg/` | — | db/redis/jwt/log 等基础设施初始化与封装，被所有业务模块依赖 |
| **公共工具（common）** | `api/common/` | — | 配置加载（Viper）、全局 DB/Redis 实例、常量定义、响应格式、模板、工具函数 |

---

## 五、模块依赖关系

### 后端模块依赖图

```
                         ┌─────────────┐
                         │   main.go   │ ← 程序入口
                         └──────┬──────┘
              ┌─────────────────┼──────────────────┐
              ▼                 ▼                  ▼
        ┌──────────┐     ┌──────────┐      ┌────────────┐
        │  router  │     │scheduler │      │  pkg/*     │
        └────┬─────┘     └──────────┘      │db/redis/jwt│
             │                             │log/cache   │
             ▼                             └─────┬──────┘
     ┌───────────────────────────────────────────┘
     │  按业务域注册路由                           │
     ▼                                            │（被所有层依赖）
┌─────────────────────────────────────────────────▼──────┐
│  Controller 层（各模块 controller/）                     │
│  system | cmdb | k8s | monitor | task | app |          │
│  configcenter | dashboard | tool                        │
└─────────────────────┬───────────────────────────────────┘
                      │ 调用
                      ▼
┌────────────────────────────────────────────────────────┐
│  Service 层（各模块 service/）                           │
│  ┌───────────┐ ┌───────────┐ ┌──────────┐ ┌─────────┐ │
│  │  系统服务  │ │ CMDB服务  │ │  K8s服务 │ │ 任务服务 │ │
│  │(LDAP/RBAC)│ │(SSH/云SDK)│ │(client-go│ │(Cron/   │ │
│  └───────────┘ └───────────┘ │Prometheus│ │Ansible) │ │
│                               └──────────┘ └─────────┘ │
└─────────────────────┬───────────────────────────────────┘
                      │ 调用
                      ▼
┌────────────────────────────────────────────────────────┐
│  DAO 层（各模块 dao/）                                   │
│  GORM 数据库操作封装                                     │
└─────────────────────┬───────────────────────────────────┘
                      │ 依赖
                      ▼
┌────────────────────────────────────────────────────────┐
│  Model 层（各模块 model/）                               │
│  数据库表结构定义（Go Struct + GORM Tag）                 │
└────────────────────────────────────────────────────────┘
```

### 前端模块依赖图

```
┌─────────────────────────────────────────────────────┐
│  App.vue + main.js（应用入口）                       │
└──────────┬──────────────────────────────────────────┘
           │
    ┌──────┴────────────────────┐
    ▼                           ▼
┌────────┐               ┌──────────────┐
│ router │               │  store(Vuex) │
│路由守卫 │               │用户信息/权限  │
└───┬────┘               └──────┬───────┘
    │                           │
    ▼                           ▼
┌────────────────────────────────────────┐
│         views/*（页面组件）             │
│  引用 components/* 中的公共组件         │
└────────────────┬───────────────────────┘
                 │ 调用
                 ▼
┌────────────────────────────────────────┐
│     api/*.js（接口封装层）              │
│  基于 utils/request.js（Axios 实例）   │
│  自动注入 JWT Token、处理响应错误       │
└────────────────────────────────────────┘
                 │ HTTP
                 ▼
         后端 Gin API 服务
```

### 外部系统依赖

```
AutoOps 后端
├── → MySQL          （数据持久化）
├── → Redis          （缓存/令牌）
├── → Prometheus     （监控数据查询）
├── → Pushgateway    （Agent 指标推送接收）
├── → Kubernetes API （多集群管理）
├── → 阿里云/腾讯云/百度云 SDK   （主机资产同步）
├── → LDAP Server    （企业域账号认证）
└── → Jenkins API    （应用发布触发）
```

---

## 六、关键配置说明

### config.yaml（后端应用配置）

```yaml
# 服务监听配置
server:
  address: 0.0.0.0:8000    # 监听地址与端口，容器内暴露 8000
  model: debug              # Gin 运行模式（debug/release），生产环境应改为 release
  enableSwagger: true       # 是否开启 /swagger/index.html，生产建议关闭
  publicUrl: "http://..."   # 对外公网地址，用于飞书 OAuth 回调，需与飞书开放平台配置一致

# 数据库连接池配置
db:
  dialects: mysql
  host/port/db/username/password  # 连接参数
  maxIdle: 50     # 最大空闲连接数，控制连接池下限
  maxOpen: 150    # 最大打开连接数，并发量上限，超出请求会排队等待

# Redis 配置
redis:
  address: 127.0.0.1:6379   # Docker 内应改为服务名（redis:6379）
  password: "123456"         # 认证密码

# 上传文件路径
imageSettings:
  uploadDir: ./upload/       # 文件上传根目录，Nginx 和后端共享此目录

# 日志配置
log:
  path: ./log                # 日志输出目录
  name: sys                  # 日志文件前缀
  model: console             # console = 输出到控制台；file = 写入文件（生产推荐 file）

# 监控配置
monitor:
  prometheus.url: "..."      # Prometheus HTTP API 地址，后端查询监控数据用
  pushgateway.url: "..."     # Pushgateway 地址，Agent 推送指标的中转站
  agent.heartbeat_server_url # Agent 上报心跳的回调地址
  agent.heartbeat_token      # Agent 鉴权 Token（防止未授权 Agent 上报）
  webhook.token              # 告警 Webhook 鉴权 Token
```

### docker-compose.yml（容器编排配置）

| 服务 | 端口映射 | 关键说明 |
|------|---------|---------|
| **devops-mysql** | `${MYSQL_PORT:-3306}:3306` | 使用 `.env` 文件注入环境变量，默认密码 `devops@2025`；挂载 `mysql/devops001.sql` 和 `init.sql` 作为初始化脚本（按文件名字母顺序执行） |
| **devops-redis** | `${REDIS_PORT:-6379}:6379` | 开启 AOF 持久化（`--appendonly yes`），密码 `zhangfan@123` |
| **devops-pushgateway** | `${PUSHGATEWAY_PORT:-9091}:9091` | 用于接收主机 Agent 主动推送的 Prometheus 指标 |
| **devops-prometheus** | `${PROMETHEUS_PORT:-9090}:9090` | 数据保留 15 天（`--storage.tsdb.retention.time=15d`），启用 Lifecycle API（`--web.enable-lifecycle`，支持热加载配置） |
| **devops-api** | `${API_PORT:-8000}:8000` | 通过环境变量覆盖 DB/Redis 连接参数；挂载 `./api/ssh_keys:/root/.ssh:ro`（只读挂载 SSH 密钥，实现 SSH 连接主机）；依赖 mysql/redis/prometheus 健康检查通过后才启动 |
| **devops-web** | `${WEB_PORT:-8080}:80` | Nginx 容器，挂载 `devops.conf` 配置文件；与 api 容器共享 `./api/upload` 卷，使 Nginx 能直接伺服上传的静态文件 |

**服务启动依赖顺序**：
```
pushgateway（健康） → prometheus（健康）
mysql（健康）+ redis（健康）+ prometheus（健康） → devops-api（启动）→ devops-web（启动）
```

**网络隔离**：所有容器处于独立的 `devops-network`（桥接网络，子网 `172.20.0.0/16`），容器间通过服务名通信，对外仅暴露必要端口。

**环境变量注入**：通过 `docker/.env` 文件统一管理密码、端口等敏感配置，`docker-compose.yml` 中使用 `${VAR:-default}` 语法提供默认值，便于不同环境（开发/生产）快速切换。

---

*文档基于源码静态分析生成，如有代码变更请同步更新本文档。*
