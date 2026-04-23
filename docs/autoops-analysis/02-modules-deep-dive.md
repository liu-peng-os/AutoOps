# 天枢 AutoOps — 模块深度解读

> 文档版本：v1.0 | 分析日期：2026-04-23 | 基于实际代码分析

---

## 目录

1. [CMDB 资产管理模块](#一cmdb-资产管理模块)
2. [Kubernetes 集群管理模块](#二kubernetes-集群管理模块)
3. [应用服务管理模块（app）](#三应用服务管理模块)
4. [任务调度模块（task）](#四任务调度模块)
5. [监控告警模块（monitor）](#五监控告警模块)
6. [配置中心模块（configcenter）](#六配置中心模块)
7. [用户权限模块（system）](#七用户权限模块)
8. [仪表盘模块（dashboard）](#八仪表盘模块)
9. [运维工具模块（tool）](#九运维工具模块)
10. [前端视图模块（web/src/views）](#十前端视图模块)

---

## 一、CMDB 资产管理模块

### 模块职责

CMDB 模块是平台的资产数据中枢，负责管理企业内所有主机（自建 + 阿里云 / 腾讯云 / 百度云）、主机分组、五种数据库（MySQL / PostgreSQL / Redis / MongoDB / Elasticsearch）的元信息录入与维护；同时提供 Web SSH 终端连接、SFTP 文件上传、SQL 在线执行及操作审计日志功能。

### 目录结构说明

```
api/api/cmdb/
├── controller/
│   ├── CmdbHostcloud.go    # 云主机批量导入（阿里云/腾讯云/百度云）
│   ├── cmdbGroup.go        # 主机分组 CRUD
│   ├── cmdbHost.go         # 主机 CRUD、Excel 批量导入、SSH 信息同步
│   ├── cmdbHostSSH.go      # WebSocket SSH 终端、命令执行、SFTP 文件上传
│   ├── cmdbSQL.go          # 数据库资产 CRUD
│   ├── cmdbSQLRecord.go    # SQL 在线执行（Select/Insert/Update/Delete）
│   └── cmdbSqlLog.go       # SQL 操作审计日志查询
├── service/                # 同名文件，实现上述 Controller 的业务逻辑
├── dao/                    # 同名文件，GORM 数据库操作
└── model/
    ├── cmdbGroup.go        # CmdbGroup / CmdbGroupHost（分组树 + 分组主机关联）
    ├── cmdbGroupHost.go    # 分组-主机关联表
    ├── cmdbHost.go         # CmdbHost（主机表）+ CreateCmdbHostDto / ExcelHostTemplate 等
    ├── cmdbSQL.go          # CmdbSQL（数据库资产表）
    └── cmdbSQLRecord.go    # SQL 执行记录 + SQL 日志
```

### 核心数据模型

| Struct | 表名 | 说明 |
|--------|------|------|
| `CmdbHost` | `cmdb_host` | 主机核心信息：SSH 连接参数（IP/端口/用户名/凭据ID）、云厂商标识（1=自建/2=阿里云/3=腾讯云）、运行状态（1=认证成功/2=未认证/3=失败）、CPU/内存/磁盘信息 |
| `CmdbGroup` | `cmdb_group` | 主机分组树（支持多级嵌套，ParentID 自关联） |
| `CmdbGroupHost` | `cmdb_group_host` | 主机与分组的多对一关联 |
| `CmdbSQL` | `cmdb_sql` | 数据库资产：Type（1=MySQL/2=PostgreSQL/3=Redis/4=MongoDB/5=ES）、AccountID（凭据）、GroupID（业务组） |
| `ExcelHostTemplate` | — | Excel 导入时的行数据结构（别名/SSH地址/端口/用户/备注） |

**DTO 设计特点**：`CreateCmdbHostDto` 只需填写 SSH 连接三要素（IP/端口/用户名/凭据ID），主机的 OS/CPU/内存/磁盘等信息通过 `SyncHostInfo` 接口自动 SSH 拉取回填。

### 对外暴露的主要接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/cmdb/grouplist` | 获取所有分组树 |
| GET | `/api/v1/cmdb/grouplistwithhosts` | 获取分组树及各组主机 |
| POST | `/api/v1/cmdb/hostcreate` | 新增主机（手动录入） |
| PUT | `/api/v1/cmdb/hostupdate` | 更新主机信息 |
| DELETE | `/api/v1/cmdb/hostdelete` | 删除主机 |
| GET | `/api/v1/cmdb/hostlist` | 分页查询主机列表 |
| GET | `/api/v1/cmdb/hostgroup?groupId=` | 按分组查询主机 |
| GET | `/api/v1/cmdb/hostbyip?ip=` | 按 IP 查询（同时匹配内/外/SSH IP） |
| POST | `/api/v1/cmdb/hostimport` | Excel 批量导入主机（multipart） |
| POST | `/api/v1/cmdb/hostsync` | 触发主机信息 SSH 自动同步 |
| POST | `/api/v1/cmdb/hostcloudcreatealiyun` | 从阿里云批量导入主机 |
| POST | `/api/v1/cmdb/hostcloudcreatetencent` | 从腾讯云批量导入主机 |
| POST | `/api/v1/cmdb/hostcloudcreatebaidu` | 从百度云批量导入主机 |
| GET | `/api/v1/cmdb/hostssh/connect/:id` | **WebSocket** SSH 终端连接 |
| GET | `/api/v1/cmdb/hostssh/command/:id` | SSH 执行单条命令 |
| POST | `/api/v1/cmdb/hostssh/upload/:id` | SFTP 上传文件到远程主机 |
| POST/PUT/DELETE/GET | `/api/v1/cmdb/database*` | 数据库资产 CRUD |
| POST | `/api/v1/cmdb/sql/execute` | 在线执行原生 SQL |
| GET | `/api/v1/cmdb/sqlLog/list` | 查询 SQL 操作审计日志 |

### 内部数据流

**主机创建流程：**
```
前端 POST /cmdb/hostcreate
  → AuthMiddleware（JWT 验证）
  → LogMiddleware（记录操作日志）
  → CmdbHostController.CreateCmdbHost()
    → 参数绑定 & 校验 (ShouldBindJSON)
    → CmdbHostService.CreateCmdbHost()
      → CmdbHostDAO.Create(host) → MySQL [cmdb_host]
    → result.Success(ctx, data)
```

**WebSocket SSH 终端流程：**
```
前端 WebSocket /cmdb/hostssh/connect/:id?token=xxx
  → CmdbHostSSHController.ConnectTerminal()
    → jwt.ValidateToken(token)           # 鉴权
    → websocket.Upgrader.Upgrade()       # HTTP升级为WS
    → CmdbHostSSHService.ConnectTerminal(id)
      → 查询主机信息 (cmdb_host)
      → 查询 SSH 凭据 (config_ecsauth)
      → golang.org/x/crypto/ssh.Dial()  # 建立 SSH 会话
    → webSSH.Connect(wsConn)             # WS ↔ SSH 双向 IO 桥接
    → select{} 保持协程阻塞（直到 WS 关闭）
```

**SQL 执行审计流程：**
```
前端 POST /cmdb/sql/execute
  → CmdbSQLRecordController.ExecuteSQL()
    → 解析目标数据库信息
    → 从 config_account 获取并解密数据库凭据
    → 建立数据库连接并执行 SQL
    → 将操作记录异步写入 cmdb_sql_log 审计表
    → 返回执行结果
```

### 与其他模块的依赖关系

- **依赖 configcenter 模块**：SSH 连接需要从 `config_ecsauth` 表读取主机凭据（密码/私钥）；数据库连接需要从 `config_account` 表读取数据库账密（AES 加密存储）。
- **被 task 模块依赖**：任务执行时通过 `HostIDs` 字段引用 `cmdb_host` 表中的主机，SSH 执行脚本也依赖 CMDB 的主机 SSH 服务。
- **被 app 模块依赖**：应用模型 `Application.Hosts` 字段存储关联的主机 ID（`cmdb_host.id`）。
- **被 dashboard 模块依赖**：仪表盘汇总查询 CMDB 的主机在线数量统计。
- **被 monitor 模块依赖**：监控 Agent 关联 `cmdb_host` 的 `HostID`。

---

## 二、Kubernetes 集群管理模块

### 模块职责

K8s 模块通过存储 kubeconfig 凭据，利用 `k8s.io/client-go` 实现对多个 Kubernetes 集群的统一纳管，覆盖集群注册/同步、节点管理（封锁/驱逐/污点/标签）、工作负载全生命周期管理（Deployment/StatefulSet/DaemonSet/Job/CronJob）、配置管理（ConfigMap/Secret）、存储管理（PV/PVC/StorageClass）、网络管理（Service/Ingress）、Web 终端直连 Pod，以及通过 `k8s.io/metrics` 采集实时资源监控数据。

### 目录结构说明

```
api/api/k8s/
├── controller/
│   ├── kubeCluster.go   # 集群 CRUD + 状态同步 + 详情
│   ├── k8snodes.go      # 节点列表/详情/封锁/驱逐/污点/标签
│   ├── k8sworkload.go   # 工作负载（Deployment等）全量 CRUD + 扩缩容 + 回滚
│   ├── k8sconfig.go     # ConfigMap / Secret 管理
│   ├── k8sstorage.go    # PV / PVC / StorageClass 管理
│   ├── k8sservice.go    # Service 管理
│   ├── k8singress.go    # Ingress 管理
│   ├── k8snamespace.go  # Namespace + ResourceQuota + LimitRange 管理
│   ├── k8sevents.go     # 集群事件查询
│   └── k8sterminal.go   # WebSocket Pod 终端连接
├── service/             # 同名文件，client-go API 调用实现
├── dao/
│   └── kubeCluster.go   # 仅集群表持久化（其余数据实时从 K8s API 获取）
└── model/
    └── kubeCluster.go   # 所有 K8s 相关结构体（极长，含集群/节点/工作负载/网络/存储等）
```

### 核心数据模型

| Struct | 存储位置 | 说明 |
|--------|---------|------|
| `KubeCluster` | MySQL `k8s_cluster` | 集群注册信息：Name/Version/Status（1=创建中/2=运行中/3=离线）/`Credential`（kubeconfig 全文）/NodeCount 等 |
| `NodeInfo` / `K8sNode` | 实时从 K8s API | 节点状态、IP、资源容量、污点等 |
| `K8sWorkload` / `K8sWorkloadDetail` | 实时从 K8s API | 工作负载基础信息 + 详情（含 Pod 列表、Events、完整 Spec） |
| `K8sPodInfo` / `K8sPodDetail` | 实时从 K8s API | Pod 状态、容器信息、日志、资源使用 |
| `K8sService` / `K8sIngress` | 实时从 K8s API | 网络对象，含端口/选择器/规则/TLS |
| `K8sPersistentVolume` / `K8sPersistentVolumeClaim` / `K8sStorageClass` | 实时从 K8s API | 存储对象，含容量/访问模式/存储源（HostPath/NFS/CSI/Local/AWS EBS等） |
| `K8sConfigMap` / `K8sSecret` | 实时从 K8s API | 配置对象（Secret 数据 Base64 编码） |
| `ClusterDetailResponse` | 聚合 | 集群全景信息：节点+工作负载统计+组件+网络+监控+运行时 |

**架构特点**：除集群注册元数据存于 MySQL，所有 Kubernetes 资源（节点/Pod/Service 等）均**实时通过 client-go 从 K8s API Server 查询**，不落库，保证数据实时性。

### 对外暴露的主要接口

| 分类 | 主要接口 |
|------|---------|
| 集群管理 | `POST/GET/PUT/DELETE /api/v1/k8s/cluster`；`GET /cluster/{id}/detail`；`POST /cluster/{id}/sync` |
| 节点管理 | `GET /k8s/{clusterID}/nodes`；`POST /nodes/{name}/cordon`（封锁）；`POST /nodes/{name}/drain`（驱逐）；污点/标签增删 |
| 工作负载 | `GET/POST/PUT/DELETE /k8s/{clusterID}/workloads/{namespace}/{type}`；扩缩容、重启、回滚（`/rollback`）、暂停/恢复 |
| 配置管理 | ConfigMap / Secret 的 CRUD + YAML 获取/更新 |
| 存储管理 | PV / PVC / StorageClass 的 CRUD |
| 网络管理 | Service / Ingress 的 CRUD |
| Namespace | Namespace CRUD + ResourceQuota + LimitRange 管理 |
| Pod 终端 | WebSocket `/k8s/{clusterID}/pods/{namespace}/{pod}/terminal` |
| 监控指标 | `GET /k8s/{clusterID}/metrics/nodes`；`/metrics/pods`；`/metrics/namespaces` |

### 内部数据流

**集群注册+同步流程：**
```
前端 POST /k8s/cluster {name, kubeconfig}
  → KubeClusterController.CreateCluster()
    → KubeClusterService.CreateCluster()
      → 解析 kubeconfig 内容
      → 使用 client-go 建立临时连接，验证 kubeconfig 有效性
      → 将集群信息写入 MySQL [k8s_cluster]
      → 触发首次同步：通过 K8s API 获取版本/节点数 → 更新数据库

POST /k8s/cluster/{id}/sync
  → KubeClusterService.SyncCluster(id)
    → 从 MySQL 取 kubeconfig
    → client-go 查询集群版本、节点列表
    → 更新 MySQL 中的 NodeCount/ReadyNodes/LastSyncAt
```

**工作负载操作数据流（以 Deployment 扩缩容为例）：**
```
前端 PUT /k8s/{clusterID}/workloads/{ns}/Deployment/{name}/scale
  → K8sWorkloadController.ScaleWorkload()
    → KubeClusterService.GetCluster(clusterID) → 获取 kubeconfig
    → client-go.NewForConfig(kubeconfig)       → 创建 K8s 客户端
    → client.AppsV1().Deployments(ns).UpdateScale()
    → result.Success(ctx, {新副本数})
```

### 与其他模块的依赖关系

- **依赖 CMDB 模块**：集群自建时通过 `NodeConfig.MasterHostIDs` / `WorkerHostIDs` 引用 CMDB 主机 ID，在目标主机上执行 K8s 安装脚本（通过 SSH 服务）。
- **被 task 模块依赖**：任务模块中的 K8s 任务类型（`CreateK8sTask`）通过 cluster ID 获取凭据后执行 K8s 相关操作。
- **被 dashboard 模块依赖**：仪表盘展示 K8s 集群总数/健康数统计。

---

## 三、应用服务管理模块

### 模块职责

应用服务管理模块以「应用」为核心资产单元，管理企业内所有软件服务的基础信息（代码仓库/负责人/技术栈/关联资源）及各环境（dev/test/prod）的 Jenkins 发布配置；提供标准化发布入口（触发 Jenkins 构建）、快速批量发布（支持串行/并行模式）以及业务线服务树视图。

### 目录结构说明

```
api/api/app/
├── controller/
│   ├── application.go  # 应用 CRUD + 服务树 + 快速发布管理
│   └── jenkins.go      # Jenkins 服务器管理 + 任务操作 + 构建日志
├── service/
│   ├── application.go  # 应用业务逻辑、发布流程协调
│   └── jenkins.go      # Jenkins HTTP API 调用封装
├── dao/
│   └── application.go  # 应用表 GORM 操作
└── model/
    ├── application.go  # Application / JenkinsEnv / QuickDeployment / QuickDeploymentTask
    └── jenkins.go      # JenkinsJob / JenkinsBuild / JenkinsServerInfo 等 Jenkins 数据模型
```

### 核心数据模型

| Struct | 表名 | 说明 |
|--------|------|------|
| `Application` | `app_application` | 应用主体：Name/Code（唯一编码）/BusinessGroupID/BusinessDeptID/DevOwners/TestOwners/OpsOwners（JSON 数组）/Domains/Hosts/Databases/OtherRes（均 JSON 存储）/StartCommand/StopCommand/HealthAPI |
| `JenkinsEnv` | `app_jenkins_env` | 应用-环境-Jenkins 绑定：AppID + EnvName（dev/test/prod）+ JenkinsServerID（关联 config_account）+ JobName |
| `QuickDeployment` | `quick_deployments` | 快速发布批次：Title/ExecutionMode（1=并行/2=串行）/Status/TaskCount/StartTime/EndTime |
| `QuickDeploymentTask` | `quick_deployment_tasks` | 快速发布子任务：ApplicationID + Environment + JenkinsEnvID + BuildNumber + Status + ExecuteOrder |
| `JenkinsJob` | — | Jenkins 侧任务对象（实时从 Jenkins API 获取） |
| `JenkinsBuild` | — | Jenkins 构建对象，含 Result/Building/Duration/ChangeSet |

**关键设计**：`JenkinsAccountType = 4`，Jenkins 服务器账密复用了 `configcenter` 的 `AccountAuth` 表（Type=4），不单独建表。

### 对外暴露的主要接口

| 分类 | 主要接口 |
|------|---------|
| 应用管理 | `GET/POST/PUT/DELETE /api/v1/app/application`；应用列表/详情/服务树 |
| Jenkins 服务器 | `GET /app/jenkins/servers`（列出 Jenkins 服务器）；测试连接 |
| 发布操作 | `POST /app/deploy`（触发单应用发布）；`GET /app/deploy/status`（查询构建状态） |
| 快速发布 | `POST /app/quickdeploy`（创建批次）；`POST /app/quickdeploy/execute`（执行）；`GET /app/quickdeploy/list` |
| 构建日志 | `GET /app/jenkins/build/log`（轮询日志）；SSE 实时推送 |
| 环境配置 | `GET/POST/PUT/DELETE /app/jenkins/env` |

### 内部数据流

**快速批量发布流程：**
```
前端 POST /app/quickdeploy {title, applications:[{appId, environment}], executionMode}
  → ApplicationController.CreateQuickDeployment()
    → 校验每个应用-环境是否有有效的 JenkinsEnv 配置
    → 创建 QuickDeployment 记录（Status=待发布）
    → 创建各子任务 QuickDeploymentTask（按数组顺序设 ExecuteOrder）

前端 POST /app/quickdeploy/execute {deploymentId, executionMode}
  → ApplicationController.ExecuteQuickDeployment()
    → ApplicationService.ExecuteQuickDeployment()
      ┌─ 并行模式（1）：goroutine 并发触发所有子任务的 Jenkins Build
      └─ 串行模式（2）：按 ExecuteOrder 顺序逐个触发，等待上一个完成
        → JenkinsService.StartBuild(serverID, jobName, params)
          → Jenkins HTTP API POST /job/{name}/build
          → 轮询 Jenkins 构建状态（或 SSE 推送）
          → 更新 QuickDeploymentTask.Status / BuildNumber / Duration
      → 全部完成后更新 QuickDeployment.Status
```

### 与其他模块的依赖关系

- **依赖 configcenter 模块**：从 `config_account` 表（Type=4）读取 Jenkins 服务器地址/账密。
- **依赖 CMDB 模块**：应用关联的 `Hosts`（主机 ID）、`Databases`（数据库 ID）均引用 CMDB 的数据。
- **依赖 system 模块**：`BusinessGroupID` 关联 CMDB 分组，`BusinessDeptID` 关联 `sys_dept`（部门表）；`DevOwners/OpsOwners` 关联 `sys_admin`（用户表）。
- **被 dashboard 模块依赖**：仪表盘展示发布统计（总次数/成功率）。

---

## 四、任务调度模块

### 模块职责

任务调度模块提供两类任务的统一调度引擎：**普通脚本任务**（Shell/Python，可立即执行或按 Cron 定时）和 **Ansible Playbook 任务**（支持多主机批量自动化运维）。核心能力包括：任务模板管理、任务创建与执行、任务队列（基于 Redis 三级优先级队列 + Worker Pool）、全局 Cron 调度器、任务日志实时推送（SSE/WebSocket 双通道），以及调度器运行状态监控。

### 目录结构说明

```
api/api/task/
├── controller/
│   ├── taskJob.go      # 任务 CRUD + 列表
│   ├── taskwork.go     # 任务执行（启动/停止/查状态/查日志）
│   ├── taskTemplate.go # 任务模板 CRUD
│   ├── taskansible.go  # Ansible 任务 CRUD + 启动 + 日志（SSE）
│   ├── taskMonitor.go  # 队列指标 + 调度器统计 + 系统状态监控
│   └── websocket.go    # WebSocket 日志推送
├── service/
│   ├── globalScheduler.go # robfig/cron 全局调度器（定时任务注册/管理）
│   ├── taskQueue.go    # Redis 任务队列（高/普通/低/重试/失败 5 条队列）
│   ├── taskwork.go     # 普通任务执行逻辑（SSH 执行脚本）
│   ├── taskansible.go  # Ansible 任务执行（生成 Playbook 并 SSH 执行）
│   ├── taskJob.go      # 任务 DAO 封装
│   └── taskTemplate.go # 模板服务
├── dao/                # 同名文件，GORM 操作
└── model/
    ├── taskJob.go      # Task（任务）+ TaskWithDetails + 状态/类型常量
    ├── taskTemplate.go # TaskTemplate（模板）
    ├── taskansible.go  # TaskAnsible（Ansible 任务）+ TaskAnsibleWork（子任务）
    └── taskwork.go     # TaskWork（执行记录）
```

### 核心数据模型

| Struct | 表名 | 说明 |
|--------|------|------|
| `Task` | `task_job` | 任务主体：Type（1=立即执行/2=定时/3=Ansible）/Shell（模板ID列表）/HostIDs（主机ID列表）/CronExpr（Cron表达式）/Status（1-5）/ExecuteCount/NextRunTime |
| `TaskTemplate` | `task_template` | 脚本模板：Name/Type（1=Shell/2=Python）/Content（脚本内容） |
| `TaskAnsible` | `task_ansible` | Ansible 任务：Name/PlaybookContent/HostIDs/Status/LogFile |
| `TaskAnsibleWork` | `task_ansible_work` | Ansible 子任务（单次执行记录）：TaskID/Status/StartTime/EndTime/Log |
| `TaskWork` | `task_work` | 普通任务执行记录：TaskID/HostID/Status/Output/StartTime/Duration |
| `TaskMessage` | Redis | Redis 队列消息体：TaskWork + Priority + RetryCount + EnqueuedAt |

**队列架构**：
```
Redis 5 条队列（FIFO）：
  dodevops:task_queue:high    # 高优先级
  dodevops:task_queue:normal  # 普通优先级
  dodevops:task_queue:low     # 低优先级
  dodevops:task_queue:retry   # 待重试队列
  dodevops:task_queue:failed  # 最终失败队列

WorkerPool：
  - 全局并发限制（globalLimit semaphore）
  - 每主机并发限制（hostLimits map[uint]*semaphore.Weighted）
  - Worker 协程池（从队列消费任务，按优先级 high→normal→low 顺序）
```

### 对外暴露的主要接口

| 分类 | 路径（前缀 /api/v1） | 说明 |
|------|---------------------|------|
| 模板管理 | `POST/GET/PUT/DELETE /template/*` | 脚本模板的增删改查 |
| 任务 CRUD | `POST/GET/PUT/DELETE /task/*` | 任务的增删改查、按名称/类型/状态查询 |
| 任务执行 | `POST /taskjob/start`；`POST /taskjob/stop` | 立即执行任务/停止任务 |
| 任务日志 | `GET /taskjob/log`；`GET /task/execution-info` | 查询日志（轮询方式） |
| Ansible 任务 | `POST/GET/DELETE /task/ansible` | Ansible 任务管理 |
| Ansible 执行 | `POST /task/ansible/:id/start` | 启动 Ansible 任务 |
| Ansible 日志 | `GET /task/ansible/:id/log/:work_id` | **SSE** 实时日志流 |
| WS 日志 | `GET /ws/task/ansible/:id/log/:work_id` | **WebSocket** 实时日志流 |
| 调度监控 | `GET /task/monitor/queue/metrics`；`/scheduler/stats`；`/system/status` | 队列/调度器/系统状态监控 |
| 定时任务控制 | `POST /task/monitor/scheduled/pause`；`/resume`；`/reset` | 定时任务暂停/恢复/重置 |

### 内部数据流

**定时任务调度流程：**
```
应用启动 → main.go
  → scheduler.GetManager().Start()         # 启动调度器管理器
  → service.InitGlobalScheduler()          # 初始化 robfig/cron 实例
  → 从 MySQL 加载 Type=2 的 Task，注册 Cron Job

Cron 触发时：
  GlobalScheduler.jobFunc(taskID)
    → 从 MySQL 加载 Task 和 HostIDs
    → 为每台主机创建 TaskWork 记录
    → 将 TaskMessage 推入 Redis 队列（high/normal/low）

WorkerPool.consumeTask()
  → 从队列 LPOP 获取 TaskMessage
  → 获取 hostLimits[hostID] semaphore
  → goroutine: TaskWorkService.Execute(taskWork)
    → 从 cmdb_host 查询主机 SSH 信息
    → 从 config_ecsauth 查询凭据
    → SSH 连接，执行脚本，收集输出
    → 更新 TaskWork 状态/日志到 MySQL
    → 释放 semaphore
```

**Ansible 任务日志实时推送（SSE）：**
```
前端 GET /task/ansible/:id/log/:work_id（SSE长连接）
  → TaskAnsibleController.GetJobLog()
    → 设置响应头 Content-Type: text/event-stream
    → goroutine: 实时读取磁盘日志文件
    → 通过 SSE channel 逐行推送到前端
    → 任务完成时发送 EOF 事件
```

### 与其他模块的依赖关系

- **依赖 CMDB 模块**：通过 `HostIDs` 获取主机 SSH 连接参数。
- **依赖 configcenter 模块**：获取 SSH 凭据（EcsAuth）执行远程脚本。
- **依赖 Redis**：任务队列的消息存储/消费，队列指标持久化。
- **依赖 scheduler 包**：`api/scheduler/manager.go` 统一管理调度器生命周期。

---

## 五、监控告警模块

### 模块职责

监控告警模块通过 Agent 体系实现主机基础资源监控，Agent 将 CPU/内存/磁盘/网络/进程指标推送到 Pushgateway，后端从 Prometheus HTTP API 查询时序数据，向前端提供主机实时指标、历史趋势（时间序列）、进程 Top 排行、端口监听状态；同时管理 Agent 的部署（通过 SSH 安装）、心跳维护及故障管理。

### 目录结构说明

```
api/api/monitor/
├── controller/
│   ├── agent.go        # Agent 部署/查询/批量部署/心跳更新
│   └── monitor.go      # 主机监控数据查询（实时/历史/进程/端口）
├── service/
│   ├── agent.go        # Agent SSH 安装、状态管理
│   └── monitorService.go # Prometheus HTTP API 查询封装
├── dao/
│   └── agentDao.go     # monitor_agent 表 GORM 操作
└── model/
    ├── agent.go        # Agent 结构体（状态/安装进度常量）
    └── monitor.go      # HostMetrics / PrometheusQueryResult / AllMetricsHistory / ProcessInfo / PortInfo
```

### 核心数据模型

| Struct | 存储位置 | 说明 |
|--------|---------|------|
| `Agent` | MySQL `monitor_agent` | Agent 注册信息：HostID/Version/Status（1=部署中→2=失败→3=运行中→4=启动异常）/InstallPath/Port/PID/LastHeartbeat/InstallProgress（0→100） |
| `HostMetrics` | 实时从 Prometheus | CPU/内存/磁盘使用率 + 在线状态 |
| `AllMetricsHistory` | 实时从 Prometheus | 7 维度时序数据：CPU/内存/磁盘读写/网络收发/负载/进程数 |
| `ProcessInfo` / `TopProcessesResult` | 实时从 Prometheus | 进程 TOP5（CPU/内存） |
| `PortInfo` / `HostPortsResult` | 实时从 Prometheus | TCP 端口监听状态 |
| `PrometheusQueryResult` | 实时从 Prometheus | Prometheus API 原始响应（支持 vector/matrix 两种类型） |

**Agent 安装进度常量**（`InstallProgress`）：0(开始) → 10(编译中) → 30(编译完成) → 50(传输中) → 70(传输完成) → 90(配置完成) → 100(启动成功)。

### 对外暴露的主要接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/monitor/agent/heartbeat` | **无需鉴权** Agent 心跳上报（Token 方式认证） |
| POST | `/api/v1/monitor/agent` | 安装 Agent 到指定主机 |
| POST | `/api/v1/monitor/agent/batch` | 批量安装 Agent |
| GET | `/api/v1/monitor/agent/list` | 分页查询 Agent 列表 |
| GET | `/api/v1/monitor/agent/:id` | Agent 详情 |
| DELETE | `/api/v1/monitor/agent/:id` | 卸载 Agent |
| GET | `/api/v1/monitor/host/metrics` | 获取主机实时监控指标（CPU/内存/磁盘/在线状态） |
| GET | `/api/v1/monitor/host/history` | 获取指标历史时序数据（7 维度） |
| GET | `/api/v1/monitor/host/processes` | 获取进程 Top 排行 |
| GET | `/api/v1/monitor/host/ports` | 获取 TCP 端口监听状态 |

### 内部数据流

**Agent 安装流程：**
```
前端 POST /monitor/agent {hostId, version}
  → AgentController.CreateAgent()
    → AgentService.DeployAgent(hostID)
      → 查询 cmdb_host（SSH 参数）
      → 查询 config_ecsauth（SSH 凭据）
      → 在 MySQL 创建 Agent 记录（Status=部署中, Progress=0）
      → goroutine（异步）:
          → SSH 连接目标主机
          → 分阶段更新 InstallProgress（10→30→50→70→90→100）
          → 传输 Agent 二进制文件（SFTP）
          → SSH 执行启动命令（配置 Prometheus 指标暴露 :9100）
          → 成功 → Status=运行中；失败 → Status=部署失败 + ErrorMsg
    → 立即返回 {agentId, status: "部署中"}（异步部署）
```

**主机监控数据查询流程：**
```
前端 GET /monitor/host/history?hostId=1&hours=24
  → MonitorController.GetHostMetricsHistory()
    → MonitorService.QueryPrometheus(promQL, timeRange)
      → HTTP GET prometheus:9090/api/v1/query_range
        PromQL: rate(node_cpu_seconds_total{instance="ip:9100"}[5m])
      → 解析 PrometheusQueryResult（matrix 类型）
      → 转换为 AllMetricsHistory（7 维度时序数组）
    → result.Success(ctx, metricsHistory)
```

### 与其他模块的依赖关系

- **依赖 CMDB 模块**：Agent 关联 `cmdb_host.id`，安装时需要主机 SSH 参数。
- **依赖 configcenter 模块**：SSH 连接需要 EcsAuth 凭据。
- **外部依赖 Prometheus + Pushgateway**：监控数据的采集和查询链路（Agent → Pushgateway → Prometheus → 后端 → 前端）。
- **被 dashboard 模块依赖**：仪表盘主机在线状态统计通过 Agent 心跳时间判断。

---

## 六、配置中心模块

### 模块职责

配置中心模块是全平台凭据管理的安全仓库，统一存储三类敏感凭据：主机 SSH 认证凭据（密码/私钥/免密）、云厂商 AK/SK 密钥、通用系统账号（数据库/Jenkins/其他工具账密）；所有密码字段使用 AES 加密存储，仅在实际使用时解密传递。此外，提供云资源定时同步调度（定期从云厂商拉取主机列表）。

### 目录结构说明

```
api/api/configcenter/
├── controller/
│   ├── accountAuth.go   # 通用账号（数据库/Jenkins等）CRUD
│   ├── ecsAuth.go       # 主机 SSH 凭据 CRUD（密码/私钥/免密三种类型）
│   ├── keyManage.go     # 云厂商 AK/SK 密钥 CRUD
│   └── syncSchedule.go  # 云资源同步调度配置
├── service/             # 同名文件，业务逻辑
├── dao/                 # 同名文件，GORM 操作
└── model/
    ├── accountAuth.go   # AccountAuth（通用账号）+ AES 加解密方法
    ├── ecsAuth.go       # EcsAuth（SSH 凭据）+ 三种创建 DTO
    ├── keyManage.go     # KeyManage（云密钥）
    └── syncSchedule.go  # SyncSchedule（同步调度配置）
```

### 核心数据模型

| Struct | 表名 | 说明 |
|--------|------|------|
| `EcsAuth` | `config_ecsauth` | SSH 凭据：Name/Type（1=密码/2=私钥/3=免认证）/Username/Password（加密）/PublicKey（私钥内容，字段名历史原因）/Port |
| `AccountAuth` | `config_account` | 通用系统账号：Alias/Host/Port/Name/Password（加密）/Type（账号类型，如 4=Jenkins）/Remark |
| `KeyManage` | `config_key` | 云厂商密钥：Vendor/AccessKey/SecretKey（加密）/Region/Alias |
| `SyncSchedule` | `config_sync_schedule` | 云资源同步调度：Vendor/Region/KeyID/CronExpr/Status |

**安全设计**：
- `EcsAuth.Password` 和 `AccountAuth.Password` 字段均调用 `util.AESEncrypt()` 加密入库。
- 使用时通过 `DecryptPassword()` 解密，明文凭据仅在内存中短暂存在，不通过 API 返回。
- SSH 凭据支持三种认证方式：密码（Type=1）、私钥文件（Type=2，存储私钥内容）、免认证/证书（Type=3）。

### 对外暴露的主要接口

| 分类 | 主要接口 |
|------|---------|
| SSH 凭据 | `GET/POST/PUT/DELETE /api/v1/ecsauth/*`；支持密码/私钥/免密三种类型创建 |
| 通用账号 | `GET/POST/PUT/DELETE /api/v1/accountauth/*`；按类型（Type）查询（如获取所有 Jenkins 服务器） |
| 云密钥 | `GET/POST/PUT/DELETE /api/v1/keymanage/*` |
| 同步调度 | `GET/POST/PUT/DELETE /api/v1/syncchedule/*`；手动触发同步 |

### 内部数据流

**SSH 凭据使用流程（被 CMDB/Task/Monitor 调用）：**
```
调用方（如 CmdbHostSSHService）:
  → EcsAuthDAO.GetById(sshKeyId)           # 查询 config_ecsauth
  → ecsAuth.DecryptPassword()              # AES 解密密码
    OR
  → 读取 ecsAuth.PublicKey                # 私钥内容（Type=2 时）
  → ssh.Dial(ip, auth)                    # 使用解密后的凭据建立 SSH 连接
  → [凭据明文在内存中使用完毕后自动 GC，不持久化也不返回前端]
```

### 与其他模块的依赖关系

- **被 CMDB 模块依赖**（最多）：主机 SSH 终端连接读取 `EcsAuth`；数据库 SQL 执行读取 `AccountAuth`。
- **被 task 模块依赖**：任务 SSH 执行脚本读取 `EcsAuth`。
- **被 monitor 模块依赖**：Agent 安装 SSH 连接读取 `EcsAuth`。
- **被 app 模块依赖**：Jenkins 服务器信息存储在 `AccountAuth`（Type=4）。
- **被 scheduler 依赖**：云资源同步定时任务从 `KeyManage` 读取 AK/SK，调用云厂商 SDK。

---

## 七、用户权限模块

### 模块职责

用户权限模块是平台的安全基础，实现基于 RBAC（Role-Based Access Control）的完整权限体系：用户管理（注册/状态/密码重置）、角色管理（角色与菜单权限绑定）、菜单管理（支持目录/菜单/按钮三级，含路由 URL 和权限标识）、部门管理（组织架构树）、岗位管理；同时提供登录（含图形验证码）、JWT 令牌签发、LDAP 企业认证集成、登录日志和操作日志审计。

### 目录结构说明

```
api/api/system/
├── controller/
│   ├── controller.go       # 通用控制器基类（登录接口在此）
│   ├── sysAdmin.go         # 用户 CRUD + 状态/密码管理 + 个人信息
│   ├── sysRole.go          # 角色 CRUD + 角色-菜单绑定
│   ├── sysMenu.go          # 菜单 CRUD + 左侧菜单树获取
│   ├── sysDept.go          # 部门 CRUD（树形结构）
│   ├── sysPost.go          # 岗位 CRUD
│   ├── sysLogininfo.go     # 登录日志查询
│   ├── sysOperationLog.go  # 操作日志查询
│   └── upload.go           # 头像等文件上传
├── service/
│   ├── sysAdmin.go         # 登录验证（本地+LDAP）、Token 签发、密码加密
│   ├── sysRole.go          # 角色-菜单权限绑定
│   ├── sysMenu.go          # 菜单树构建、用户菜单权限查询
│   ├── captcha.go          # 图形验证码生成（base64Captcha）
│   └── ...
├── dao/                    # 同名文件，GORM 操作
└── model/
    ├── sysAdmin.go         # SysAdmin + LoginDto + JwtAdmin + 各 DTO
    ├── sysRole.go          # SysRole + RoleMenu（角色-菜单绑定）
    ├── sysMenu.go          # SysMenu（三级树） + LeftMenuVo + ValueVo（权限值）
    ├── sysAdmin_Role.go    # SysAdminRole（用户-角色关联）
    ├── sysRole_Menu.go     # SysRoleMenu（角色-菜单关联）
    ├── sysDept.go          # SysDept（部门树）
    ├── sysLogininfo.go     # SysLogininfo（登录日志）
    └── sysOperationLog.go  # SysOperationLog（操作日志）
```

### 核心数据模型

| Struct | 表名 | 说明 |
|--------|------|------|
| `SysAdmin` | `sys_admin` | 用户：Username/Password（bcrypt加密）/Nickname/DeptId/PostId/Status（1=启用/2=禁用） |
| `SysRole` | `sys_role` | 角色：RoleName/RoleKey（权限字符串）/Status |
| `SysMenu` | `sys_menu` | 菜单：ParentId（支持树形嵌套）/MenuType（1=目录/2=菜单/3=按钮）/Url（路由）/Value（权限标识）/Sort |
| `SysAdminRole` | `sys_admin_role` | 用户-角色多对一关联 |
| `SysRoleMenu` | `sys_role_menu` | 角色-菜单多对多关联 |
| `SysDept` | `sys_dept` | 部门树（ParentID 自关联） |
| `SysLogininfo` | `sys_logininfo` | 登录日志：Username/IP/Browser/OS/Status/LoginTime |
| `SysOperationLog` | `sys_operation_log` | 操作日志：Username/Method/URL/Params/Status/IP/RequestTime（由 LogMiddleware 写入） |
| `JwtAdmin` | — | JWT 载荷：ID/Username/Nickname/Icon/Email/Phone（不含密码） |

**RBAC 权限链**：`SysAdmin` → `SysAdminRole` → `SysRole` → `SysRoleMenu` → `SysMenu`（叶子节点的 `Value` 即按钮级权限标识，如 `cmdb:host:add`）。

### 对外暴露的主要接口

| 方法 | 路径（前缀 /api/v1） | 说明 |
|------|---------------------|------|
| GET | `/captcha` | **无需鉴权** 获取图形验证码（UUID + base64 图片） |
| POST | `/login` | **无需鉴权** 登录（验证码 + 账密，支持 LDAP） |
| GET/POST/PUT/DELETE | `/system/admin*` | 用户 CRUD + 状态设置 + 密码重置 |
| GET/POST/PUT/DELETE | `/system/role*` | 角色 CRUD |
| POST | `/system/role/menu` | 角色绑定菜单权限 |
| GET/POST/PUT/DELETE | `/system/menu*` | 菜单 CRUD；`GET /system/nav` 返回当前用户左侧菜单树 |
| GET | `/system/permission` | 返回当前用户的按钮级权限标识列表 |
| GET/POST/PUT/DELETE | `/system/dept*` | 部门 CRUD |
| GET/POST/PUT/DELETE | `/system/post*` | 岗位 CRUD |
| GET | `/system/logininfo` | 登录日志分页查询 |
| GET | `/system/operationlog` | 操作日志分页查询 |
| GET/PUT | `/system/personal` | 个人信息查看/修改 |

### 内部数据流

**登录流程（含验证码和 LDAP）：**
```
前端 GET /captcha
  → CaptchaService.GenerateCaptcha()
    → base64Captcha 生成图片
    → 将验证码答案存入 Redis（key=uuid, ttl=5min）
    → 返回 {uuid, base64图片}

前端 POST /login {username, password, image, idKey}
  → SysAdminController.Login()
    → 从 Redis 验证 image == captcha[idKey]（验证码）
    → SysAdminService.Login()
      ┌─ 本地认证：查询 sys_admin，bcrypt.Compare(password, hash)
      └─ LDAP认证：go-ldap Bind(username, password) 验证
    → 生成 JwtAdmin 并签发 JWT Token（dgrijalva/jwt-go）
    → 将 Token 存入 Redis（支持多端登录管理）
    → 写入登录日志 sys_logininfo（IP/Browser/OS/状态）
    → 返回 {token, userInfo, menus, permissions}
```

### 与其他模块的依赖关系

- **被所有模块依赖**（隐式）：所有业务接口均经过 `AuthMiddleware` → `jwt.ValidateToken()` 验证，间接依赖 system 模块颁发的 JWT。
- **被所有模块依赖**（操作日志）：`LogMiddleware` 将每次操作写入 `sys_operation_log`，该表属于 system 模块。
- **依赖 Redis**：验证码存储、Token 缓存（多端登录控制）。

---

## 八、仪表盘模块

### 模块职责

仪表盘模块是平台的数据汇总层，从 CMDB、K8s、App、Task 各模块数据库中聚合统计数据，为运维人员提供全局概览：主机在线状态统计、K8s 集群健康度、应用发布成功率、任务执行统计、服务/业务线分布、数据库资产按类型分布，无复杂业务逻辑，以聚合查询为主。

### 目录结构说明

```
api/api/dashboard/
├── controller/
│   └── dashboard.go   # 单个控制器，处理所有统计查询
├── service/
│   └── dashboard.go   # 跨模块聚合查询逻辑
└── model1/            # 注意：目录名为 model1（避免与其他模块 model 包冲突）
    └── dashboard.go   # 所有统计结构体
```

### 核心数据模型

| Struct | 说明 |
|--------|------|
| `DashboardStats` | 总览统计：HostStats + K8sClusterStats + DeploymentStats + TaskStats + ServiceStats + DatabaseStats |
| `HostStats` | 主机统计：Total/Online/Offline（通过 Agent.LastHeartbeat 判断在线） |
| `K8sClusterStats` | 集群统计：Total/Healthy/Offline（通过 KubeCluster.Status 判断） |
| `DeploymentStats` | 发布统计：Total/Success/Failed/SuccessRate（查 quick_deployments 表） |
| `TaskStats` | 任务统计：Total/Success/Failed/SuccessRate（查 task_job 表） |
| `ServiceStats` | 服务统计：Total/BusinessLines（查 app_application 表） |
| `DatabaseStats` | 数据库统计：Total/ByType（按 cmdb_sql.type 分组计数） |
| `BusinessDistributionStats` | 业务线分布：各 BusinessGroup 的服务数量及占比 |
| `AssetStats` | 资产分布：主机（自建/阿里云/腾讯云）+ 数据库（各类型）+ K8s 分类统计 |

### 对外暴露的主要接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/dashboard/stats` | 全局统计数据（主机/集群/发布/任务汇总） |
| GET | `/api/v1/dashboard/business` | 业务线服务分布统计 |
| GET | `/api/v1/dashboard/assets` | 资产分类统计（按云厂商/数据库类型） |

### 内部数据流

```
前端 GET /dashboard/stats
  → DashboardController.GetStats()
    → DashboardService.GetDashboardStats()
      → 并发查询（goroutine）：
        ├─ MySQL COUNT(cmdb_host)                 → HostStats.Total
        ├─ Redis/MySQL monitor_agent last_heartbeat → HostStats.Online/Offline
        ├─ MySQL COUNT(k8s_cluster) BY status     → K8sClusterStats
        ├─ MySQL COUNT/GROUP(quick_deployments)   → DeploymentStats
        ├─ MySQL COUNT/GROUP(task_job)            → TaskStats
        ├─ MySQL COUNT(app_application)           → ServiceStats
        └─ MySQL COUNT(cmdb_sql) GROUP BY type    → DatabaseStats
      → 组装 DashboardStats
    → result.Success(ctx, stats)
```

### 与其他模块的依赖关系

- **依赖 CMDB 模块**：查询 `cmdb_host`（主机统计）、`cmdb_sql`（数据库统计）。
- **依赖 K8s 模块**：查询 `k8s_cluster`（集群统计）。
- **依赖 app 模块**：查询 `quick_deployments`（发布统计）、`app_application`（服务统计）。
- **依赖 task 模块**：查询 `task_job`（任务统计）。
- **依赖 monitor 模块**：查询 `monitor_agent.last_heartbeat`（主机在线状态）。
- 单向依赖，不被其他模块依赖。

---

## 九、运维工具模块

### 模块职责

运维工具模块提供两类功能：**快捷导航**（自定义运维工具链接，如 Grafana/Nightingale/JumpServer 等外部系统入口）和**服务一键部署**（从预定义的服务目录中选择服务版本，通过 SSH 在目标主机上执行安装脚本，支持 Docker 容器化和二进制两种部署方式，并追踪部署进度）。

### 目录结构说明

```
api/api/tool/
├── controller/
│   ├── knowledge.go     # 知识库（Markdown 文档）CRUD
│   ├── serviceDeploy.go # 服务一键部署（列表/创建/状态查询/日志）
│   └── tool.go          # 快捷导航 CRUD
├── service/
│   ├── knowledge.go     # 知识库服务
│   ├── serviceDeploy.go # 部署执行逻辑（SSH + Docker/二进制安装脚本）
│   └── tool.go          # 导航工具服务
├── dao1/                # 注意：目录名为 dao1
└── model/
    ├── knowledge.go     # Knowledge（知识库条目）
    ├── tool.go          # Tool（导航链接）+ ServiceDeploy（部署记录）+ ServiceInfo（服务目录定义）
```

### 核心数据模型

| Struct | 表名 | 说明 |
|--------|------|------|
| `Tool` | `tool_link` | 导航链接：Title/Icon/Link（URL）/Sort |
| `ServiceDeploy` | `tool_service_deploy` | 部署记录：ServiceName/ServiceID/Version/HostID/HostIP/InstallDir/ContainerName/Ports/EnvVars/Status（0=部署中→1=运行→2=已停止→3=失败）/DeployLog |
| `ServiceInfo` | 来自 `services.json` 文件 | 服务目录定义：ID/Name/Category/Versions（含 DeployType=container/binary）/DefaultPort/EnvVars/MinMemory/MinCPU |
| `ServiceVersion` | 嵌套在 ServiceInfo | 版本信息：ID/Name/Stable/Recommended/DeployType |

**服务目录**通过读取磁盘上的 `services.json` 配置文件（不存数据库），新增可部署服务只需更新 JSON 文件，无需修改代码。

### 对外暴露的主要接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET/POST/PUT/DELETE | `/api/v1/tool/link*` | 快捷导航的增删改查 |
| GET | `/api/v1/tool/services` | 获取可部署服务目录（来自 services.json） |
| POST | `/api/v1/tool/deploy` | 创建并触发服务部署 |
| GET | `/api/v1/tool/deploy/list` | 部署记录列表 |
| GET | `/api/v1/tool/deploy/:id/log` | 查询部署日志（SSE） |
| GET/POST/PUT/DELETE | `/api/v1/tool/knowledge*` | 知识库文章管理 |

### 内部数据流

**服务一键部署流程：**
```
前端 POST /tool/deploy {serviceId, version, hostId, installDir, envVars}
  → ServiceDeployController.CreateDeploy()
    → ServiceDeployService.Deploy()
      → 从 services.json 获取服务定义（ServiceInfo）
      → 创建 ServiceDeploy 记录（Status=部署中）
      → goroutine（异步）:
          ┌─ container 类型：
          │   → SSH 连接主机
          │   → 执行 Docker pull + run 命令（带端口映射/环境变量）
          └─ binary 类型：
              → SFTP 上传二进制文件
              → SSH 执行安装脚本
          → 实时追加 DeployLog
          → 成功 → Status=1；失败 → Status=3 + 错误信息
    → 立即返回 {deployId, status: "部署中"}
```

### 与其他模块的依赖关系

- **依赖 CMDB 模块**：查询 `cmdb_host` 获取部署目标主机信息。
- **依赖 configcenter 模块**：获取 `EcsAuth` SSH 凭据建立连接。
- 快捷导航功能独立，不依赖其他业务模块。

---

## 十、前端视图模块

### 整体说明

前端采用 Vue3 + Vue Router（SPA）+ Vuex 架构，`web/src/views/` 下的子目录与后端 API 模块一一对应，每个业务页面通过 `web/src/api/` 下对应的 JS 文件调用后端接口。

### 各视图模块详解

#### 1. `views/dashboard/`
**职责**：运维总览页，展示主机/集群/发布/任务统计数字卡片，用 **ECharts** 绘制业务分布饼图和发布趋势折线图。
**交互**：调用 `api/dashboard.js` → 后端 `/dashboard/*`。

#### 2. `views/cmdb/`
**职责**：资产管理中心，包含以下子页面：

| 组件文件 | 说明 |
|---------|------|
| `cmdbHost.vue` | 主机列表（分组树 + 主机表格），支持分组过滤、按 IP/名称搜索 |
| `cmdbGroup.vue` | 分组树管理（新增/编辑/删除分组节点） |
| `cmdbDB.vue` | 数据库资产列表（5 种数据库图标区分） |
| `DBdetails.vue` | 数据库详情 + SQL 在线执行（Monaco 代码编辑器）+ SQL 操作日志 |
| `Host/HostSsh.vue` | **xterm.js** Web SSH 终端（建立 WebSocket 连接 `/cmdb/hostssh/connect/:id`） |
| `Host/Terminal.vue` | 终端组件封装（含窗口自适应 fit addon） |
| `Host/SSH.vue` | SSH 快捷命令/文件上传操作面板 |
| `Host/MonitorDialog.vue` | 主机监控弹窗（CPU/内存/磁盘折线图，调用 monitor 接口） |
| `Host/CreateHost.vue` / `EditHost.vue` | 手动添加/编辑主机表单 |
| `Host/CreateCloud.vue` | 云厂商批量导入（选择厂商/区域/AK/SK） |
| `Host/CreateExcel.vue` | Excel 批量导入向导 |

**关键技术**：xterm.js 与 WebSocket 结合，`addon-attach` 将 WS 消息直接桥接到终端；`addon-fit` 实现终端尺寸自适应并向服务端发送 resize 消息。

#### 3. `views/K8s/`
**职责**：K8s 多集群管理门户。

| 文件 | 说明 |
|------|------|
| `k8s-clusters.vue` | 集群列表，状态指示灯（运行中/离线），进入集群详情 |
| `k8s-workloads.vue` | 工作负载管理（Deployment/StatefulSet/DaemonSet 切换Tab，副本扩缩容滑块） |
| `k8s-nodes.vue` | 节点列表（CPU/内存使用率进度条，污点标签管理） |
| `k8s-namespace.vue` | Namespace 管理（ResourceQuota/LimitRange 配置） |
| `k8s-config.vue` | ConfigMap/Secret 管理（Monaco YAML 编辑器直接编辑） |
| `k8s-storage.vue` | PV/PVC/StorageClass 管理 |
| `k8s-network.vue` | Service/Ingress 管理（路由规则可视化） |
| `clusters/` / `namespaces/` / `nodes/` / `pods/` | 各资源类型的详情页组件（含 Events/Logs/Terminal） |

**关键技术**：K8s YAML 编辑使用 `@guolao/vue-monaco-editor` + `js-yaml` 库实时解析校验；`@xterm/xterm` 实现 Pod Web Terminal。

#### 4. `views/app/`
**职责**：应用发布中心。

| 文件 | 说明 |
|------|------|
| `application.vue` | 应用台账管理（CRUD，关联主机/数据库/域名/负责人） |
| `app-list.vue` | 应用+环境矩阵视图（行=应用，列=dev/test/prod，格=最新发布状态） |
| `app_quick_release.vue` | 快速发布面板（选择应用+环境，拖拽排序，一键触发批量发布） |
| `app_quick_temp.vue` | 快速发布模板管理 |

**关键技术**：快速发布使用 `@antv/x6` 绘制发布流程拓扑图，展示各应用发布进度（节点颜色变化表示状态）；`@logicflow/core` 用于工单流程图绘制。

#### 5. `views/task/`
**职责**：任务调度中心。

| 文件 | 说明 |
|------|------|
| `TaskTemplate.vue` | 脚本模板库（Monaco 代码编辑器，支持 Shell/Python 语法高亮） |
| `TaskJob.vue` | 定时任务列表（Cron 表达式配置，暂停/恢复/立即执行） |
| `TaskConfig.vue` | 任务配置（关联模板+主机+执行计划） |
| `TaskAnsible.vue` | Ansible 任务管理（Playbook 编辑 + 执行历史） |
| `Job/AnsibleTaskHistory.vue` | Ansible 执行历史（SSE 实时日志流，通过 `utils/sseLogManager.js` 管理） |

**关键技术**：日志实时推送通过 `EventSource`（SSE）或 WebSocket 两路并存，`utils/sseLogManager.js` 封装 SSE 连接管理（重连/心跳/多路复用）。

#### 6. `views/monitor/`
**职责**：监控告警中心。

| 文件 | 说明 |
|------|------|
| `base.vue` | 主机基础监控（CPU/内存/磁盘/网络历史折线图，ECharts 多系列）  |
| `https.vue` | 域名 HTTPS 证书监控（到期天数/状态） |
| `Alarm-rules.vue` | 告警规则配置 |
| `Alarm-notify.vue` | 告警通知渠道配置 |
| `alarm-history.vue` | 告警历史记录 |
| `DBLog.vue` | 数据库操作审计日志查看 |
| `LoginLog.vue` | 登录日志审计 |
| `Operator.vue` | 操作日志审计 |

**关键技术**：ECharts 时序折线图（多系列、时间轴联动），数据从 Prometheus 时序接口获取后格式化为 ECharts 的 `dataset` 格式。

#### 7. `views/system/`
**职责**：系统管理后台。

| 文件 | 说明 |
|------|------|
| `Admin.vue` | 用户管理（搜索/新增/编辑/启禁用/重置密码） |
| `Role.vue` | 角色管理 + 菜单权限树勾选 |
| `Menu.vue` | 菜单管理（目录/菜单/按钮三级树，拖拽排序） |
| `Dept.vue` | 部门树管理 |
| `Post.vue` | 岗位管理 |
| `Personal.vue` | 个人信息/头像/密码修改 |

#### 8. `views/configcenter/`
**职责**：凭据管理，支持 SSH 凭据（密码/私钥/免密）、云厂商 AK/SK、通用账号的增删改查，密码字段前端显示为掩码。

#### 9. `views/Tools/`
**职责**：快捷导航配置（添加/排序工具链接卡片）+ 服务部署向导（选服务→选版本→选主机→填配置→查看安装进度日志）。

#### 10. `views/work/`
**职责**：运营工单（发起/处理/关闭），工单状态流转可视化（logicflow 流程图）。

### 前端公共层说明

| 目录/文件 | 说明 |
|---------|------|
| `utils/request.js` | Axios 实例：统一 baseURL、请求拦截器（注入 `Authorization: Bearer token`）、响应拦截器（401 自动跳转登录、错误提示） |
| `utils/authority.js` | 按钮级权限工具函数：`hasPermission(value)` 判断当前用户权限标识列表中是否包含指定权限 |
| `utils/sseLogManager.js` | SSE 连接管理器：维护多路 SSE 连接生命周期（创建/重连/关闭/心跳） |
| `permission/index.js` | 路由守卫：`beforeEach` 检查 localStorage Token，无 Token 跳转 `/login`；有 Token 但 Vuex 中无菜单则重新拉取用户信息和菜单 |
| `store/index.js` | Vuex Store：`user`（token/userInfo）、`menus`（动态菜单树）、`permissions`（按钮权限列表） |
| `store/mutations.js` | Vuex Mutations：`SET_TOKEN` / `SET_USER` / `SET_MENUS` / `SET_PERMISSIONS` |
| `components/CodeEditor.vue` | Monaco 编辑器封装（语言/主题可配置） |
| `components/HostSelector.vue` | 主机选择器（分组树 + 主机 Checkbox，用于任务/部署选择目标主机） |
| `components/Authority.js` | 自定义指令 `v-auth` 实现按钮级权限控制（无权限则隐藏/禁用按钮） |

---

*文档基于源码静态分析生成。如代码有重大变更，请同步更新本文档对应章节。*
